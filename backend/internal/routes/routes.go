package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yourname/MarketEase/internal/controllers"
	"github.com/yourname/MarketEase/internal/middlewares"
	"github.com/yourname/MarketEase/internal/utils"
	"time"
)

// SetupRouter настраивает все маршруты приложения
func SetupRouter() *gin.Engine {
	router := gin.New()

	// Добавляем middleware для восстановления после паники
	router.Use(gin.Recovery())

	// Добавляем middleware для логирования
	router.Use(utils.LoggerMiddleware())

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Статические файлы для веб-интерфейса
	router.Static("/web", "./web")

	api := router.Group("/api")
	{
		api.GET("/ping", controllers.Ping)

		// Аутентификация
		api.POST("/register", controllers.RegisterUser)
		api.POST("/login", controllers.LoginUser)
		api.POST("/refresh", controllers.RefreshToken)

		// Профиль пользователя
		profile := api.Group("/profile")
		profile.Use(middlewares.AuthMiddleware()) // Защищённая зона профиля
		{
			profile.GET("", controllers.GetProfile)
			profile.PUT("/password", controllers.UpdatePassword)
		}

		// Администрирование пользователей
		admin := api.Group("/admin")
		admin.Use(middlewares.AuthMiddleware())
		{
			admin.GET("/users", middlewares.DirectorMiddleware(), controllers.GetAllUsers)
			admin.GET("/deleted-users", middlewares.DirectorMiddleware(), controllers.GetDeletedUsers)
			admin.POST("/assign-role", middlewares.DirectorMiddleware(), controllers.AssignRole)
			admin.POST("/ban-user", middlewares.DirectorMiddleware(), controllers.BanUser)
			admin.POST("/unban-user", middlewares.DirectorMiddleware(), controllers.UnbanUser)
			admin.POST("/delete-user", middlewares.DirectorMiddleware(), controllers.DeleteUser)
			admin.POST("/restore-user", middlewares.DirectorMiddleware(), controllers.RestoreUser)
		}

		// Управление товарами
		products := api.Group("/products")
		products.Use(middlewares.AuthMiddleware())
		{
			products.GET("", middlewares.ViewerMiddleware(), controllers.GetProducts)
			products.GET("/deleted", middlewares.ManagerMiddleware(), controllers.GetDeletedProducts)
			products.GET("/:id", middlewares.ViewerMiddleware(), controllers.GetProduct)
			products.POST("", middlewares.ManagerMiddleware(), controllers.CreateProduct)
			products.PUT("/:id", middlewares.ManagerMiddleware(), controllers.UpdateProduct)
			products.PUT("/:id/restore", middlewares.ManagerMiddleware(), controllers.RestoreProduct)
			products.DELETE("/:id", middlewares.ManagerMiddleware(), controllers.DeleteProduct)
		}
	}

	// Swagger документация
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
