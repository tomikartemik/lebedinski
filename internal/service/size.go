package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"strconv"
)

type SizeService struct {
	repo repository.Size
}

func NewSizeService(repo repository.Size) *SizeService {
	return &SizeService{repo: repo}
}

func (s *SizeService) AddNewSizes(sizes []model.Size) error {
	return s.repo.AddNewSizes(sizes)
}

func (s *SizeService) UpdateSize(sizeIDStr string, updateData map[string]interface{}) error {
	sizeID, err := strconv.Atoi(sizeIDStr)
	if err != nil {
		return err
	}

	// Опционально: здесь можно добавить валидацию данных в updateData,
	// например, пересчет скидки, если меняется цена.

	return s.repo.UpdateSize(sizeID, updateData)
}

func (s *SizeService) DeleteSize(sizeIDStr string) error {
	sizeID, err := strconv.Atoi(sizeIDStr)

	if err != nil {
		return err
	}

	return s.repo.DeleteSize(sizeID)
}
