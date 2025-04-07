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
	id := c.Query("id")
	if id == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "item id is required")
		return
	}

	var updateData map[string]interface{}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Удаляем поле "id" из данных для обновления, так как оно уже есть в URL
	delete(updateData, "id")

	if err := h.services.UpdateSize(id, updateData); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Size updated!"})
}

func (h *Handler) DeleteSize(c *gin.Context) {
	id := c.Query("id")
	err := h.services.DeleteSize(id)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Size successfully deleted")
}
