package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Константы для ролей пользователей
const (
	RoleViewer   = "viewer"
	RoleManager  = "manager"
	RoleDirector = "director"
)

// Константы для статусов пользователей
const (
	StatusActive  = "active"
	StatusBanned  = "banned"
	StatusDeleted = "deleted"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey" example:"e7b1c2f1-567b-4b14-a8d7-cd5b08bfc9d0"`
	Username  string     `json:"username" gorm:"unique;not null" example:"john_doe"`
	Email     string     `json:"email" gorm:"unique;not null" example:"john@example.com"`
	Password  string     `json:"-" gorm:"not null"` // В ответах на API пароль скрываем!
	Role      string     `json:"role" gorm:"type:varchar(20);default:viewer" example:"viewer"`
	Status    string     `json:"status" gorm:"type:varchar(20);default:active" example:"active"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}

// BeforeCreate будет вызываться перед вставкой новой записи
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

// IsValidRole проверяет, является ли роль допустимой
func IsValidRole(role string) bool {
	return role == RoleViewer || role == RoleManager || role == RoleDirector
}

// IsValidStatus проверяет, является ли статус допустимым
func IsValidStatus(status string) bool {
	return status == StatusActive || status == StatusBanned || status == StatusDeleted
}
