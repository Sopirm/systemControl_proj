<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'

interface Project {
  id: number
  name: string
  description: string
  status: string
  defectsCount: number
  completionRate: number
}

// Временные данные для демонстрации
const projects = ref<Project[]>([
  {
    id: 1,
    name: 'ЖК "Солнечный"',
    description: 'Жилой комплекс из 3 корпусов на ул. Ленина',
    status: 'active',
    defectsCount: 12,
    completionRate: 65
  },
  {
    id: 2,
    name: 'Бизнес-центр "Меркурий"',
    description: 'Офисное здание класса А в центре города',
    status: 'active',
    defectsCount: 8,
    completionRate: 40
  },
  {
    id: 3,
    name: 'Торговый центр "Галактика"',
    description: 'Многофункциональный торговый комплекс',
    status: 'planning',
    defectsCount: 0,
    completionRate: 0
  }
])

const isLoading = ref(false)
</script>

<template>
  <div class="projects-page">
    <div class="page-header flex-between">
      <h1>Проекты</h1>
      <BaseButton>Новый проект</BaseButton>
    </div>

    <div v-if="isLoading" class="loading-indicator">Загрузка проектов...</div>

    <div v-else-if="projects.length === 0" class="empty-state">
      <p>Нет доступных проектов</p>
      <BaseButton>Создать проект</BaseButton>
    </div>

    <div v-else class="projects-grid">
      <BaseCard v-for="project in projects" :key="project.id" class="project-card">
        <h3 class="project-title">
          <RouterLink :to="`/projects/${project.id}`">{{ project.name }}</RouterLink>
        </h3>
        <p class="project-description">{{ project.description }}</p>
        
        <div class="project-stats">
          <div class="stat-item">
            <span class="stat-label">Дефектов:</span>
            <span class="stat-value">{{ project.defectsCount }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Готовность:</span>
            <span class="stat-value">{{ project.completionRate }}%</span>
          </div>
        </div>
        
        <div class="progress-bar">
          <div
            class="progress-bar-fill"
            :style="{ width: `${project.completionRate}%` }"
          ></div>
        </div>
        
        <div class="project-status">
          <span
            class="status-badge"
            :class="project.status === 'active' ? 'status-active' : 'status-planning'"
          >
            {{ project.status === 'active' ? 'Активный' : 'Планирование' }}
          </span>
        </div>
        
        <template #footer>
          <div class="card-actions">
            <RouterLink :to="`/projects/${project.id}`" class="btn btn-outline">Подробнее</RouterLink>
          </div>
        </template>
      </BaseCard>
    </div>
  </div>
</template>

<style scoped>
.projects-page {
  padding: 1rem 0;
}

.page-header {
  margin-bottom: 1.5rem;
}

h1 {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--color-primary);
  margin: 0;
}

.projects-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.project-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.project-title {
  margin: 0;
  margin-bottom: 0.5rem;
  font-size: 1.2rem;
}

.project-title a {
  color: var(--color-primary);
  text-decoration: none;
}

.project-description {
  color: var(--color-text-light);
  font-size: 0.9rem;
  margin-bottom: 1rem;
  flex-grow: 1;
}

.project-stats {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.stat-label {
  color: var(--color-text-light);
}

.stat-value {
  font-weight: 500;
  margin-left: 0.25rem;
}

.progress-bar {
  height: 6px;
  background-color: var(--color-border);
  border-radius: 3px;
  margin-bottom: 1rem;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background-color: var(--color-primary);
}

.project-status {
  margin-bottom: 0.5rem;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-active {
  background-color: rgba(76, 175, 80, 0.15);
  color: #2e7d32;
}

.status-planning {
  background-color: rgba(3, 169, 244, 0.15);
  color: #0277bd;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
}

.loading-indicator {
  text-align: center;
  padding: 2rem;
  color: var(--color-text-light);
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
}

.empty-state p {
  margin-bottom: 1rem;
  color: var(--color-text-light);
}
</style>
