package controllers

import (
	"net/http"
	"systemControl_proj/config"
	"systemControl_proj/database"
	"systemControl_proj/middleware"
	"systemControl_proj/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// контроллер запросов связанных с пользователями
type UserController struct {
	DB     *gorm.DB
	Config *config.Config
}

// создание нового экземпляра контроллера пользователей
func NewUserController(config *config.Config) *UserController {
	return &UserController{
		DB:     database.DB,
		Config: config,
	}
}

// регистрация нового пользователя
func (uc *UserController) Register(c *gin.Context) {
	var userReg models.UserRegistration

	if err := c.ShouldBindJSON(&userReg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка, что пользователь с таким именем не существует
	var existingUser models.User
	if result := uc.DB.Where("username = ? OR email = ?", userReg.Username, userReg.Email).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь с таким именем или email уже существует"})
		return
	}

	// Хеширование пароля
	hashedPassword, err := models.HashPassword(userReg.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при хешировании пароля"})
		return
	}

	// Создание нового пользователя
	user := models.User{
		Username:     userReg.Username,
		Email:        userReg.Email,
		PasswordHash: hashedPassword,
		FullName:     userReg.FullName,
		Role:         userReg.Role,
	}

	// Сохранение пользователя в базе данных
	if result := uc.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при сохранении пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "пользователь успешно зарегистрирован",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// вход пользователя в систему
func (uc *UserController) Login(c *gin.Context) {
	var userLogin models.UserLogin

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Поиск пользователя по имени
	var user models.User
	if result := uc.DB.Where("username = ?", userLogin.Username).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверное имя пользователя или пароль"})
		return
	}

	// Проверка пароля
	if err := user.CheckPassword(userLogin.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверное имя пользователя или пароль"})
		return
	}

	// Создание JWT токена
	token, err := middleware.GenerateToken(&user, uc.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "успешная авторизация",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

// получение профиля текущего пользователя
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	var user models.User
	if result := uc.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}
