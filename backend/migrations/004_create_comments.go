package migrations

import (
	"gorm.io/gorm"
)

// CreateCommentsTable миграция для создания таблицы комментариев
type CreateCommentsTable struct{}

// Up создает таблицу комментариев
func (m *CreateCommentsTable) Up(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			defect_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE,
			FOREIGN KEY (defect_id) REFERENCES defects(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`).Error
}

// Down удаляет таблицу комментариев
func (m *CreateCommentsTable) Down(tx *gorm.DB) error {
	return tx.Exec(`DROP TABLE IF EXISTS comments`).Error
}

// Name возвращает имя миграции
func (m *CreateCommentsTable) Name() string {
	return "004_create_comments_table"
}
