package migrations

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// представляет структуру миграции
type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null;unique"`
	Applied   bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// интерфейс для миграций
type Migrator interface {
	Up(tx *gorm.DB) error
	Down(tx *gorm.DB) error
	Name() string
}

// сервис для управления миграциями
type MigrationService struct {
	DB         *gorm.DB
	Migrations []Migrator
}

// создает новый сервис миграций
func NewMigrationService(db *gorm.DB, migrations []Migrator) *MigrationService {
	return &MigrationService{
		DB:         db,
		Migrations: migrations,
	}
}

// создает таблицу для отслеживания миграций
func (ms *MigrationService) SetupMigrationTable() error {
	return ms.DB.AutoMigrate(&Migration{})
}

// запускает все неприменённые миграции
func (ms *MigrationService) RunMigrations() error {
	if err := ms.SetupMigrationTable(); err != nil {
		return fmt.Errorf("ошибка создания таблицы миграций: %w", err)
	}

	for _, migrator := range ms.Migrations {
		var migration Migration
		migrationName := migrator.Name()

		// Проверка, была ли уже применена миграция
		result := ms.DB.Where("name = ?", migrationName).First(&migration)

		if result.Error == nil && migration.Applied {
			log.Printf("Миграция '%s' уже применена", migrationName)
			continue
		}

		log.Printf("Применение миграции '%s'", migrationName)

		// Начало транзакции
		tx := ms.DB.Begin()

		// Применение миграции
		if err := migrator.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("ошибка применения миграции '%s': %w", migrationName, err)
		}

		// Сохранение информации о миграции
		if result.Error != nil {
			// Миграция ещё не существует в таблице
			newMigration := Migration{
				Name:    migrationName,
				Applied: true,
			}
			if err := tx.Create(&newMigration).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("ошибка сохранения информации о миграции '%s': %w", migrationName, err)
			}
		} else {
			// Обновление существующей миграции
			if err := tx.Model(&migration).Update("applied", true).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("ошибка обновления информации о миграции '%s': %w", migrationName, err)
			}
		}

		// Завершение транзакции
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("ошибка применения транзакции для миграции '%s': %w", migrationName, err)
		}

		log.Printf("Миграция '%s' успешно применена", migrationName)
	}

	return nil
}

// отменяет последнюю миграцию
func (ms *MigrationService) RollbackMigration() error {
	var lastMigration Migration
	if err := ms.DB.Where("applied = ?", true).Order("id DESC").First(&lastMigration).Error; err != nil {
		return fmt.Errorf("нет миграций для отмены: %w", err)
	}

	// Поиск соответствующей миграции
	var migrator Migrator
	for _, m := range ms.Migrations {
		if m.Name() == lastMigration.Name {
			migrator = m
			break
		}
	}

	if migrator == nil {
		return fmt.Errorf("не найден код миграции '%s' для отмены", lastMigration.Name)
	}

	// Начало транзакции
	tx := ms.DB.Begin()

	// Отмена миграции
	if err := migrator.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка отмены миграции '%s': %w", lastMigration.Name, err)
	}

	// Обновление информации о миграции
	if err := tx.Model(&lastMigration).Update("applied", false).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ошибка обновления информации о миграции '%s': %w", lastMigration.Name, err)
	}

	// Завершение транзакции
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("ошибка применения транзакции для отмены миграции '%s': %w", lastMigration.Name, err)
	}

	log.Printf("Миграция '%s' успешно отменена", lastMigration.Name)
	return nil
}
