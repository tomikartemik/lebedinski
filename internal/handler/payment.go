package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

	// Проверяем наличие подписи
	signature := c.GetHeader("Signature")
	if signature == "" {
		log.Println("Missing Signature header")
		c.Status(http.StatusBadRequest)
		return
	}

	// Читаем тело
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	// Проверяем подпись
	mac := hmac.New(sha256.New, []byte(os.Getenv("SECRET_KEY")))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Подпись от ЮКассы приходит в формате "v1 <signature>"
	// Извлекаем только саму подпись
	parts := strings.Split(signature, " ")
	if len(parts) != 3 {
		log.Println("Invalid signature format")
		c.Status(http.StatusBadRequest)
		return
	}
	actualSignature := parts[1]

	if !hmac.Equal([]byte(actualSignature), []byte(expectedSignature)) {
		log.Println("Invalid signature received")
		c.Status(http.StatusForbidden)
		return
	}

	// Парсим тело
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

	// Логируем информацию
	log.Printf(
		"Webhook received: Event=%s, ID=%s, Status=%s, Amount=%s\n",
		notification.Event,
		notification.Object.ID,
		notification.Object.Status,
		notification.Object.Amount.Value,
	)

	// Возвращаем 200 OK
	c.Status(http.StatusOK)

}
