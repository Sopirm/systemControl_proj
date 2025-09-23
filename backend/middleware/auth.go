package middleware

import (
	"errors"
	"net/http"
	"strings"
	"systemControl_proj/config"
	"systemControl_proj/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// структура для JWT токена
type JWTClaims struct {
	UserID   uint        `json:"user_id"`
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
	jwt.RegisteredClaims
}

// создание нового JWT токена для пользователя
func GenerateToken(user *models.User, cfg *config.Config) (string, error) {
	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.Auth.JWTExpirationHrs))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.Auth.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// проверка JWT токена и возвращает claims
func ValidateToken(tokenString string, cfg *config.Config) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Auth.JWTSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("некорректный токен")
}

// проверка JWT токена в заголовке запроса
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "токен авторизации не предоставлен"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "некорректный формат токена"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(tokenParts[1], cfg)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен: " + err.Error()})
			c.Abort()
			return
		}

		// сохранение данных пользователя в контексте
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// проверка роли пользователя
func RoleMiddleware(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
			c.Abort()
			return
		}

		role := userRole.(models.Role)
		authorized := false

		for _, r := range roles {
			if role == r {
				authorized = true
				break
			}
		}

		if !authorized {
			c.JSON(http.StatusForbidden, gin.H{"error": "недостаточно прав для доступа"})
			c.Abort()
			return
		}

		c.Next()
	}
}
