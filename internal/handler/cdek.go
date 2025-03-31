package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	var order model.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cdekUUID, err := h.services.CreateCdekOrder(order)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.ProcessOrder(order, cdekUUID)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: cdekUUID})
}
