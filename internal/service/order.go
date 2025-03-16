package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"strconv"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(orderRepo repository.Order) *OrderService {
	return &OrderService{repo: orderRepo}
}

func (s *OrderService) CreateOrder(order model.Order) (int, error) {
	id, err := s.repo.CreateOrder(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *OrderService) GetOrderByID(idStr string) (model.Order, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return model.Order{}, err
	}

	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
