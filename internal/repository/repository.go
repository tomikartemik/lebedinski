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
	Cart
	Order
	PromoCode
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Item:      NewItemRepository(db),
		Photo:     NewPhotoRepository(db),
		Size:      NewSizeRepository(db),
		Category:  NewCategoryRepository(db),
		Cart:      NewCartRepository(db),
		Order:     NewOrderRepository(db),
		PromoCode: NewPromoCodeRepository(db),
	}
}

type Item interface {
	CreateItem(item model.Item) (int, error)
	GetAllItems() ([]model.Item, error)
	GetItemByID(id int) (model.Item, error)
	UpdateItem(itemID int, updateData map[string]interface{}) error
	DeleteItem(itemID int) error
	GetTopItems() ([]model.Top, error)
	ChangeTopItem(position, itemID int) error
}

type Photo interface {
	NewPhoto(photo model.Photo) error
	DeletePhoto(photoID int) error
	GetPhotoByID(photoID int) (model.Photo, error)
}

type Size interface {
	AddNewSizes(sizes []model.Size) error
	DecreaseStock(itemID int, size string, quantity int) error
	UpdateSize(sizeID int, updateData map[string]interface{}) error
	DeleteSize(sizeID int) error
}

type Category interface {
	AddCategory(category model.Category) error
	GetAllCategories() ([]model.Category, error)
	UpdateCategory(category model.Category) error
	DeleteCategory(categoryID int) error
}

type Cart interface {
	CreateCart(cart *model.Cart) error
	GetCartByID(cartID int) (model.Cart, error)
}

type Order interface {
	SaveOrder(order model.Order) error
	GetCartItemsByCartID(cartID int) ([]model.CartItem, error)
	GetAllOrders() ([]model.Order, error)
	GetOrderByCartID(id int) (model.Order, error)
	UpdateOrder(order model.Order) error
	DeleteOrder(cartID int) error
	ChangeStatus(orderID int, status string) error
}

type PromoCode interface {
	CreatePromoCode(promocode model.PromoCode) error
	GetPromoCodeByCode(code string) (model.PromoCode, error)
	GetAllPromocodes() ([]model.PromoCode, error)
	DeletePromoCodeByCode(code string) error
	UpdatePromoCode(promocode model.PromoCode) error
}
