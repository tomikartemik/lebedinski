package service

import (
	"fmt"
	"io"
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
)

type PhotoService struct {
	repo repository.Photo
}

func NewPhotoService(repo repository.Photo) *PhotoService {
	return &PhotoService{repo: repo}
}

func (s *PhotoService) SavePhoto(itemIDStr string, file *multipart.FileHeader) error {
	itemID, err := strconv.Atoi(itemIDStr)

	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("photo_%d%s", itemID, ext)
	filePath := filepath.Join("uploads", newFileName)

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return fmt.Errorf("unable to create upload directory: %v", err)
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create file: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("unable to save file: %v", err)
	}

	photo := model.Photo{
		Link:   filePath,
		ItemID: itemID,
	}

	if err := s.repo.NewPhoto(photo); err != nil {
		return fmt.Errorf("unable to save photo to database: %v", err)
	}

	return nil
}
