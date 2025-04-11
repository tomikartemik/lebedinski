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
	if err := r.db.Preload("Category").Preload("Photos").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemRepository) GetItemByID(id int) (model.Item, error) {
	var item model.Item
	if err := r.db.Where("id = ?", id).Preload("Category").Preload("Photos").Preload("Sizes").Find(&item).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *ItemRepository) UpdateItem(itemID int, updateData map[string]interface{}) error {
	return r.db.Model(&model.Item{}).Where("id = ?", itemID).Updates(updateData).Error
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
