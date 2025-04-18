package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) UploadPhoto(c *gin.Context) {
	itemID := c.Query("item_id")

	photo, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get photo from form"})
		return
	}

	err = h.services.SavePhoto(itemID, photo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Saved photo successfully")
}

func (h *Handler) DeletePhoto(c *gin.Context) {
	photoID := c.Query("photo_id")
	if photoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo_id is required"})
		return
	}

	err := h.services.DeletePhoto(photoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Photo deleted successfully")
}
