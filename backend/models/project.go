package models

import (
	"time"

	"gorm.io/gorm"
)

// статус проекта
type ProjectStatus string

const (
	ProjectStatusActive    ProjectStatus = "active"
	ProjectStatusCompleted ProjectStatus = "completed"
	ProjectStatusSuspended ProjectStatus = "suspended"
)

// модель строительного проекта/объекта
type Project struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Location    string         `json:"location"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	Status      ProjectStatus  `json:"status" gorm:"type:varchar(20);default:'active'"`
	ManagerID   uint           `json:"manager_id"`
	Manager     User           `json:"manager" gorm:"foreignKey:ManagerID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// данные для создания проекта
type ProjectCreate struct {
	Name        string        `json:"name" binding:"required"`
	Description string        `json:"description"`
	Location    string        `json:"location" binding:"required"`
	StartDate   time.Time     `json:"start_date" binding:"required"`
	EndDate     time.Time     `json:"end_date" binding:"required"`
	Status      ProjectStatus `json:"status"`
	ManagerID   uint          `json:"manager_id" binding:"required"`
}

// данные для обновления проекта
type ProjectUpdate struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Location    string        `json:"location"`
	StartDate   time.Time     `json:"start_date"`
	EndDate     time.Time     `json:"end_date"`
	Status      ProjectStatus `json:"status"`
	ManagerID   uint          `json:"manager_id"`
}
