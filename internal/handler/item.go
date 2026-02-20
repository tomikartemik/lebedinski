package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"lebedinski/internal/model"
	"lebedinski/internal/utils"
	"net/http"
	"strconv"
)

type createItemRequest struct {
	model.Item
	CategoryIDs []int `json:"category_ids"`
}

func (h *Handler) CreateItem(c *gin.Context) {
	var req createItemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemID, err := h.services.CreateItem(req.Item)
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if len(req.CategoryIDs) > 0 {
		if err := h.services.UpdateItemCategories(itemID, req.CategoryIDs); err != nil {
			fmt.Printf("Warning: Failed to set categories for item %d: %v\n", itemID, err)
		}
	}

	c.JSON(http.StatusOK, itemID)
}

func (h *Handler) AllItems(c *gin.Context) {
	items, err := h.services.GetAllItems()
	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) ItemByID(c *gin.Context) {
	id := c.Query("id")

	item, err := h.services.GetItemByID(id)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
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
		return
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
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) ChangeTopItem(c *gin.Context) {
	positionStr := c.Query("position")
	position, err := strconv.Atoi(positionStr)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	itemIDStr := c.Query("item_id")
	itemID, err := strconv.Atoi(itemIDStr)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.ChangeTopItem(position, itemID)

	if err != nil {
		utils.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Item changed successfully!"})
}
