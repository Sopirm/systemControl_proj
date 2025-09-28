package migrations

import (
	"gorm.io/gorm"
)

// CreateDefectsTable миграция для создания таблицы дефектов
type CreateDefectsTable struct{}

// Up создает таблицу дефектов
func (m *CreateDefectsTable) Up(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE IF NOT EXISTS defects (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			project_id INTEGER NOT NULL,
			status VARCHAR(20) DEFAULT 'new',
			priority VARCHAR(10) DEFAULT 'medium',
			reporter_id INTEGER NOT NULL,
			assignee_id INTEGER,
			due_date TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE,
			FOREIGN KEY (project_id) REFERENCES projects(id),
			FOREIGN KEY (reporter_id) REFERENCES users(id),
			FOREIGN KEY (assignee_id) REFERENCES users(id)
		)
	`).Error
}

// Down удаляет таблицу дефектов
func (m *CreateDefectsTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS defects`).Error
}

// Name возвращает имя миграции
func (m *CreateDefectsTable) Name() string {
	return "003_create_defects_table"
}
