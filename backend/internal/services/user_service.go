package services

import (
	"errors"
	"github.com/yourname/MarketEase/internal/models"
	"github.com/yourname/MarketEase/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetUserByID(userID string) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("пользователь не найден")
	}
	return user, nil
}

func (s *UserService) GetAllUsers(includeDeleted, includeBanned bool) ([]models.User, error) {
	users, err := s.userRepo.FindAll(includeDeleted, includeBanned)
	if err != nil {
		return nil, errors.New("ошибка получения пользователей")
	}
	return users, nil
}

func (s *UserService) GetDeletedUsers() ([]models.User, error) {
	users, err := s.userRepo.FindDeleted()
	if err != nil {
		return nil, errors.New("ошибка получения удаленных пользователей")
	}
	return users, nil
}

func (s *UserService) UpdateUserRole(userID, role string) error {
	if !models.IsValidRole(role) {
		return errors.New("недопустимая роль")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	user.Role = role
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("ошибка обновления роли пользователя")
	}

	return nil
}

func (s *UserService) BanUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	user.Status = models.StatusBanned
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("ошибка блокировки пользователя")
	}

	return nil
}

func (s *UserService) UnbanUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	user.Status = models.StatusActive
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("ошибка разблокировки пользователя")
	}

	return nil
}

func (s *UserService) DeleteUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	user.Status = models.StatusDeleted
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("ошибка удаления пользователя")
	}

	return nil
}

func (s *UserService) RestoreUser(userID string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	user.Status = models.StatusActive
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("ошибка восстановления пользователя")
	}

	return nil
}

func (s *UserService) UpdatePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("пользователь не найден")
	}

	// Проверяем старый пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("старый пароль неверный")
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка шифрования пароля")
	}

	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("ошибка обновления пароля")
	}

	return nil
}
