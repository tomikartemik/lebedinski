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

// GetPvzList обрабатывает запрос на получение списка ПВЗ СДЭК по региону
func (h *Handler) GetPvzList(c *gin.Context) {
	countryName := c.Query("country") // Например: Россия
	regionName := c.Query("region")   // Например: Московская область

	if countryName == "" || regionName == "" {
		log.Println("Ошибка: не указаны обязательные параметры 'country' и 'region'")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameters 'country' and 'region' are required"})
		return
	}

	// Преобразуем название страны в код
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

	// Собираем параметры для *сервиса*
	params := map[string]string{
		"country_codes": countryCode, // Передаем код страны
		"region":        regionName,  // Передаем название региона
	}

	log.Printf("Запрос списка ПВЗ СДЭК от фронтенда: страна=%s, регион=%s -> параметры для сервиса: %+v", countryName, regionName, params)

	// Вызываем сервис для получения списка ПВЗ
	pvzList, err := h.services.Cdek.GetPvzList(params) // Сервис теперь сам найдет код региона
	if err != nil {
		log.Printf("Ошибка получения списка ПВЗ из сервиса: %v", err)
		// Проверяем, не связана ли ошибка с тем, что регион не найден
		if strings.Contains(err.Error(), "region not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Region not found in CDEK database", "region": regionName, "country": countryCode})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get PVZ list", "details": err.Error()})
		}
		return
	}

	log.Printf("Получено %d ПВЗ для региона %s, страна %s", len(pvzList), regionName, countryName)
	c.JSON(http.StatusOK, pvzList)
}
