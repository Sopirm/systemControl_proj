package controllers

import (
	"net/http"
	"strconv"
	"systemControl_proj/database"
	"systemControl_proj/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// контроллер запросов связанных с проектами
type ProjectController struct {
	DB *gorm.DB
}

// создание нового экземпляра контроллера проектов
func NewProjectController() *ProjectController {
	return &ProjectController{
		DB: database.DB,
	}
}

// создание нового проекта
func (pc *ProjectController) CreateProject(c *gin.Context) {
	var projectCreate models.ProjectCreate

	if err := c.ShouldBindJSON(&projectCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования менеджера
	var manager models.User
	if result := pc.DB.First(&manager, projectCreate.ManagerID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "указанный менеджер не найден"})
		return
	}

	// Создание нового проекта
	project := models.Project{
		Name:        projectCreate.Name,
		Description: projectCreate.Description,
		Location:    projectCreate.Location,
		StartDate:   projectCreate.StartDate,
		EndDate:     projectCreate.EndDate,
		Status:      projectCreate.Status,
		ManagerID:   projectCreate.ManagerID,
	}

	// Установка статуса по умолчанию, если не указан
	if project.Status == "" {
		project.Status = models.ProjectStatusActive
	}

	// Сохранение проекта в базе данных
	if result := pc.DB.Create(&project); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при сохранении проекта"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "проект успешно создан",
		"project": project,
	})
}

// получение списка всех проектов
func (pc *ProjectController) GetAllProjects(c *gin.Context) {
	var projects []models.Project

	// Получение параметров запроса для фильтрации
	status := c.Query("status")
	managerID := c.Query("manager_id")

	query := pc.DB.Preload("Manager")

	// Применение фильтров
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if managerID != "" {
		query = query.Where("manager_id = ?", managerID)
	}

	// Получение проектов из базы данных
	if result := query.Find(&projects); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении проектов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

// получение конкретного проекта по ID
func (pc *ProjectController) GetProject(c *gin.Context) {
	id := c.Param("id")

	var project models.Project
	if result := pc.DB.Preload("Manager").First(&project, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "проект не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project": project,
	})
}

// обновление существующего проекта
func (pc *ProjectController) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	projectID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID проекта"})
		return
	}

	var projectUpdate models.ProjectUpdate
	if err := c.ShouldBindJSON(&projectUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка существования проекта
	var project models.Project
	if result := pc.DB.First(&project, projectID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "проект не найден"})
		return
	}

	// Обновление полей проекта
	if projectUpdate.Name != "" {
		project.Name = projectUpdate.Name
	}
	if projectUpdate.Description != "" {
		project.Description = projectUpdate.Description
	}
	if projectUpdate.Location != "" {
		project.Location = projectUpdate.Location
	}
	if !projectUpdate.StartDate.IsZero() {
		project.StartDate = projectUpdate.StartDate
	}
	if !projectUpdate.EndDate.IsZero() {
		project.EndDate = projectUpdate.EndDate
	}
	if projectUpdate.Status != "" {
		project.Status = projectUpdate.Status
	}
	if projectUpdate.ManagerID != 0 {
		// Проверка существования менеджера
		var manager models.User
		if result := pc.DB.First(&manager, projectUpdate.ManagerID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "указанный менеджер не найден"})
			return
		}
		project.ManagerID = projectUpdate.ManagerID
	}

	if result := pc.DB.Save(&project); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обновлении проекта"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "проект успешно обновлен",
		"project": project,
	})
}

// удаление проекта
func (pc *ProjectController) DeleteProject(c *gin.Context) {
	id := c.Param("id")

	// Проверка существования проекта
	var project models.Project
	if result := pc.DB.First(&project, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "проект не найден"})
		return
	}

	// Удаление проекта из базы данных (мягкое удаление с помощью DeletedAt)
	if result := pc.DB.Delete(&project); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при удалении проекта"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "проект успешно удален",
	})
}
