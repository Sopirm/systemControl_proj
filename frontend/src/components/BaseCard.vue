<script setup lang="ts">
defineProps<{
  title?: string
  border?: boolean
}>()
</script>

<template>
  <div class="base-card" :class="{ 'with-border': border }">
    <div v-if="title" class="card-header">
      <h3 class="card-title">{{ title }}</h3>
    </div>
    <div class="card-content">
      <slot />
    </div>
    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<style scoped>
.base-card {
  background-color: var(--color-background);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow-default);
  overflow: hidden;
}

.with-border {
  border: 1px solid var(--color-border);
}

.card-header {
  padding: 1rem;
  border-bottom: 1px solid var(--color-border);
}

.card-title {
  margin: 0;
  color: var(--color-primary);
  font-size: 1.2rem;
  font-weight: 500;
}

.card-content {
  padding: 1rem;
}

.card-footer {
  padding: 1rem;
  border-top: 1px solid var(--color-border);
  background-color: rgba(0, 0, 0, 0.02);
}
</style>
