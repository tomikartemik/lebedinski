package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lebedinski/internal/service"
	"log"
	"net/http"
	"strings"
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

	//UPLOADS
	////////////////////////////////////////////////////////////
	router.Static("/uploads", "./uploads")
	////////////////////////////////////////////////////////////

	//ITEM
	////////////////////////////////////////////////////////////
	item := router.Group("/item")
	{
		item.POST("/new", h.CreateItem)
		item.GET("/all", h.AllItems)
		item.GET("", h.ItemByID)
		item.PUT("", h.UpdateItem)
	}
	////////////////////////////////////////////////////////////

	//ITEM
	////////////////////////////////////////////////////////////
	photo := router.Group("/photo")
	{
		photo.POST("/new", h.UploadPhoto)
	}
	////////////////////////////////////////////////////////////

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
		//payment.POST("/new", h.CreatePayment)
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
	countryName := c.Query("country") // Например: Россия
	cityName := c.Query("city")       // Например: Москва

	if countryName == "" || cityName == "" {
		log.Println("Ошибка: не указаны обязательные параметры 'country' и 'city'")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameters 'country' and 'city' are required"})
		return
	}

	// Преобразуем название страны в код (простой пример, можно расширить)
	countryCode := ""
	switch strings.ToLower(countryName) {
	case "россия":
		countryCode = "RU"
	case "беларусь":
		countryCode = "BY"
	case "казахстан":
		countryCode = "KZ"
	// Добавьте другие страны по необходимости
	default:
		// Попытаемся использовать переданное значение как код напрямую,
		// если оно короткое (например, "RU", "BY")
		if len(countryName) == 2 {
			countryCode = strings.ToUpper(countryName)
		} else {
			log.Printf("Неизвестное название страны: %s", countryName)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported country name", "country": countryName})
			return
		}
	}

	// Собираем параметры для *сервиса* (текстовые/кодовые)
	params := map[string]string{
		"country_codes": countryCode, // Передаем код страны
		"city":          cityName,    // Передаем название города
	}

	log.Printf("Запрос списка ПВЗ СДЭК от фронтенда: страна=%s, город=%s -> параметры для сервиса: %+v", countryName, cityName, params)

	// Вызываем сервис для получения списка ПВЗ
	pvzList, err := h.services.Cdek.GetPvzList(params) // Сервис теперь сам разберется с кодами
	if err != nil {
		log.Printf("Ошибка получения списка ПВЗ из сервиса: %v", err)
		// Проверяем, не связана ли ошибка с тем, что город не найден
		if strings.Contains(err.Error(), "city not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "City not found in CDEK database", "city": cityName, "country": countryCode})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PVZ list", "details": err.Error()})
		}
		return
	}

	log.Printf("Получено %d ПВЗ для города %s, страна %s", len(pvzList), cityName, countryName)
	c.JSON(http.StatusOK, pvzList)
}
