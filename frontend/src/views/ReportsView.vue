<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'
import { defectService, type Defect } from '../services/defectService'
import { Chart, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, BarController, DoughnutController } from 'chart.js'
Chart.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, BarController, DoughnutController)
Chart.defaults.responsive = true
Chart.defaults.maintainAspectRatio = false

const selectedReport = ref<'defects-by-status' | 'defects-by-priority' | 'defects-by-project'>('defects-by-status')
const isLoading = ref(false)
const exportFormat = ref<'excel' | 'csv'>('excel')

// Canvas refs
const statusCanvas = ref<HTMLCanvasElement | null>(null)
const priorityCanvas = ref<HTMLCanvasElement | null>(null)
const projectCanvas = ref<HTMLCanvasElement | null>(null)

// Chart instances
let statusChart: Chart | null = null
let priorityChart: Chart | null = null
let projectChart: Chart | null = null

// Data
const defects = ref<Defect[]>([])
const statusAgg = ref<{ labels: string[]; data: number[] }>({ labels: [], data: [] })
const priorityAgg = ref<{ labels: string[]; data: number[] }>({ labels: [], data: [] })
const projectAgg = ref<{ labels: string[]; data: number[] }>({ labels: [], data: [] })
const sum = (arr: number[]) => arr.reduce((s, n) => s + n, 0)

// Helpers for labels
const statusLabel = (s: string) => ({
  new: 'Новый',
  in_progress: 'В работе',
  review: 'На проверке',
  closed: 'Закрыт',
  cancelled: 'Отменен', // фронт
  canceled: 'Отменен'   // бэк
} as Record<string, string>)[s] || s

const priorityLabel = (p: string) => ({
  low: 'Низкий',
  medium: 'Средний',
  high: 'Высокий'
} as Record<string, string>)[p] || p

// Aggregations
function aggregateByStatus(items: Defect[]) {
  // Бэкенд использует 'canceled' (без второй L)
  const keys = ['new', 'in_progress', 'review', 'closed', 'canceled']
  const counts = keys.map(k => items.filter(d => d.status === k).length)
  return { labels: keys.map(statusLabel), data: counts }
}

function aggregateByPriority(items: Defect[]) {
  const keys = ['low', 'medium', 'high']
  const counts = keys.map(k => items.filter(d => d.priority === k).length)
  return { labels: keys.map(priorityLabel), data: counts }
}

function aggregateByProject(items: Defect[]) {
  const map = new Map<string, number>()
  for (const d of items) {
    const name = d.project?.name || `Проект #${d.project_id}`
    map.set(name, (map.get(name) || 0) + 1)
  }
  const labels = Array.from(map.keys())
  const data = Array.from(map.values())
  return { labels, data }
}

// Chart builders
function buildStatusChart(ctx: CanvasRenderingContext2D, labels: string[], data: number[]) {
  const chart = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels,
      datasets: [{
        label: 'Дефекты по статусу',
        data,
        backgroundColor: ['#e3f2fd', '#e8f5e9', '#fff3e0', '#f5f5f5', '#ffebee'],
        borderColor: ['#1976d2', '#2e7d32', '#e65100', '#616161', '#c62828']
      }]
    },
    options: { responsive: true, maintainAspectRatio: false }
  })
  chart.resize();
  chart.update();
  return chart
}

function buildPriorityChart(ctx: CanvasRenderingContext2D, labels: string[], data: number[]) {
  const chart = new Chart(ctx, {
    type: 'doughnut',
    data: {
      labels,
      datasets: [{
        label: 'Дефекты по приоритету',
        data,
        backgroundColor: ['#e8f5e9', '#fff3e0', '#ffebee'],
        borderColor: ['#2e7d32', '#e65100', '#c62828']
      }]
    },
    options: { responsive: true, maintainAspectRatio: false }
  })
  chart.resize();
  chart.update();
  return chart
}

function buildProjectChart(ctx: CanvasRenderingContext2D, labels: string[], data: number[]) {
  const chart = new Chart(ctx, {
    type: 'bar',
    data: {
      labels,
      datasets: [{
        label: 'Дефекты по объектам',
        data,
        backgroundColor: '#90caf9',
        borderColor: '#1976d2'
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: { y: { beginAtZero: true, ticks: { precision: 0 } } }
    }
  })
  chart.resize();
  chart.update();
  return chart
}

function destroyCharts() {
  statusChart?.destroy(); statusChart = null
  priorityChart?.destroy(); priorityChart = null
  projectChart?.destroy(); projectChart = null
}

function renderSelectedChart() {
  // Перерисовываем только активный график, остальные уничтожаем
  if (selectedReport.value === 'defects-by-status') {
    priorityChart?.destroy(); priorityChart = null
    projectChart?.destroy(); projectChart = null
    if (statusCanvas.value && sum(statusAgg.value.data) > 0 && statusCanvas.value.offsetParent !== null) {
      const ctx = statusCanvas.value.getContext('2d')!
      statusChart?.destroy();
      statusChart = buildStatusChart(ctx, statusAgg.value.labels, statusAgg.value.data)
    }
    return
  }
  if (selectedReport.value === 'defects-by-priority') {
    statusChart?.destroy(); statusChart = null
    projectChart?.destroy(); projectChart = null
    if (priorityCanvas.value && sum(priorityAgg.value.data) > 0 && priorityCanvas.value.offsetParent !== null) {
      const ctx = priorityCanvas.value.getContext('2d')!
      priorityChart?.destroy();
      priorityChart = buildPriorityChart(ctx, priorityAgg.value.labels, priorityAgg.value.data)
    }
    return
  }
  if (selectedReport.value === 'defects-by-project') {
    statusChart?.destroy(); statusChart = null
    priorityChart?.destroy(); priorityChart = null
    if (projectCanvas.value && sum(projectAgg.value.data) > 0 && projectCanvas.value.offsetParent !== null) {
      const ctx = projectCanvas.value.getContext('2d')!
      projectChart?.destroy();
      projectChart = buildProjectChart(ctx, projectAgg.value.labels, projectAgg.value.data)
    }
  }
}

async function loadAndRender() {
  isLoading.value = true
  try {
    defects.value = await defectService.getAllDefects()
    const byStatus = aggregateByStatus(defects.value)
    const byPriority = aggregateByPriority(defects.value)
    const byProject = aggregateByProject(defects.value)
    statusAgg.value = byStatus
    priorityAgg.value = byPriority
    projectAgg.value = byProject
    // Перерисовываем только активный график, когда DOM обновлен и контейнер видим
    await nextTick()
    // debug размеров контейнера
    const container = selectedReport.value === 'defects-by-status' ? statusCanvas.value?.parentElement
      : selectedReport.value === 'defects-by-priority' ? priorityCanvas.value?.parentElement
      : projectCanvas.value?.parentElement
    console.debug('ReportsView agg:', { byStatus, byPriority, byProject, size: { w: container?.clientWidth, h: container?.clientHeight } })
    requestAnimationFrame(() => renderSelectedChart())
  } catch (e) {
    console.error('Ошибка загрузки данных отчетов:', e)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  loadAndRender()
})

onBeforeUnmount(() => destroyCharts())

// Перерисовка при переключении вкладок отчета, чтобы избегать отрисовки в скрытом контейнере
watch(selectedReport, async () => {
  await nextTick()
  // В некоторых случаях стили применяются кадром позже
  requestAnimationFrame(() => renderSelectedChart())
})

// ===== Export helpers =====
type Row = { label: string; value: number }

function getCurrentAgg(): { title: string; rows: Row[]; file: string; sheet: string } | null {
  if (selectedReport.value === 'defects-by-status') {
    const rows = statusAgg.value.labels.map((l, i) => ({ label: l, value: statusAgg.value.data[i] || 0 }))
    return { title: 'Дефекты по статусам', rows, file: 'defects_by_status', sheet: 'ByStatus' }
  }
  if (selectedReport.value === 'defects-by-priority') {
    const rows = priorityAgg.value.labels.map((l, i) => ({ label: l, value: priorityAgg.value.data[i] || 0 }))
    return { title: 'Дефекты по приоритету', rows, file: 'defects_by_priority', sheet: 'ByPriority' }
  }
  if (selectedReport.value === 'defects-by-project') {
    const rows = projectAgg.value.labels.map((l, i) => ({ label: l, value: projectAgg.value.data[i] || 0 }))
    return { title: 'Дефекты по объектам', rows, file: 'defects_by_project', sheet: 'ByProject' }
  }
  return null
}

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

function exportCsv() {
  const agg = getCurrentAgg()
  if (!agg) return
  const total = agg.rows.reduce((s, r) => s + r.value, 0)
  if (total === 0) {
    console.warn('Нет данных для экспорта CSV')
    return
  }
  const sep = ';'
  const header = ['Название', 'Количество']
  const lines = [header.join(sep), ...agg.rows.map(r => [`"${r.label.replaceAll('"', '""')}"`, String(r.value)].join(sep))]
  const csv = '\ufeff' + lines.join('\n')
  downloadBlob(new Blob([csv], { type: 'text/csv;charset=utf-8;' }), `${agg.file}.csv`)
}

async function exportExcel() {
  const agg = getCurrentAgg()
  if (!agg) return
  const total = agg.rows.reduce((s, r) => s + r.value, 0)
  if (total === 0) {
    console.warn('Нет данных для экспорта Excel')
    return
  }
  // динамически импортируем xlsx, чтобы не раздувать бандл
  const XLSX: any = await import('xlsx')
  const aoa = [[agg.title], [], ['Название', 'Количество'], ...agg.rows.map(r => [r.label, r.value])]
  const ws = XLSX.utils.aoa_to_sheet(aoa)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, agg.sheet)
  const wbout = XLSX.write(wb, { bookType: 'xlsx', type: 'array' })
  downloadBlob(new Blob([wbout], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }), `${agg.file}.xlsx`)
}

function onExport() {
  if (exportFormat.value === 'csv') return exportCsv()
  return exportExcel()
}
</script>

<template>
  <div class="reports-page">
    <div class="page-header flex-between">
      <h1>Отчеты</h1>
      <BaseButton :disabled="isLoading" @click="onExport">Экспорт</BaseButton>
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
            <BaseButton variant="outline" :disabled="isLoading" @click="onExport">Скачать</BaseButton>
          </div>
        </div>
      </aside>

      <div class="report-content">
        <BaseCard :title="selectedReport === 'defects-by-status' ? 'Дефекты по статусам' : 
                          selectedReport === 'defects-by-priority' ? 'Дефекты по приоритету' : 
                          'Дефекты по объектам'">
          <div v-if="isLoading" class="loading-indicator">
            Загрузка данных...
          </div>
          <div v-else>
            <div v-show="selectedReport === 'defects-by-status'" class="chart-container">
              <template v-if="sum(statusAgg.data) > 0">
                <canvas ref="statusCanvas"></canvas>
              </template>
              <template v-else>
                <div class="report-placeholder"><p class="placeholder-text">Нет данных для отображения</p></div>
              </template>
            </div>
            <div v-show="selectedReport === 'defects-by-priority'" class="chart-container">
              <template v-if="sum(priorityAgg.data) > 0">
                <canvas ref="priorityCanvas"></canvas>
              </template>
              <template v-else>
                <div class="report-placeholder"><p class="placeholder-text">Нет данных для отображения</p></div>
              </template>
            </div>
            <div v-show="selectedReport === 'defects-by-project'" class="chart-container">
              <template v-if="sum(projectAgg.data) > 0">
                <canvas ref="projectCanvas"></canvas>
              </template>
              <template v-else>
                <div class="report-placeholder"><p class="placeholder-text">Нет данных для отображения</p></div>
              </template>
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

.chart-container {
  position: relative;
  width: 100%;
  height: 380px;
}

canvas {
  width: 100% !important;
  height: 100% !important;
  display: block;
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
