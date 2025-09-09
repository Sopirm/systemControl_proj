# «СистемаКонтроля» по работе с строительными объектами
Этот проект разработан на Golang (Gin) и Vue.js.

## Используемые технологии

### Backend (GoLang Gin)

*   **Go (Golang)** 
*   **Gin Web Framework** 

### Frontend (Vue.js)

*   **Vue.js**
*   **Vite** 
*   **TypeScript** 
*   **Vue Router**

### Тестирование

* **Playwright**

## Запуск проекта на вашем компьютере

Для запуска этого проекта вам потребуется установить Go и Node.js (включая npm) на вашей системе.

### 1. Клонирование репозитория

```bash
git clone <https://github.com/Sopirm/systemControl_proj>
cd systemControl_proj
```

### 2. Запуск Backend (GoLang Gin)

Перейдите в директорию `backend`, инициализируйте Go модуль, установите зависимости и запустите сервер:

```bash
cd backend
go mod tidy
go run main.go
```

Сервер будет запущен на `http://localhost:8080`.

### 3. Запуск Frontend (Vue.js)

В новом терминале перейдите в директорию `frontend`, установите зависимости и запустите следующие команды:

```bash
cd frontend
npm install
npm run dev
```

Фронтенд будет доступен по адресу, который будет указан в консоли (обычно `http://localhost:5173`).
