package model

type Photo struct {
	ID     int    `gorm:"primaryKey;uniqueIndex;autoIncrement;primaryKey" json:"id"`
	Link   string `gorm:"not null" json:"link"`
	ItemID int    `gorm:"not null" json:"item_id"`
}
