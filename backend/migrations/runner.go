package migrations

import (
	"log"

	"gorm.io/gorm"
)

// возвращает список всех миграций в нужном порядке
func GetMigrations() []Migrator {
	return []Migrator{
		&CreateUsersTable{},
		&CreateProjectsTable{},
		&CreateDefectsTable{},
		&CreateCommentsTable{},
		&AddIndices{},
	}
}

// запускает все доступные миграции
func RunAllMigrations(db *gorm.DB) error {
	migrations := GetMigrations()
	service := NewMigrationService(db, migrations)

	log.Println("Запуск миграций...")
	if err := service.RunMigrations(); err != nil {
		log.Printf("Ошибка при выполнении миграций: %v", err)
		return err
	}
	log.Println("Миграции успешно выполнены")

	return nil
}

// отменяет последнюю применённую миграцию
func RollbackLastMigration(db *gorm.DB) error {
	migrations := GetMigrations()
	service := NewMigrationService(db, migrations)

	log.Println("Отмена последней миграции...")
	if err := service.RollbackMigration(); err != nil {
		log.Printf("Ошибка при отмене миграции: %v", err)
		return err
	}
	log.Println("Миграция успешно отменена")

	return nil
}
