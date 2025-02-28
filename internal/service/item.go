package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"lebedinski/internal/utils"
	"strconv"
)

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(item model.Item) error {
	return s.repo.CreateItem(item)
}

func (s *ItemService) GetAllItems() ([]model.ItemShortInfo, error) {
	var itemsShortInfo []model.ItemShortInfo

	items, err := s.repo.GetAllItems()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		itemsShortInfo = append(itemsShortInfo, utils.ConvertItemToShortInfo(item))
	}

	return itemsShortInfo, nil
}

func (s *ItemService) GetItemByID(idStr string) (model.Item, error) {
	id, err := strconv.Atoi(idStr)

	if err != nil {
		return model.Item{}, err
	}

	item, err := s.repo.GetItemByID(id)

	if err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (s *ItemService) UpdateItem(item model.Item) error {
	return s.repo.UpdateItem(item)
}
