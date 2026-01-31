package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) Create(product *models.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetById(id string) (*models.Product, error) {
	return s.repo.GetById(id)
}

func (s *ProductService) DeleteById(id string) error {
	return s.repo.DeleteById(id)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}
