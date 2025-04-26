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

func (h *Handler) HandleWebhook(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Error reading body:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	mac := hmac.New(sha256.New, []byte(os.Getenv("SECRET_KEY")))
	mac.Write(body)

	var notification struct {
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

	if err := json.Unmarshal(body, &notification); err != nil {
		log.Println("JSON parse error:", err)
		c.Status(http.StatusBadRequest)
		return
	}

	log.Printf(
		"Webhook received: Event=%s, ID=%s, Status=%s, Amount=%s, Description=%s",
		notification.Event,
		notification.Object.ID,
		notification.Object.Status,
		notification.Object.Amount.Value,
		notification.Object.Description,
	)

	if notification.Object.Status == "succeeded" {
		h.services.CreateCdekOrder(notification.Object.Description)
		h.services.SendOrderConfirmation(notification.Object.Description)
	}

	log.Printf(notification.Object.Status)

	c.Status(http.StatusOK)
}
