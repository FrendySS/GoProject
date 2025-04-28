package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yourname/MarketEase/internal/models"
	"github.com/yourname/MarketEase/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) RegisterUser(username, email, password string) (*models.User, error) {
	// Проверяем, существует ли пользователь с таким email
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("ошибка шифрования пароля")
	}

	// Создаем нового пользователя
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     models.RoleViewer,
		Status:   models.StatusActive,
	}

	// Сохраняем пользователя в БД
	if err := s.userRepo.Create(&user); err != nil {
		return nil, errors.New("ошибка создания пользователя")
	}

	return &user, nil
}

func (s *AuthService) LoginUser(email, password string) (string, string, error) {
	// Ищем пользователя по email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("неверный email или пароль")
	}

	// Проверяем статус пользователя
	if user.Status == models.StatusBanned {
		return "", "", errors.New("пользователь заблокирован")
	}

	if user.Status == models.StatusDeleted {
		return "", "", errors.New("пользователь удален")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("неверный email или пароль")
	}

	// Генерируем JWT токен
	accessToken, err := generateAccessToken(*user)
	if err != nil {
		return "", "", errors.New("ошибка генерации токена доступа")
	}

	// Генерируем рефреш токен
	refreshToken, err := generateRefreshToken(*user)
	if err != nil {
		return "", "", errors.New("ошибка генерации рефреш токена")
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	// Проверяем валидность рефреш токена
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		secret := os.Getenv("JWT_REFRESH_SECRET")
		if secret == "" {
			secret = "refresh_secret"
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("неверный или просроченный рефреш токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("не удалось распознать токен")
	}

	// Получаем пользователя по ID из токена
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("некорректный формат токена")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", "", errors.New("пользователь не найден")
	}

	// Проверяем статус пользователя
	if user.Status == models.StatusBanned {
		return "", "", errors.New("пользователь заблокирован")
	}

	if user.Status == models.StatusDeleted {
		return "", "", errors.New("пользователь удален")
	}

	// Генерируем новый токен доступа
	accessToken, err := generateAccessToken(*user)
	if err != nil {
		return "", "", errors.New("ошибка генерации токена доступа")
	}

	// Генерируем новый рефреш токен
	newRefreshToken, err := generateRefreshToken(*user)
	if err != nil {
		return "", "", errors.New("ошибка генерации рефреш токена")
	}

	return accessToken, newRefreshToken, nil
}

func generateAccessToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // Токен действителен 15 минут
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret"
	}

	return token.SignedString([]byte(secret))
}

func generateRefreshToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // Рефреш токен действителен 7 дней
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "refresh_secret"
	}

	return token.SignedString([]byte(secret))
}
