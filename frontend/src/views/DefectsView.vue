<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import DefectItem from '../components/DefectItem.vue'
import BaseButton from '../components/BaseButton.vue'

const router = useRouter()

interface Defect {
  id: number
  title: string
  description: string
  priority: 'low' | 'medium' | 'high'
  status: 'new' | 'in_progress' | 'review' | 'closed' | 'cancelled'
  project: {
    id: number
    name: string
  }
  assignee?: string
  dueDate?: string
}

// Временные данные для демонстрации
const defects = ref<Defect[]>([
  {
    id: 1,
    title: 'Трещина в несущей стене',
    description: 'На 3 этаже, 2 подъезд обнаружена трещина в несущей стене размером 1.5м',
    priority: 'high',
    status: 'new',
    project: {
      id: 1,
      name: 'ЖК "Солнечный"'
    },
    assignee: 'Иванов И.И.',
    dueDate: '2023-10-15'
  },
  {
    id: 2,
    title: 'Протечка в потолке',
    description: 'Кв. 42, 5 этаж - следы протечки на потолке в ванной комнате',
    priority: 'medium',
    status: 'in_progress',
    project: {
      id: 1,
      name: 'ЖК "Солнечный"'
    },
    assignee: 'Петров П.П.'
  },
  {
    id: 3,
    title: 'Неправильная укладка плитки',
    description: 'Холл 1 этажа - неровная укладка напольной плитки на площади 5м²',
    priority: 'low',
    status: 'review',
    project: {
      id: 2,
      name: 'Бизнес-центр "Меркурий"'
    },
    assignee: 'Сидоров С.С.',
    dueDate: '2023-10-10'
  }
])

const isLoading = ref(false)
const filterStatus = ref('all')

const filteredDefects = computed(() => {
  if (filterStatus.value === 'all') {
    return defects.value
  }
  return defects.value.filter(defect => defect.status === filterStatus.value)
})

function openDefect(id: number) {
  router.push(`/defects/${id}`)
}
</script>

<template>
  <div class="defects-page">
    <div class="page-header flex-between">
      <h1>Дефекты</h1>
      <BaseButton>Создать дефект</BaseButton>
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
        @click="openDefect(defect.id)"
      >
        <template #actions>
          <BaseButton variant="outline">Подробнее</BaseButton>
        </template>
      </DefectItem>
    </div>
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
</style>