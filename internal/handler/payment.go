package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func (h *Handler) CreatePayment(c *gin.Context) {
	var request struct {
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentResponse, err := h.services.CreatePayment(request.Amount, request.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"payment_id":  paymentResponse.ID,
		"payment_url": paymentResponse.Confirmation.ConfirmationURL,
	})

}

func (h *Handler) HandleWebhook(c *gin.Context) {
	signature := c.GetHeader("Webhook-Signature")
	body, _ := c.GetRawData()

	// Генерация HMAC-SHA256 подписи
	mac := hmac.New(sha256.New, []byte(os.Getenv("SECRET_KEY")))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	if signature != expectedSignature {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	// Парсинг уведомления
	var notification struct {
		Event  string `json:"event"`
		Object struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"object"`
	}

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	fmt.Println(notification.Object.Status)
	log.Println(notification.Object.Status)

	// Обработка события
	switch notification.Event {
	case "payment.succeeded":
		fmt.Println(notification.Object.Status)
	case "payment.canceled":
		fmt.Println(notification.Object.Status)
	}

	c.Status(http.StatusOK) // Важно вернуть 200 OK!

}
