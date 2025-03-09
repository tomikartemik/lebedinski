package repository

import (
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
