package model

import "time"

type Order struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	Telegram  string    `json:"telegram"`
	CreatedAt time.Time `json:"created_at"`
	CartID    int       `gorm:"not null" json:"cart_id"`
	Cart      Cart      `gorm:"foreignKey:CartID" json:"cart"`
}
