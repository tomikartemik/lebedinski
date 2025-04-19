package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lebedinski/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Static("/uploads", "./uploads")

	banner := router.Group("/banner")
	{
		banner.POST("/upload", h.UploadBanner)
		banner.POST("/upload_mobile", h.UploadMobileBanner)
		banner.GET("", h.GetBanner)
		banner.GET("/mobile", h.GetMobileBanner)
	}

	item := router.Group("/item")
	{
		item.POST("/new", h.CreateItem)
		item.POST("/change-top-item", h.ChangeTopItem)
		item.GET("/all", h.AllItems)
		item.GET("", h.ItemByID)
		item.GET("/top", h.GetTopItems)
		item.PUT("", h.UpdateItem)
		item.DELETE("", h.DeleteItem)
	}

	photo := router.Group("/photo")
	{
		photo.POST("/new", h.UploadPhoto)
		photo.DELETE("", h.DeletePhoto)
	}

	size := router.Group("size")
	{
		size.POST("/add", h.AddNewSizes)
		size.PUT("", h.UpdateSize)
		size.DELETE("", h.DeleteSize)
	}

	category := router.Group("category")
	{
		category.GET("/all", h.GetAllCategorise)
		category.POST("/new", h.AddNewCategory)
		category.PUT("", h.UpdateCategory)
		category.DELETE("", h.DeleteCategory)
	}

	payment := router.Group("payment")
	{
		payment.POST("/response", h.HandleWebhook)
	}

	cart := router.Group("cart")
	{
		cart.POST("/create", h.CreateCart)
		cart.GET("", h.GetCartById)
	}

	order := router.Group("order")
	{
		order.POST("/new", h.CreateOrder)
		order.GET("/all", h.GetAllOrders)
		order.GET("/by-cart-id", h.GetCartById)
	}

	cdek := router.Group("cdek")
	{
		cdek.GET("/pvz", h.GetPvzList)
	}

	return router
}
