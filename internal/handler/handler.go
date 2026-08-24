package handler

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"lebedinski/internal/service"
	"os"
	"strings"
	"time"
)

type Handler struct {
	services       *service.Service
	auth           *adminAuth
	dadata         *dadataProxy
	allowedOrigins []string
}

func NewHandler(services *service.Service) (*Handler, error) {
	auth, err := newAdminAuthFromEnv()
	if err != nil {
		return nil, err
	}
	dadata, err := newDadataProxyFromEnv()
	if err != nil {
		return nil, err
	}

	origins, err := configuredOrigins()
	if err != nil {
		return nil, err
	}

	return &Handler{
		services:       services,
		auth:           auth,
		dadata:         dadata,
		allowedOrigins: origins,
	}, nil
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	_ = router.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	router.Use(cors.New(cors.Config{
		AllowOrigins:     h.allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Static("/uploads", "./uploads")

	dadata := router.Group("/dadata")
	{
		dadata.POST("/suggest/address", h.dadataSuggestAddress)
		dadata.POST("/find-delivery", h.dadataFindDelivery)
	}

	auth := router.Group("/auth")
	{
		auth.POST("/login", h.login)
		auth.POST("/logout", h.logout)
		auth.GET("/session", h.requireAdmin(), h.sessionStatus)
	}

	banner := router.Group("/banner", h.requireAdmin())
	{
		banner.POST("/upload", h.UploadBanner)
		banner.POST("/upload_mobile", h.UploadMobileBanner)
	}

	item := router.Group("/item")
	{
		item.GET("/all", h.AllItems)
		item.GET("", h.ItemByID)
		item.GET("/top", h.GetTopItems)
	}
	itemAdmin := router.Group("/item", h.requireAdmin())
	{
		itemAdmin.POST("/new", h.CreateItem)
		itemAdmin.POST("/change-top-item", h.ChangeTopItem)
		itemAdmin.PUT("", h.UpdateItem)
		itemAdmin.DELETE("", h.DeleteItem)
	}

	photo := router.Group("/photo", h.requireAdmin())
	{
		photo.POST("/new", h.UploadPhoto)
		photo.DELETE("", h.DeletePhoto)
	}

	size := router.Group("/size", h.requireAdmin())
	{
		size.POST("/add", h.AddNewSizes)
		size.PUT("", h.UpdateSize)
		size.DELETE("", h.DeleteSize)
	}

	category := router.Group("/category")
	{
		category.GET("/all", h.GetAllCategorise)
	}
	categoryAdmin := router.Group("/category", h.requireAdmin())
	{
		categoryAdmin.POST("/new", h.AddNewCategory)
		categoryAdmin.PUT("", h.UpdateCategory)
		categoryAdmin.PUT("/update", h.UpdateCategory)
		categoryAdmin.DELETE("", h.DeleteCategory)
	}

	payment := router.Group("/payment")
	{
		payment.POST("/response", h.HandleWebhook)
		payment.POST("/send-message-if-failed", h.SendMessageIfFailed)
	}

	cart := router.Group("/cart")
	{
		cart.POST("/create", h.CreateCart)
		cart.GET("", h.GetCartById)
	}

	order := router.Group("/order")
	{
		order.POST("/new", h.CreateOrder)
	}
	orderAdmin := router.Group("/order", h.requireAdmin())
	{
		orderAdmin.GET("/all", h.GetAllOrders)
		orderAdmin.GET("/by-cart-id", h.GetCartById)
		orderAdmin.POST("/status", h.ChangeStatusToSent)
		orderAdmin.POST("/sent", h.ChangeStatusToSent)
		orderAdmin.POST("/new-status", h.NewStatus)
		orderAdmin.PUT("", h.UpdateOrder)
		orderAdmin.DELETE("", h.DeleteOrder)
	}

	cdek := router.Group("/cdek")
	{
		cdek.GET("/pvz", h.GetPvzList)
	}

	promocode := router.Group("/promocode")
	{
		promocode.GET("", h.GetPromocodeByCode)
	}
	promocodeAdmin := router.Group("/promocode", h.requireAdmin())
	{
		promocodeAdmin.POST("", h.CreatePromoCode)
		promocodeAdmin.GET("/all", h.GetPromocodeList)
		promocodeAdmin.DELETE("", h.DeletePromocode)
		promocodeAdmin.PUT("", h.UpdatePromocode)
	}
	return router
}

func configuredOrigins() ([]string, error) {
	raw := strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	if raw == "" {
		raw = "https://lebedinski.shop,https://www.lebedinski.shop,https://admin.lebedinski.shop"
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin == "" {
			continue
		}
		if origin == "*" {
			return nil, fmt.Errorf("ALLOWED_ORIGINS must not contain a wildcard")
		}
		origins = append(origins, origin)
	}
	if len(origins) == 0 {
		return nil, fmt.Errorf("ALLOWED_ORIGINS does not contain any valid origins")
	}
	return origins, nil
}
