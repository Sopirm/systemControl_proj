# 🏗️ «СистемаКонтроля» для строительных объектов

![Golang](https://img.shields.io/badge/Backend-Golang%2FGin-00ADD8)
![Vue.js](https://img.shields.io/badge/Frontend-Vue.js-4FC08D)
![PostgreSQL](https://img.shields.io/badge/БД-PostgreSQL-336791)

Монолитное веб-приложение для централизованного управления дефектами на строительных объектах. Система обеспечивает полный цикл работы: от регистрации дефекта и назначения исполнителя до контроля статусов и формирования отчётности.

## 📋 Назначение системы

Система предназначена для:
* **Инженеров** — регистрация дефектов, обновление информации
* **Менеджеров** — назначение задач, контроль сроков, формирование отчётов
* **Руководителей и заказчиков** — просмотр прогресса и отчётности

## ⚙️ Функциональные возможности

- ✅ Регистрация пользователей и аутентификация
- 👥 Разграничение прав доступа (менеджер, инженер, наблюдатель)
- 🏢 Управление проектами/объектами и их этапами
- 🔍 Создание и редактирование дефектов (заголовок, описание, приоритет, исполнитель, сроки, вложения)
- 📊 Управление статусами дефектов: Новая → В работе → На проверке → Закрыта/Отменена
- 💬 Ведение комментариев и истории изменений
- 🔎 Поиск, сортировка и фильтрация дефектов
- 📈 Экспорт отчётности в CSV/Excel
- 📊 Просмотр аналитических отчётов (графики, статистика)

## 🧰 Технологический стек

### Backend
- **Golang** с **Gin Framework** — основа серверной части
- **PostgreSQL** — реляционная СУБД
- **GORM** — ORM для работы с базой данных
- **JWT** — токены для авторизации
- **Docker** — контейнеризация

### Frontend
- **Vue.js** — основной фреймворк
- **TypeScript** — типизация JavaScript
- **Vite** — быстрая сборка
- **Vue Router** — маршрутизация
- **Chart.js** — визуализация данных (планируется)

### Тестирование
- **Playwright** — e2e тестирование

## 🗄️ Структура базы данных

### Таблица Users
- id (PK)
- username
- email
- password_hash
- full_name
- role
- created_at
- updated_at

### Таблица Projects
- id (PK)
- name
- description
- location
- start_date
- end_date
- status
- manager_id (FK -> Users)
- created_at
- updated_at

### Таблица Defects
- id (PK)
- title
- description
- project_id (FK -> Projects)
- status (enum: new, in_progress, review, closed, canceled)
- priority (enum: low, medium, high)
- reporter_id (FK -> Users)
- assignee_id (FK -> Users)
- due_date
- created_at
- updated_at

### Таблица Comments
- id (PK)
- defect_id (FK -> Defects)
- user_id (FK -> Users)
- content
- created_at

## 🔒 Безопасность

- JWT токены для авторизации
- Хеширование паролей через bcrypt
- RBAC (Role-Based Access Control)
- Защита от SQL-инъекций (параметризованные запросы)
- Защита от XSS (санитация входных данных)
- Валидация всех входящих данных

## 🚀 Запуск проекта

Для запуска проекта вам потребуется установить Go и Node.js на вашей системе.

### 1. Клонирование репозитория

```bash
git clone https://github.com/Sopirm/systemControl_proj
cd systemControl_proj
```

### 2. Запуск Backend (GoLang Gin)

```bash
cd backend
go mod tidy
go run main.go
```

Сервер будет запущен на `http://localhost:8080`.

### 3. Запуск Frontend (Vue.js)

```bash
cd frontend
npm install
npm run dev
```

Фронтенд будет доступен по адресу, указанному в консоли (обычно `http://localhost:5173`).

---

Ссылка на GitHub: [https://github.com/Sopirm/systemControl_proj](https://github.com/Sopirm/systemControl_proj)