package repositories

import (
	"github.com/yourname/MarketEase/internal/database"
	"github.com/yourname/MarketEase/internal/models"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) FindByID(id string) (*models.Product, error) {
	var product models.Product
	result := database.DB.First(&product, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return database.DB.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return database.DB.Save(product).Error
}

func (r *ProductRepository) FindAll(filter models.ProductFilter) ([]models.Product, error) {
	var products []models.Product
	query := database.DB

	// Применяем фильтры
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}

	if filter.MinPrice > 0 {
		query = query.Where("price >= ?", filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		query = query.Where("price <= ?", filter.MaxPrice)
	}

	if filter.MinStock > 0 {
		query = query.Where("stock >= ?", filter.MinStock)
	}

	if filter.MaxStock > 0 {
		query = query.Where("stock <= ?", filter.MaxStock)
	}

	if !filter.IncludeDeleted {
		query = query.Where("status = ?", models.ProductStatusActive)
	}

	// Применяем сортировку
	if filter.SortBy != "" {
		order := "asc"
		if filter.SortOrder == "desc" {
			order = "desc"
		}

		// Проверяем допустимые поля для сортировки
		switch filter.SortBy {
		case "name", "price", "stock", "created_at":
			query = query.Order(filter.SortBy + " " + order)
		}
	}

	result := query.Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}

	return products, nil
}

func (r *ProductRepository) FindDeleted() ([]models.Product, error) {
	var products []models.Product
	result := database.DB.Where("status = ?", models.ProductStatusDeleted).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}
