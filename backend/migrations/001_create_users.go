package migrations

import (
	"gorm.io/gorm"
)

// миграция для создания таблицы пользователей
type CreateUsersTable struct{}

// создает таблицу пользователей
func (m *CreateUsersTable) Up(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) NOT NULL UNIQUE,
			email VARCHAR(100) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			full_name VARCHAR(100),
			role VARCHAR(20) NOT NULL DEFAULT 'observer',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		)
	`).Error
}

// удаляет таблицу пользователей
func (m *CreateUsersTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS users`).Error
}

// возвращает имя миграции
func (m *CreateUsersTable) Name() string {
	return "001_create_users_table"
}
