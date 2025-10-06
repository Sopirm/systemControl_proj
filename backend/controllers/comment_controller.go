package controllers

import (
	"net/http"
	"systemControl_proj/database"
	"systemControl_proj/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// контроллер запросов связанных с комментариями
type CommentController struct {
	DB *gorm.DB
}

// создание нового экземпляра контроллера комментариев
func NewCommentController() *CommentController {
	return &CommentController{
		DB: database.DB,
	}
}

// создание нового комментария к дефекту
func (cc *CommentController) CreateComment(c *gin.Context) {
	var commentCreate models.CommentCreate

	if err := c.ShouldBindJSON(&commentCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получение ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	// Проверка существования дефекта
	var defect models.Defect
	if result := cc.DB.First(&defect, commentCreate.DefectID); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "указанный дефект не найден"})
		return
	}

	// Создание нового комментария
	comment := models.Comment{
		DefectID: commentCreate.DefectID,
		UserID:   userID.(uint),
		Content:  commentCreate.Content,
	}

	// Сохранение комментария в базе данных
	if result := cc.DB.Create(&comment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при сохранении комментария"})
		return
	}

	cc.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "комментарий успешно создан",
		"comment": comment,
	})
}

// получение всех комментариев для конкретного дефекта
func (cc *CommentController) GetDefectComments(c *gin.Context) {
	defectID := c.Param("id")

	var comments []models.Comment
	if result := cc.DB.Where("defect_id = ?", defectID).Preload("User").Find(&comments); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении комментариев"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}

// удаление комментария
func (cc *CommentController) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	// получение ID пользователя из контекста
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	// Получение роли пользователя из контекста
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "роль пользователя не определена"})
		return
	}

	// Проверка существования комментария
	var comment models.Comment
	if result := cc.DB.First(&comment, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "комментарий не найден"})
		return
	}

	// Проверка прав на удаление комментария (только автор или менеджер)
	if comment.UserID != userID.(uint) && userRole.(models.Role) != models.RoleManager {
		c.JSON(http.StatusForbidden, gin.H{"error": "у вас нет прав на удаление этого комментария"})
		return
	}

	// Удаление комментария из базы данных
	if result := cc.DB.Delete(&comment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при удалении комментария"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "комментарий успешно удален",
	})
}
