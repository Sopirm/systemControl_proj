<script setup lang="ts">
import { ref } from 'vue'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'

const selectedReport = ref('defects-by-status')
const isLoading = ref(false)
const exportFormat = ref('excel')
</script>

<template>
  <div class="reports-page">
    <div class="page-header flex-between">
      <h1>Отчеты</h1>
      <BaseButton>Экспорт</BaseButton>
    </div>

    <div class="reports-container">
      <aside class="reports-sidebar">
        <h2>Список отчетов</h2>
        <ul class="reports-list">
          <li>
            <button 
              :class="{ active: selectedReport === 'defects-by-status' }"
              @click="selectedReport = 'defects-by-status'"
            >
              Дефекты по статусам
            </button>
          </li>
          <li>
            <button 
              :class="{ active: selectedReport === 'defects-by-priority' }"
              @click="selectedReport = 'defects-by-priority'"
            >
              Дефекты по приоритету
            </button>
          </li>
          <li>
            <button 
              :class="{ active: selectedReport === 'defects-by-project' }"
              @click="selectedReport = 'defects-by-project'"
            >
              Дефекты по объектам
            </button>
          </li>
          <li>
            <button 
              :class="{ active: selectedReport === 'completion-rate' }"
              @click="selectedReport = 'completion-rate'"
            >
              Степень готовности объектов
            </button>
          </li>
        </ul>

        <div class="export-section">
          <h3>Экспорт данных</h3>
          <div class="export-options">
            <div class="radio-group">
              <label>
                <input type="radio" v-model="exportFormat" value="excel"> Excel
              </label>
              <label>
                <input type="radio" v-model="exportFormat" value="csv"> CSV
              </label>
            </div>
            <BaseButton variant="outline">Скачать</BaseButton>
          </div>
        </div>
      </aside>

      <div class="report-content">
        <BaseCard :title="selectedReport === 'defects-by-status' ? 'Дефекты по статусам' : 
                          selectedReport === 'defects-by-priority' ? 'Дефекты по приоритету' : 
                          selectedReport === 'defects-by-project' ? 'Дефекты по объектам' : 
                          'Степень готовности объектов'">
          <div v-if="isLoading" class="loading-indicator">
            Загрузка данных...
          </div>
          <div v-else>
            <!-- Здесь будет содержимое выбранного отчета -->
            <div class="report-placeholder">
              <p>Диаграмма для выбранного отчета будет отображена здесь</p>
              <div class="placeholder-chart"></div>
              <p class="placeholder-text">
                В будущих версиях здесь будут интерактивные графики, отображающие аналитические данные 
                по дефектам и объектам строительства.
              </p>
            </div>
          </div>
        </BaseCard>
      </div>
    </div>
  </div>
</template>

<style scoped>
.reports-page {
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

h2 {
  font-size: 1.2rem;
  color: var(--color-text);
  margin-bottom: 1rem;
}

h3 {
  font-size: 1rem;
  margin-bottom: 0.5rem;
}

.reports-container {
  display: grid;
  grid-template-columns: 250px 1fr;
  gap: 1.5rem;
}

.reports-sidebar {
  background-color: #f8f9fa;
  border-radius: var(--border-radius);
  padding: 1rem;
}

.reports-list {
  list-style: none;
  padding: 0;
  margin-bottom: 2rem;
}

.reports-list li {
  margin-bottom: 0.5rem;
}

.reports-list button {
  display: block;
  width: 100%;
  text-align: left;
  padding: 0.6rem 0.8rem;
  border: none;
  background-color: transparent;
  color: var(--color-text);
  cursor: pointer;
  border-radius: var(--border-radius);
  transition: var(--transition-default);
}

.reports-list button:hover {
  background-color: rgba(30, 144, 255, 0.1);
}

.reports-list button.active {
  background-color: var(--color-primary);
  color: white;
}

.report-content {
  min-height: 400px;
}

.export-section {
  padding-top: 1rem;
  border-top: 1px solid var(--color-border);
}

.export-options {
  margin-top: 0.5rem;
}

.radio-group {
  margin-bottom: 1rem;
}

.radio-group label {
  display: block;
  margin-bottom: 0.5rem;
  cursor: pointer;
}

.loading-indicator {
  text-align: center;
  padding: 2rem;
  color: var(--color-text-light);
}

.report-placeholder {
  text-align: center;
  padding: 2rem 1rem;
}

.placeholder-chart {
  height: 200px;
  background-color: #f5f5f5;
  border-radius: var(--border-radius);
  margin: 1.5rem auto;
  position: relative;
}

.placeholder-chart:after {
  content: "";
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 100px;
  height: 100px;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='%231e90ff'%3E%3Cpath d='M3 13h2v7c0 1.103.897 2 2 2h12c1.103 0 2-.897 2-2v-7h2v-2H3v2zm4 7v-9h10v9H7zm6-11V7H11v2H9V5h2V3h2v2h2v2h-2z'%3E%3C/path%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: center;
  opacity: 0.2;
}

.placeholder-text {
  color: var(--color-text-light);
  font-size: 0.9rem;
  max-width: 600px;
  margin: 0 auto;
}

@media (max-width: 768px) {
  .reports-container {
    grid-template-columns: 1fr;
  }
  
  .reports-sidebar {
    margin-bottom: 1rem;
  }
}
</style>
