package controllers

import (
	"log"
	"net/http"
	"strconv"
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

	// Отладочная информация
	log.Printf("Получен запрос на регистрацию: username=%s, email=%s, role=%s",
		userReg.Username, userReg.Email, userReg.Role)

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

	log.Printf("Создание пользователя: %+v", user)

	// Сохранение пользователя в базе данных
	if result := uc.DB.Create(&user); result.Error != nil {
		log.Printf("Ошибка при сохранении пользователя: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при сохранении пользователя"})
		return
	}

	log.Printf("Пользователь успешно создан: ID=%d, username=%s, role=%s",
		user.ID, user.Username, user.Role)

	c.JSON(http.StatusCreated, gin.H{
		"message": "пользователь успешно зарегистрирован",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"full_name": user.FullName,
			"role":     user.Role,
		},
	})
}

// вход пользователя в систему
func (uc *UserController) Login(c *gin.Context) {
	var userLogin models.UserLogin

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		log.Printf("Ошибка при разборе JSON запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Попытка входа пользователя: %s", userLogin.Username)

	// Поиск пользователя по имени пользователя или email
	var user models.User
	if result := uc.DB.Where("username = ? OR email = ?", userLogin.Username, userLogin.Username).First(&user); result.Error != nil {
		log.Printf("Пользователь не найден: %s, ошибка: %v", userLogin.Username, result.Error)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверное имя пользователя или пароль"})
		return
	}

	log.Printf("Пользователь найден: %s, email: %s, ID: %d, роль: %s", user.Username, user.Email, user.ID, user.Role)

	// Проверка пароля
	if err := user.CheckPassword(userLogin.Password); err != nil {
		log.Printf("Неверный пароль для пользователя %s: %v", userLogin.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверное имя пользователя или пароль"})
		return
	}

	log.Printf("Пароль верный для пользователя: %s", userLogin.Username)

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
			"full_name": user.FullName,
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

// получение списка всех пользователей (только для менеджеров)
func (uc *UserController) GetAllUsers(c *gin.Context) {
	// Проверка роли выполняется в middleware

	var users []models.User
	if result := uc.DB.Find(&users); result.Error != nil {
		log.Printf("Ошибка при получении списка пользователей: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении списка пользователей"})
		return
	}

	// Создаем безопасные для передачи объекты (без хешей паролей)
	safeUsers := make([]gin.H, len(users))
	for i, user := range users {
		safeUsers[i] = gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"full_name":  user.FullName,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users": safeUsers,
		"count": len(users),
	})
}

// получение списка инженеров (доступно менеджерам и инженерам)
func (uc *UserController) GetEngineers(c *gin.Context) {
    // Получаем роль и ID пользователя из контекста
    userRole, _ := c.Get("role")
    userID, _ := c.Get("userID")

    var users []models.User
    query := uc.DB.Where("role = ?", models.RoleEngineer)

    // Если запрос делает инженер, исключаем его самого из результата
    if userRole == models.RoleEngineer {
        if uid, ok := userID.(uint); ok {
            query = query.Where("id <> ?", uid)
        }
    }

    if result := query.Find(&users); result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении списка инженеров"})
        return
    }

    // Создаем безопасные для передачи объекты (без хешей паролей)
    safeUsers := make([]gin.H, len(users))
    for i, user := range users {
        safeUsers[i] = gin.H{
            "id":         user.ID,
            "username":   user.Username,
            "email":      user.Email,
            "full_name":  user.FullName,
            "role":       user.Role,
            "created_at": user.CreatedAt,
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "users": safeUsers,
        "count": len(users),
    })
}

// обновление роли пользователя (только для менеджеров)
func (uc *UserController) UpdateUserRole(c *gin.Context) {
	// Проверка роли выполняется в middleware

	// ID цели из маршрута и проверка корректности
	userIDStr := c.Param("id")
	targetID, convErr := strconv.Atoi(userIDStr)
	if convErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID пользователя"})
		return
	}

	// Запрещаем менять роль у самого себя
	if ctxUserID, ok := c.Get("userID"); ok {
		if uid, ok2 := ctxUserID.(uint); ok2 && uid == uint(targetID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "нельзя изменять роль собственного аккаунта"})
			return
		}
	}

	var roleUpdate struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&roleUpdate); err != nil {
		log.Printf("Ошибка при разборе JSON запроса обновления роли: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка, что роль допустима
	if roleUpdate.Role != string(models.RoleManager) &&
		roleUpdate.Role != string(models.RoleEngineer) &&
		roleUpdate.Role != string(models.RoleObserver) {
		log.Printf("Недопустимая роль: %s", roleUpdate.Role)
		c.JSON(http.StatusBadRequest, gin.H{"error": "недопустимая роль"})
		return
	}

	// Поиск пользователя
	var user models.User
	if result := uc.DB.First(&user, targetID); result.Error != nil {
		log.Printf("Пользователь для обновления роли не найден: ID=%d", targetID)
		c.JSON(http.StatusNotFound, gin.H{"error": "пользователь не найден"})
		return
	}

	// Обновление роли
	user.Role = models.Role(roleUpdate.Role)
	if result := uc.DB.Save(&user); result.Error != nil {
		log.Printf("Ошибка при обновлении роли пользователя: %v", result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обновлении роли пользователя"})
		return
	}

	log.Printf("Роль пользователя %s (ID=%d) изменена на %s", user.Username, user.ID, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"message": "роль пользователя успешно обновлена",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}
