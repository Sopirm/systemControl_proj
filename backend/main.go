package main

import (
	"log"
	"systemControl_proj/config"
	"systemControl_proj/database"
	"systemControl_proj/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.GetConfig()

	db, err := database.SetupDatabase(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка получения соединения с базой данных: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ошибка проверки соединения с базой данных: %v", err)
	}

	log.Println("Соединение с базой данных установлено")

	router := gin.Default()

	routes.SetupRoutes(router, cfg)

	log.Printf("Сервер запущен на порту %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
