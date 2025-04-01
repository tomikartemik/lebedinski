package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
)

func (h *Handler) AddNewSizes(c *gin.Context) {
	var sizes []model.Size

	if err := c.ShouldBindJSON(&sizes); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.AddNewSizes(sizes)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Size successfully added")
}

func (h *Handler) UpdateSize(c *gin.Context) {
	var size model.Size

	if err := c.ShouldBindJSON(&size); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.services.UpdateSize(size)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Size successfully updated")
}
