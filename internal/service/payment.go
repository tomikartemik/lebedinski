package service

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (s *PaymentService) CreatePayment(amount float64, currency, description string) (*model.PaymentResponse, error) {
	paymentRequest := &model.PaymentRequest{
		Amount: model.Amount{
			Value:    formatAmount(amount),
			Currency: currency,
		},
		Description: description,
		RedirectURL: "https://tomikartemik.ru",
		Capture:     true,
	}

	requestBody, err := json.Marshal(paymentRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(os.Getenv("SHOP_ID"), os.Getenv("SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", "unique-key")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to create payment")
	}

	var paymentResponse model.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		return nil, err
	}

	return &paymentResponse, nil
}

func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
