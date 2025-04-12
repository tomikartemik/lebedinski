package model

type Item struct {
	ID              int      `gorm:"primaryKey;uniqueIndex;autoIncrement;primaryKey" json:"id"`
	Name            string   `gorm:"not null" json:"name"`
	Description     string   `gorm:"not null" json:"description"`
	Details         string   `gorm:"not null" json:"details"`
	Price           int      `gorm:"not null" json:"price"`
	ActualPrice     int      `gorm:"not null" json:"actual_price"`
	Discount        int      `gorm:"not null;default:0" json:"discount"`
	CustomTailoring bool     `json:"custom_tailoring"`
	SoldOut         bool     `json:"sold_out"`
	CategoryID      int      `gorm:"not null" json:"category_id"`
	Category        Category `gorm:"foreignKey:CategoryID" json:"category"`
	Sizes           []Size   `gorm:"foreignKey:ItemID" json:"sizes"`
	Photos          []Photo  `gorm:"foreignKey:ItemID" json:"photos"`
}

type ItemShortInfo struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	ActualPrice int      `json:"actual_price"`
	Discount    int      `json:"discount"`
	SoldOut     bool     `json:"sold_out"`
	CategoryID  int      `json:"category_id"`
	Category    Category `json:"category"`
	Sizes       []Size   `json:"sizes"`
	Photos      []Photo  `json:"photos"`
}

type CartItem struct {
	ID       int    `gorm:"primaryKey;autoIncrement" json:"id"`
	ItemID   int    `gorm:"not null" json:"item_id"`
	Quantity int    `gorm:"not null" json:"quantity"`
	Size     string `gorm:"not null" json:"size"`
	CartID   int    `gorm:"not null" json:"cart_id"`
}
