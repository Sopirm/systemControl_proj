<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'

const isLoggedIn = ref(false)
</script>

<template>
  <header class="app-header">
    <div class="container flex-between">
      <div class="logo">
        <RouterLink to="/">Система контроля</RouterLink>
      </div>
      <nav class="main-nav">
        <ul class="nav-list">
          <li v-if="isLoggedIn">
            <RouterLink to="/projects">Проекты</RouterLink>
          </li>
          <li v-if="isLoggedIn">
            <RouterLink to="/defects">Дефекты</RouterLink>
          </li>
          <li v-if="isLoggedIn">
            <RouterLink to="/reports">Отчеты</RouterLink>
          </li>
          <li v-if="isLoggedIn">
            <button class="btn-outline">Выйти</button>
          </li>
        </ul>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.app-header {
  background-color: var(--color-background);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1rem 0;
}

.logo a {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--color-primary);
}

.nav-list {
  display: flex;
  list-style: none;
  gap: 1.5rem;
  align-items: center;
}

.nav-list a {
  color: var(--color-text);
  font-weight: 500;
  transition: var(--transition-default);
}

.nav-list a:hover,
.nav-list a.router-link-active {
  color: var(--color-primary);
}

@media (max-width: 768px) {
  .app-header .container {
    flex-direction: column;
    gap: 1rem;
  }
  
  .nav-list {
    flex-wrap: wrap;
    justify-content: center;
    gap: 1rem;
  }
}
</style>
