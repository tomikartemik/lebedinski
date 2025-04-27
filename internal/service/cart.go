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

		for _, size := range item.Sizes {
			if size.Size == ci.Size && size.Stock < ci.Quantity || item.CustomTailoring != true {
				return 0, fmt.Errorf("itemID: %d size: %s stock: %d", ci.ItemID, ci.Size, size.Stock)
			}
		}

	}

	cart := model.Cart{Items: items}
	if err := s.cartRepo.CreateCart(&cart); err != nil {
		return 0, err
	}

	return cart.ID, nil
}

func (s *CartService) GetCartByID(id int) (model.Cart, error) {
	cart, err := s.cartRepo.GetCartByID(id)
	if err != nil {
		return model.Cart{}, err
	}
	return cart, nil
}
