<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import DefectForm from './DefectForm.vue'
import type { Defect } from '../services/defectService'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  projectId: {
    type: [Number, String, null],
    default: null
  }
})

const emit = defineEmits(['close', 'created'])

// Обработка нажатия клавиши Escape для закрытия модального окна
const handleEscKey = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && props.show) {
    emit('close')
  }
}

// Обработка сохранения дефекта
const handleSaved = (defect: Defect) => {
  emit('created', defect)
}

// Обработка отмены создания
const handleCancelled = () => {
  emit('close')
}

// Добавляем и удаляем обработчик клавиши Escape
onMounted(() => {
  document.addEventListener('keydown', handleEscKey)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscKey)
})
</script>

<template>
  <Teleport to="body">
    <div v-if="show" class="modal-overlay" @click="emit('close')">
      <div class="modal-container" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">Создание нового дефекта</h3>
          <button class="modal-close" @click="emit('close')">&times;</button>
        </div>
        <div class="modal-body">
          <DefectForm 
            :projectId="projectId" 
            @saved="handleSaved" 
            @cancelled="handleCancelled" 
          />
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-container {
  background-color: var(--color-background);
  border-radius: var(--border-radius);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  width: 100%;
  max-width: 800px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 1.5rem;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.modal-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-primary);
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: var(--color-text-light);
  padding: 0;
  margin: 0;
  line-height: 1;
}

.modal-body {
  margin-bottom: 1rem;
}
</style>
