<script setup lang="ts">
import { RouterView } from 'vue-router'
import AppHeader from './components/AppHeader.vue'
</script>

<template>
  <AppHeader />
  <main class="main-content">
    <div class="container">
      <RouterView />
    </div>
  </main>
</template>

<style>
@import '@/assets/main.css';

.main-content {
  padding: 2rem 0;
  min-height: calc(100vh - 70px);
}
</style>