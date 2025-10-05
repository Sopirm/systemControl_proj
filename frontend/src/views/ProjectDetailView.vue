<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'
import DefectItem from '../components/DefectItem.vue'
import { useAuth } from '../composables/useAuth'
import { projectService, type Project } from '../services/projectService'
import { defectService, type Defect, type DefectStats } from '../services/defectService'

const route = useRoute()
const router = useRouter()
const { isManager } = useAuth()
const projectId = route.params.id as string
const project = ref<Project | null>(null)
const isLoading = ref(true)
const error = ref('')
const defects = ref<Defect[]>([])
const defectStats = ref<DefectStats>({ active: 0, resolved: 0, total: 0 })
const isLoadingDefects = ref(false)

// Функция для форматирования даты
const formatDate = (dateString: string) => {
  if (!dateString) return 'Не указана'
  
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ru-RU').format(date)
}

// Вспомогательные функции для статусов проекта
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

// Загрузка статистики дефектов
const loadDefectStats = async () => {
  try {
    const stats = await defectService.getDefectStatsByProjectId(Number(projectId))
    defectStats.value = stats
  } catch (err) {
    console.error(`Ошибка при загрузке статистики дефектов для проекта ${projectId}:`, err)
    // Используем демо-данные в случае ошибки
    defectStats.value = defectService.generateDemoDefectStats(Number(projectId))
  }
}

// Загрузка дефектов проекта
const loadDefects = async () => {
  try {
    isLoadingDefects.value = true
    // Попробуем загрузить реальные дефекты через API
    const projectDefects = await defectService.getDefectsByProjectId(Number(projectId))
    defects.value = projectDefects
  } catch (err) {
    console.error(`Ошибка при загрузке дефектов для проекта ${projectId}:`, err)
    // В случае ошибки используем демо-данные
    defects.value = [
      {
        id: 1,
        title: 'Трещина в несущей стене',
        description: 'На 3 этаже, 2 подъезд обнаружена трещина в несущей стене размером 1.5м',
        priority: 'high',
        status: 'new',
        project_id: Number(projectId),
        assignee: { id: 1, username: 'ivanov', full_name: 'Иванов И.И.', email: 'ivanov@example.com' },
        due_date: '2023-10-15'
      },
      {
        id: 2,
        title: 'Протечка в потолке',
        description: 'Кв. 42, 5 этаж - следы протечки на потолке в ванной комнате',
        priority: 'medium',
        status: 'in_progress',
        project_id: Number(projectId),
        assignee: { id: 2, username: 'petrov', full_name: 'Петров П.П.', email: 'petrov@example.com' }
      }
    ] as Defect[]
  } finally {
    isLoadingDefects.value = false
  }
}

// Загрузка данных проекта
const loadProject = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const data = await projectService.getProjectById(Number(projectId))
    project.value = data
    
    // Загружаем статистику дефектов и дефекты
    await Promise.all([
      loadDefectStats(),
      loadDefects()
    ])
  } catch (err) {
    console.error('Ошибка при загрузке проекта:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить проект'
  } finally {
    isLoading.value = false
  }
}

// Редактирование проекта
const editProject = () => {
  router.push(`/projects/${projectId}/edit`)
}

// Загрузка проекта при монтировании компонента
onMounted(() => {
  loadProject()
})
</script>

<template>
  <div class="project-detail-page">
    <div v-if="error" class="alert alert-error">{{ error }}</div>

    <div v-if="isLoading" class="loading-indicator">
      Загрузка информации о проекте...
    </div>

    <template v-else-if="project">
      <div class="project-header">
        <div class="project-title-section">
          <h1>{{ project.name }}</h1>
          <span 
            class="status-badge"
            :class="getStatusClass(project.status)"
          >
            {{ getStatusLabel(project.status) }}
          </span>
        </div>
        <div class="project-actions">
          <BaseButton v-if="isManager" variant="outline" @click="editProject">Редактировать</BaseButton>
        </div>
      </div>

      <div class="project-info-section">
        <BaseCard title="Информация о проекте">
          <div class="project-details">
            <div class="detail-item">
              <span class="detail-label">Описание:</span>
              <span class="detail-value">{{ project.description || 'Описание отсутствует' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Дата начала:</span>
              <span class="detail-value">{{ formatDate(project.start_date) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Дата окончания:</span>
              <span class="detail-value">{{ formatDate(project.end_date) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Местоположение:</span>
              <span class="detail-value">{{ project.location || 'Не указано' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Руководитель:</span>
              <span class="detail-value">{{ project.manager?.full_name || 'Не назначен' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Активные дефекты:</span>
              <span class="detail-value">{{ defectStats.active }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Решенные дефекты:</span>
              <span class="detail-value">{{ defectStats.resolved }}</span>
            </div>
          </div>
          
          <div class="defect-stats-section">
            <div class="defect-stats-header">
              <div class="defect-stats-title">Статистика дефектов:</div>
              <div class="defect-stats-total">Всего дефектов: {{ defectStats.total }}</div>
            </div>
            <div v-if="defectStats.total > 0" class="defect-ratio">
              <div 
                class="ratio-fill ratio-active" 
                :style="{ width: `${defectStats.active / defectStats.total * 100}%` }"
                :title="`Активные дефекты: ${defectStats.active}`"
              ></div>
              <div 
                class="ratio-fill ratio-resolved" 
                :style="{ width: `${defectStats.resolved / defectStats.total * 100}%` }"
                :title="`Решенные дефекты: ${defectStats.resolved}`"
              ></div>
            </div>
            <div v-else class="defect-ratio empty-ratio">
              <span class="empty-ratio-text">Нет дефектов</span>
            </div>
            <div class="defect-stats-legend">
              <div class="legend-item">
                <div class="legend-color active-color"></div>
                <div class="legend-label">Активные</div>
              </div>
              <div class="legend-item">
                <div class="legend-color resolved-color"></div>
                <div class="legend-label">Решенные</div>
              </div>
            </div>
          </div>
        </BaseCard>
      </div>

      <div class="project-defects-section">
        <div class="section-header flex-between">
          <h2>Дефекты проекта</h2>
          <BaseButton>Добавить дефект</BaseButton>
        </div>

        <div v-if="defects.length === 0" class="empty-state">
          <p>Для этого проекта не зарегистрировано дефектов</p>
        </div>

        <div v-else class="defects-list">
          <DefectItem 
            v-for="defect in defects" 
            :key="defect.id" 
            :defect="defect"
          >
            <template #actions>
              <BaseButton variant="outline">Подробнее</BaseButton>
            </template>
          </DefectItem>
        </div>
      </div>
    </template>
    <div v-else-if="!isLoading && !error" class="not-found">
      <h2>Проект не найден</h2>
      <p>Проект с ID {{ projectId }} не существует или был удален.</p>
      <RouterLink to="/projects" class="btn">Вернуться к списку проектов</RouterLink>
    </div>
  </div>
</template>

<style scoped>
.project-detail-page {
  padding: 1rem 0;
}

.loading-indicator {
  text-align: center;
  padding: 3rem 0;
  color: var(--color-text-light);
}

.project-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.project-title-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

h1 {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--color-primary);
  margin: 0;
}

h2 {
  font-size: 1.5rem;
  color: var(--color-primary);
  margin: 0;
}

.status-badge {
  display: inline-block;
  padding: 0.3rem 0.6rem;
  border-radius: 12px;
  font-size: 0.8rem;
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

.project-info-section {
  margin-bottom: 2rem;
}

.project-details {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.detail-label {
  color: var(--color-text-light);
  font-size: 0.9rem;
}

.detail-value {
  font-weight: 500;
}

.defect-stats-section {
  margin-top: 1.5rem;
}

.defect-stats-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.defect-stats-title {
  font-weight: 500;
}

.defect-stats-total {
  font-weight: 500;
}

.defect-ratio {
  height: 8px;
  background-color: var(--color-border);
  border-radius: 4px;
  overflow: hidden;
  display: flex;
  margin-bottom: 0.75rem;
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
  height: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  font-size: 0.7rem;
  color: var(--color-text-light);
  margin-bottom: 0.75rem;
}

.empty-ratio-text {
  display: none;
}

.defect-stats-legend {
  display: flex;
  gap: 1.5rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}

.active-color {
  background-color: #ff9800;
}

.resolved-color {
  background-color: #4caf50;
}

.project-defects-section {
  margin-top: 2rem;
}

.section-header {
  margin-bottom: 1.5rem;
}

.defects-list {
  margin-top: 1.5rem;
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  background-color: #f8f9fa;
  border-radius: var(--border-radius);
}

.empty-state p {
  color: var(--color-text-light);
}

.not-found {
  text-align: center;
  padding: 3rem 1rem;
}

.not-found h2 {
  margin-bottom: 1rem;
}

.not-found p {
  margin-bottom: 2rem;
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