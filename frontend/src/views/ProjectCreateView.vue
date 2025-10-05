<script setup lang="ts">
import { useRouter } from 'vue-router'
import ProjectForm from '../components/ProjectForm.vue'
import type { Project } from '../services/projectService'

const router = useRouter()

const handleSaved = (project: Project) => {
  router.push(`/projects/${project.id}`)
}

const handleCancelled = () => {
  router.push('/projects')
}
</script>

<template>
  <div class="project-create-page">
    <div class="page-header">
      <h1>Создание нового проекта</h1>
    </div>
    
    <ProjectForm 
      @saved="handleSaved" 
      @cancelled="handleCancelled" 
    />
  </div>
</template>

<style scoped>
.project-create-page {
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
</style>
