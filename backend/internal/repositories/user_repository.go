package repositories

import (
	"github.com/yourname/MarketEase/internal/database"
	"github.com/yourname/MarketEase/internal/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	return database.DB.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error {
	return database.DB.Save(user).Error
}

func (r *UserRepository) FindAll(includeDeleted, includeBanned bool) ([]models.User, error) {
	var users []models.User
	query := database.DB

	if !includeDeleted {
		query = query.Where("status != ?", models.StatusDeleted)
	}

	if !includeBanned {
		query = query.Where("status != ?", models.StatusBanned)
	}

	result := query.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (r *UserRepository) FindDeleted() ([]models.User, error) {
	var users []models.User
	result := database.DB.Where("status = ?", models.StatusDeleted).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
