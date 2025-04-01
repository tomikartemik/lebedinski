package service

import (
	"lebedinski/internal/model"
	"lebedinski/internal/repository"
	"strconv"
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

func (s *CategoryService) UpdateCategory(category model.Category) error {
	return s.repo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(categoryIDStr string) error {
	categoryID, err := strconv.Atoi(categoryIDStr)

	if err != nil {
		return err
	}

	return s.repo.DeleteCategory(categoryID)
}
