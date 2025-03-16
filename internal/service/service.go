package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"mime/multipart"
)

type Service struct {
	Item
	Photo
	Size
	Category
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Item:     NewItemService(repos),
		Photo:    NewPhotoService(repos),
		Size:     NewSizeService(repos),
		Category: NewCategoryService(repos),
		Order:    NewOrderService(repos),
	}
}

type Item interface {
	CreateItem(item model.Item) (int, error)
	GetAllItems() ([]model.ItemShortInfo, error)
	GetItemByID(id string) (model.Item, error)
	UpdateItem(item model.Item) error
}

type Photo interface {
	SavePhoto(itemIDStr string, file *multipart.FileHeader) error
}

type Size interface {
	AddNewSizes(sizes []model.Size) error
}

type Category interface {
	AddCategory(category model.Category) error
	GetAllCategories() ([]model.Category, error)
}

type Order interface {
	CreateOrder(order model.Order) (int, error)
	GetOrderByID(idStr string) (model.Order, error)
}
