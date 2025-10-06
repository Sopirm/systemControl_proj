package migrations

import (
	"log"
	"systemControl_proj/models"

	"gorm.io/gorm"
)

// миграция для создания учетной записи администратора по умолчанию
// Создает пользователя admin/admin с ролью manager при пустой таблице users
// Выполняется только один раз (если таблица пользователей пуста)
type SeedAdminUser struct{}

func (m *SeedAdminUser) Name() string { return "006_seed_admin_user" }

func (m *SeedAdminUser) Up(tx *gorm.DB) error {
	// Проверяем, существует ли уже пользователь admin
	var exists int64
	if err := tx.Table("users").Where("username = ?", "admin").Count(&exists).Error; err != nil {
		return err
	}
	if exists > 0 {
		log.Println("SeedAdminUser: пользователь 'admin' уже существует, пропуск")
		return nil
	}

	// Хешируем пароль для admin
	hash, err := models.HashPassword("admin")
	if err != nil {
		return err
	}

	// Создаем запись пользователя admin
	username := "admin"
	email := "admin@example.com"
	fullName := "Admin"
	role := string(models.RoleManager)

	return tx.Exec(`
		INSERT INTO users (username, email, password_hash, full_name, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, username, email, hash, fullName, role).Error
}

func (m *SeedAdminUser) Down(tx *gorm.DB) error {
	// Удаляем только пользователя с username=admin
	return tx.Exec(`DELETE FROM users WHERE username = ?`, "admin").Error
}
