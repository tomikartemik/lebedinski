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
	if signature == "" {
		log.Println("Missing Webhook-Signature header")
		c.Status(http.StatusBadRequest)
		return
	}

	// 2. Читаем тело ОДИН раз
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	// 3. Проверяем подпись
	mac := hmac.New(sha256.New, []byte(os.Getenv("SECRET_KEY")))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		log.Println("Invalid signature received")
		c.Status(http.StatusForbidden)
		return
	}

	// 4. Парсим тело заново
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

	// 5. Логируем всю информацию
	log.Printf(
		"Webhook received: Event=%s, ID=%s, Status=%s, Amount=%s\n",
		notification.Event,
		notification.Object.ID,
		notification.Object.Status,
		notification.Object.Amount.Value,
	)

	// 6. Всегда возвращаем 200 OK
	c.Status(http.StatusOK)

}
