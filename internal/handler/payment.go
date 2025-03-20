package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
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
	log.Printf("Request from IP: %s", c.Request.RemoteAddr)
	log.Printf("Headers: %+v", c.Request.Header)

	// Проверка заголовка
	signatureHeader := c.GetHeader("Signature")
	if signatureHeader == "" {
		log.Println("Missing Signature header")
		c.Status(http.StatusBadRequest)
		return
	}

	// Парсинг заголовка
	parts := strings.Split(signatureHeader, " ")
	if len(parts) < 3 || parts[0] != "v1" {
		log.Println("Invalid signature format")
		c.Status(http.StatusBadRequest)
		return
	}
	signature := strings.Join(parts[1:len(parts)-1], " ")
	keyID := parts[len(parts)-1]

	// Логируем keyID
	log.Printf("Key ID: %s", keyID)

	// Чтение тела
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	// Декодирование подписи
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		log.Println("Base64 decode error:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	// Проверка подписи
	mac := hmac.New(sha256.New, []byte(os.Getenv("SECRET_KEY")))
	mac.Write(body)
	expectedSignature := mac.Sum(nil)

	if !hmac.Equal(decodedSignature, expectedSignature) {
		log.Println("Invalid signature")
		c.Status(http.StatusForbidden)
		return
	}

	// Парсинг тела
	var notification struct {
		Event  string `json:"event"`
		Object struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Amount struct {
				Value string `json:"value"`
			} `json:"amount"`
		} `json:"object"`
	}

	if err := json.Unmarshal(body, &notification); err != nil {
		log.Println("JSON parse error:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Printf(
		"Webhook received: Event=%s, ID=%s, Status=%s, Amount=%s",
		notification.Event,
		notification.Object.ID,
		notification.Object.Status,
		notification.Object.Amount.Value,
	)

	c.Status(http.StatusOK)
}
