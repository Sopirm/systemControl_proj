package routes

import (
	"systemControl_proj/config"
	"systemControl_proj/controllers"
	"systemControl_proj/database"
	"systemControl_proj/middleware"
	"systemControl_proj/models"

	"github.com/gin-gonic/gin"
)

// настраивает маршруты для API
func SetupRoutes(router *gin.Engine, cfg *config.Config) {
	userController := controllers.NewUserController(cfg)
	projectController := controllers.NewProjectController()
	defectController := controllers.NewDefectController()
	commentController := controllers.NewCommentController()
	debugController := controllers.NewDebugController(cfg) // Отладочный контроллер

	// Middleware для CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API СистемаКонтроля",
		})
	})

	// маршруты для аутентификации
	auth := router.Group("/auth")
	{
		auth.POST("/register", userController.Register)
		auth.POST("/login", userController.Login)
	}

	// Отладочные маршруты
	debug := router.Group("/debug")
	{
		// Получение списка всех пользователей
		debug.GET("/users", func(c *gin.Context) {
			var users []models.User
			if err := database.DB.Find(&users).Error; err != nil {
				c.JSON(500, gin.H{"error": "Ошибка получения пользователей", "details": err.Error()})
				return
			}

			// Создаем безопасные для передачи объекты (без хешей паролей)
			safeUsers := make([]gin.H, len(users))
			for i, user := range users {
				safeUsers[i] = gin.H{
					"id":          user.ID,
					"username":    user.Username,
					"email":       user.Email,
					"full_name":   user.FullName,
					"role":        user.Role,
					"created_at":  user.CreatedAt,
					"hash_length": len(user.PasswordHash),
				}
			}

			c.JSON(200, gin.H{
				"users": safeUsers,
				"count": len(users),
			})
		})

		// Сброс пароля для пользователя
		debug.POST("/reset-password", debugController.ResetUserPassword)

		// Создание тестового пользователя
		debug.POST("/create-test-user", debugController.CreateTestUser)
	}

	// маршруты, требующие аутентификации
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		api.GET("/profile", userController.GetProfile)

		// Маршруты для управления пользователями (только для менеджеров)
		users := api.Group("/users")
		{
			users.GET("", middleware.RoleMiddleware(models.RoleManager), userController.GetAllUsers)
			users.PUT("/:id/role", middleware.RoleMiddleware(models.RoleManager), userController.UpdateUserRole)
		}

		projects := api.Group("/projects")
		{
			projects.GET("", projectController.GetAllProjects)
			projects.GET("/:id", projectController.GetProject)
			projects.POST("", middleware.RoleMiddleware(models.RoleManager), projectController.CreateProject)
			projects.PUT("/:id", middleware.RoleMiddleware(models.RoleManager), projectController.UpdateProject)
			projects.DELETE("/:id", middleware.RoleMiddleware(models.RoleManager), projectController.DeleteProject)
		}

		// маршруты для дефектов
		defects := api.Group("/defects")
		{
			// комментарии к дефектам
			defects.GET("/:id/comments", commentController.GetDefectComments)

			defects.GET("", defectController.GetAllDefects)
			defects.GET("/:id", defectController.GetDefect)
			defects.POST("", defectController.CreateDefect)
			defects.PUT("/:id", defectController.UpdateDefect)
			defects.DELETE("/:id", middleware.RoleMiddleware(models.RoleManager, models.RoleEngineer), defectController.DeleteDefect)

			// комментарии к дефектам
			defects.POST("/comments", commentController.CreateComment)
			defects.DELETE("/comments/:id", commentController.DeleteComment)
		}
	}
}
