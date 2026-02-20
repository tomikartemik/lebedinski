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

func (r *OrderRepository) GetAllOrders() ([]model.Order, error) {
	var orders []model.Order
	if err := r.db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderByCartID(id int) (model.Order, error) {
	var order model.Order
	if err := r.db.Where("cart_id = ?", id).First(&order).Error; err != nil {
		return order, err
	}
	return order, nil
}

func (r *OrderRepository) UpdateOrder(order model.Order) error {
	return r.db.Model(&model.Order{}).
		Where("cart_id = ?", order.CartID).
		Updates(map[string]interface{}{
			"full_name":       order.FullName,
			"email":           order.Email,
			"phone":           order.Phone,
			"additional_info": order.AdditionalInfo,
			"point_code":      order.PointCode,
			"delivery_city":   order.DeliveryCity,
			"promocode":       order.Promocode,
			"status":          order.Status,
			"archive":         order.Archive,
			"marked":          order.Marked,
			"telegram_id":     order.TelegramID,
			"cdek_order_uuid": order.CdekOrderUUID,
		}).Error
}

func (r *OrderRepository) DeleteOrder(cartID int) error {
	return r.db.Delete(&model.Order{}, "cart_id = ?", cartID).Error
}

func (r *OrderRepository) ChangeStatus(orderID int, status string) error {
	return r.db.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status).Error
}
