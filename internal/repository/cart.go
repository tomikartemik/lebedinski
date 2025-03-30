package repository

import (
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) CreateCart(cart *model.Cart) error {
	return r.db.Create(cart).Error
}

func (r *CartRepository) GetCartByID(cartID int) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Preload("Items").First(&cart, cartID).Error
	return &cart, err
}
