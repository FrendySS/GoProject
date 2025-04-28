package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/MarketEase/internal/utils"
	"net/http"
)

// GetProfile godoc
// @Summary Профиль пользователя
// @Description Получить информацию о своём профиле
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /profile [get]
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.HandleError(c, utils.NewAppError(http.StatusUnauthorized, "Необходима авторизация"))
		return
	}

	user, err := userService.GetUserByID(userID.(string))
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdatePassword godoc
// @Summary Смена пароля
// @Description Изменить свой пароль
// @Tags profile
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body UpdatePasswordInput true "Старый и новый пароль"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /profile/password [put]
func UpdatePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.HandleError(c, utils.NewAppError(http.StatusUnauthorized, "Необходима авторизация"))
		return
	}

	var input UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	if err := userService.UpdatePassword(userID.(string), input.OldPassword, input.NewPassword); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменен"})
}

// UpdatePasswordInput структура тела запроса на смену пароля
type UpdatePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required" example:"oldpassword123"`
	NewPassword string `json:"newPassword" binding:"required,min=6" example:"newpassword456"`
}
