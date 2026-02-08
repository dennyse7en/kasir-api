package service

import (
	"errors"
	"kasir-api/internal/model"
	"kasir-api/internal/repository"
)

type ProductService interface {
	Create(product model.Product) (model.Product, error)
	GetAll() ([]model.Product, error)
	GetByID(id int) (model.Product, error)
	Update(id int, product model.Product) (model.Product, error)
	Delete(id int) error
	SearchByName(name string) ([]model.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(product model.Product) (model.Product, error) {
	if product.Name == "" {
		return model.Product{}, errors.New("product name is required")
	}
	if product.CategoryID == 0 {
		return model.Product{}, errors.New("category id is required")
	}
	if product.Price < 0 {
		return model.Product{}, errors.New("price cannot be negative")
	}
	return s.repo.Create(product)
}

func (s *productService) GetAll() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) GetByID(id int) (model.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) Update(id int, product model.Product) (model.Product, error) {
	if product.Name == "" {
		return model.Product{}, errors.New("product name is required")
	}
	if product.Price < 0 {
		return model.Product{}, errors.New("price cannot be negative")
	}
	return s.repo.Update(id, product)
}

func (s *productService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *productService) SearchByName(name string) ([]model.Product, error) {
	if name == "" {
		return s.GetAll()
	}
	return s.repo.SearchByName(name)
}
