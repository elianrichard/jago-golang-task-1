package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) Create(product *models.Category) error {
	return s.repo.Create(product)
}

func (s *CategoryService) GetById(id string) (*models.Category, error) {
	return s.repo.GetById(id)
}

func (s *CategoryService) DeleteById(id string) error {
	return s.repo.DeleteById(id)
}

func (s *CategoryService) Update(product *models.Category) error {
	return s.repo.Update(product)
}
