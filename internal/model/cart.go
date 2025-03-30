package model

type Cart struct {
	ID    int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Items []CartItem `gorm:"foreignKey:CartID" json:"items"`
}
