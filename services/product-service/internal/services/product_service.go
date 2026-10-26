package services

import (
	"database/sql"
	"errors"
	"strings"

	"product-service/internal/models"
	"product-service/internal/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) validateProductFields(name string, price float64, quantity int) error {

	if strings.TrimSpace(name) == "" {
		return errors.New("Missing required field: name")
	}

	if len(strings.TrimSpace(name)) < 3 {
		return errors.New("name must be at least 3 characters")
	}

	if price <= 0 {
		return errors.New("price must be greater than 0")
	}

	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	return nil
}

func (s *ProductService) Create(req models.CreateProductRequest) (*models.Product, error) {

	if err := s.validateProductFields(req.Name, req.Price, req.Quantity); err != nil {
		return nil, err
	}

	return s.repo.Create(req)
}

func (s *ProductService) GetAll(keyword string) ([]models.Product, error) {
	return s.repo.GetAll(keyword)
}

func (s *ProductService) GetByID(id int) (*models.Product, error) {

	product, err := s.repo.GetByID(id)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}

		return nil, err
	}

	return product, nil
}

func (s *ProductService) Update(id int, req models.UpdateProductRequest) (*models.Product, error) {

	if err := s.validateProductFields(req.Name, req.Price, req.Quantity); err != nil {
		return nil, err
	}

	product, err := s.repo.Update(id, req)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Delete(id int) error {

	_, err := s.repo.GetByID(id)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("product not found")
		}

		return err
	}

	return s.repo.Delete(id)
}
