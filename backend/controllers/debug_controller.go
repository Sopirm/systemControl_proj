package controllers

import (
	"log"
	"net/http"
	"systemControl_proj/config"
	"systemControl_proj/database"
	"systemControl_proj/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DebugController - контроллер для отладочных функций
type DebugController struct {
	DB     *gorm.DB
	Config *config.Config
}

// NewDebugController - создает новый экземпляр отладочного контроллера
func NewDebugController(config *config.Config) *DebugController {
	return &DebugController{
		DB:     database.DB,
		Config: config,
	}
}

// ResetUserPassword - сбрасывает пароль пользователя для отладки
func (dc *DebugController) ResetUserPassword(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Ошибка при разборе JSON в ResetUserPassword: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Запрос на сброс пароля для пользователя: %s", req.Username)

	var user models.User
	if result := dc.DB.Where("username = ?", req.Username).First(&user); result.Error != nil {
		log.Printf("Пользователь для сброса пароля не найден: %s", req.Username)
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	hashedPassword, err := models.HashPassword(req.NewPassword)
	if err != nil {
		log.Printf("Ошибка при хешировании нового пароля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при хешировании пароля"})
		return
	}

	// Обновляем пароль
	user.PasswordHash = hashedPassword
	if result := dc.DB.Save(&user); result.Error != nil {
		log.Printf("Ошибка при сохранении нового пароля: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при сохранении пароля"})
		return
	}

	log.Printf("Пароль успешно сброшен для пользователя %s", user.Username)

	c.JSON(http.StatusOK, gin.H{
		"message": "пароль успешно сброшен",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}

// CreateTestUser - создает тестового пользователя для отладки
func (dc *DebugController) CreateTestUser(c *gin.Context) {
	// Стандартные данные для тестового пользователя
	username := "testuser"
	email := "test@example.com"
	password := "test123"
	fullName := "Тестовый Пользователь"
	role := models.RoleManager

	// Проверяем, существует ли уже такой пользователь
	var existingUser models.User
	if result := dc.DB.Where("username = ? OR email = ?", username, email).First(&existingUser); result.Error == nil {
		log.Printf("Тестовый пользователь уже существует: %s", username)
		c.JSON(http.StatusOK, gin.H{
			"message": "тестовый пользователь уже существует",
			"user": gin.H{
				"id":        existingUser.ID,
				"username":  existingUser.Username,
				"email":     existingUser.Email,
				"full_name": existingUser.FullName,
				"role":      existingUser.Role,
			},
		})
		return
	}

	// Хешируем пароль
	hashedPassword, err := models.HashPassword(password)
	if err != nil {
		log.Printf("Ошибка при хешировании пароля для тестового пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при хешировании пароля"})
		return
	}

	// Создаем тестового пользователя
	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		Role:         role,
	}

	// Сохраняем в базу
	if result := dc.DB.Create(&user); result.Error != nil {
		log.Printf("Ошибка при создании тестового пользователя: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании тестового пользователя"})
		return
	}

	log.Printf("Тестовый пользователь успешно создан: ID=%d, username=%s", user.ID, user.Username)

	c.JSON(http.StatusCreated, gin.H{
		"message": "тестовый пользователь успешно создан",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
			"password":  password, // Отправляем пароль в открытом виде, так как это тестовый пользователь
		},
	})
}
