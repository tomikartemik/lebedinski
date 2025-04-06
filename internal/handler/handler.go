package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lebedinski/internal/service"
	"net/http"
	"strings"
	"log"
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

	item := router.Group("/item")
	{
		item.POST("/new", h.CreateItem)
		item.GET("/all", h.AllItems)
		item.GET("", h.ItemByID)
		item.PUT("", h.UpdateItem)
	}

	photo := router.Group("/photo")
	{
		photo.POST("/new", h.UploadPhoto)
	}

	size := router.Group("size")
	{
		size.POST("/add", h.AddNewSizes)
		size.PUT("/", h.UpdateSize)
		size.DELETE("/", h.DeleteSize)
	}

	category := router.Group("category")
	{
		category.GET("/all", h.GetAllCategorise)
		category.POST("/new", h.AddNewCategory)
		category.PUT("/", h.UpdateCategory)
		category.DELETE("/", h.DeleteCategory)
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

func (h *Handler) GetPvzList(c *gin.Context) {
	country := c.Query("country")
	cityCode := c.Query("city_code")

	if country == "" || cityCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "country and city_code are required"})
		return
	}

	// Преобразование названия страны в код
	countryCode := ""
	switch strings.ToLower(country) {
	case "россия":
		countryCode = "RU"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported country"})
		return
	}

	params := map[string]string{
		"country_codes": countryCode,
		"city_code":     cityCode,
	}

	log.Printf("Запрос списка ПВЗ СДЭК от фронтенда: страна=%s, код города=%s -> параметры для сервиса: %+v", country, cityCode, params)

	pvzList, err := h.services.Cdek.GetPvzList(params)
	if err != nil {
		log.Printf("Ошибка получения списка ПВЗ из сервиса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Получено %d ПВЗ от API СДЭК", len(pvzList))
	c.JSON(http.StatusOK, pvzList)
}
