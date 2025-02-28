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

	return router
}
