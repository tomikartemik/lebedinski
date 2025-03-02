package model

type Photo struct {
	ID     int    `gorm:"primaryKey;uniqueIndex;autoIncrement" json:"id"`
	Link   string `gorm:"not null" json:"link"`
	ItemID int    `gorm:"not null" json:"item_id"`
}
