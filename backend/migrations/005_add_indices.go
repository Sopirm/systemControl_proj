package migrations

import (
	"gorm.io/gorm"
)

// AddIndices миграция для добавления индексов к таблицам
type AddIndices struct{}

// Up создает индексы для таблиц
func (m *AddIndices) Up(tx *gorm.DB) error {
	// Индексы для пользователей
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`).Error; err != nil {
		return err
	}

	// Индексы для проектов
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_projects_manager_id ON projects(manager_id)`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status)`).Error; err != nil {
		return err
	}

	// Индексы для дефектов
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_defects_project_id ON defects(project_id)`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_defects_status ON defects(status)`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_defects_priority ON defects(priority)`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_defects_reporter_id ON defects(reporter_id)`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_defects_assignee_id ON defects(assignee_id)`).Error; err != nil {
		return err
	}

	// Индексы для комментариев
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_comments_defect_id ON comments(defect_id)`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id)`).Error; err != nil {
		return err
	}

	return nil
}

// Down удаляет созданные индексы
func (m *AddIndices) Down(tx *gorm.DB) error {
	// Удаление индексов для пользователей
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_users_username`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_users_email`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_users_role`).Error; err != nil {
		return err
	}

	// Удаление индексов для проектов
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_projects_manager_id`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_projects_status`).Error; err != nil {
		return err
	}

	// Удаление индексов для дефектов
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_defects_project_id`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_defects_status`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_defects_priority`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_defects_reporter_id`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_defects_assignee_id`).Error; err != nil {
		return err
	}

	// Удаление индексов для комментариев
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_comments_defect_id`).Error; err != nil {
		return err
	}
	if err := tx.Exec(`DROP INDEX IF EXISTS idx_comments_user_id`).Error; err != nil {
		return err
	}

	return nil
}

// Name возвращает имя миграции
func (m *AddIndices) Name() string {
	return "005_add_indices"
}
