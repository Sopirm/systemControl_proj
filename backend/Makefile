.PHONY: run build migrate rollback

# Запуск сервера
run:
	go run main.go

# Сборка приложения
build:
	go build -o server main.go

# Миграции
migrate:
	go run cmd/migrate/main.go

# Отмена последней миграции
rollback:
	go run cmd/migrate/main.go --rollback

# Установка зависимостей
deps:
	go mod tidy

# Очистка сгенерированных файлов
clean:
	rm -f server
