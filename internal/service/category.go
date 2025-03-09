package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) AddCategory(category model.Category) error {
	return s.repo.AddCategory(category)
}

func (s *CategoryService) GetAllCategories() ([]model.Category, error) {
	return s.repo.GetAllCategories()
}
