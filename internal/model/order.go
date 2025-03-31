package model

type Order struct {
	CartID         int        `gorm:"unique;not null" json:"cart_id"`
	FullName       string     `gorm:"not null" json:"full_name"`
	Email          string     `gorm:"not null" json:"email"`
	Phone          string     `gorm:"not null" json:"phone"`
	AdditionalInfo string     `json:"additional_info"`
	PointCode      string     `gorm:"not null" json:"point_code"`
	Address        string     `gorm:"not null" json:"address"`
	City           string     `gorm:"not null" json:"city"`
	DeliveryCost   float64    `gorm:"not null" json:"delivery_cost"`
	Promocode      string     `json:"promocode"`
	CdekOrderUUID  string     `gorm:"unique"`
	Status         string     `gorm:"default:'created';not null"`
}

type CdekOrderRequest struct {
	Number     string      `json:"number"`
	TariffCode int         `json:"tariff_code"`
	Recipient  struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email,omitempty"`
	} `json:"recipient"`
	DeliveryPoint string      `json:"delivery_point"`
	Packages      []CdekPackage `json:"packages"`
}

type CdekPackage struct {
	Number string            `json:"number"`
	Weight int               `json:"weight"`
	Length int               `json:"length,omitempty"`
	Width  int               `json:"width,omitempty"`
	Height int               `json:"height,omitempty"`
	Items  []CdekPackageItem `json:"items"`
}

type CdekPackageItem struct {
	Name    string      `json:"name"`
	WareKey string      `json:"ware_key"`
	Payment CdekPayment `json:"payment"`
	Cost    float64       `json:"cost"`
	Weight  int           `json:"weight"`
	Amount  int           `json:"amount"`
}

type CdekPayment struct {
	Value float64 `json:"value"`
}
