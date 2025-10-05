<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseCard from '../components/BaseCard.vue'
import { authService } from '../services/api'
import { useRouter } from 'vue-router'

const router = useRouter()
const isLoggedIn = ref(authService.isLoggedIn())

// Функция для проверки авторизации
const checkAuth = () => {
  isLoggedIn.value = authService.isLoggedIn()
  
  // Если пользователь авторизован, перенаправляем на страницу проектов
  if (isLoggedIn.value) {
    router.push('/projects')
  }
}

// Проверяем при монтировании компонента
onMounted(() => {
  checkAuth()
  
  // Слушаем события авторизации
  window.addEventListener('auth-change', checkAuth)
})
</script>

<template>
  <div class="home">
    <section class="hero-section">
      <h1>Система контроля</h1>
      <p class="subtitle">Управление дефектами строительных объектов</p>
      <div class="hero-actions">
        <RouterLink v-if="!isLoggedIn" to="/login" class="btn">Войти</RouterLink>
        <RouterLink v-if="!isLoggedIn" to="/register" class="btn btn-outline">Регистрация</RouterLink>
        <RouterLink v-else to="/projects" class="btn">Проекты</RouterLink>
      </div>
    </section>

    <section class="features-section">
      <h2 class="section-title">Возможности системы</h2>
      <div class="features-grid">
        <BaseCard title="Управление проектами">
          <p>Централизованное управление проектами и объектами строительства.</p>
        </BaseCard>

        <BaseCard title="Учет дефектов">
          <p>Регистрация, классификация и отслеживание дефектов на объектах.</p>
        </BaseCard>

        <BaseCard title="Назначение задач">
          <p>Распределение заданий по исполнителям с контролем сроков.</p>
        </BaseCard>

        <BaseCard title="Аналитика и отчеты">
          <p>Формирование отчетов и аналитических данных для руководства.</p>
        </BaseCard>
      </div>
    </section>
  </div>
</template>

<style scoped>
.home {
  text-align: center;
}

.hero-section {
  padding: 3rem 0;
  margin-bottom: 2rem;
}

h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--color-primary);
  margin-bottom: 0.5rem;
}

.subtitle {
  font-size: 1.2rem;
  color: var(--color-text-light);
  margin-bottom: 2rem;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.section-title {
  font-size: 1.8rem;
  color: var(--color-primary);
  margin-bottom: 2rem;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
}

@media (max-width: 768px) {
  .hero-actions {
    flex-direction: column;
    align-items: center;
  }

  .hero-actions .btn {
    width: 100%;
    max-width: 300px;
  }
}
</style>