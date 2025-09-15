<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseInput from '../components/BaseInput.vue'
import BaseButton from '../components/BaseButton.vue'

const router = useRouter()
const email = ref('')
const password = ref('')
const error = ref('')
const isLoading = ref(false)

const handleLogin = async () => {
  error.value = ''
  
  if (!email.value || !password.value) {
    error.value = 'Пожалуйста, заполните все поля'
    return
  }
  
  isLoading.value = true
  
  try {
    // Здесь будет реальный запрос к API
    // Временная заглушка для демонстрации
    setTimeout(() => {
      console.log('Попытка входа:', { email: email.value, password: '***' })
      
      // После успешной авторизации перенаправляем на страницу проектов
      router.push('/projects')
      
      isLoading.value = false
    }, 1000)
  } catch (err) {
    error.value = 'Ошибка при входе в систему'
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
          v-model="email"
          label="Email"
          type="email"
          placeholder="Введите ваш email"
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
