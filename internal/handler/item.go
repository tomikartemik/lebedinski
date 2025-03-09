package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
)

func (h *Handler) CreateItem(c *gin.Context) {
	var item model.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemID, err := h.services.CreateItem(item)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, itemID)
}

func (h *Handler) AllItems(c *gin.Context) {
	items, err := h.services.GetAllItems()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) ItemByID(c *gin.Context) {
	id := c.Query("id")

	item, err := h.services.GetItemByID(id)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) UpdateItem(c *gin.Context) {
	var item model.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.UpdateItem(item); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Item updated!"})
}
