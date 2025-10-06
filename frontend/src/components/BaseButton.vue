<script setup lang="ts">
defineProps<{
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'outline' | 'text' | 'danger'
  disabled?: boolean
}>()
</script>

<template>
  <button
    :type="type || 'button'"
    :class="['base-button', `btn-${variant || 'primary'}`]"
    :disabled="disabled"
  >
    <slot />
  </button>
</template>

<style scoped>
.base-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.6rem 1.2rem;
  border-radius: var(--border-radius);
  font-weight: 500;
  transition: var(--transition-default);
  cursor: pointer;
}

.btn-primary {
  background-color: var(--color-primary);
  color: white;
  border: 1px solid var(--color-primary);
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--color-primary-dark);
  border-color: var(--color-primary-dark);
}

.btn-outline {
  background-color: transparent;
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
}

.btn-outline:hover:not(:disabled) {
  background-color: var(--color-primary);
  color: white;
}

.btn-text {
  background-color: transparent;
  border: none;
  color: var(--color-primary);
  padding: 0.4rem 0.8rem;
}

.btn-text:hover:not(:disabled) {
  background-color: rgba(30, 144, 255, 0.1);
}

.btn-danger {
  background-color: #f44336;
  color: white;
  border: 1px solid #f44336;
}

.btn-danger:hover:not(:disabled) {
  background-color: #d32f2f;
  border-color: #d32f2f;
}

.base-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
