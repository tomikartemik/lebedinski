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
			Value:    formatAmount(amount), // Форматируем сумму (например, "100.00")
			Currency: currency,
		},
		Description: description,
		RedirectURL: "https://tomikartemik.ru", // Куда перенаправить после оплаты
		Capture:     true,                      // Автоматическое подтверждение платежа
	}

	// Сериализация запроса в JSON
	requestBody, err := json.Marshal(paymentRequest)
	if err != nil {
		return nil, err
	}

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// Установка заголовков
	req.SetBasicAuth(os.Getenv("SHOP_ID"), os.Getenv("SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", "unique-key") // Уникальный ключ для идемпотентности

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Логирование статуса ответа и тела ошибки
	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			return nil, errors.New(errorResponse["description"].(string))
		}
		return nil, errors.New("failed to create payment")
	}

	// Десериализация ответа
	var paymentResponse model.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResponse); err != nil {
		return nil, err
	}

	return &paymentResponse, nil
}

func formatAmount(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 2, 64)
}
