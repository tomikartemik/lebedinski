package repository

import (
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) SaveOrder(order model.Order) error {
	return r.db.Create(&order).Error
}

func (r *OrderRepository) GetCartItemsByCartID(cartID int) ([]model.CartItem, error) {
	var cartItems []model.CartItem

	if err := r.db.Where("cart_id = ?", cartID).Find(&cartItems).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}
