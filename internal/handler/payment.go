package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

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
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	mac := hmac.New(sha256.New, []byte(os.Getenv("SECRET_KEY")))
	mac.Write(body)

	var notification paymentNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Println("JSON parse error:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if notification.Object.Status == "succeeded" {
		h.processSuccessfulPayment(notification.Object.Description, notification.Object.Amount.Value)
	}

	c.Status(http.StatusOK)
}

func (h *Handler) SendMessageIfFailed(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	var notification paymentNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Println("JSON parse error:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	if notification.Object.Status == "succeeded" {
		h.processSuccessfulPayment(notification.Object.Description, notification.Object.Amount.Value)
	}

	c.Status(http.StatusOK)
}

// processSuccessfulPayment runs the one-time side effects of a paid order: it
// creates the CDEK shipment and, only if this call actually processed the order
// (i.e. was not a duplicate webhook), sends the confirmation email. CreateCdekOrder
// owns the idempotency claim, so the confirmation — which decrements stock and
// promocode uses — must be gated on its `created` result to avoid double-running.
// Both /response and /send-message-if-failed funnel here, and YooKassa retries
// deliveries, so this can be invoked several times for one payment.
func (h *Handler) processSuccessfulPayment(cartID, amount string) {
	_, created, err := h.services.CreateCdekOrder(cartID)
	if err != nil {
		log.Printf("processSuccessfulPayment: CreateCdekOrder failed for cart %s: %v", cartID, err)
		return
	}
	if !created {
		// Duplicate delivery — order already processed. Skip email/stock updates.
		return
	}
	if err := h.services.SendOrderConfirmation(cartID, amount); err != nil {
		log.Printf("processSuccessfulPayment: SendOrderConfirmation failed for cart %s: %v", cartID, err)
	}
}
