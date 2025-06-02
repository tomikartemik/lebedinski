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
	PromoCode
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Item:      NewItemService(repos),
		Photo:     NewPhotoService(repos),
		Size:      NewSizeService(repos),
		Category:  NewCategoryService(repos),
		Payment:   NewPaymentService(repos.Item, repos.Cart, repos.PromoCode),
		Cart:      NewCartService(repos.Cart, repos.Item),
		Cdek:      NewCdekService(repos.Item, repos.Order),
		Order:     NewOrderService(repos.Item, repos.Order, repos.Size, repos.Cart, repos.PromoCode),
		PromoCode: NewPromoCodeService(repos.PromoCode),
	}
}

type Item interface {
	CreateItem(item model.Item) (int, error)
	GetAllItems() ([]model.ItemShortInfo, error)
	GetItemByID(id string) (model.Item, error)
	UpdateItem(itemIDStr string, updateData map[string]interface{}) error
	DeleteItem(itemIDStr string) error
	GetTopItems() ([]model.ItemShortInfo, error)
	ChangeTopItem(position, itemID int) error
}

type Photo interface {
	SavePhoto(itemIDStr string, file *multipart.FileHeader) error
	DeletePhoto(photoIDStr string) error
}

type Size interface {
	AddNewSizes(sizes []model.Size) error
	UpdateSize(sizeIDStr string, updateData map[string]interface{}) error
	DeleteSize(sizeIDStr string) error
}

type Category interface {
	AddCategory(category model.Category) error
	GetAllCategories() ([]model.Category, error)
	UpdateCategory(category model.Category) error
	DeleteCategory(categoryID string) error
}

type Order interface {
	ProcessOrder(order model.Order, paymentID string) error
	GetAllOrders() ([]model.Order, error)
	GetOrderByCartID(id int) (model.Order, error)
	SendOrderConfirmation(cartIDStr, total string) error
	SendOrderShippedNotification(cartIDStr string) error
	DeleteOrder(cartID int) error
	UpdateOrder(order model.Order) error
	ChangeStatus(orderID int, status string) error
}

type Payment interface {
	CreatePayment(order model.Order) (*model.PaymentResponse, error)
}

type Cart interface {
	CreateValidCart(items []model.CartItem) (int, error)
	GetCartByID(id int) (model.Cart, error)
}

type Cdek interface {
	GetToken() (string, error)
	CreateCdekOrder(cartIDStr string) (string, error)
	GetPvzList(params map[string]string) ([]model.Pvz, error)
}

type PromoCode interface {
	CreatePromoCode(promocode model.PromoCode) error
	GetPromoCodeByCode(code string) (model.PromoCode, error)
	GetAllPromoCodes() ([]model.PromoCode, error)
	DeletePromoCodeByCode(code string) error
	UpdatePromoCode(promocode model.PromoCode) error
}
