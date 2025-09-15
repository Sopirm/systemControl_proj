<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'

const route = useRoute()
const defectId = route.params.id
const isLoading = ref(true)

// Временные данные для демонстрации
const defect = ref({
  id: Number(defectId),
  title: 'Трещина в несущей стене',
  description: 'На 3 этаже, 2 подъезд обнаружена трещина в несущей стене размером 1.5м. Необходимо провести экспертизу и определить причину. Возможно, потребуется усиление конструкции.',
  priority: 'high',
  status: 'new',
  createdAt: '2023-09-20',
  project: {
    id: 1,
    name: 'ЖК "Солнечный"'
  },
  assignee: 'Иванов И.И.',
  reporter: 'Сидоров С.С.',
  dueDate: '2023-10-15',
  attachments: [
    { id: 1, name: 'Фото 1.jpg', type: 'image' },
    { id: 2, name: 'Экспертное заключение.pdf', type: 'document' }
  ]
})

const comments = ref([
  {
    id: 1,
    author: 'Петров П.П.',
    text: 'Осмотрел дефект. Требуется дополнительная диагностика.',
    createdAt: '2023-09-21 14:32'
  },
  {
    id: 2,
    author: 'Иванов И.И.',
    text: 'Запросил экспертизу у специалистов. Ожидаем результат через 2 дня.',
    createdAt: '2023-09-22 09:15'
  }
])

const newComment = ref('')
const statusOptions = ['new', 'in_progress', 'review', 'closed', 'cancelled']

onMounted(() => {
  // Имитация загрузки данных с сервера
  setTimeout(() => {
    isLoading.value = false
  }, 500)
})

const getStatusLabel = (status: string) => {
  const statusLabels: Record<string, string> = {
    new: 'Новый',
    in_progress: 'В работе',
    review: 'На проверке',
    closed: 'Закрыт',
    cancelled: 'Отменен'
  }
  return statusLabels[status] || 'Неизвестно'
}

const getPriorityLabel = (priority: string) => {
  const priorityLabels: Record<string, string> = {
    low: 'Низкий',
    medium: 'Средний',
    high: 'Высокий'
  }
  return priorityLabels[priority] || 'Неизвестно'
}

const addComment = () => {
  if (!newComment.value.trim()) return
  
  comments.value.push({
    id: comments.value.length + 1,
    author: 'Текущий пользователь',
    text: newComment.value.trim(),
    createdAt: new Date().toLocaleString()
  })
  
  newComment.value = ''
}

const updateStatus = (status: string) => {
  defect.value.status = status
}
</script>

<template>
  <div class="defect-detail-page">
    <div v-if="isLoading" class="loading-indicator">
      Загрузка информации о дефекте...
    </div>

    <template v-else>
      <div class="defect-header">
        <RouterLink :to="`/projects/${defect.project.id}`" class="project-link">
          {{ defect.project.name }}
        </RouterLink>
        <h1>{{ defect.title }}</h1>
        <div class="defect-badges">
          <span class="badge status-badge" :class="`status-${defect.status}`">
            {{ getStatusLabel(defect.status) }}
          </span>
          <span class="badge priority-badge" :class="`priority-${defect.priority}`">
            {{ getPriorityLabel(defect.priority) }}
          </span>
        </div>
      </div>

      <div class="defect-content">
        <div class="defect-info-section">
          <BaseCard title="Информация о дефекте">
            <p class="defect-description">{{ defect.description }}</p>
            
            <div class="defect-details">
              <div class="detail-item">
                <span class="detail-label">Создан:</span>
                <span class="detail-value">{{ defect.createdAt }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Создатель:</span>
                <span class="detail-value">{{ defect.reporter }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Исполнитель:</span>
                <span class="detail-value">{{ defect.assignee }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Срок:</span>
                <span class="detail-value">{{ defect.dueDate }}</span>
              </div>
            </div>
            
            <div class="attachments-section" v-if="defect.attachments.length">
              <h3>Вложения:</h3>
              <ul class="attachments-list">
                <li v-for="attachment in defect.attachments" :key="attachment.id" class="attachment-item">
                  <a href="#" class="attachment-link">
                    <span class="attachment-icon" :class="attachment.type">
                      {{ attachment.type === 'image' ? '🖼️' : '📄' }}
                    </span>
                    <span class="attachment-name">{{ attachment.name }}</span>
                  </a>
                </li>
              </ul>
            </div>
          </BaseCard>

          <div class="status-update-section">
            <BaseCard title="Обновить статус">
              <div class="status-buttons">
                <button
                  v-for="status in statusOptions"
                  :key="status"
                  class="status-button"
                  :class="[`status-${status}`, { active: defect.status === status }]"
                  @click="updateStatus(status)"
                >
                  {{ getStatusLabel(status) }}
                </button>
              </div>
            </BaseCard>
          </div>

          <div class="comments-section">
            <BaseCard title="Обсуждение">
              <div v-if="comments.length === 0" class="no-comments">
                Нет комментариев
              </div>
              
              <div v-else class="comments-list">
                <div v-for="comment in comments" :key="comment.id" class="comment-item">
                  <div class="comment-header">
                    <span class="comment-author">{{ comment.author }}</span>
                    <span class="comment-date">{{ comment.createdAt }}</span>
                  </div>
                  <p class="comment-text">{{ comment.text }}</p>
                </div>
              </div>
              
              <div class="add-comment">
                <h4>Добавить комментарий:</h4>
                <textarea
                  v-model="newComment"
                  rows="3"
                  class="comment-textarea"
                  placeholder="Введите комментарий..."
                ></textarea>
                <BaseButton @click="addComment">Отправить</BaseButton>
              </div>
            </BaseCard>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.defect-detail-page {
  padding: 1rem 0;
}

.loading-indicator {
  text-align: center;
  padding: 3rem 0;
  color: var(--color-text-light);
}

.defect-header {
  margin-bottom: 2rem;
}

.project-link {
  display: block;
  color: var(--color-primary);
  font-size: 0.9rem;
  margin-bottom: 0.5rem;
}

h1 {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--color-text);
  margin: 0 0 0.5rem;
}

.defect-badges {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
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

.status-in_progress {
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

.defect-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

.defect-description {
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.defect-details {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
}

.detail-label {
  font-size: 0.85rem;
  color: var(--color-text-light);
}

.detail-value {
  font-weight: 500;
}

.attachments-section {
  margin-top: 1.5rem;
}

h3 {
  font-size: 1.1rem;
  margin-bottom: 0.5rem;
  color: var(--color-text);
}

.attachments-list {
  list-style: none;
  padding: 0;
}

.attachment-item {
  margin-bottom: 0.5rem;
}

.attachment-link {
  display: flex;
  align-items: center;
  padding: 0.5rem;
  background-color: #f5f5f5;
  border-radius: var(--border-radius);
  transition: var(--transition-default);
}

.attachment-link:hover {
  background-color: #e0e0e0;
}

.attachment-icon {
  margin-right: 0.5rem;
}

.status-update-section {
  margin: 1.5rem 0;
}

.status-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.status-button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: var(--border-radius);
  cursor: pointer;
  font-size: 0.85rem;
  transition: var(--transition-default);
  opacity: 0.7;
}

.status-button:hover {
  opacity: 1;
}

.status-button.active {
  opacity: 1;
  font-weight: 500;
}

.comments-section {
  margin-top: 1.5rem;
}

.no-comments {
  text-align: center;
  padding: 1rem;
  color: var(--color-text-light);
}

.comments-list {
  margin-bottom: 2rem;
}

.comment-item {
  padding: 1rem;
  border-bottom: 1px solid var(--color-border);
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.comment-author {
  font-weight: 500;
}

.comment-date {
  font-size: 0.85rem;
  color: var(--color-text-light);
}

.comment-text {
  line-height: 1.5;
}

.add-comment {
  margin-top: 1.5rem;
}

h4 {
  font-size: 1rem;
  margin-bottom: 0.5rem;
  color: var(--color-text);
}

.comment-textarea {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  resize: vertical;
  margin-bottom: 0.8rem;
}

.comment-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
}

@media (max-width: 768px) {
  .defect-content {
    grid-template-columns: 1fr;
  }
  
  .defect-details {
    grid-template-columns: 1fr;
  }
  
  .status-buttons {
    flex-direction: column;
  }
}
</style>
