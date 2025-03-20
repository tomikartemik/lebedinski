package model

type PaymentRequest struct {
	Amount      Amount `json:"amount"`
	Description string `json:"description"`
	RedirectURL string `json:"confirmation_url"` // Куда перенаправить после оплаты
	Capture     bool   `json:"capture"`          // Автоматическое подтверждение платежа
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
