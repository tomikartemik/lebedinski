package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
)

func (h *Handler) CreatePromoCode(c *gin.Context) {
	var promocode model.PromoCode

	if err := c.ShouldBindJSON(&promocode); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.CreatePromoCode(promocode)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Promocode successfully added")
}

func (h *Handler) GetPromocodeByCode(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "code is required")
		return
	}

	promocode, err := h.services.GetPromoCodeByCode(code)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, promocode)
}

func (h *Handler) GetPromocodeList(c *gin.Context) {
	promocodes, err := h.services.GetAllPromoCodes()

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, promocodes)
}

func (h *Handler) DeletePromocode(c *gin.Context) {
	code := c.Query("code")

	if code == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "code is required")
		return
	}

	if err := h.services.DeletePromoCodeByCode(code); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Promocode successfully deleted")
}

func (h *Handler) UpdatePromocode(c *gin.Context) {
	var promocode model.PromoCode

	if err := c.ShouldBindJSON(&promocode); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.services.UpdatePromoCode(promocode)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Promocode successfully updated")
}
