package model

type Category struct {
	ID   int    `gorm:"primaryKey;uniqueIndex;autoIncrement" json:"id"`
	Name string `gorm:"not null" json:"name"`
}
