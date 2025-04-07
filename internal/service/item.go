package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"lebedinski/internal/utils"
	"os"
	"strconv"
)

type ItemService struct {
	repo repository.Item
}

func NewItemService(repo repository.Item) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) CreateItem(item model.Item) (int, error) {
	item.Discount = (1 - int(item.ActualPrice/item.Price)) * 100
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

func (s *ItemService) DeleteItem(itemIDStr string) error {
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		return err
	}

	// Получаем информацию о товаре для удаления фотографий
	item, err := s.repo.GetItemByID(itemID)
	if err != nil {
		return err
	}

	// Удаляем все фотографии товара с сервера
	for _, photo := range item.Photos {
		if err := os.Remove(photo.Link); err != nil {
			return err
		}
	}

	// Удаляем товар и связанные данные из базы
	return s.repo.DeleteItem(itemID)
}
