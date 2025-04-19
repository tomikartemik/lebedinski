package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) UploadBanner(c *gin.Context) {
	// Получаем файл из формы
	file, err := c.FormFile("banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get banner file from form"})
		return
	}

	// Проверяем допустимые расширения файлов
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".webp"}
	validExtension := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			validExtension = true
			break
		}
	}

	if !validExtension {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .jpg, .jpeg, .png, .gif, .mp4, .webp files are allowed"})
		return
	}

	// Создаем директорию uploads, если она не существует
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
		return
	}

	// Путь к файлу баннера (сохраняем оригинальное расширение)
	bannerPath := filepath.Join("uploads", "banner"+ext)

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

func (h *Handler) UploadMobileBanner(c *gin.Context) {
	// Получаем файл из формы
	file, err := c.FormFile("mobile_banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get mobile banner file from form"})
		return
	}

	// Проверяем допустимые расширения файлов
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".mp4", ".webp"}
	validExtension := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			validExtension = true
			break
		}
	}

	if !validExtension {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .jpg, .jpeg, .png, .gif, .mp4, .webp files are allowed"})
		return
	}

	// Создаем директорию uploads, если она не существует
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
		return
	}

	// Путь к файлу мобильного баннера (сохраняем оригинальное расширение)
	bannerPath := filepath.Join("uploads", "mobile_banner"+ext)

	// Если файл существует, удаляем его
	if _, err := os.Stat(bannerPath); err == nil {
		if err := os.Remove(bannerPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove existing mobile banner"})
			return
		}
	}

	// Сохраняем новый файл
	if err := c.SaveUploadedFile(file, bannerPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save mobile banner file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mobile banner uploaded successfully"})
}
