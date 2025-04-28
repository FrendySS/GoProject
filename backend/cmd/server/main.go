// @title MarketEase API
// @version 2.0
// @description API для управления складом товаров с ролевой моделью доступа
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"github.com/joho/godotenv"
	"github.com/yourname/MarketEase/config"
	_ "github.com/yourname/MarketEase/docs"
	"github.com/yourname/MarketEase/internal/database"
	"github.com/yourname/MarketEase/internal/routes"
	"github.com/yourname/MarketEase/internal/utils"
)

func main() {
	// Инициализация логгера
	utils.InitLogger()
	defer utils.Logger.Sync()

	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Info("Предупреждение: .env файл не найден, используются переменные окружения")
	}

	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Подключение к базе данных
	database.ConnectDatabase(cfg)

	// Заполнение базы данных тестовыми данными
	database.SeedData()

	// Настройка и запуск сервера
	router := routes.SetupRouter()
	utils.Logger.Info("Сервер запущен на порту :8080")
	if err := router.Run(":8080"); err != nil {
		utils.Logger.Fatal("Ошибка запуска сервера: " + err.Error())
	}
}
