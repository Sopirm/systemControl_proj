package main

import (
	"flag"
	"fmt"
	"log"
	"systemControl_proj/config"
	"systemControl_proj/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Парсинг флагов командной строки
	var rollback bool
	flag.BoolVar(&rollback, "rollback", false, "отменить последнюю миграцию")
	flag.Parse()

	// Загрузка конфигурации
	cfg := config.GetConfig()

	// Подключение к базе данных
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	// Получение соединения для проверки подключения
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("ошибка получения соединения с базой данных: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("ошибка проверки подключения к базе данных: %v", err)
	}

	log.Println("Соединение с базой данных установлено")

	// Выполнение миграций или отмена последней миграции
	if rollback {
		if err := migrations.RollbackLastMigration(db); err != nil {
			log.Fatalf("ошибка при отмене миграции: %v", err)
		}
	} else {
		if err := migrations.RunAllMigrations(db); err != nil {
			log.Fatalf("ошибка при выполнении миграций: %v", err)
		}
	}
}
