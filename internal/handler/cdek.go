package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"log"
	"net/http"
	"strings"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	var order model.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//cdekUUID, err := h.services.CreateCdekOrder(order)
	//if err != nil {
	//	utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	paymentResponse, err := h.services.CreatePayment(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.services.ProcessOrder(order, paymentResponse.ID)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusCreated, gin.H{
		"payment_url": paymentResponse.Confirmation.ConfirmationURL,
	})
}

func (h *Handler) GetPvzList(c *gin.Context) {
	countryName := c.Query("country")
	cityName := c.Query("city")

	if countryName == "" || cityName == "" {
		log.Println("Ошибка: не указаны обязательные параметры 'country' и 'city'")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameters 'country' and 'city' are required"})
		return
	}

	countryCode := ""
	switch strings.ToLower(countryName) {
	case "россия":
		countryCode = "RU"
	case "беларусь":
		countryCode = "BY"
	case "казахстан":
		countryCode = "KZ"
	default:
		if len(countryName) == 2 {
			countryCode = strings.ToUpper(countryName)
		} else {
			log.Printf("Неизвестное название страны: %s", countryName)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported country name", "country": countryName})
			return
		}
	}

	params := map[string]string{
		"country_codes": countryCode,
		"city":          cityName,
	}

	log.Printf("Запрос списка ПВЗ СДЭК от фронтенда: страна=%s, город=%s -> параметры для сервиса: %+v", countryName, cityName, params)

	pvzList, err := h.services.Cdek.GetPvzList(params)
	if err != nil {
		log.Printf("Ошибка получения списка ПВЗ из сервиса: %v", err)
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
