package service

import (
	"fmt"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
)

type CartService struct {
	cartRepo repository.Cart
	itemRepo repository.Item
}

func NewCartService(cartRepo repository.Cart, itemRepo repository.Item) *CartService {
	return &CartService{
		cartRepo: cartRepo,
		itemRepo: itemRepo,
	}
}

func (s *CartService) CreateValidCart(items []model.CartItem) (int, error) {
	for _, ci := range items {
		item, err := s.itemRepo.GetItemByID(ci.ItemID)
		if err != nil {
			return 0, fmt.Errorf("item %d not found", ci.ItemID)
		}

		if item.SoldOut {
			return 0, fmt.Errorf("item %d is sold out", ci.ItemID)
		}

		validSize := false
		for _, size := range item.Sizes {
			if size.Size == ci.Size && size.Stock >= ci.Quantity {
				validSize = true
				break
			}
		}

		if !validSize {
			return 0, fmt.Errorf("invalid size or quantity for item %d", ci.ItemID)
		}
	}

	cart := &model.Cart{Items: items}
	if err := s.cartRepo.CreateCart(cart); err != nil {
		return 0, err
	}

	return cart.ID, nil
}
