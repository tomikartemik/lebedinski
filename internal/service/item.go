package service

import (
	"fmt"
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
		intIDs, err := parseCategoryIDs(categoryIDs)
		if err != nil {
			return err
		}

		if err := s.repo.UpdateItemCategories(itemID, intIDs); err != nil {
			return err
		}

		// Keep the primary category in sync with the many-to-many relation.
		if len(intIDs) > 0 {
			updateData["category_id"] = intIDs[0]
		}
	}

	// Remove categories from update data as they're handled via associations
	delete(updateData, "category")
	delete(updateData, "categories")

	if len(updateData) > 0 {
		return s.repo.UpdateItem(itemID, updateData)
	}
	return nil
}

func parseCategoryIDs(raw interface{}) ([]int, error) {
	switch v := raw.(type) {
	case []int:
		return v, nil
	case []interface{}:
		result := make([]int, 0, len(v))
		for _, id := range v {
			switch typedID := id.(type) {
			case float64:
				result = append(result, int(typedID))
			case int:
				result = append(result, typedID)
			default:
				return nil, fmt.Errorf("category_ids must contain only numbers")
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("category_ids must be an array of numbers")
	}
}

func (s *ItemService) UpdateItemCategories(itemID int, categoryIDs []int) error {
	return s.repo.UpdateItemCategories(itemID, categoryIDs)
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
