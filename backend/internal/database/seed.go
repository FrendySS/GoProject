package database

import (
	"github.com/yourname/MarketEase/internal/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

// SeedData заполняет базу данных тестовыми данными
func SeedData() {
	// Проверяем, есть ли уже пользователи в базе данных
	var userCount int64
	DB.Model(&models.User{}).Count(&userCount)

	if userCount == 0 {
		log.Println("Заполнение базы данных тестовыми пользователями...")

		// Создаем директора
		directorPassword, _ := bcrypt.GenerateFromPassword([]byte("director123"), bcrypt.DefaultCost)
		director := models.User{
			Username: "director",
			Email:    "director@marketease.com",
			Password: string(directorPassword),
			Role:     models.RoleDirector,
			Status:   models.StatusActive,
		}
		DB.Create(&director)

		// Создаем менеджера
		managerPassword, _ := bcrypt.GenerateFromPassword([]byte("manager123"), bcrypt.DefaultCost)
		manager := models.User{
			Username: "manager",
			Email:    "manager@marketease.com",
			Password: string(managerPassword),
			Role:     models.RoleManager,
			Status:   models.StatusActive,
		}
		DB.Create(&manager)

		// Создаем обычного пользователя
		userPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
		user := models.User{
			Username: "user",
			Email:    "user@marketease.com",
			Password: string(userPassword),
			Role:     models.RoleViewer,
			Status:   models.StatusActive,
		}
		DB.Create(&user)

		log.Println("Тестовые пользователи созданы успешно")
	}

	// Проверяем, есть ли уже товары в базе данных
	var productCount int64
	DB.Model(&models.Product{}).Count(&productCount)

	if productCount == 0 {
		log.Println("Заполнение базы данных тестовыми товарами...")

		products := []models.Product{
			{
				Name:        "Ноутбук HP Pavilion",
				Description: "Мощный ноутбук для работы и развлечений",
				Price:       59999.99,
				Stock:       15,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Смартфон Samsung Galaxy S21",
				Description: "Флагманский смартфон с отличной камерой",
				Price:       49999.99,
				Stock:       25,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Планшет Apple iPad Pro",
				Description: "Профессиональный планшет для творчества",
				Price:       79999.99,
				Stock:       10,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Умные часы Xiaomi Mi Band 6",
				Description: "Фитнес-трекер с множеством функций",
				Price:       3999.99,
				Stock:       50,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Беспроводные наушники Sony WH-1000XM4",
				Description: "Наушники с шумоподавлением премиум-класса",
				Price:       29999.99,
				Stock:       20,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Игровая консоль PlayStation 5",
				Description: "Новейшая игровая консоль от Sony",
				Price:       49999.99,
				Stock:       5,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Фотоаппарат Canon EOS R5",
				Description: "Профессиональная беззеркальная камера",
				Price:       299999.99,
				Stock:       3,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Умная колонка Яндекс Станция",
				Description: "Умная колонка с голосовым помощником Алиса",
				Price:       12999.99,
				Stock:       30,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Электросамокат Xiaomi Mi Electric Scooter Pro 2",
				Description: "Мощный электросамокат для городских поездок",
				Price:       39999.99,
				Stock:       8,
				Status:      models.ProductStatusActive,
			},
			{
				Name:        "Кофемашина DeLonghi Magnifica S",
				Description: "Автоматическая кофемашина для дома",
				Price:       34999.99,
				Stock:       12,
				Status:      models.ProductStatusActive,
			},
		}

		for _, product := range products {
			DB.Create(&product)
		}

		log.Println("Тестовые товары созданы успешно")
	}
}
