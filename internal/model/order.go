package model

type Order struct {
	CartID         int     `gorm:"unique;not null" json:"cart_id"`
	FullName       string  `gorm:"not null" json:"full_name"`
	Email          string  `gorm:"not null" json:"email"`
	Phone          string  `gorm:"not null" json:"phone"`
	AdditionalInfo string  `json:"additional_info"`
	PointCode      string  `gorm:"not null" json:"point_code"`
	Address        string  `gorm:"not null" json:"address"`
	City           string  `gorm:"not null" json:"city"`
	DeliveryCost   float64 `gorm:"not null" json:"delivery_cost"`
	Promocode      string  `json:"promocode"`
	CdekOrderUUID  string  `gorm:"unique"`
	Status         string  `gorm:"default:'created';not null"`
}

type CdekOrderRequest struct {
	Number     string `json:"number"`
	TariffCode int    `json:"tariff_code"`
	Recipient  struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	} `json:"recipient"`
	ToLocation struct {
		Code    string `json:"code"`
		Address string `json:"address"`
		City    string `json:"city"`
		Country string `json:"country"`
	} `json:"to_location"`
	DeliveryCost float64 `json:"delivery_cost"`
}
