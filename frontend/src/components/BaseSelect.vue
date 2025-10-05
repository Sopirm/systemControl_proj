<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: string
  label?: string
  options: { value: string, label: string }[]
  placeholder?: string
  required?: boolean
  disabled?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}>()
</script>

<template>
  <div class="form-group">
    <label v-if="label" class="form-label">
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>
    <select
      class="form-select"
      :class="{ 'has-error': error }"
      :value="modelValue"
      :required="required"
      :disabled="disabled"
      @change="(e) => {
        const value = (e.target as HTMLSelectElement).value;
        emit('update:modelValue', value);
        emit('change', value);
      }"
    >
      <option value="" disabled selected v-if="placeholder">{{ placeholder }}</option>
      <option v-for="option in options" :key="option.value" :value="option.value">
        {{ option.label }}
      </option>
    </select>
    <p v-if="error" class="error-message">{{ error }}</p>
  </div>
</template>

<style scoped>
.form-group {
  margin-bottom: 1rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--color-text);
}

.required {
  color: var(--color-error);
}

.form-select {
  width: 100%;
  padding: 0.6rem 0.8rem;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  transition: var(--transition-default);
  font-size: 1rem;
  background-color: #fff;
  height: 41px;
}

.form-select:focus {
  border-color: var(--color-primary);
  outline: none;
  box-shadow: 0 0 0 3px rgba(30, 144, 255, 0.2);
}

.form-select.has-error {
  border-color: var(--color-error);
}

.error-message {
  color: var(--color-error);
  font-size: 0.85rem;
  margin-top: 0.25rem;
}

.form-select:disabled {
  background-color: #f8f9fa;
  cursor: not-allowed;
  opacity: 0.7;
}
</style>
