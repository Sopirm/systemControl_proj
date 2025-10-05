<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseInput from '../components/BaseInput.vue'
import BaseButton from '../components/BaseButton.vue'
import { authService } from '../services/api'
import { useAuth } from '../composables/useAuth'

// Создаем событие для оповещения о изменении состояния авторизации
const authEvent = new CustomEvent('auth-change', { detail: { action: 'login' } })

const router = useRouter()
const username = ref('')
const password = ref('')
const error = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  error.value = ''
  
  if (!username.value || !password.value) {
    error.value = 'Пожалуйста, заполните все поля'
    return
  }
  
  isLoading.value = true
  
  try {
    // Авторизация пользователя через API сервис
    const response = await authService.login({
      username: username.value,
      password: password.value
    })
    
    console.log('Успешная авторизация:', response.user)
    
    // Обновляем состояние авторизации
    window.dispatchEvent(authEvent)
    
    // После успешной авторизации перенаправляем на страницу проектов
    await nextTick()
    router.push('/projects')
  } catch (err) {
    console.error('Ошибка авторизации:', err)
    error.value = err instanceof Error ? err.message : 'Ошибка при входе в систему'
    isLoading.value = false
  }
}
</script>

<template>
  <div class="login-page">
    <BaseCard title="Вход в систему" class="login-card">
      <form @submit.prevent="handleLogin">
        <div v-if="error" class="alert alert-error">{{ error }}</div>

        <BaseInput
          v-model="username"
          label="Имя пользователя или Email"
          placeholder="Введите имя пользователя или email"
          required
        />

        <BaseInput
          v-model="password"
          label="Пароль"
          type="password"
          placeholder="Введите ваш пароль"
          required
        />

        <div class="form-actions">
          <BaseButton type="submit" :disabled="isLoading">
            {{ isLoading ? 'Вход...' : 'Войти' }}
          </BaseButton>
        </div>

        <div class="login-footer">
          <p>
            Еще нет аккаунта?
            <RouterLink to="/register">Зарегистрироваться</RouterLink>
          </p>
        </div>
      </form>
    </BaseCard>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
}

.login-card {
  width: 100%;
  max-width: 420px;
}

.form-actions {
  margin-top: 1.5rem;
}

.login-footer {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.9rem;
}

.login-footer a {
  color: var(--color-primary);
  font-weight: 500;
}
</style>
