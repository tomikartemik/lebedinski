package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
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
