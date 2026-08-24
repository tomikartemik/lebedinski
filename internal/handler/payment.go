package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

const maxPaymentNotificationBytes = 1 << 20

type paymentNotification struct {
	Event  string `json:"event"`
	Object struct {
		ID          string `json:"id"`
		Status      string `json:"status"`
		Description string `json:"description"`
		Amount      struct {
			Value string `json:"value"`
		} `json:"amount"`
	} `json:"object"`
}

func (h *Handler) HandleWebhook(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxPaymentNotificationBytes)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("payment webhook body error:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	var notification paymentNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Println("payment webhook JSON error:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if notification.Event != "payment.succeeded" || notification.Object.Status != "succeeded" {
		c.Status(http.StatusOK)
		return
	}
	if notification.Object.ID == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	// The webhook endpoint is public, so verify the payment directly with
	// YooKassa and use only the authoritative response below.
	payment, err := h.services.GetPayment(notification.Object.ID)
	if err != nil {
		log.Printf("payment webhook verification failed for %s: %v", notification.Object.ID, err)
		c.Status(http.StatusBadGateway)
		return
	}
	if payment.ID != notification.Object.ID || payment.Status != "succeeded" || !payment.Paid {
		log.Printf("payment webhook rejected: provider did not confirm payment %s as succeeded", notification.Object.ID)
		c.Status(http.StatusBadRequest)
		return
	}

	order, err := h.services.GetOrderByPaymentID(payment.ID)
	if err != nil {
		log.Printf("payment webhook has no local order for payment %s: %v", payment.ID, err)
		c.Status(http.StatusInternalServerError)
		return
	}
	if payment.Description != strconv.Itoa(order.CartID) {
		log.Printf("payment webhook order mismatch for payment %s", payment.ID)
		c.Status(http.StatusBadRequest)
		return
	}

	if _, err := h.services.MarkPaymentSucceeded(payment.ID); err != nil {
		log.Printf("failed to persist successful payment %s: %v", payment.ID, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	claimed, err := h.services.ClaimFulfillment(payment.ID)
	if err != nil {
		log.Printf("failed to claim fulfillment for payment %s: %v", payment.ID, err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Payment truth is now durable. CDEK and email processing continue outside
	// the webhook response so a pickup-point problem cannot cause a 502 or turn
	// a captured payment back into Not Paid.
	if claimed {
		go h.processPaidOrder(payment.ID, payment.Amount.Value)
	}
	c.Status(http.StatusOK)
}

func (h *Handler) SendMessageIfFailed(c *gin.Context) {
	h.HandleWebhook(c)
}

func (h *Handler) processPaidOrder(paymentID, amount string) {
	defer func() {
		if recovered := recover(); recovered != nil {
			message := "Внутренняя ошибка при создании отправления СДЭК. Обратитесь к разработчику."
			if err := h.services.SetFulfillmentFailed(paymentID, message); err != nil {
				log.Printf("failed to persist fulfillment panic for payment %s: %v", paymentID, err)
			}
			log.Printf("payment fulfillment panic for %s: %v\n%s", paymentID, recovered, debug.Stack())
		}
	}()

	if _, err := h.services.CreateCdekOrder(paymentID); err != nil {
		message := fulfillmentMessage(err)
		if persistErr := h.services.SetFulfillmentFailed(paymentID, message); persistErr != nil {
			log.Printf("failed to persist fulfillment error for payment %s: %v", paymentID, persistErr)
		}
		log.Printf("CDEK fulfillment failed for payment %s: %v", paymentID, err)
		return
	}
	if err := h.services.SendOrderConfirmation(paymentID, amount); err != nil {
		log.Printf("order confirmation failed for payment %s: %v", paymentID, err)
	}
}

func fulfillmentMessage(err error) string {
	if err == nil {
		return ""
	}
	message := strings.Join(strings.Fields(err.Error()), " ")
	if strings.Contains(message, "appropriate_pickup_point_not_found") {
		return "СДЭК не принял указанный пункт выдачи. Уточните код или адрес ПВЗ через поддержку. Оплата подтверждена."
	}
	if len(message) > 600 {
		message = message[:600]
	}
	return fmt.Sprintf("Не удалось создать отправление СДЭК: %s", message)
}
