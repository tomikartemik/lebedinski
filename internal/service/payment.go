package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"net/http"
	"os"
	"strconv"
)

type PaymentService struct {
	repo repository.Order
}

func NewPaymentService(orderRepo repository.Order) *PaymentService {
	return &PaymentService{repo: orderRepo}
}

func (s *PaymentService) CreatePayment(amount float64, description string) (*model.PaymentResponse, error) {
	idempotenceKey := uuid.New().String()

	paymentRequest := model.PaymentRequest{
		Amount: model.Amount{
			Value:    formatAmount(amount),
			Currency: "RUB",
		},
		Description: description,
		Capture:     true,
		Confirmation: model.Confirmation{
			Type:      "redirect",
			ReturnURL: "https://your-site.com/payment-callback",
		},
		PaymentMethod: model.PaymentMethod{
			Type: "bank_card",
		},
	}

	requestBody, err := json.Marshal(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.yookassa.ru/v3/payments",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	req.SetBasicAuth(os.Getenv("SHOP_ID"), os.Getenv("SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", idempotenceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var paymentResponse model.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		return nil, fmt.Errorf("response parsing failed: %w", err)
	}

	return &paymentResponse, nil
}

func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
