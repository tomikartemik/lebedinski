package repository

import (
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
	if err := r.db.Preload("Photos").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemRepository) GetItemByID(id int) (model.Item, error) {
	var item model.Item
	if err := r.db.Where("id = ?", id).Preload("Photos").Find(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *ItemRepository) UpdateItem(item model.Item) error {
	return r.db.Save(&item).Error
}
