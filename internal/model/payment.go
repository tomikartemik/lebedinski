package model

type PaymentRequest struct {
	Amount       Amount       `json:"amount"`
	Description  string       `json:"description"`
	Capture      bool         `json:"capture"`
	Confirmation Confirmation `json:"confirmation"` // Добавьте это поле
}

type Confirmation struct {
	Type      string `json:"type"`       // Тип подтверждения (например, "redirect")
	ReturnURL string `json:"return_url"` // Куда перенаправить после оплаты
}

type Amount struct {
	Value    string `json:"value"`    // Сумма платежа
	Currency string `json:"currency"` // Валюта (например, "RUB")
}

type PaymentResponse struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	PaymentURL string `json:"confirmation_url"` // Ссылка для оплаты
}
