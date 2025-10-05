<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import BaseInput from './BaseInput.vue'
import BaseButton from './BaseButton.vue'
import { projectService, type Project, type ProjectCreate, type ProjectUpdate } from '../services/projectService'
import { userService } from '../services/userService'
import type { User } from '../services/userService'

const props = defineProps({
  projectId: {
    type: [Number, String],
    default: null
  },
  project: {
    type: Object as () => Project | null,
    default: null
  }
})

const emit = defineEmits(['saved', 'cancelled'])

// Состояние формы
const name = ref('')
const description = ref('')
const location = ref('')
const startDate = ref('')
const endDate = ref('')
const status = ref('active')
const managerId = ref<number | null>(null)
const isLoading = ref(false)
const isSaving = ref(false)
const error = ref('')
const managers = ref<User[]>([])

// Опции для выбора статуса проекта
const statusOptions = [
  { value: 'active', label: 'Активный' },
  { value: 'completed', label: 'Завершен' },
  { value: 'suspended', label: 'Приостановлен' }
]

// Функция для загрузки менеджеров
const loadManagers = async () => {
  try {
    const users = await userService.getAllUsers()
    managers.value = users.filter(user => user.role === 'manager')
  } catch (err) {
    console.error('Ошибка при загрузке списка менеджеров:', err)
    error.value = 'Не удалось загрузить список менеджеров'
  }
}

// Функция для загрузки проекта при редактировании
const loadProject = async () => {
  if (!props.projectId) return
  
  try {
    isLoading.value = true
    error.value = ''
    
    const projectData = await projectService.getProjectById(Number(props.projectId))
    name.value = projectData.name
    description.value = projectData.description || ''
    location.value = projectData.location || ''
    startDate.value = projectData.start_date ? new Date(projectData.start_date).toISOString().split('T')[0] : ''
    endDate.value = projectData.end_date ? new Date(projectData.end_date).toISOString().split('T')[0] : ''
    status.value = projectData.status
    managerId.value = projectData.manager_id
  } catch (err) {
    console.error('Ошибка при загрузке проекта:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить данные проекта'
  } finally {
    isLoading.value = false
  }
}

// Наблюдаем за изменением props.project
watch(() => props.project, (newProject) => {
  if (newProject) {
    name.value = newProject.name
    description.value = newProject.description || ''
    location.value = newProject.location || ''
    startDate.value = newProject.start_date ? new Date(newProject.start_date).toISOString().split('T')[0] : ''
    endDate.value = newProject.end_date ? new Date(newProject.end_date).toISOString().split('T')[0] : ''
    status.value = newProject.status
    managerId.value = newProject.manager_id
  }
}, { immediate: true })

// Валидация формы
const validateForm = () => {
  if (!name.value) {
    error.value = 'Имя проекта обязательно для заполнения'
    return false
  }
  if (!location.value) {
    error.value = 'Местоположение обязательно для заполнения'
    return false
  }
  if (!startDate.value) {
    error.value = 'Дата начала обязательна для заполнения'
    return false
  }
  if (!endDate.value) {
    error.value = 'Дата окончания обязательна для заполнения'
    return false
  }
  if (!managerId.value) {
    error.value = 'Необходимо выбрать менеджера проекта'
    return false
  }
  return true
}

// Сохранение проекта
const saveProject = async () => {
  error.value = ''
  
  if (!validateForm()) {
    return
  }
  
  try {
    isSaving.value = true
    
    // Форматирование дат в ISO формат с временем
    const formatDateForBackend = (dateStr: string) => {
      if (!dateStr) return ''
      // Преобразуем строку даты в формат ISO с временем 00:00:00 в UTC
      const date = new Date(dateStr)
      return date.toISOString()
    }
    
    const projectData: ProjectCreate | ProjectUpdate = {
      name: name.value,
      description: description.value,
      location: location.value,
      start_date: formatDateForBackend(startDate.value),
      end_date: formatDateForBackend(endDate.value),
      status: status.value,
      manager_id: managerId.value!
    }
    
    let result
    if (props.projectId) {
      // Обновление проекта
      result = await projectService.updateProject(Number(props.projectId), projectData)
    } else {
      // Создание нового проекта
      result = await projectService.createProject(projectData as ProjectCreate)
    }
    
    emit('saved', result)
  } catch (err) {
    console.error('Ошибка при сохранении проекта:', err)
    error.value = err instanceof Error ? err.message : 'Ошибка при сохранении проекта'
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
  loadManagers()
  if (props.projectId) {
    loadProject()
  }
})
</script>

<template>
  <div class="project-form">
    <div v-if="isLoading" class="loading">
      Загрузка данных проекта...
    </div>
    <form v-else @submit.prevent="saveProject">
      <div v-if="error" class="form-error">{{ error }}</div>

      <div class="form-group">
        <BaseInput
          v-model="name"
          label="Название проекта"
          placeholder="Введите название проекта"
          required
        />
      </div>

      <div class="form-group">
        <BaseInput
          v-model="description"
          label="Описание"
          placeholder="Введите описание проекта"
          type="textarea"
        />
      </div>

      <div class="form-group">
        <BaseInput
          v-model="location"
          label="Местоположение"
          placeholder="Введите адрес объекта"
          required
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <BaseInput
            v-model="startDate"
            label="Дата начала"
            type="date"
            required
          />
        </div>
        <div class="form-group">
          <BaseInput
            v-model="endDate"
            label="Дата окончания"
            type="date"
            required
          />
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label for="status" class="form-label">Статус</label>
          <select
            id="status"
            v-model="status"
            class="form-select"
            required
          >
            <option v-for="option in statusOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
        <div class="form-group">
          <label for="manager" class="form-label">Менеджер проекта</label>
          <select
            id="manager"
            v-model="managerId"
            class="form-select"
            required
          >
            <option value="">Выберите менеджера</option>
            <option v-for="manager in managers" :key="manager.id" :value="manager.id">
              {{ manager.full_name }}
            </option>
          </select>
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
          {{ props.projectId ? 'Сохранить' : 'Создать проект' }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<style scoped>
.project-form {
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
