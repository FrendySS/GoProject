package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/MarketEase/internal/models"
	"github.com/yourname/MarketEase/internal/repositories"
	"github.com/yourname/MarketEase/internal/services"
	"github.com/yourname/MarketEase/internal/utils"
	"net/http"
	"strconv"
)

var (
	productRepo    = repositories.NewProductRepository()
	productService = services.NewProductService(productRepo)
)

// GetProducts godoc
// @Summary Получить все товары
// @Description Получить список всех товаров с возможностью фильтрации и сортировки
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param name query string false "Фильтр по названию"
// @Param minPrice query number false "Минимальная цена"
// @Param maxPrice query number false "Максимальная цена"
// @Param minStock query integer false "Минимальное количество на складе"
// @Param maxStock query integer false "Максимальное количество на складе"
// @Param sortBy query string false "Поле для сортировки (name, price, stock, created_at)"
// @Param sortOrder query string false "Порядок сортировки (asc, desc)"
// @Param includeDeleted query boolean false "Включать удаленные товары"
// @Success 200 {array} models.Product
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products [get]
func GetProducts(c *gin.Context) {
	var filter models.ProductFilter

	// Получаем параметры фильтрации из запроса
	filter.Name = c.Query("name")

	if minPrice, err := strconv.ParseFloat(c.Query("minPrice"), 64); err == nil {
		filter.MinPrice = minPrice
	}

	if maxPrice, err := strconv.ParseFloat(c.Query("maxPrice"), 64); err == nil {
		filter.MaxPrice = maxPrice
	}

	if minStock, err := strconv.Atoi(c.Query("minStock")); err == nil {
		filter.MinStock = minStock
	}

	if maxStock, err := strconv.Atoi(c.Query("maxStock")); err == nil {
		filter.MaxStock = maxStock
	}

	filter.SortBy = c.Query("sortBy")
	filter.SortOrder = c.Query("sortOrder")
	filter.IncludeDeleted = c.Query("includeDeleted") == "true"

	products, err := productService.GetProducts(filter)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Получить товар по ID
// @Description Получить информацию о товаре по его ID
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID товара"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := productService.GetProductByID(id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Создать новый товар
// @Description Добавить новый товар (требуется роль manager или director)
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body models.Product true "Данные товара"
// @Success 201 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	if err := productService.CreateProduct(&product); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Обновить товар
// @Description Изменить информацию о товаре (требуется роль manager или director)
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID товара"
// @Param product body models.Product true "Обновленные данные товара"
// @Success 200 {object} models.Product
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var updatedProduct models.Product

	if err := c.ShouldBindJSON(&updatedProduct); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	product, err := productService.UpdateProduct(id, &updatedProduct)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Удалить товар
// @Description Пометить товар как удаленный (требуется роль manager или director)
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID товара"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id} [delete]
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := productService.DeleteProduct(id); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар успешно удалён"})
}

// RestoreProduct godoc
// @Summary Восстановить удаленный товар
// @Description Восстановить ранее удаленный товар (требуется роль manager или director)
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "ID товара"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/{id}/restore [put]
func RestoreProduct(c *gin.Context) {
	id := c.Param("id")

	if err := productService.RestoreProduct(id); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар успешно восстановлен"})
}

// GetDeletedProducts godoc
// @Summary Получить список удаленных товаров
// @Description Получить список всех удаленных товаров (требуется роль manager или director)
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /products/deleted [get]
func GetDeletedProducts(c *gin.Context) {
	products, err := productService.GetDeletedProducts()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, products)
}
