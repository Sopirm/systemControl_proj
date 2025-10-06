<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import BaseCard from '../components/BaseCard.vue'
import BaseButton from '../components/BaseButton.vue'
import { defectService, type Defect } from '../services/defectService'
import { Chart, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, BarController, DoughnutController } from 'chart.js'
import { projectService, type Project } from '../services/projectService'
Chart.register(ArcElement, Tooltip, Legend, CategoryScale, LinearScale, BarElement, BarController, DoughnutController)
Chart.defaults.responsive = true
Chart.defaults.maintainAspectRatio = false

const selectedReport = ref<'defects-by-status' | 'defects-by-priority' | 'defects-by-project' | 'overview'>('defects-by-status')
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
const projects = ref<Project[]>([])
const statusAgg = ref<{ labels: string[]; data: number[] }>({ labels: [], data: [] })
const priorityAgg = ref<{ labels: string[]; data: number[] }>({ labels: [], data: [] })
const projectAgg = ref<{ labels: string[]; data: number[] }>({ labels: [], data: [] })
const sum = (arr: number[]) => arr.reduce((s, n) => s + n, 0)

// форматирование даты: YYYY-MM-DD HH:MM:SS
function fmtDateTime(v?: string) {
  if (!v) return '—'
  const d = new Date(v)
  if (isNaN(d.getTime())) return '—'
  const pad = (n: number) => n.toString().padStart(2, '0')
  const y = d.getFullYear()
  const m = pad(d.getMonth() + 1)
  const day = pad(d.getDate())
  const hh = pad(d.getHours())
  const mm = pad(d.getMinutes())
  const ss = pad(d.getSeconds())
  // По требованию: HH:MM:SS (через двоеточие)
  return `${y}-${m}-${day} ${hh}:${mm}:${ss}`
}

type OverviewRow = {
  project_id: number
  project_name: string
  project_status?: string
  project_location?: string
  project_start_date?: string
  project_end_date?: string
  manager_name?: string
  defect_id?: number
  defect_title?: string
  defect_status?: string
  defect_priority?: string
  reporter?: string
  assignee?: string
  due_date?: string
  defect_created_at?: string
}
const overviewRows = ref<OverviewRow[]>([])

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
    // Загружаем проекты и дефекты
    projects.value = await projectService.getAllProjects()
    defects.value = await defectService.getAllDefects()
    const byStatus = aggregateByStatus(defects.value)
    const byPriority = aggregateByPriority(defects.value)
    const byProject = aggregateByProject(defects.value)
    statusAgg.value = byStatus
    priorityAgg.value = byPriority
    projectAgg.value = byProject
    // Формируем строки общего отчета (все проекты + их дефекты)
    const map = new Map<number, Defect[]>()
    for (const d of defects.value) {
      const arr = map.get(d.project_id) || []
      arr.push(d)
      map.set(d.project_id, arr)
    }
    const rows: OverviewRow[] = []
    for (const p of projects.value) {
      const ds = map.get(p.id) || []
      if (ds.length === 0) {
        rows.push({
          project_id: p.id,
          project_name: p.name,
          project_status: p.status,
          project_location: p.location,
          project_start_date: p.start_date,
          project_end_date: p.end_date,
          manager_name: p.manager?.full_name || p.manager?.username,
        })
      } else {
        for (const d of ds) {
          rows.push({
            project_id: p.id,
            project_name: p.name,
            project_status: p.status,
            project_location: p.location,
            project_start_date: p.start_date,
            project_end_date: p.end_date,
            manager_name: p.manager?.full_name || p.manager?.username,
            defect_id: d.id,
            defect_title: d.title,
            defect_status: d.status,
            defect_priority: d.priority,
            reporter: (d.reporter?.full_name || d.reporter?.username),
            assignee: (d.assignee?.full_name || d.assignee?.username),
            due_date: d.due_date,
            defect_created_at: d.created_at,
          })
        }
      }
    }
    overviewRows.value = rows
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
  if (selectedReport.value === 'overview') {
    return null
  }
  if (selectedReport.value === 'defects-by-status') {
    const rows = statusAgg.value.labels.map((l, i) => ({ label: l, value: statusAgg.value.data[i] || 0 }))
    return { title: 'Дефекты по статусам', rows, file: 'defects_by_status', sheet: 'По статусам' }
  }
  if (selectedReport.value === 'defects-by-priority') {
    const rows = priorityAgg.value.labels.map((l, i) => ({ label: l, value: priorityAgg.value.data[i] || 0 }))
    return { title: 'Дефекты по приоритету', rows, file: 'defects_by_priority', sheet: 'По приоритету' }
  }
  if (selectedReport.value === 'defects-by-project') {
    const rows = projectAgg.value.labels.map((l, i) => ({ label: l, value: projectAgg.value.data[i] || 0 }))
    return { title: 'Дефекты по объектам', rows, file: 'defects_by_project', sheet: 'По объектам' }
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

// CSV escape for quotes (compatible with older JS targets)
function csvEsc(val: any): string {
  return '"' + String(val).replace(/"/g, '""') + '"'
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
  const lines = [header.join(sep), ...agg.rows.map(r => [csvEsc(r.label), String(r.value)].join(sep))]
  const csv = '\ufeff' + lines.join('\n')
  downloadBlob(new Blob([csv], { type: 'text/csv;charset=utf-8;' }), `${agg.file}.csv`)
}

function exportOverviewCsv() {
  const sep = ';'
  const header = ['ID проекта', 'Название проекта', 'Статус проекта', 'Локация', 'Дата начала', 'Дата окончания', 'Менеджер', 'ID дефекта', 'Заголовок дефекта', 'Статус дефекта', 'Приоритет', 'Автор', 'Исполнитель', 'Срок', 'Создан']
  const rows = overviewRows.value.map(r => [
    r.project_id,
    r.project_name,
    r.project_status || '',
    r.project_location || '',
    fmtDateTime(r.project_start_date),
    fmtDateTime(r.project_end_date),
    r.manager_name || '',
    r.defect_id ?? '',
    r.defect_title || '',
    statusLabel(r.defect_status || ''),
    priorityLabel(r.defect_priority || ''),
    r.reporter || '',
    r.assignee || '',
    fmtDateTime(r.due_date),
    fmtDateTime(r.defect_created_at)
  ])
  if (rows.length === 0) {
    console.warn('Нет данных для экспорта CSV (overview)')
    return
  }
  const lines = [header.join(sep), ...rows.map(r => r.map(csvEsc).join(sep))]
  const csv = '\ufeff' + lines.join('\n')
  downloadBlob(new Blob([csv], { type: 'text/csv;charset=utf-8;' }), `Общий_отчет_проекты_дефекты.csv`)
}

async function exportExcel() {
  const agg = getCurrentAgg()
  if (!agg) {
    // Обработка общего отчета (таблица)
    if (selectedReport.value === 'overview') {
      const XLSX: any = await import('xlsx')
      const header = ['ID проекта', 'Название проекта', 'Статус проекта', 'Локация', 'Дата начала', 'Дата окончания', 'Менеджер', 'ID дефекта', 'Заголовок дефекта', 'Статус дефекта', 'Приоритет', 'Автор', 'Исполнитель', 'Срок', 'Создан']
      const rows = overviewRows.value.map(r => [
        r.project_id,
        r.project_name,
        r.project_status || '',
        r.project_location || '',
        fmtDateTime(r.project_start_date),
        fmtDateTime(r.project_end_date),
        r.manager_name || '',
        r.defect_id ?? '',
        r.defect_title || '',
        statusLabel(r.defect_status || ''),
        priorityLabel(r.defect_priority || ''),
        r.reporter || '',
        r.assignee || '',
        fmtDateTime(r.due_date),
        fmtDateTime(r.defect_created_at)
      ])
      if (rows.length === 0) {
        console.warn('Нет данных для экспорта Excel (overview)')
        return
      }
      const aoa = [['Общий отчет (проекты и дефекты)'], [], header, ...rows]
      const ws = XLSX.utils.aoa_to_sheet(aoa)
      const wb = XLSX.utils.book_new()
      XLSX.utils.book_append_sheet(wb, ws, 'Общий отчет')
      const wbout = XLSX.write(wb, { bookType: 'xlsx', type: 'array' })
      downloadBlob(new Blob([wbout], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }), 'Общий_отчет_проекты_дефекты.xlsx')
    }
    return
  }
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
  if (selectedReport.value === 'overview') {
    if (overviewRows.value.length === 0) {
      console.warn('Нет данных для экспорта (overview)')
      return
    }
    if (exportFormat.value === 'csv') return exportOverviewCsv()
    return exportExcel()
  }
  if (exportFormat.value === 'csv') return exportCsv()
  return exportExcel()
}
</script>

<template>
  <div class="reports-page" :style="{ overflowX: 'hidden' }">
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
          <li>
            <button 
              :class="{ active: selectedReport === 'overview' }"
              @click="selectedReport = 'overview'"
            >
              Общий отчет (проекты + дефекты)
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

      <div class="report-content" :style="{ minWidth: '0' }">
        <BaseCard :title="selectedReport === 'defects-by-status' ? 'Дефекты по статусам' : 
                          selectedReport === 'defects-by-priority' ? 'Дефекты по приоритету' : 
                          selectedReport === 'defects-by-project' ? 'Дефекты по объектам' :
                          'Общий отчет (проекты и дефекты)'">
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
            <div v-show="selectedReport === 'overview'" class="table-container" :style="{ overflowX: 'auto', maxWidth: '100%', display: 'block' }">
              <template v-if="overviewRows.length > 0">
                <table class="report-table">
                  <thead>
                    <tr>
                      <th>ID проекта</th>
                      <th>Название проекта</th>
                      <th>Статус проекта</th>
                      <th>Локация</th>
                      <th>Дата начала</th>
                      <th>Дата окончания</th>
                      <th>Менеджер</th>
                      <th>ID дефекта</th>
                      <th>Заголовок дефекта</th>
                      <th>Статус дефекта</th>
                      <th>Приоритет</th>
                      <th>Автор</th>
                      <th>Исполнитель</th>
                      <th>Срок</th>
                      <th>Создан</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(r, idx) in overviewRows" :key="idx">
                      <td>{{ r.project_id }}</td>
                      <td>{{ r.project_name }}</td>
                      <td>{{ r.project_status || '—' }}</td>
                      <td>{{ r.project_location || '—' }}</td>
                      <td>{{ fmtDateTime(r.project_start_date) }}</td>
                      <td>{{ fmtDateTime(r.project_end_date) }}</td>
                      <td>{{ r.manager_name || '—' }}</td>
                      <td>{{ r.defect_id ?? '—' }}</td>
                      <td>{{ r.defect_title || '—' }}</td>
                      <td>{{ statusLabel(r.defect_status || '') || '—' }}</td>
                      <td>{{ priorityLabel(r.defect_priority || '') || '—' }}</td>
                      <td>{{ r.reporter || '—' }}</td>
                      <td>{{ r.assignee || '—' }}</td>
                      <td>{{ fmtDateTime(r.due_date) }}</td>
                      <td>{{ fmtDateTime(r.defect_created_at) }}</td>
                    </tr>
                  </tbody>
                </table>
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
  /* не позволяем странице скроллиться горизонтально, скролл должен быть на контейнере таблицы */
  overflow-x: hidden;
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
  /* скрываем переполнение, чтобы прокрутка была внутри контента */
  overflow: hidden;
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
  /* критично для grid: разрешаем ужиматься по ширине, чтобы внутренний блок мог прокручиваться */
  min-width: 0;
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

/* Табличный отчет: горизонтальная прокрутка, чтобы не выходил за экран */
.table-container {
  display: block;
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
  scrollbar-gutter: stable both-edges;
}

.report-table {
  width: max-content; /* ширина = ширина содержимого, чтобы появился горизонтальный скролл */
  min-width: 100%;
  border-collapse: collapse;
}

.report-table th,
.report-table td {
  padding: 0.5rem 0.75rem;
  border-bottom: 1px solid var(--color-border);
  white-space: nowrap; /* чтобы длинные заголовки не ломали сетку */
}

.report-table thead th {
  position: sticky;
  top: 0;
  background: var(--color-background-soft);
  z-index: 1;
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
