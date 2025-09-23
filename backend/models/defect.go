package models

import (
	"time"

	"gorm.io/gorm"
)

// статус дефекта
type DefectStatus string

const (
	DefectStatusNew        DefectStatus = "new"
	DefectStatusInProgress DefectStatus = "in_progress"
	DefectStatusReview     DefectStatus = "review"
	DefectStatusClosed     DefectStatus = "closed"
	DefectStatusCanceled   DefectStatus = "canceled"
)

// приоритет дефекта
type DefectPriority string

const (
	DefectPriorityLow    DefectPriority = "low"
	DefectPriorityMedium DefectPriority = "medium"
	DefectPriorityHigh   DefectPriority = "high"
)

// модель дефекта на строительном объекте
type Defect struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	ProjectID   uint           `json:"project_id"`
	Project     Project        `json:"project" gorm:"foreignKey:ProjectID"`
	Status      DefectStatus   `json:"status" gorm:"type:varchar(20);default:'new'"`
	Priority    DefectPriority `json:"priority" gorm:"type:varchar(10);default:'medium'"`
	ReporterID  uint           `json:"reporter_id"`
	Reporter    User           `json:"reporter" gorm:"foreignKey:ReporterID"`
	AssigneeID  uint           `json:"assignee_id"`
	Assignee    User           `json:"assignee" gorm:"foreignKey:AssigneeID"`
	DueDate     time.Time      `json:"due_date"`
	Comments    []Comment      `json:"comments" gorm:"foreignKey:DefectID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// данные для создания дефекта
type DefectCreate struct {
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description"`
	ProjectID   uint           `json:"project_id" binding:"required"`
	Priority    DefectPriority `json:"priority"`
	AssigneeID  uint           `json:"assignee_id"`
	DueDate     time.Time      `json:"due_date"`
}

// данные для обновления дефекта
type DefectUpdate struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      DefectStatus   `json:"status"`
	Priority    DefectPriority `json:"priority"`
	AssigneeID  uint           `json:"assignee_id"`
	DueDate     time.Time      `json:"due_date"`
}
