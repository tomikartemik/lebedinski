package service

import (
	"fmt"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
)

type OrderService struct {
	repoItem  repository.Item
	repoOrder repository.Order
	repoSize  repository.Size
}

func NewOrderService(repoItem repository.Item, repoOrder repository.Order, repoSize repository.Size) *OrderService {
	return &OrderService{
		repoItem:  repoItem,
		repoOrder: repoOrder,
		repoSize:  repoSize,
	}
}

func (s *OrderService) ProcessOrder(order model.Order, paymentID string) error {
	order.PaymentID = paymentID
	order.Status = "Paid"

	cartItems, err := s.repoOrder.GetCartItemsByCartID(order.CartID)

	if err != nil {
		return fmt.Errorf("не удалось получить товары для CartID %d: %w", order.CartID, err)
	}

	if len(cartItems) == 0 {
		return fmt.Errorf("корзина с ID %d пуста или не найдена", order.CartID)
	}

	for _, item := range cartItems {
		err := s.repoSize.DecreaseStock(item.ItemID, item.Size, item.Quantity)
		if err != nil {
			return fmt.Errorf("не удалось списать остаток для ItemID %d, Size %s: %w", item.ItemID, item.Size, err)
		}
	}

	err = s.repoOrder.SaveOrder(order)
	if err != nil {
		return err
	}
	return nil
}

func (s *OrderService) GetAllOrders() ([]model.Order, error) {
	return s.repoOrder.GetAllOrders()
}

func (s *OrderService) GetOrderByCartID(id int) (model.Order, error) {
	order, err := s.repoOrder.GetOrderByCartID(id)

	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
