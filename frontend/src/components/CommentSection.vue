<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import BaseCard from './BaseCard.vue'
import BaseButton from './BaseButton.vue'
import ConfirmModal from './ConfirmModal.vue'
import { commentService, type Comment } from '../services/commentService'
import { authService } from '../services/api'

const props = defineProps<{
  defectId: number
}>()

const comments = ref<Comment[]>([])
const newCommentContent = ref('')
const isLoading = ref(true)
const isSubmitting = ref(false)
const error = ref('')
const showDeleteConfirm = ref(false)
const commentToDelete = ref<number | null>(null)

const currentUser = computed(() => authService.getCurrentUser())

// Форматирование даты и времени
const formatDateTime = (dateString: string) => {
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

// Загрузка комментариев
const loadComments = async () => {
  try {
    isLoading.value = true
    error.value = ''
    comments.value = await commentService.getDefectComments(props.defectId)
  } catch (err) {
    console.error('Ошибка при загрузке комментариев:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить комментарии'
  } finally {
    isLoading.value = false
  }
}

// Создание нового комментария
const submitComment = async () => {
  if (!newCommentContent.value.trim()) {
    error.value = 'Комментарий не может быть пустым'
    return
  }

  try {
    isSubmitting.value = true
    error.value = ''
    
    const newComment = await commentService.createComment({
      defect_id: props.defectId,
      content: newCommentContent.value.trim()
    })
    
    comments.value.push(newComment)
    newCommentContent.value = ''
  } catch (err) {
    console.error('Ошибка при создании комментария:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось создать комментарий'
  } finally {
    isSubmitting.value = false
  }
}

// Открытие диалога подтверждения удаления
const openDeleteConfirm = (commentId: number) => {
  commentToDelete.value = commentId
  showDeleteConfirm.value = true
}

// Удаление комментария
const deleteComment = async () => {
  if (!commentToDelete.value) return

  try {
    await commentService.deleteComment(commentToDelete.value)
    comments.value = comments.value.filter(c => c.id !== commentToDelete.value)
    commentToDelete.value = null
  } catch (err) {
    console.error('Ошибка при удалении комментария:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось удалить комментарий'
  } finally {
    showDeleteConfirm.value = false
  }
}

// Проверка, может ли пользователь удалить комментарий
const canDeleteComment = (comment: Comment) => {
  return commentService.canDeleteComment(comment)
}

// Загрузка комментариев при монтировании компонента
onMounted(() => {
  loadComments()
})
</script>

<template>
  <BaseCard title="Комментарии">
    <div v-if="error" class="alert alert-error">{{ error }}</div>

    <!-- Форма добавления комментария -->
    <div class="comment-form">
      <textarea
        v-model="newCommentContent"
        class="comment-textarea"
        placeholder="Добавить комментарий..."
        rows="3"
        :disabled="isSubmitting"
      ></textarea>
      <div class="comment-form-actions">
        <BaseButton 
          @click="submitComment" 
          :disabled="isSubmitting || !newCommentContent.trim()"
        >
          {{ isSubmitting ? 'Отправка...' : 'Добавить комментарий' }}
        </BaseButton>
      </div>
    </div>

    <!-- Список комментариев -->
    <div v-if="isLoading" class="loading-indicator">
      Загрузка комментариев...
    </div>

    <div v-else-if="comments.length === 0" class="no-comments">
      Комментариев пока нет. Будьте первым!
    </div>

    <div v-else class="comments-list">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <div class="comment-header">
          <div class="comment-author">
            <span class="author-name">{{ comment.user?.full_name || 'Неизвестный пользователь' }}</span>
            <span class="author-role" v-if="comment.user?.role === 'manager'">Менеджер</span>
            <span class="author-role" v-else-if="comment.user?.role === 'engineer'">Инженер</span>
          </div>
          <div class="comment-meta">
            <span class="comment-date">{{ formatDateTime(comment.created_at) }}</span>
            <button 
              v-if="canDeleteComment(comment)"
              @click="openDeleteConfirm(comment.id)"
              class="delete-btn"
              title="Удалить комментарий"
            >
              ✕
            </button>
          </div>
        </div>
        <div class="comment-content">
          {{ comment.content }}
        </div>
      </div>
    </div>

    <!-- Модальное окно подтверждения удаления -->
    <ConfirmModal
      :show="showDeleteConfirm"
      title="Удаление комментария"
      message="Вы уверены, что хотите удалить этот комментарий? Это действие нельзя отменить."
      confirm-text="Удалить"
      cancel-text="Отмена"
      confirm-variant="danger"
      @confirm="deleteComment"
      @cancel="showDeleteConfirm = false"
    />
  </BaseCard>
</template>

<style scoped>
.alert {
  padding: 1rem;
  margin-bottom: 1rem;
  border-radius: var(--border-radius);
}

.alert-error {
  background-color: rgba(244, 67, 54, 0.1);
  color: #d32f2f;
  border: 1px solid rgba(244, 67, 54, 0.3);
}

.comment-form {
  margin-bottom: 2rem;
}

.comment-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  font-family: inherit;
  font-size: 0.95rem;
  resize: vertical;
  transition: border-color 0.2s;
}

.comment-textarea:focus {
  outline: none;
  border-color: var(--color-primary);
}

.comment-textarea:disabled {
  background-color: #f5f5f5;
  cursor: not-allowed;
}

.comment-form-actions {
  margin-top: 0.75rem;
  display: flex;
  justify-content: flex-end;
}

.loading-indicator {
  text-align: center;
  padding: 2rem 0;
  color: var(--color-text-light);
}

.no-comments {
  text-align: center;
  padding: 2rem 0;
  color: var(--color-text-light);
  font-style: italic;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.comment-item {
  padding: 1rem;
  background-color: rgba(0, 0, 0, 0.02);
  border-radius: var(--border-radius);
  border: 1px solid var(--color-border);
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.comment-author {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.author-name {
  font-weight: 600;
  color: var(--color-primary);
}

.author-role {
  padding: 0.125rem 0.5rem;
  background-color: var(--color-primary);
  color: white;
  font-size: 0.7rem;
  border-radius: 12px;
  font-weight: 500;
}

.comment-meta {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.comment-date {
  font-size: 0.85rem;
  color: var(--color-text-light);
}

.delete-btn {
  background: none;
  border: none;
  color: #f44336;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.delete-btn:hover {
  background-color: rgba(244, 67, 54, 0.1);
}

.comment-content {
  color: var(--color-text);
  line-height: 1.6;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
