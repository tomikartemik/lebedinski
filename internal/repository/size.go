package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type SizeRepository struct {
	db *gorm.DB
}

func NewSizeRepository(db *gorm.DB) *SizeRepository {
	return &SizeRepository{db: db}
}

func (r *SizeRepository) AddNewSizes(sizes []model.Size) error {
	tx := r.db.Create(&sizes)
	return tx.Error
}

func (r *SizeRepository) DecreaseStock(itemID int, size string, quantity int) error {
	if quantity <= 0 {
		return errors.New("количество для списания должно быть положительным")
	}

	result := r.db.Model(&model.Size{}).
		Where("item_id = ? AND size = ? AND stock >= ?", itemID, size, quantity).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity))

	if result.Error != nil {
		return fmt.Errorf("ошибка при обновлении остатка для item_id %d, size %s: %w", itemID, size, result.Error)
	}

	if result.RowsAffected == 0 {
		var count int64
		r.db.Model(&model.Size{}).Where("item_id = ? AND size = ?", itemID, size).Count(&count)
		if count == 0 {
			return fmt.Errorf("запись для item_id %d и size %s не найдена", itemID, size)
		}
		return fmt.Errorf("недостаточно остатка для item_id %d, size %s (нужно %d)", itemID, size, quantity)
	}

	return nil
}

func (r *SizeRepository) UpdateSize(size model.Size) error {
	return r.db.Where("id = ?", size.ID).Updates(&size).Error
}

func (r *SizeRepository) DeleteSize(sizeID int) error {
	return r.db.Where("id = ?", sizeID).Delete(&model.Size{}).Error
}
