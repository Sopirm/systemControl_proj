package migrations

import (
	"gorm.io/gorm"
)

// миграция для создания таблицы проектов
type CreateProjectsTable struct{}

// создает таблицу проектов
func (m *CreateProjectsTable) Up(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE IF NOT EXISTS projects (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			location VARCHAR(255),
			start_date TIMESTAMP WITH TIME ZONE,
			end_date TIMESTAMP WITH TIME ZONE,
			status VARCHAR(20) DEFAULT 'active',
			manager_id INTEGER NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE,
			FOREIGN KEY (manager_id) REFERENCES users(id)
		)
	`).Error
}

// удаляет таблицу проектов
func (m *CreateProjectsTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS projects`).Error
}

// возвращает имя миграции
func (m *CreateProjectsTable) Name() string {
	return "002_create_projects_table"
}
