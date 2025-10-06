<script setup lang="ts">
import { computed } from 'vue'
import type { Defect } from '../services/defectService'

const props = defineProps<{
  defect: Defect
}>()

// Форматирование даты
const formatDate = (dateString?: string) => {
  if (!dateString) return 'Не указана'
  
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ru-RU').format(date)
}

const getStatusClass = (status: string) => {
  const statusMap: Record<string, string> = {
    new: 'status-new',
    in_progress: 'status-progress',
    review: 'status-review',
    closed: 'status-closed',
    cancelled: 'status-cancelled'
  }
  return statusMap[status] || 'status-new'
}

const getPriorityClass = (priority: string) => {
  const priorityMap: Record<string, string> = {
    low: 'priority-low',
    medium: 'priority-medium',
    high: 'priority-high'
  }
  return priorityMap[priority] || 'priority-medium'
}

const getStatusLabel = (status: string) => {
  const statusLabels: Record<string, string> = {
    new: 'Новый',
    in_progress: 'В работе',
    review: 'На проверке',
    closed: 'Закрыт',
    cancelled: 'Отменен'
  }
  return statusLabels[status] || 'Новый'
}

const getPriorityLabel = (priority: string) => {
  const priorityLabels: Record<string, string> = {
    low: 'Низкий',
    medium: 'Средний',
    high: 'Высокий'
  }
  return priorityLabels[priority] || 'Средний'
}

// Вычисляемое свойство для имени исполнителя
const assigneeName = computed(() => {
  return props.defect.assignee?.full_name || 'Не назначен'
})
</script>

<template>
  <div class="defect-item">
    <div class="defect-header">
      <h3 class="defect-title">{{ defect.title }}</h3>
      <div class="defect-badges">
        <span class="badge status-badge" :class="getStatusClass(defect.status)">
          {{ getStatusLabel(defect.status) }}
        </span>
        <span class="badge priority-badge" :class="getPriorityClass(defect.priority)">
          {{ getPriorityLabel(defect.priority) }}
        </span>
      </div>
    </div>
    <p class="defect-description">{{ defect.description }}</p>
    <div class="defect-details">
      <div class="detail-item">
        <span class="detail-label">Исполнитель:</span>
        <span class="detail-value">{{ assigneeName }}</span>
      </div>
      <div v-if="defect.due_date" class="detail-item">
        <span class="detail-label">Срок:</span>
        <span class="detail-value">{{ formatDate(defect.due_date) }}</span>
      </div>
    </div>
    <div class="defect-actions">
      <slot name="actions"></slot>
    </div>
  </div>
</template>

<style scoped>
.defect-item {
  background-color: var(--color-background);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-default);
  padding: 1rem;
  margin-bottom: 1rem;
  border-left: 4px solid var(--color-primary);
}

.defect-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.5rem;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.defect-title {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 500;
  color: var(--color-text);
}

.defect-badges {
  display: flex;
  gap: 0.5rem;
}

.badge {
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.status-badge {
  background-color: #e0e0e0;
  color: #333;
}

.status-new {
  background-color: #e3f2fd;
  color: #1976d2;
}

.status-progress {
  background-color: #e8f5e9;
  color: #2e7d32;
}

.status-review {
  background-color: #fff3e0;
  color: #e65100;
}

.status-closed {
  background-color: #f5f5f5;
  color: #616161;
}

.status-cancelled {
  background-color: #ffebee;
  color: #c62828;
}

.priority-badge {
  border-radius: 4px;
}

.priority-low {
  background-color: #e8f5e9;
  color: #2e7d32;
}

.priority-medium {
  background-color: #fff3e0;
  color: #e65100;
}

.priority-high {
  background-color: #ffebee;
  color: #c62828;
}

.defect-description {
  margin: 0.5rem 0;
  color: var(--color-text-light);
  font-size: 0.95rem;
}

.defect-details {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin: 0.8rem 0;
  font-size: 0.9rem;
}

.detail-item {
  display: flex;
  align-items: center;
}

.detail-label {
  color: var(--color-text-light);
  margin-right: 0.3rem;
}

.detail-value {
  font-weight: 500;
}

.defect-actions {
  margin-top: 1rem;
  display: flex;
  gap: 0.5rem;
}
</style>