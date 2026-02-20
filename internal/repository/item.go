package repository

import (
	"fmt"
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) CreateItem(item model.Item) (int, error) {
	err := r.db.Create(&item).Error
	if err != nil {
		return 0, err
	}
	return item.ID, nil
}

func (r *ItemRepository) GetAllItems() ([]model.Item, error) {
	var items []model.Item
	if err := r.db.Preload("Category").Preload("Categories").Preload("Photos").Preload("Sizes").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemRepository) GetItemByID(id int) (model.Item, error) {
	var item model.Item
	if err := r.db.Where("id = ?", id).Preload("Category").Preload("Categories").Preload("Photos").Preload("Sizes").Find(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *ItemRepository) UpdateItemCategories(itemID int, categoryIDs []int) error {
	var item model.Item
	item.ID = itemID

	uniqueCategoryIDs := make(map[int]struct{}, len(categoryIDs))
	normalizedCategoryIDs := make([]int, 0, len(categoryIDs))
	for _, categoryID := range categoryIDs {
		if categoryID <= 0 {
			return fmt.Errorf("category_ids must contain positive ids")
		}
		if _, exists := uniqueCategoryIDs[categoryID]; exists {
			continue
		}
		uniqueCategoryIDs[categoryID] = struct{}{}
		normalizedCategoryIDs = append(normalizedCategoryIDs, categoryID)
	}

	var categories []model.Category
	if len(normalizedCategoryIDs) > 0 {
		if err := r.db.Where("id IN ?", normalizedCategoryIDs).Find(&categories).Error; err != nil {
			return err
		}

		if len(categories) != len(normalizedCategoryIDs) {
			return fmt.Errorf("one or more category_ids do not exist")
		}
	}

	return r.db.Model(&item).Association("Categories").Replace(categories)
}

func (r *ItemRepository) UpdateItem(itemID int, updateData map[string]interface{}) error {
	if rawCategoryID, ok := updateData["category_id"]; ok {
		categoryID, err := parseCategoryID(rawCategoryID)
		if err != nil {
			return err
		}

		var count int64
		if err := r.db.Model(&model.Category{}).Where("id = ?", categoryID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return fmt.Errorf("category_id %d does not exist", categoryID)
		}

		updateData["category_id"] = categoryID
	}

	return r.db.Model(&model.Item{}).Where("id = ?", itemID).Updates(updateData).Error
}

func parseCategoryID(raw interface{}) (int, error) {
	switch v := raw.(type) {
	case int:
		if v <= 0 {
			return 0, fmt.Errorf("category_id must be positive")
		}
		return v, nil
	case float64:
		categoryID := int(v)
		if categoryID <= 0 {
			return 0, fmt.Errorf("category_id must be positive")
		}
		return categoryID, nil
	default:
		return 0, fmt.Errorf("category_id must be a number")
	}
}

func (r *ItemRepository) DeleteItem(itemID int) error {
	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Получаем информацию о товаре для удаления фотографий
	var item model.Item
	if err := tx.Where("id = ?", itemID).Preload("Photos").First(&item).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Удаляем все размеры товара
	if err := tx.Where("item_id = ?", itemID).Delete(&model.Size{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Удаляем все фотографии товара из базы данных
	if err := tx.Where("item_id = ?", itemID).Delete(&model.Photo{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Удаляем связи с категориями
	if err := tx.Model(&item).Association("Categories").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// Удаляем сам товар
	if err := tx.Delete(&item).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Завершаем транзакцию
	return tx.Commit().Error
}

func (r *ItemRepository) GetTopItems() ([]model.Top, error) {
	var tops []model.Top
	if err := r.db.Model(&model.Top{}).Find(&tops).Error; err != nil {
		return nil, err
	}
	return tops, nil
}

func (r *ItemRepository) ChangeTopItem(position, itemID int) error {
	return r.db.Model(&model.Top{}).Where("position = ?", position).Update("item_id", itemID).Error
}
