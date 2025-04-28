package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/MarketEase/internal/repositories"
	"github.com/yourname/MarketEase/internal/services"
	"github.com/yourname/MarketEase/internal/utils"
	"net/http"
)

var (
	userRepo    = repositories.NewUserRepository()
	authService = services.NewAuthService(userRepo)
)

// RegisterUser godoc
// @Summary Регистрация пользователя
// @Description Создать нового пользователя с ролью viewer
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterInput true "Данные пользователя для регистрации"
// @Success 201 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /register [post]
func RegisterUser(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	user, err := authService.RegisterUser(input.Username, input.Email, input.Password)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// LoginUser godoc
// @Summary Вход пользователя
// @Description Аутентификация пользователя по email и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginInput true "Данные для входа"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /login [post]
func LoginUser(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	accessToken, refreshToken, err := authService.LoginUser(input.Email, input.Password)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	})
}

// RefreshToken godoc
// @Summary Обновление токена
// @Description Получение нового токена доступа с помощью рефреш токена
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body RefreshTokenInput true "Рефреш токен"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /refresh [post]
func RefreshToken(c *gin.Context) {
	var input RefreshTokenInput

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	accessToken, refreshToken, err := authService.RefreshToken(input.RefreshToken)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	})
}

// Структуры запросов и ответа

type RegisterInput struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"supersecret"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"supersecret"`
}

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string `json:"token_type" example:"Bearer"`
}
