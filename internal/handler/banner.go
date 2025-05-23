package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

	// Сохраняем тип контента и расширение в отдельном файле
	meta := struct {
		ContentType string `json:"content_type"`
		Extension   string `json:"extension"`
	}{
		ContentType: contentType,
		Extension:   ext,
	}

	metaJson, _ := json.Marshal(meta)
	if err := os.WriteFile(filepath.Join("uploads", "banner.meta"), metaJson, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save meta data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Banner uploaded successfully"})
}

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format. Allowed: .jpg, .jpeg, .png, .gif, .mp4, .webp"})
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

	// Сохраняем метаданные
	meta := struct {
		ContentType string `json:"content_type"`
		Extension   string `json:"extension"`
	}{
		ContentType: contentType,
		Extension:   ext,
	}

	metaJson, _ := json.Marshal(meta)
	if err := os.WriteFile(filepath.Join("uploads", "mobile_banner.meta"), metaJson, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save mobile banner meta data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Mobile banner uploaded successfully",
		"extension": ext,
	})
}
