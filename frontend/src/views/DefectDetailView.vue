<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'
import DefectForm from '../components/DefectForm.vue'
import ConfirmModal from '../components/ConfirmModal.vue'
import CommentSection from '../components/CommentSection.vue'
import { useAuth } from '../composables/useAuth'
import { defectService, type Defect } from '../services/defectService'

const route = useRoute()
const router = useRouter()
const defectId = route.params.id as string
const defect = ref<Defect | null>(null)
const isLoading = ref(true)
const error = ref('')
const isEditing = ref(false)
const isDeleting = ref(false)
const showDeleteConfirm = ref(false)

// Получаем права доступа
const permissions = defectService.checkPermissions()

// Форматирование даты
const formatDate = (dateString?: string) => {
  if (!dateString) return 'Не указана'
  
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ru-RU').format(date)
}

// Получение метки статуса
const getStatusLabel = (status: string) => {
  switch (status) {
    case 'new': return 'Новый'
    case 'in_progress': return 'В работе'
    case 'review': return 'На проверке'
    case 'closed': return 'Закрыт'
    case 'cancelled': return 'Отменен'
    default: return status
  }
}

// Получение класса для статуса
const getStatusClass = (status: string) => {
  switch (status) {
    case 'new': return 'status-new'
    case 'in_progress': return 'status-progress'
    case 'review': return 'status-review'
    case 'closed': return 'status-closed'
    case 'cancelled': return 'status-cancelled'
    default: return ''
  }
}

// Получение метки приоритета
const getPriorityLabel = (priority: string) => {
  switch (priority) {
    case 'low': return 'Низкий'
    case 'medium': return 'Средний'
    case 'high': return 'Высокий'
    default: return priority
  }
}

// Получение класса для приоритета
const getPriorityClass = (priority: string) => {
  switch (priority) {
    case 'low': return 'priority-low'
    case 'medium': return 'priority-medium'
    case 'high': return 'priority-high'
    default: return ''
  }
}

// Загрузка данных дефекта
const loadDefect = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const data = await defectService.getDefectById(Number(defectId))
    defect.value = data
  } catch (err) {
    console.error('Ошибка при загрузке дефекта:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить дефект'
  } finally {
    isLoading.value = false
  }
}

// Переключение режима редактирования
const toggleEditMode = () => {
  isEditing.value = !isEditing.value
}

// Обработка сохранения дефекта
const handleSaved = (savedDefect: Defect) => {
  defect.value = savedDefect
  isEditing.value = false
}

// Обработка отмены редактирования
const handleCancelled = () => {
  isEditing.value = false
}

// Открытие диалога подтверждения удаления
const openDeleteConfirm = () => {
  showDeleteConfirm.value = true
}

// Удаление дефекта
const deleteDefect = async () => {
  try {
    isDeleting.value = true
    await defectService.deleteDefect(Number(defectId))
    // После успешного удаления перенаправляем на страницу проекта
    if (defect.value?.project_id) {
      router.push(`/projects/${defect.value.project_id}`)
    } else {
      router.push('/projects')
    }
  } catch (err) {
    console.error('Ошибка при удалении дефекта:', err)
    error.value = err instanceof Error ? err.message : 'Ошибка при удалении дефекта'
  } finally {
    isDeleting.value = false
    showDeleteConfirm.value = false
  }
}

// Возврат к проекту
const backToProject = () => {
  if (defect.value?.project_id) {
    router.push(`/projects/${defect.value.project_id}`)
  } else {
    router.push('/projects')
  }
}

// Загрузка дефекта при монтировании компонента
onMounted(() => {
  loadDefect()
})
</script>

<template>
  <div class="defect-detail-page">
    <div v-if="error" class="alert alert-error">{{ error }}</div>

    <div v-if="isLoading" class="loading-indicator">
      Загрузка информации о дефекте...
    </div>

    <template v-else-if="defect">
      <div class="page-header">
        <div class="navigation-buttons">
          <button class="back-button" @click="backToProject">
            &larr; Вернуться к проекту
          </button>
          <RouterLink to="/defects" class="back-button">
            &larr; К списку дефектов
          </RouterLink>
        </div>
      </div>

      <template v-if="isEditing">
        <h1>Редактирование дефекта</h1>
        <DefectForm 
          :defectId="defectId" 
          :defect="defect"
          @saved="handleSaved" 
          @cancelled="handleCancelled" 
        />
      </template>

      <template v-else>
        <div class="defect-header">
          <div class="defect-title-section">
            <h1>{{ defect.title }}</h1>
            <div class="defect-badges">
              <span class="badge status-badge" :class="getStatusClass(defect.status)">
                {{ getStatusLabel(defect.status) }}
              </span>
              <span class="badge priority-badge" :class="getPriorityClass(defect.priority)">
                {{ getPriorityLabel(defect.priority) }}
              </span>
            </div>
          </div>
          <div class="defect-actions">
            <div v-if="permissions.canEdit || permissions.canDelete" class="button-group">
              <BaseButton v-if="permissions.canEdit" variant="outline" @click="toggleEditMode">
                Редактировать
              </BaseButton>
              <BaseButton v-if="permissions.canDelete" variant="danger" @click="openDeleteConfirm">
                Удалить
              </BaseButton>
            </div>
          </div>
        </div>

        <BaseCard title="Информация о дефекте">
          <div class="defect-details">
            <div class="detail-item">
              <span class="detail-label">Описание:</span>
              <span class="detail-value">{{ defect.description || 'Описание отсутствует' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Проект:</span>
              <span class="detail-value">
                <RouterLink :to="`/projects/${defect.project_id}`">
                  {{ defect.project?.name || `Проект #${defect.project_id}` }}
                </RouterLink>
              </span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Исполнитель:</span>
              <span class="detail-value">{{ defect.assignee?.full_name || 'Не назначен' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Автор:</span>
              <span class="detail-value">{{ defect.reporter?.full_name || 'Неизвестно' }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Срок выполнения:</span>
              <span class="detail-value">{{ formatDate(defect.due_date) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Дата создания:</span>
              <span class="detail-value">{{ formatDate(defect.created_at) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Последнее обновление:</span>
              <span class="detail-value">{{ formatDate(defect.updated_at) }}</span>
            </div>
          </div>
        </BaseCard>

        <!-- Секция комментариев -->
        <div class="comments-section">
          <CommentSection :defect-id="Number(defectId)" />
        </div>
      </template>
    </template>
    
    <div v-else-if="!isLoading && !error" class="not-found">
      <h2>Дефект не найден</h2>
      <p>Дефект с ID {{ defectId }} не существует или был удален.</p>
      <RouterLink to="/projects" class="btn">Вернуться к списку проектов</RouterLink>
    </div>
    
    <!-- Модальное окно подтверждения удаления -->
    <ConfirmModal
      :show="showDeleteConfirm"
      title="Удаление дефекта"
      message="Вы уверены, что хотите удалить этот дефект? Это действие нельзя отменить."
      confirm-text="Удалить"
      cancel-text="Отмена"
      confirm-variant="danger"
      @confirm="deleteDefect"
      @cancel="showDeleteConfirm = false"
    />
  </div>
</template>

<style scoped>
.defect-detail-page {
  padding: 1rem 0;
}

.page-header {
  margin-bottom: 1.5rem;
}

.navigation-buttons {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.back-button {
  background: none;
  border: none;
  color: var(--color-primary);
  cursor: pointer;
  padding: 0.5rem 0;
  font-size: 0.95rem;
}

.back-button:hover {
  text-decoration: underline;
}

.loading-indicator {
  text-align: center;
  padding: 3rem 0;
  color: var(--color-text-light);
}

.defect-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.defect-title-section {
  flex: 1;
}

h1 {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--color-primary);
  margin: 0 0 0.5rem 0;
}

.defect-badges {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
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

.defect-actions {
  display: flex;
  justify-content: flex-end;
}

.button-group {
  display: flex;
  gap: 0.75rem;
}

.defect-details {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  margin-bottom: 1rem;
}

.detail-label {
  color: var(--color-text-light);
  font-size: 0.9rem;
}

.detail-value {
  font-weight: 500;
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