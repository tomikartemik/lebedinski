package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func (h *Handler) UploadBanner(c *gin.Context) {
	// Получаем файл из формы
	file, err := c.FormFile("banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get banner file from form"})
		return
	}

	// Проверяем расширение файла
	ext := filepath.Ext(file.Filename)
	if ext != ".mp4" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .mp4 files are allowed"})
		return
	}

	// Создаем директорию uploads, если она не существует
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
		return
	}

	// Путь к файлу баннера
	bannerPath := filepath.Join("uploads", "banner.mp4")

	// Если файл существует, удаляем его
	if _, err := os.Stat(bannerPath); err == nil {
		if err := os.Remove(bannerPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove existing banner"})
			return
		}
	}

	// Сохраняем новый файл
	if err := c.SaveUploadedFile(file, bannerPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save banner file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Banner uploaded successfully"})
} 