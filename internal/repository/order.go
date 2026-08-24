package repository

import (
	"gorm.io/gorm"
	"lebedinski/internal/model"
	"strings"
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

func (r *OrderRepository) GetOrderByPaymentID(paymentID string) (model.Order, error) {
	var order model.Order
	if err := r.db.Where("payment_id = ?", paymentID).First(&order).Error; err != nil {
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

// SetStatusByCartID sets an order's status, keyed by its CartID.
func (r *OrderRepository) SetStatusByCartID(cartID int, status string) error {
	return r.db.Model(&model.Order{}).Where("cart_id = ?", cartID).Update("status", status).Error
}

// ClaimOrderForProcessing atomically transitions an order into the "Processing"
// state, but only if it has not already been paid or claimed. It returns true
// only for the single caller that wins the claim; concurrent/duplicate webhook
// deliveries get false and must skip processing. This closes the race where two
// webhooks both read a not-yet-paid order before either marks it done.
func (r *OrderRepository) ClaimOrderForProcessing(cartID int) (bool, error) {
	res := r.db.Model(&model.Order{}).
		Where("cart_id = ? AND status NOT IN ?", cartID, []string{"Processing", "Paid", "Sent"}).
		Update("status", "Processing")
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

// MarkPaymentSucceeded records YooKassa's result on the exact payment row.
// Delivery failures must never move this status back to Not Paid.
func (r *OrderRepository) MarkPaymentSucceeded(paymentID string) (model.Order, error) {
	order, err := r.GetOrderByPaymentID(paymentID)
	if err != nil {
		return order, err
	}

	updates := map[string]interface{}{"payment_status": "Succeeded"}
	switch strings.ToLower(strings.TrimSpace(order.Status)) {
	case "", "created", "pending", "processing", "not paid":
		updates["status"] = "Paid"
	}
	if err := r.db.Model(&model.Order{}).Where("payment_id = ?", paymentID).Updates(updates).Error; err != nil {
		return order, err
	}
	return r.GetOrderByPaymentID(paymentID)
}

// ClaimFulfillment ensures duplicate YooKassa notifications cannot create
// duplicate shipments or repeat stock/email side effects.
func (r *OrderRepository) ClaimFulfillment(paymentID string) (bool, error) {
	res := r.db.Model(&model.Order{}).
		Where("payment_id = ? AND COALESCE(cdek_order_uuid, '') = '' AND fulfillment_status = ?", paymentID, "Pending").
		Updates(map[string]interface{}{
			"fulfillment_status": "Processing",
			"fulfillment_error":  "",
		})
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected > 0, nil
}

func (r *OrderRepository) CompleteFulfillment(paymentID, cdekOrderNumber string) error {
	return r.db.Model(&model.Order{}).Where("payment_id = ?", paymentID).Updates(map[string]interface{}{
		"cdek_order_uuid":    cdekOrderNumber,
		"fulfillment_status": "Ready",
		"fulfillment_error":  "",
	}).Error
}

func (r *OrderRepository) SetFulfillmentFailed(paymentID, message string) error {
	return r.db.Model(&model.Order{}).Where("payment_id = ?", paymentID).Updates(map[string]interface{}{
		"fulfillment_status": "Needs Attention",
		"fulfillment_error":  message,
	}).Error
}
