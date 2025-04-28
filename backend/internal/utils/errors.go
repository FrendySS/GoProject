package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/MarketEase/internal/models"
	"net/http"
)

// AppError представляет ошибку приложения с HTTP-статусом
type AppError struct {
	StatusCode int
	Message    string
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError создает новую ошибку приложения
func NewAppError(statusCode int, message string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// HandleError обрабатывает ошибку и отправляет соответствующий ответ
func HandleError(c *gin.Context, err error) {
	// Проверяем, является ли ошибка AppError
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.StatusCode, models.ErrorResponse{Message: appErr.Message})
		return
	}

	// Обрабатываем обычные ошибки
	switch err.Error() {
	case "пользователь с таким email уже существует":
		c.JSON(http.StatusConflict, models.ErrorResponse{Message: err.Error()})
	case "неверный email или пароль", "пользователь заблокирован", "пользователь удален", "неверный или просроченный токен", "необходима авторизация":
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: err.Error()})
	case "пользователь не найден", "товар не найден":
		c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
	case "недопустимая роль", "старый пароль неверный", "неверные данные":
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	case "требуется роль manager или director", "требуется роль director", "нет доступа":
		c.JSON(http.StatusForbidden, models.ErrorResponse{Message: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
}
