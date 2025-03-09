package model

type Size struct {
	ID     int    `gorm:"primaryKey;uniqueIndex;autoIncrement" json:"id"`
	Size   string `json:"size"`
	Stock  int    `json:"stock"`
	ItemID int    `gorm:"not null" json:"item_id"`
}
