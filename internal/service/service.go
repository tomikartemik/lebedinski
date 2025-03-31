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
	Payment
	Cart
	Cdek
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Item:     NewItemService(repos),
		Photo:    NewPhotoService(repos),
		Size:     NewSizeService(repos),
		Category: NewCategoryService(repos),
		Payment:  NewPaymentService(repos),
		Cart:     NewCartService(repos.Cart, repos.Item),
		Cdek:     NewCdekService(repos.Item),
		Order:    NewOrderService(repos.Item, repos.Order, repos.Size),
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
	ProcessOrder(order model.Order, cdekUUID string) error
}

type Payment interface {
	CreatePayment(amount float64, description string) (*model.PaymentResponse, error)
}

type Cart interface {
	CreateValidCart(items []model.CartItem) (int, error)
	GetCartByID(id int) (model.Cart, error)
}

type Cdek interface {
	GetToken() (string, error)
	CreateCdekOrder(order model.Order) (string, error)
}
