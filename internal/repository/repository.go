package repository

import (
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type Repository struct {
	Item
	Photo
	Size
	Category
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Item:     NewItemRepository(db),
		Photo:    NewPhotoRepository(db),
		Size:     NewSizeRepository(db),
		Category: NewCategoryRepository(db),
	}
}

type Item interface {
	CreateItem(item model.Item) (int, error)
	GetAllItems() ([]model.Item, error)
	GetItemByID(id int) (model.Item, error)
	UpdateItem(item model.Item) error
}

type Photo interface {
	NewPhoto(photo model.Photo) error
}

type Size interface {
	AddNewSizes(sizes []model.Size) error
}

type Category interface {
	AddCategory(category model.Category) error
	GetAllCategories() ([]model.Category, error)
}
