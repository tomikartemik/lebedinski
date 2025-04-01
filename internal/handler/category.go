package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
)

func (h *Handler) AddNewCategory(c *gin.Context) {
	var category model.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.AddCategory(category)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Category successfully added")
}

func (h *Handler) GetAllCategorise(c *gin.Context) {
	categories, err := h.services.GetAllCategories()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, categories)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	var category model.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.services.UpdateCategory(category)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Category successfully updated")
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Query("id")

	err := h.services.DeleteCategory(id)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, "Category successfully deleted")
}
