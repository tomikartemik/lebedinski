package repository

import (
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) AddCategory(category model.Category) error {
	return r.db.Create(&category).Error
}

func (r *CategoryRepository) GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Find(&categories).Error
	return categories, err
}
