package model

import "time"

type PromoCode struct {
	ID                 int       `gorm:"primaryKey;uniqueIndex;autoIncrement" json:"id"`
	Code               string    `gorm:"unique" json:"code"`
	MinAmount          float64   `json:"min_amount"`
	DiscountPercentage float64   `json:"discount_percentage"`
	MaxDiscount        float64   `json:"maxDiscount"`
	NumberOfUses       int       `json:"numberOfUses"`
	EndDate            time.Time `json:"endDate"`
}
