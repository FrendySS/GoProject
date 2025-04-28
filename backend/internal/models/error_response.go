package models

// ErrorResponse стандартная структура ошибок
type ErrorResponse struct {
	Message string `json:"message" example:"Произошла ошибка"`
}
