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

func (h *Handler) GetBanner(c *gin.Context) {
	bannerPath := filepath.Join("uploads", "banner")
	metaPath := filepath.Join("uploads", "banner.meta")

	// Проверяем существование файла
	if _, err := os.Stat(bannerPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Banner not found"})
		return
	}

	// Читаем метаданные
	metaJson, err := os.ReadFile(metaPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read meta data"})
		return
	}

	var meta struct {
		ContentType string `json:"content_type"`
		Extension   string `json:"extension"`
	}
	if err := json.Unmarshal(metaJson, &meta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse meta data"})
		return
	}

	// Если запрос пришел с правильным расширением - отдаем файл
	if strings.HasSuffix(c.Request.URL.Path, meta.Extension) {
		c.Header("Content-Type", meta.ContentType)
		c.File(bannerPath)
		return
	}

	// Если запрос без расширения - редиректим на URL с расширением
	newUrl := c.Request.URL.Path + meta.Extension
	if c.Request.URL.RawQuery != "" {
		newUrl += "?" + c.Request.URL.RawQuery
	}
	c.Redirect(http.StatusMovedPermanently, newUrl)
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

func (h *Handler) GetMobileBanner(c *gin.Context) {
	bannerPath := filepath.Join("uploads", "mobile_banner")
	metaPath := filepath.Join("uploads", "mobile_banner.meta")

	// Проверяем существование файла
	if _, err := os.Stat(bannerPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mobile banner not found"})
		return
	}

	// Читаем метаданные
	metaJson, err := os.ReadFile(metaPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to read mobile banner meta data"})
		return
	}

	var meta struct {
		ContentType string `json:"content_type"`
		Extension   string `json:"extension"`
	}
	if err := json.Unmarshal(metaJson, &meta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse mobile banner meta data"})
		return
	}

	// Если запрос уже содержит правильное расширение
	if strings.HasSuffix(c.Request.URL.Path, meta.Extension) {
		c.Header("Content-Type", meta.ContentType)
		c.Header("Cache-Control", "public, max-age=86400") // Кеширование на 1 день
		c.File(bannerPath)
		return
	}

	// Редирект на URL с правильным расширением
	newUrl := strings.TrimSuffix(c.Request.URL.Path, "/") + meta.Extension
	if c.Request.URL.RawQuery != "" {
		newUrl += "?" + c.Request.URL.RawQuery
	}
	c.Redirect(http.StatusMovedPermanently, newUrl)
}
