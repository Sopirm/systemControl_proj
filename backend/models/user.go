package models

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleObserver Role = "observer"
	RoleEngineer Role = "engineer"
	RoleManager  Role = "manager"
)

// модель пользователя системы
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"unique;not null"`
	Email        string         `json:"email" gorm:"unique;not null"`
	PasswordHash string         `json:"-" gorm:"not null"`
	FullName     string         `json:"full_name"`
	Role         Role           `json:"role" gorm:"type:varchar(20);default:'observer'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// данные для регистрации пользователя
type UserRegistration struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// данные для входа пользователя
type UserLogin struct {
	Username string `json:"username" binding:"required"` // Username может быть именем пользователя или email
	Password string `json:"password" binding:"required"`
}

// хеширует пароль пользователя используя bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хеширования пароля: %v", err)
		return "", err
	}
	log.Printf("Пароль успешно хеширован, длина хеша: %d", len(string(bytes)))
	return string(bytes), nil
}

// проверяет пароль пользователя
func (user *User) CheckPassword(password string) error {
	log.Printf("Проверка пароля, длина хеша в БД: %d", len(user.PasswordHash))
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("Ошибка сравнения пароля: %v", err)
	} else {
		log.Printf("Пароли совпадают")
	}
	return err
}
