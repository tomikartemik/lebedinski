package handler

import (
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
	"strconv"
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

	if err := h.services.UpdateItem(id, updateData); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Item updated!"})
}

func (h *Handler) DeleteItem(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.NewErrorResponse(c, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.services.DeleteItem(id); err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Item deleted successfully"})
}

func (h *Handler) GetTopItems(c *gin.Context) {
	items, err := h.services.GetTopItems()

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) ChangeTopItem(c *gin.Context) {
	positionStr := c.Query("position")
	position, err := strconv.Atoi(positionStr)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	itemIDStr := c.Query("item_id")
	itemID, err := strconv.Atoi(itemIDStr)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err = h.services.ChangeTopItem(position, itemID)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Item changed successfully!"})
}
