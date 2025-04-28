package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yourname/MarketEase/internal/models"
	"net/http"
	"os"
	"strings"
)

// AuthMiddleware проверяет наличие и валидность JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Необходим токен авторизации"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Неверный формат токена"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				secret = "default_secret"
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Неверный или просроченный токен"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Message: "Не удалось распознать токен"})
			c.Abort()
			return
		}

		// Сохраняем user_id и роль в контекст запроса
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

// ViewerMiddleware проверяет, что пользователь имеет хотя бы роль viewer
func ViewerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("role"); !exists {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "Нет доступа"})
			c.Abort()
			return
		}

		// Любая роль имеет доступ к просмотру
		c.Next()
	}
}

// ManagerMiddleware проверяет, что пользователь имеет роль manager или director
func ManagerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "Нет доступа"})
			c.Abort()
			return
		}

		if role != models.RoleManager && role != models.RoleDirector {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "Требуется роль manager или director"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// DirectorMiddleware проверяет, что пользователь имеет роль director
func DirectorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "Нет доступа"})
			c.Abort()
			return
		}

		if role != models.RoleDirector {
			c.JSON(http.StatusForbidden, models.ErrorResponse{Message: "Требуется роль director"})
			c.Abort()
			return
		}

		c.Next()
	}
}
