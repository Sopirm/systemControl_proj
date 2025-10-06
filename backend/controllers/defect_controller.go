package controllers

import (
	"net/http"
	"strconv"
	"systemControl_proj/database"
	"systemControl_proj/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// контроллер запросов связанных с дефектами
type DefectController struct {
	DB *gorm.DB
}

// создание нового экземпляра контроллера дефектов
func NewDefectController() *DefectController {
	return &DefectController{
		DB: database.DB,
	}
}

// создание нового дефекта
func (dc *DefectController) CreateDefect(c *gin.Context) {
	var defectCreate models.DefectCreate

	if err := c.ShouldBindJSON(&defectCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	// Получение роли пользователя из контекста
	userRole, _ := c.Get("role")

	// Проверка существования проекта
	var project models.Project
	if result := dc.DB.First(&project, defectCreate.ProjectID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "указанный проект не найден"})
		return
	}

	// Проверка существования исполнителя, если он указан
	if defectCreate.AssigneeID != 0 {
		var assignee models.User
		if result := dc.DB.First(&assignee, defectCreate.AssigneeID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "указанный исполнитель не найден"})
			return
		}
		// Если текущий пользователь инженер, он может назначать только инженера
		if role, ok := userRole.(models.Role); ok && role == models.RoleEngineer {
			if assignee.Role != models.RoleEngineer {
				c.JSON(http.StatusForbidden, gin.H{"error": "инженер может назначать исполнителем только инженера"})
				return
			}
		}
	}

	// Создание нового дефекта
	defect := models.Defect{
		Title:       defectCreate.Title,
		Description: defectCreate.Description,
		ProjectID:   defectCreate.ProjectID,
		Status:      models.DefectStatusNew,
		ReporterID:  userID.(uint),
		AssigneeID:  defectCreate.AssigneeID,
		DueDate:     defectCreate.DueDate,
	}

	// Установка приоритета по умолчанию, если не указан
	if defectCreate.Priority != "" {
		defect.Priority = defectCreate.Priority
	} else {
		defect.Priority = models.DefectPriorityMedium
	}

	// Сохранение дефекта в базе данных
	if result := dc.DB.Create(&defect); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при сохранении дефекта"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "дефект успешно создан",
		"defect":  defect,
	})
}

// получение списка всех дефектов
func (dc *DefectController) GetAllDefects(c *gin.Context) {
	var defects []models.Defect

	// получение параметров запроса для фильтрации
	// Сначала проверяем параметр из URL пути (для маршрута /projects/:id/defects)
	projectID := c.Param("id")
	// Если параметра в пути нет, проверяем query-параметр
	if projectID == "" {
		projectID = c.Query("project_id")
	}

	status := c.Query("status")
	priority := c.Query("priority")
	assigneeID := c.Query("assignee_id")
	reporterID := c.Query("reporter_id")

	query := dc.DB.Preload("Project").Preload("Reporter").Preload("Assignee")

	// применение фильтров
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if assigneeID != "" {
		query = query.Where("assignee_id = ?", assigneeID)
	}
	if reporterID != "" {
		query = query.Where("reporter_id = ?", reporterID)
	}

	// Получение дефектов из базы данных
	if result := query.Find(&defects); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении дефектов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"defects": defects,
	})
}

// получение конкретного дефекта по ID
func (dc *DefectController) GetDefect(c *gin.Context) {
	id := c.Param("id")

	var defect models.Defect
	if result := dc.DB.Preload("Project").Preload("Reporter").Preload("Assignee").Preload("Comments").Preload("Comments.User").First(&defect, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "дефект не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"defect": defect,
	})
}

// обновление существующего дефекта
func (dc *DefectController) UpdateDefect(c *gin.Context) {
	id := c.Param("id")
	defectID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID дефекта"})
		return
	}

	var defectUpdate models.DefectUpdate
	if err := c.ShouldBindJSON(&defectUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка существования дефекта
	var defect models.Defect
	if result := dc.DB.First(&defect, defectID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "дефект не найден"})
		return
	}

	// Обновление полей дефекта
	if defectUpdate.Title != "" {
		defect.Title = defectUpdate.Title
	}
	if defectUpdate.Description != "" {
		defect.Description = defectUpdate.Description
	}
	if defectUpdate.Status != "" {
		defect.Status = defectUpdate.Status
	}
	if defectUpdate.Priority != "" {
		defect.Priority = defectUpdate.Priority
	}
	if defectUpdate.AssigneeID != 0 {
		var assignee models.User
		if result := dc.DB.First(&assignee, defectUpdate.AssigneeID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "указанный исполнитель не найден"})
			return
		}
		// Если текущий пользователь инженер, он может назначать только инженера
		userRole, _ := c.Get("role")
		if role, ok := userRole.(models.Role); ok && role == models.RoleEngineer {
			if assignee.Role != models.RoleEngineer {
				c.JSON(http.StatusForbidden, gin.H{"error": "инженер может назначать исполнителем только инженера"})
				return
			}
		}
		defect.AssigneeID = defectUpdate.AssigneeID
	}
	if !defectUpdate.DueDate.IsZero() {
		defect.DueDate = defectUpdate.DueDate
	}

	if result := dc.DB.Save(&defect); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обновлении дефекта"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "дефект успешно обновлен",
		"defect":  defect,
	})
}

// удаление дефекта
func (dc *DefectController) DeleteDefect(c *gin.Context) {
	id := c.Param("id")

	// проверка существования дефекта
	var defect models.Defect
	if result := dc.DB.First(&defect, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "дефект не найден"})
		return
	}

	// удаление дефекта из базы данных (мягкое удаление с помощью DeletedAt)
	if result := dc.DB.Delete(&defect); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при удалении дефекта"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "дефект успешно удален",
	})
}
