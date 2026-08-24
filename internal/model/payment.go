package model

type PaymentRequest struct {
	Amount        Amount        `json:"amount"`
	Description   string        `json:"description"`
	Capture       bool          `json:"capture"`
	Confirmation  Confirmation  `json:"confirmation"`
	PaymentMethod PaymentMethod `json:"payment_method_data"`
}

type PaymentMethod struct {
	Type string `json:"type"`
}

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type Confirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type PaymentResponse struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	Paid         bool   `json:"paid"`
	Amount       Amount `json:"amount"`
	Description  string `json:"description"`
	CapturedAt   string `json:"captured_at"`
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}
