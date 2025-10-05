<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'
import { useAuth } from '../composables/useAuth'
import { onMounted, ref, watch } from 'vue'

const router = useRouter()
const { logout, isAuthenticated: authStatus, isManager: managerStatus, authService } = useAuth()
const isAuthenticated = ref(authStatus.value)
const isManager = ref(managerStatus.value)

// Обновление состояния аутентификации при монтировании компонента
onMounted(() => {
  checkAuthentication()
  
  // Слушаем глобальное событие изменения авторизации
  window.addEventListener('auth-change', () => {
    console.log('Получено событие изменения авторизации')
    checkAuthentication()
  })
})

// Обновление при изменении маршрута
watch(() => router.currentRoute.value.path, () => {
  checkAuthentication()
})

function checkAuthentication() {
  // Проверяем текущее состояние авторизации
  const isLoggedIn = authService.isLoggedIn()
  console.log('Проверка авторизации:', isLoggedIn)
  isAuthenticated.value = isLoggedIn
  
  if (isLoggedIn) {
    const user = authService.getCurrentUser()
    console.log('Текущий пользователь:', user)
    isManager.value = user?.role === 'manager'
  } else {
    isManager.value = false
  }
}

const handleLogout = () => {
  logout()
  // Создаем событие выхода
  const logoutEvent = new CustomEvent('auth-change', { detail: { action: 'logout' } })
  window.dispatchEvent(logoutEvent)
  
  // Перенаправляем на страницу входа
  router.push('/login')
}
</script>

<template>
  <header class="app-header">
    <div class="container flex-between">
      <div class="logo">
        <RouterLink to="/">Система контроля</RouterLink>
      </div>
      <nav class="main-nav">
        <ul class="nav-list">
          <li v-if="isAuthenticated">
            <RouterLink to="/projects">Проекты</RouterLink>
          </li>
          <li v-if="isAuthenticated">
            <RouterLink to="/defects">Дефекты</RouterLink>
          </li>
          <li v-if="isAuthenticated && isManager">
            <RouterLink to="/reports">Отчеты</RouterLink>
          </li>
          <li v-if="isAuthenticated && isManager">
            <RouterLink to="/users">Пользователи</RouterLink>
          </li>
          <li v-if="!isAuthenticated">
            <RouterLink to="/login">Войти</RouterLink>
          </li>
          <li v-if="!isAuthenticated">
            <RouterLink to="/register">Регистрация</RouterLink>
          </li>
          <li v-if="isAuthenticated">
            <button class="logout-button" @click="handleLogout">Выйти</button>
          </li>
        </ul>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  background-color: var(--color-background);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1rem 0;
}

.logo a {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--color-primary);
}

.nav-list {
  display: flex;
  list-style: none;
  gap: 1.5rem;
  align-items: center;
}

.nav-list a {
  color: var(--color-text);
  font-weight: 500;
  transition: var(--transition-default);
}

.nav-list a:hover,
.nav-list a.router-link-active {
  color: var(--color-primary);
}

.logout-button {
  background: none;
  border: 1px solid var(--color-border);
  color: var(--color-text);
  padding: 0.4rem 0.8rem;
  border-radius: 4px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.logout-button:hover {
  background-color: #f0f0f0;
  color: var(--color-primary);
  border-color: var(--color-primary);
}

@media (max-width: 768px) {
  .app-header .container {
    flex-direction: column;
    gap: 1rem;
  }
  
  .nav-list {
    flex-wrap: wrap;
    justify-content: center;
    gap: 1rem;
  }
}
</style>
