package services

import (
	"errors"
	"github.com/yourname/MarketEase/internal/models"
	"github.com/yourname/MarketEase/internal/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) GetProducts(filter models.ProductFilter) ([]models.Product, error) {
	products, err := s.productRepo.FindAll(filter)
	if err != nil {
		return nil, errors.New("ошибка получения списка товаров")
	}
	return products, nil
}

func (s *ProductService) GetProductByID(productID string) (*models.Product, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("товар не найден")
	}
	return product, nil
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	product.Status = models.ProductStatusActive
	if err := s.productRepo.Create(product); err != nil {
		return errors.New("ошибка при создании товара")
	}
	return nil
}

func (s *ProductService) UpdateProduct(productID string, updatedProduct *models.Product) (*models.Product, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("товар не найден")
	}

	// Обновляем поля
	product.Name = updatedProduct.Name
	product.Description = updatedProduct.Description
	product.Price = updatedProduct.Price
	product.Stock = updatedProduct.Stock

	if err := s.productRepo.Update(product); err != nil {
		return nil, errors.New("ошибка при обновлении товара")
	}

	return product, nil
}

func (s *ProductService) DeleteProduct(productID string) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("товар не найден")
	}

	product.Status = models.ProductStatusDeleted
	if err := s.productRepo.Update(product); err != nil {
		return errors.New("ошибка при удалении товара")
	}

	return nil
}

func (s *ProductService) RestoreProduct(productID string) error {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("товар не найден")
	}

	product.Status = models.ProductStatusActive
	if err := s.productRepo.Update(product); err != nil {
		return errors.New("ошибка при восстановлении товара")
	}

	return nil
}

func (s *ProductService) GetDeletedProducts() ([]models.Product, error) {
	products, err := s.productRepo.FindDeleted()
	if err != nil {
		return nil, errors.New("ошибка получения списка удаленных товаров")
	}
	return products, nil
}
