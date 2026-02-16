package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"lebedinski/internal/utils"
	"os"
	"sort"
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

func (s *ItemService) UpdateItem(itemIDStr string, updateData map[string]interface{}) error {
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		return err
	}

	// Handle category_ids if present
	if categoryIDs, ok := updateData["category_ids"]; ok {
		delete(updateData, "category_ids")
		if ids, ok := categoryIDs.([]interface{}); ok {
			var intIDs []int
			for _, id := range ids {
				switch v := id.(type) {
				case float64:
					intIDs = append(intIDs, int(v))
				case int:
					intIDs = append(intIDs, v)
				}
			}
			if err := s.repo.UpdateItemCategories(itemID, intIDs); err != nil {
				return err
			}
		}
	}

	// Remove categories from update data as they're handled via associations
	delete(updateData, "categories")

	if len(updateData) > 0 {
		return s.repo.UpdateItem(itemID, updateData)
	}
	return nil
}

func (s *ItemService) DeleteItem(itemIDStr string) error {
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		return err
	}

	item, err := s.repo.GetItemByID(itemID)
	if err != nil {
		return err
	}

	for _, photo := range item.Photos {
		if err := os.Remove(photo.Link); err != nil {
			return err
		}
	}

	return s.repo.DeleteItem(itemID)
}

func (s *ItemService) GetTopItems() ([]model.ItemShortInfo, error) {
	var items []model.ItemShortInfo

	tops, err := s.repo.GetTopItems()
	if err != nil {
		return nil, err
	}

	sort.Slice(tops, func(i, j int) bool {
		return tops[i].Position < tops[j].Position
	})

	for _, top := range tops {
		item, err := s.repo.GetItemByID(top.ItemID)
		if err != nil {
			return nil, err
		}
		items = append(items, utils.ConvertItemToShortInfo(item))
	}

	return items, nil
}

func (s *ItemService) ChangeTopItem(position, itemID int) error {
	return s.repo.ChangeTopItem(position, itemID)
}
