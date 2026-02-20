package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
	"strconv"
)

func (h *Handler) CreateCart(c *gin.Context) {
	var request struct {
		Items []model.CartItem `json:"items"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cartID, err := h.services.Cart.CreateValidCart(request.Items)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: strconv.Itoa(cartID)})
}

func (h *Handler) GetCartById(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))

	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cart, err := h.services.GetCartByID(id)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, cart)
}
