<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import BaseInput from './BaseInput.vue'
import BaseButton from './BaseButton.vue'
import { defectService, type Defect, type DefectCreate, type DefectUpdate } from '../services/defectService'
import { projectService, type Project } from '../services/projectService'
import { userService, type User } from '../services/userService'
import { useAuth } from '../composables/useAuth'

const props = defineProps({
  defectId: {
    type: [Number, String],
    default: null
  },
  projectId: {
    type: [Number, String, null],
    default: null
  },
  defect: {
    type: Object as () => Defect | null,
    default: null
  }
})

const emit = defineEmits(['saved', 'cancelled'])

// Состояние формы
const title = ref('')
const description = ref('')
const priority = ref('medium')
const status = ref('new')
const assigneeId = ref<number | null>(null)
const dueDate = ref('')
const projectId = ref<number | null>(props.projectId ? Number(props.projectId) : null)
const isLoading = ref(false)
const isSaving = ref(false)
const error = ref('')
const engineers = ref<User[]>([])
const projects = ref<Project[]>([])

// Получаем текущего пользователя и его права
const { currentUser } = useAuth()
const permissions = defectService.checkPermissions()

// Опции для выбора приоритета
const priorityOptions = [
  { value: 'low', label: 'Низкий' },
  { value: 'medium', label: 'Средний' },
  { value: 'high', label: 'Высокий' }
]

// Опции для выбора статуса
const statusOptions = [
  { value: 'new', label: 'Новый' },
  { value: 'in_progress', label: 'В работе' },
  { value: 'review', label: 'На проверке' },
  { value: 'closed', label: 'Закрыт' },
  { value: 'cancelled', label: 'Отменен' }
]

// Функция для загрузки инженеров
const loadEngineers = async () => {
  try {
    const users = await userService.getAllUsers()
    engineers.value = users.filter(user => 
      user.role === 'engineer' || user.role === 'manager'
    )
  } catch (err) {
    console.error('Ошибка при загрузке списка инженеров:', err)
    error.value = 'Не удалось загрузить список исполнителей'
  }
}

// Функция для загрузки проектов
const loadProjects = async () => {
  try {
    const data = await projectService.getAllProjects()
    // Отображаем только активные проекты
    projects.value = data.filter(project => project.status === 'active')
  } catch (err) {
    console.error('Ошибка при загрузке списка проектов:', err)
    error.value = 'Не удалось загрузить список проектов'
  }
}

// Функция для загрузки дефекта при редактировании
const loadDefect = async () => {
  if (!props.defectId) return
  
  try {
    isLoading.value = true
    error.value = ''
    
    const defectData = await defectService.getDefectById(Number(props.defectId))
    title.value = defectData.title
    description.value = defectData.description || ''
    priority.value = defectData.priority
    status.value = defectData.status
    assigneeId.value = defectData.assignee_id || null
    dueDate.value = defectData.due_date ? new Date(defectData.due_date).toISOString().split('T')[0] : ''
  } catch (err) {
    console.error('Ошибка при загрузке дефекта:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить данные дефекта'
  } finally {
    isLoading.value = false
  }
}

// Наблюдаем за изменением props.defect
watch(() => props.defect, (newDefect) => {
  if (newDefect) {
    title.value = newDefect.title
    description.value = newDefect.description || ''
    priority.value = newDefect.priority
    status.value = newDefect.status
    assigneeId.value = newDefect.assignee_id || null
    dueDate.value = newDefect.due_date ? new Date(newDefect.due_date).toISOString().split('T')[0] : ''
  }
}, { immediate: true })

// Валидация формы
const validateForm = () => {
  if (!title.value) {
    error.value = 'Название дефекта обязательно для заполнения'
    return false
  }
  if (!projectId.value) {
    error.value = 'Необходимо выбрать проект'
    return false
  }
  return true
}

// Проверка, может ли пользователь изменить статус на выбранный
const canChangeToStatus = (statusValue: string) => {
  return defectService.canChangeStatusTo(statusValue)
}

// Сохранение дефекта
const saveDefect = async () => {
  error.value = ''
  
  if (!validateForm()) {
    return
  }
  
  try {
    isSaving.value = true
    
    // Проверка прав на изменение статуса
    if (props.defectId && !canChangeToStatus(status.value)) {
      error.value = 'У вас нет прав для изменения статуса на выбранный'
      return
    }
    
    let result
    if (props.defectId) {
      // Обновление дефекта
      const defectData: DefectUpdate = {
        title: title.value,
        description: description.value,
        priority: priority.value as 'low' | 'medium' | 'high',
        status: status.value as 'new' | 'in_progress' | 'review' | 'closed' | 'cancelled',
        assignee_id: assigneeId.value || undefined,
        due_date: dueDate.value || undefined
      }
      
      result = await defectService.updateDefect(Number(props.defectId), defectData)
    } else {
      // Создание нового дефекта
      const defectData: DefectCreate = {
        title: title.value,
        description: description.value,
        project_id: projectId.value!,
        priority: priority.value as 'low' | 'medium' | 'high',
        assignee_id: assigneeId.value || undefined,
        due_date: dueDate.value || undefined
      }
      
      result = await defectService.createDefect(defectData)
    }
    
    emit('saved', result)
  } catch (err) {
    console.error('Ошибка при сохранении дефекта:', err)
    error.value = err instanceof Error ? err.message : 'Ошибка при сохранении дефекта'
  } finally {
    isSaving.value = false
  }
}

// Отмена
const cancel = () => {
  emit('cancelled')
}

// Инициализация
onMounted(() => {
  loadEngineers()
  loadProjects()
  if (props.defectId) {
    loadDefect()
  }
})
</script>

<template>
  <div class="defect-form">
    <div v-if="isLoading" class="loading">
      Загрузка данных дефекта...
    </div>
    <form v-else @submit.prevent="saveDefect">
      <div v-if="error" class="form-error">{{ error }}</div>

      <div class="form-group">
        <BaseInput
          v-model="title"
          label="Название дефекта"
          placeholder="Введите название дефекта"
          required
        />
      </div>

      <div class="form-group">
        <BaseInput
          v-model="description"
          label="Описание"
          placeholder="Введите описание дефекта"
          type="textarea"
        />
      </div>

      <div v-if="!props.projectId" class="form-group">
        <label for="project" class="form-label">Проект</label>
        <select
          id="project"
          v-model="projectId"
          class="form-select"
          required
        >
          <option :value="null">Выберите проект</option>
          <option v-for="project in projects" :key="project.id" :value="project.id">
            {{ project.name }}
          </option>
        </select>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="priority" class="form-label">Приоритет</label>
          <select
            id="priority"
            v-model="priority"
            class="form-select"
            required
          >
            <option v-for="option in priorityOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
        
        <div v-if="props.defectId" class="form-group">
          <label for="status" class="form-label">Статус</label>
          <select
            id="status"
            v-model="status"
            class="form-select"
            required
          >
            <option 
              v-for="option in statusOptions" 
              :key="option.value" 
              :value="option.value"
              :disabled="!canChangeToStatus(option.value)"
            >
              {{ option.label }}
            </option>
          </select>
          <small v-if="!permissions.canCloseOrCancel" class="form-text">
            Только менеджер может закрыть или отменить дефект
          </small>
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="assignee" class="form-label">Исполнитель</label>
          <select
            id="assignee"
            v-model="assigneeId"
            class="form-select"
          >
            <option :value="null">Не назначен</option>
            <option v-for="engineer in engineers" :key="engineer.id" :value="engineer.id">
              {{ engineer.full_name }} ({{ engineer.role === 'manager' ? 'Менеджер' : 'Инженер' }})
            </option>
          </select>
        </div>
        <div class="form-group">
          <BaseInput
            v-model="dueDate"
            label="Срок исполнения"
            type="date"
          />
        </div>
      </div>

      <div class="form-actions">
        <BaseButton
          type="button"
          variant="outline"
          @click="cancel"
          :disabled="isSaving"
        >
          Отмена
        </BaseButton>
        <BaseButton
          type="submit"
          :loading="isSaving"
        >
          {{ props.defectId ? 'Сохранить' : 'Создать дефект' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<style scoped>
.defect-form {
  max-width: 800px;
  margin: 0 auto;
}

.loading {
  text-align: center;
  padding: 2rem;
  color: var(--color-text-light);
}

.form-group {
  margin-bottom: 1.5rem;
  width: 100%;
}

.form-row {
  display: flex;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.form-row .form-group {
  margin-bottom: 0;
  flex: 1;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.form-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  background-color: var(--color-background);
  font-size: 1rem;
}

.form-text {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.8rem;
  color: var(--color-text-light);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.form-error {
  padding: 1rem;
  margin-bottom: 1.5rem;
  background-color: rgba(244, 67, 54, 0.1);
  color: #d32f2f;
  border: 1px solid rgba(244, 67, 54, 0.3);
  border-radius: var(--border-radius);
}

@media (max-width: 768px) {
  .form-row {
    flex-direction: column;
    gap: 1.5rem;
  }
}
</style>
