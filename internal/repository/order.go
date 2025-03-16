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

func (r *OrderRepository) CreateOrder(order model.Order) (int, error) {
	err := r.db.Create(&order).Error
	if err != nil {
		return 0, err
	}

	return order.ID, nil
}

func (r *OrderRepository) GetOrderByID(id int) (model.Order, error) {
	var order model.Order

	err := r.db.Preload("OrderItems").First(&order, id).Error
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
