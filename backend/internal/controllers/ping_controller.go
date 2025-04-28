package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping отвечает на тестовый запрос
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
