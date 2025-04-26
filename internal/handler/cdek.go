package handler

import (
	"fmt"
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

	fmt.Println(order)
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

	pvzList, err := h.services.Cdek.GetPvzList(params)
	if err != nil {
		log.Printf("Ошибка получения списка ПВЗ из сервиса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Получено %d ПВЗ от API СДЭК", len(pvzList))
	c.JSON(http.StatusOK, pvzList)
}
