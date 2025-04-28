package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/MarketEase/internal/services"
	"github.com/yourname/MarketEase/internal/utils"
	"net/http"
)

var (
	userService = services.NewUserService(userRepo)
)

// GetAllUsers godoc
// @Summary Получить список всех пользователей
// @Description Только для director
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Param includeDeleted query bool false "Включать удаленных пользователей"
// @Param includeBanned query bool false "Включать заблокированных пользователей"
// @Success 200 {array} models.User
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/users [get]
func GetAllUsers(c *gin.Context) {
	includeDeleted := c.Query("includeDeleted") == "true"
	includeBanned := c.Query("includeBanned") == "true"

	users, err := userService.GetAllUsers(includeDeleted, includeBanned)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetDeletedUsers godoc
// @Summary Получить список удаленных пользователей
// @Description Только для director
// @Tags admin
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/deleted-users [get]
func GetDeletedUsers(c *gin.Context) {
	users, err := userService.GetDeletedUsers()
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, users)
}

// AssignRole godoc
// @Summary Изменить роль пользователя
// @Description Только для director
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body AssignRoleInput true "Данные для изменения роли"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/assign-role [post]
func AssignRole(c *gin.Context) {
	var input AssignRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	if err := userService.UpdateUserRole(input.UserID, input.Role); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Роль пользователя успешно обновлена на " + input.Role})
}

// BanUser godoc
// @Summary Заблокировать пользователя
// @Description Только для director
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body UserIDInput true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/ban-user [post]
func BanUser(c *gin.Context) {
	var input UserIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	if err := userService.BanUser(input.UserID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно заблокирован"})
}

// UnbanUser godoc
// @Summary Разблокировать пользователя
// @Description Только для director
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body UserIDInput true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/unban-user [post]
func UnbanUser(c *gin.Context) {
	var input UserIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	if err := userService.UnbanUser(input.UserID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно разблокирован"})
}

// DeleteUser godoc
// @Summary Удалить пользователя
// @Description Только для director
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body UserIDInput true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/delete-user [post]
func DeleteUser(c *gin.Context) {
	var input UserIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	if err := userService.DeleteUser(input.UserID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
}

// RestoreUser godoc
// @Summary Восстановить удаленного пользователя
// @Description Только для director
// @Tags admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body UserIDInput true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/restore-user [post]
func RestoreUser(c *gin.Context) {
	var input UserIDInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewAppError(http.StatusBadRequest, "Неверные данные"))
		return
	}

	if err := userService.RestoreUser(input.UserID); err != nil {
		utils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно восстановлен"})
}

// Структуры запросов

type AssignRoleInput struct {
	UserID string `json:"userId" binding:"required" example:"e7b1c2f1-567b-4b14-a8d7-cd5b08bfc9d0"`
	Role   string `json:"role" binding:"required" example:"manager"`
}

type UserIDInput struct {
	UserID string `json:"userId" binding:"required" example:"e7b1c2f1-567b-4b14-a8d7-cd5b08bfc9d0"`
}
