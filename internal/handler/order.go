package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
	"strconv"
)

func (h *Handler) GetAllOrders(c *gin.Context) {
	orders, err := h.services.GetAllOrders()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrderByCartID(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.services.GetOrderByCartID(id)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) ChangeStatusToSent(c *gin.Context) {
	cartID := c.Query("cart_id")

	err := h.services.SendOrderShippedNotification(cartID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *Handler) DeleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("cart_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.DeleteOrder(id)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Deleted!")
}

func (h *Handler) UpdateOrder(c *gin.Context) {
	order := model.Order{}
	err := c.ShouldBindJSON(&order)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.UpdateOrder(order)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *Handler) NewStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("cart_id"))
	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	status := c.Query("status")

	err = h.services.ChangeStatus(id, status)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Status: "+status)
}
