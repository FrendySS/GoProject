package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Logger *zap.Logger

// InitLogger инициализирует логгер
func InitLogger() {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Logger, err = config.Build()
	if err != nil {
		fmt.Printf("Ошибка инициализации логгера: %v\n", err)
		os.Exit(1)
	}
}

// LoggerMiddleware middleware для логирования запросов
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Время начала запроса
		startTime := time.Now()

		// Обработка запроса
		c.Next()

		// Время окончания запроса
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Логирование информации о запросе
		Logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// Логирование ошибок
		for _, err := range c.Errors {
			Logger.Error("Request error",
				zap.String("error", err.Error()),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
			)
		}
	}
}
