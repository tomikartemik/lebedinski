package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UploadBanner загружает основной баннер
func (h *Handler) UploadBanner(c *gin.Context) {
	file, err := c.FormFile("banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get banner file from form"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".mp4":  "video/mp4",
		".webp": "image/webp",
	}

	contentType, ok := allowedExtensions[ext]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
		return
	}

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
		return
	}

	// Сохраняем файл без расширения
	bannerPath := filepath.Join("uploads", "banner")

	// Удаляем старый файл если существует
	if _, err := os.Stat(bannerPath); err == nil {
		if err := os.Remove(bannerPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove existing banner"})
			return
		}
	}

	// Сохраняем файл
	if err := c.SaveUploadedFile(file, bannerPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save banner file"})
		return
	}

	// Сохраняем тип контента в отдельном файле
	if err := os.WriteFile(filepath.Join("uploads", "banner.content-type"), []byte(contentType), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save content type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Banner uploaded successfully"})
}

// UploadMobileBanner загружает мобильный баннер
func (h *Handler) UploadMobileBanner(c *gin.Context) {
	file, err := c.FormFile("mobile_banner")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to get mobile banner file from form"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".mp4":  "video/mp4",
		".webp": "image/webp",
	}

	contentType, ok := allowedExtensions[ext]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format"})
		return
	}

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
		return
	}

	// Сохраняем файл без расширения
	bannerPath := filepath.Join("uploads", "mobile_banner")

	// Удаляем старый файл если существует
	if _, err := os.Stat(bannerPath); err == nil {
		if err := os.Remove(bannerPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove existing mobile banner"})
			return
		}
	}

	// Сохраняем файл
	if err := c.SaveUploadedFile(file, bannerPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save mobile banner file"})
		return
	}

	// Сохраняем тип контента в отдельном файле
	if err := os.WriteFile(filepath.Join("uploads", "mobile_banner.content-type"), []byte(contentType), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save content type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mobile banner uploaded successfully"})
}

// GetBanner возвращает основной баннер с правильным Content-Type
func (h *Handler) GetBanner(c *gin.Context) {
	bannerPath := filepath.Join("uploads", "banner")
	contentTypePath := filepath.Join("uploads", "banner.content-type")

	// Проверяем существование файла
	if _, err := os.Stat(bannerPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
		return
	}

	// Читаем тип контента
	contentType, err := os.ReadFile(contentTypePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read content type"})
		return
	}

	// Отправляем файл с правильным Content-Type
	c.Header("Content-Type", string(contentType))
	c.File(bannerPath)
}

// GetMobileBanner возвращает мобильный баннер с правильным Content-Type
func (h *Handler) GetMobileBanner(c *gin.Context) {
	bannerPath := filepath.Join("uploads", "mobile_banner")
	contentTypePath := filepath.Join("uploads", "mobile_banner.content-type")

	// Проверяем существование файла
	if _, err := os.Stat(bannerPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mobile banner not found"})
		return
	}

	// Читаем тип контента
	contentType, err := os.ReadFile(contentTypePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read content type"})
		return
	}

	// Отправляем файл с правильным Content-Type
	c.Header("Content-Type", string(contentType))
	c.File(bannerPath)
}
