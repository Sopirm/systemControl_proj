package models

import (
	"time"

	"gorm.io/gorm"
)

// модель комментария к дефекту
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	DefectID  uint           `json:"defect_id"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Content   string         `json:"content" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// данные для создания комментария
type CommentCreate struct {
	DefectID uint   `json:"defect_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
}
