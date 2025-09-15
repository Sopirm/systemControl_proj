<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: string | number
  label?: string
  type?: string
  placeholder?: string
  required?: boolean
  disabled?: boolean
  error?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const inputType = computed(() => props.type || 'text')
</script>

<template>
  <div class="form-group">
    <label v-if="label" class="form-label">
      {{ label }}
      <span v-if="required" class="required">*</span>
    </label>
    <input
      :type="inputType"
      class="form-input"
      :class="{ 'has-error': error }"
      :value="modelValue"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
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

.form-input {
  width: 100%;
  padding: 0.6rem 0.8rem;
  border: 1px solid var(--color-border);
  border-radius: var(--border-radius);
  transition: var(--transition-default);
  font-size: 1rem;
}

.form-input:focus {
  border-color: var(--color-primary);
  outline: none;
  box-shadow: 0 0 0 3px rgba(30, 144, 255, 0.2);
}

.form-input.has-error {
  border-color: var(--color-error);
}

.error-message {
  color: var(--color-error);
  font-size: 0.85rem;
  margin-top: 0.25rem;
}

.form-input:disabled {
  background-color: #f8f9fa;
  cursor: not-allowed;
  opacity: 0.7;
}
</style>
