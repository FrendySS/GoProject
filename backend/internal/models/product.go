package models

import (
	"github.com/google/uuid"
	"time"
)

// Константы для статусов товаров
const (
	ProductStatusActive  = "active"
	ProductStatusDeleted = "deleted"
)

type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	Status      string    `json:"status" gorm:"type:varchar(20);default:active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

// ProductFilter структура для фильтрации товаров
type ProductFilter struct {
	Name           string  `form:"name"`
	MinPrice       float64 `form:"minPrice"`
	MaxPrice       float64 `form:"maxPrice"`
	MinStock       int     `form:"minStock"`
	MaxStock       int     `form:"maxStock"`
	SortBy         string  `form:"sortBy"`
	SortOrder      string  `form:"sortOrder"`
	IncludeDeleted bool    `form:"includeDeleted"`
}
