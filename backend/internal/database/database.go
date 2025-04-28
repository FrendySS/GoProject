package database

import (
	"github.com/yourname/MarketEase/config"
	"github.com/yourname/MarketEase/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDatabase(cfg *config.Config) {
	connectionString := cfg.GetDBConnectionString()

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Не удалось подключиться к базе данных: %v", err)
	}

	// 🚀 Создаем расширение uuid-ossp, если его нет
	if err := ensureUUIDExtension(db); err != nil {
		log.Fatalf("❌ Не удалось создать расширение uuid-ossp: %v", err)
	}

	// 🔥 Автомиграция моделей
	err = db.AutoMigrate(
		&models.User{},
		&models.Product{},
	)
	if err != nil {
		log.Fatalf("❌ Ошибка миграции моделей: %v", err)
	}

	DB = db
	log.Println("✅ Успешное подключение к базе данных и миграция моделей")
}

// ensureUUIDExtension создает расширение uuid-ossp, если его нет
func ensureUUIDExtension(db *gorm.DB) error {
	return db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
}
