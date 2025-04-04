package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/utils"
	"net/http"
	"strconv"
)

func (h *Handler) GetAllOrders(c *gin.Context) {
	orders, err := h.services.GetAllOrders()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrderByCartID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	order, err := h.services.GetOrderByCartID(id)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, order)
}
