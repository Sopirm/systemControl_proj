<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'
import { useAuth } from '../composables/useAuth'
import { projectService, type Project } from '../services/projectService'
import { defectService, type DefectStats } from '../services/defectService'

const router = useRouter()
const { isManager } = useAuth()
const projects = ref<Project[]>([])
const isLoading = ref(true)
const error = ref('')
const projectDefectStats = ref<Record<number, DefectStats>>({})

// Загрузка статистики дефектов для проекта
const loadDefectStats = async (projectId: number) => {
  try {
    const stats = await defectService.getDefectStatsByProjectId(projectId)
    projectDefectStats.value[projectId] = stats
    console.log(`Статистика дефектов для проекта ${projectId}:`, stats)
  } catch (err) {
    console.error(`Ошибка при загрузке статистики дефектов для проекта ${projectId}:`, err)
    // В случае ошибки устанавливаем нулевые значения
    projectDefectStats.value[projectId] = { active: 0, resolved: 0, total: 0 }
  }
}

// Получение статистики дефектов для проекта
const getDefectStats = (project: Project): DefectStats => {
  return projectDefectStats.value[project.id] || { active: 0, resolved: 0, total: 0 }
}

// Загрузка списка проектов
const loadProjects = async () => {
  try {
    isLoading.value = true
    error.value = ''
    const data = await projectService.getAllProjects()
    projects.value = data
    
    // Загружаем статистику дефектов для каждого проекта
    for (const project of data) {
      await loadDefectStats(project.id)
    }
  } catch (err) {
    console.error('Ошибка при загрузке проектов:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить проекты'
  } finally {
    isLoading.value = false
  }
}

// Создание нового проекта
const createProject = () => {
  router.push('/projects/create')
}

// Отображаемые статусы
const getStatusLabel = (status: string) => {
  switch (status) {
    case 'active':
      return 'Активный'
    case 'completed':
      return 'Завершен'
    case 'suspended':
      return 'Приостановлен'
    default:
      return 'Планирование'
  }
}

// Классы для статусов
const getStatusClass = (status: string) => {
  switch (status) {
    case 'active':
      return 'status-active'
    case 'completed':
      return 'status-completed'
    case 'suspended':
      return 'status-suspended'
    default:
      return 'status-planning'
  }
}

// Загрузка проектов при монтировании компонента
onMounted(() => {
  loadProjects()
})
</script>

<template>
  <div class="projects-page">
    <div class="page-header flex-between">
      <h1>Проекты</h1>
      <BaseButton v-if="isManager" @click="createProject">Новый проект</BaseButton>
    </div>

    <div v-if="error" class="alert alert-error">{{ error }}</div>

    <div v-if="isLoading" class="loading-indicator">Загрузка проектов...</div>

    <div v-else-if="projects.length === 0" class="empty-state">
      <p>Нет доступных проектов</p>
      <BaseButton v-if="isManager" @click="createProject">Создать проект</BaseButton>
    </div>

    <div v-else class="projects-grid">
      <BaseCard v-for="project in projects" :key="project.id" class="project-card">
        <h3 class="project-title">
          <RouterLink :to="`/projects/${project.id}`">{{ project.name }}</RouterLink>
        </h3>
        <p class="project-description">{{ project.description }}</p>
        
        <div class="project-stats">
          <div class="stat-item">
            <span class="stat-label">Активные дефекты:</span>
            <span class="stat-value">{{ getDefectStats(project).active }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">Решенные дефекты:</span>
            <span class="stat-value">{{ getDefectStats(project).resolved }}</span>
          </div>
        </div>
        
        <div class="defect-progress">
          <div class="defect-count">
            <span class="count-label">Всего дефектов:</span>
            <span class="count-value">{{ getDefectStats(project).total }}</span>
          </div>
          <div v-if="getDefectStats(project).total > 0" class="defect-ratio">
            <div 
              class="ratio-fill ratio-active" 
              :style="{ width: `${getDefectStats(project).active / getDefectStats(project).total * 100}%` }"
            ></div>
            <div 
              class="ratio-fill ratio-resolved" 
              :style="{ width: `${getDefectStats(project).resolved / getDefectStats(project).total * 100}%` }"
            ></div>
          </div>
          <div v-else class="defect-ratio empty-ratio">
            <span class="empty-ratio-text">Нет дефектов</span>
          </div>
        </div>
        
        <div class="project-info">
          <div class="info-item">
            <span class="info-label">Менеджер:</span>
            <span class="info-value">{{ project.manager?.full_name || 'Не назначен' }}</span>
          </div>
          <div class="project-status">
            <span
              class="status-badge"
              :class="getStatusClass(project.status)"
            >
              {{ getStatusLabel(project.status) }}
            </span>
          </div>
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

.defect-progress {
  margin-bottom: 1rem;
}

.defect-count {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.25rem;
  font-size: 0.9rem;
}

.count-label {
  color: var(--color-text-light);
}

.count-value {
  font-weight: 500;
}

.defect-ratio {
  height: 6px;
  background-color: var(--color-border);
  border-radius: 3px;
  overflow: hidden;
  display: flex;
}

.ratio-fill {
  height: 100%;
}

.ratio-active {
  background-color: #ff9800;
}

.ratio-resolved {
  background-color: #4caf50;
}

.empty-ratio {
  height: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  font-size: 0.7rem;
  color: var(--color-text-light);
}

.empty-ratio-text {
  display: none;
}

.project-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.info-label {
  color: var(--color-text-light);
}

.info-value {
  font-weight: 500;
  margin-left: 0.25rem;
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

.status-completed {
  background-color: rgba(96, 125, 139, 0.15);
  color: #455a64;
}

.status-suspended {
  background-color: rgba(255, 152, 0, 0.15);
  color: #ef6c00;
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

.alert {
  padding: 1rem;
  margin-bottom: 1.5rem;
  border-radius: var(--border-radius);
}

.alert-error {
  background-color: rgba(244, 67, 54, 0.1);
  color: #d32f2f;
  border: 1px solid rgba(244, 67, 54, 0.3);
}
</style>