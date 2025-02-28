package repository

import (
	"gorm.io/gorm"
	"lebedinski/internal/model"
)

type Repository struct {
	Item
	Photo
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Item:  NewItemRepository(db),
		Photo: NewPhotoRepository(db),
	}
}

type Item interface {
	CreateItem(item model.Item) error
	GetAllItems() ([]model.Item, error)
	GetItemByID(id int) (model.Item, error)
	UpdateItem(item model.Item) error
}

type Photo interface {
	NewPhoto(photo model.Photo) error
}
