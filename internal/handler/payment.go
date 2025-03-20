package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	var notification struct {
		Event  string `json:"event"`
		Object struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"object"`
	}
	if err := c.ShouldBindJSON(&notification); err != nil {
		fmt.Println("invalid notification")
		return
	}

	switch notification.Event {
	case "payment.succeeded":
		fmt.Println("Payment succeeded")
	case "payment.canceled":
		fmt.Println("Payment canceled")
	default:
		fmt.Println("Unknown event")
	}
}
