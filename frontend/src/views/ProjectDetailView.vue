<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'
import DefectItem from '../components/DefectItem.vue'

const route = useRoute()
const projectId = route.params.id
const isLoading = ref(true)

// Временные данные для демонстрации
const project = ref({
  id: Number(projectId),
  name: 'ЖК "Солнечный"',
  description: 'Жилой комплекс из 3 корпусов на ул. Ленина',
  status: 'active',
  completionRate: 65,
  startDate: '2022-03-15',
  endDate: '2023-12-31',
  location: 'г. Москва, ул. Ленина, д. 12',
  manager: 'Смирнов А.В.'
})

const defects = ref([
  {
    id: 1,
    title: 'Трещина в несущей стене',
    description: 'На 3 этаже, 2 подъезд обнаружена трещина в несущей стене размером 1.5м',
    priority: 'high',
    status: 'new',
    assignee: 'Иванов И.И.',
    dueDate: '2023-10-15'
  },
  {
    id: 2,
    title: 'Протечка в потолке',
    description: 'Кв. 42, 5 этаж - следы протечки на потолке в ванной комнате',
    priority: 'medium',
    status: 'in_progress',
    assignee: 'Петров П.П.'
  }
])

onMounted(() => {
  // Имитация загрузки данных с сервера
  setTimeout(() => {
    isLoading.value = false
  }, 500)
})
</script>

<template>
  <div class="project-detail-page">
    <div v-if="isLoading" class="loading-indicator">
      Загрузка информации о проекте...
    </div>

    <template v-else>
      <div class="project-header">
        <div class="project-title-section">
          <h1>{{ project.name }}</h1>
          <span 
            class="status-badge"
            :class="project.status === 'active' ? 'status-active' : 'status-planning'"
          >
            {{ project.status === 'active' ? 'Активный' : 'Планирование' }}
          </span>
        </div>
        <div class="project-actions">
          <BaseButton variant="outline">Редактировать</BaseButton>
        </div>
      </div>

      <div class="project-info-section">
        <BaseCard title="Информация о проекте">
          <div class="project-details">
            <div class="detail-item">
              <span class="detail-label">Описание:</span>
              <span class="detail-value">{{ project.description }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Дата начала:</span>
              <span class="detail-value">{{ project.startDate }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Дата окончания:</span>
              <span class="detail-value">{{ project.endDate }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Местоположение:</span>
              <span class="detail-value">{{ project.location }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Руководитель:</span>
              <span class="detail-value">{{ project.manager }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Степень готовности:</span>
              <span class="detail-value">{{ project.completionRate }}%</span>
            </div>
          </div>
          
          <div class="progress-section">
            <div class="progress-label">Прогресс проекта:</div>
            <div class="progress-bar">
              <div
                class="progress-bar-fill"
                :style="{ width: `${project.completionRate}%` }"
              ></div>
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

.progress-section {
  margin-top: 1rem;
}

.progress-label {
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.progress-bar {
  height: 8px;
  background-color: var(--color-border);
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar-fill {
  height: 100%;
  background-color: var(--color-primary);
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
</style>
