<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ProjectForm from '../components/ProjectForm.vue'
import type { Project } from '../services/projectService'
import { projectService } from '../services/projectService'

const route = useRoute()
const router = useRouter()
const projectId = route.params.id as string
const project = ref<Project | null>(null)
const isLoading = ref(true)
const error = ref('')

const loadProject = async () => {
  try {
    isLoading.value = true
    error.value = ''
    
    const data = await projectService.getProjectById(Number(projectId))
    project.value = data
  } catch (err) {
    console.error('Ошибка при загрузке проекта:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить проект'
  } finally {
    isLoading.value = false
  }
}

const handleSaved = (project: Project) => {
  router.push(`/projects/${project.id}`)
}

const handleCancelled = () => {
  router.push(`/projects/${projectId}`)
}

onMounted(() => {
  loadProject()
})
</script>

<template>
  <div class="project-edit-page">
    <div class="page-header">
      <h1>Редактирование проекта</h1>
    </div>
    
    <div v-if="error" class="alert alert-error">{{ error }}</div>

    <div v-if="isLoading" class="loading-indicator">
      Загрузка данных проекта...
    </div>
    
    <ProjectForm 
      v-else
      :projectId="projectId" 
      :project="project"
      @saved="handleSaved" 
      @cancelled="handleCancelled" 
    />
  </div>
</template>

<style scoped>
.project-edit-page {
  padding: 1rem 0;
}

.page-header {
  margin-bottom: 2rem;
}

h1 {
  font-size: 1.8rem;
  font-weight: 600;
  color: var(--color-primary);
  margin: 0;
}

.loading-indicator {
  text-align: center;
  padding: 2rem;
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
