<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import DefectItem from '../components/DefectItem.vue'
import BaseButton from '../components/BaseButton.vue'
import CreateDefectModal from '../components/CreateDefectModal.vue'
import { defectService, type Defect } from '../services/defectService'
import { useAuth } from '../composables/useAuth'

const router = useRouter()
const { isManager, isEngineer } = useAuth()
const defects = ref<Defect[]>([])
const isLoading = ref(true)
const error = ref('')
const filterStatus = ref('all')
const showCreateDefectModal = ref(false)

// Получаем права доступа
const permissions = defectService.checkPermissions()

// Загрузка всех дефектов
const loadDefects = async () => {
  try {
    isLoading.value = true
    error.value = ''
    const data = await defectService.getAllDefects()
    defects.value = data
    console.log('Загружено дефектов:', data.length)
  } catch (err) {
    console.error('Ошибка при загрузке дефектов:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить дефекты'
  } finally {
    isLoading.value = false
  }
}

// Фильтрация дефектов по статусу
const filteredDefects = computed(() => {
  if (filterStatus.value === 'all') {
    return defects.value
  }
  return defects.value.filter(defect => defect.status === filterStatus.value)
})

// Открытие страницы дефекта
function openDefect(id: number) {
  router.push(`/defects/${id}`)
}

// Открытие модального окна создания дефекта
const openCreateDefectModal = () => {
  showCreateDefectModal.value = true
}

// Закрытие модального окна создания дефекта
const closeCreateDefectModal = () => {
  showCreateDefectModal.value = false
}

// Обработка создания дефекта
const handleDefectCreated = (defect: Defect) => {
  // Добавляем новый дефект в список
  defects.value.push(defect)
  
  // Закрываем модальное окно
  closeCreateDefectModal()
}

// Загрузка дефектов при монтировании компонента
onMounted(() => {
  loadDefects()
})
</script>

<template>
  <div class="defects-page">
    <div v-if="error" class="alert alert-error">{{ error }}</div>
    
    <div class="page-header flex-between">
      <h1>Дефекты</h1>
      <BaseButton 
        v-if="permissions.canCreate" 
        @click="openCreateDefectModal"
      >
        Создать дефект
      </BaseButton>
    </div>

    <div class="filters">
      <div class="filter-group">
        <label for="status-filter">Статус:</label>
        <select id="status-filter" v-model="filterStatus">
          <option value="all">Все</option>
          <option value="new">Новые</option>
          <option value="in_progress">В работе</option>
          <option value="review">На проверке</option>
          <option value="closed">Закрытые</option>
          <option value="cancelled">Отмененные</option>
        </select>
      </div>
    </div>

    <div v-if="isLoading" class="loading-indicator">Загрузка дефектов...</div>

    <div v-else-if="filteredDefects.length === 0" class="empty-state">
      <p>Нет дефектов по заданным критериям</p>
    </div>

    <div v-else class="defects-list">
      <DefectItem 
        v-for="defect in filteredDefects" 
        :key="defect.id" 
        :defect="defect"
      >
        <template #actions>
          <RouterLink :to="`/defects/${defect.id}`">
            <BaseButton variant="outline">Подробнее</BaseButton>
          </RouterLink>
        </template>
      </DefectItem>
    </div>
    
    <!-- Модальное окно создания дефекта -->
    <CreateDefectModal
      :show="showCreateDefectModal"
      :projectId="null"
      @close="closeCreateDefectModal"
      @created="handleDefectCreated"
    />
  </div>
</template>

<style scoped>
.defects-page {
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

.filters {
  margin-bottom: 1.5rem;
  padding: 1rem;
  background-color: #f8f9fa;
  border-radius: var(--border-radius);
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-group label {
  font-weight: 500;
  color: var(--color-text);
}

.filter-group select {
  padding: 0.4rem 0.8rem;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  font-size: 0.9rem;
  width: auto;
}

.defects-list {
  margin-top: 1rem;
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