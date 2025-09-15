<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseInput from '../components/BaseInput.vue'
import BaseButton from '../components/BaseButton.vue'

const router = useRouter()
const fullName = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const error = ref('')
const isLoading = ref(false)

const handleRegister = async () => {
  error.value = ''
  
  if (!fullName.value || !email.value || !password.value || !confirmPassword.value) {
    error.value = 'Пожалуйста, заполните все поля'
    return
  }
  
  if (password.value !== confirmPassword.value) {
    error.value = 'Пароли не совпадают'
    return
  }
  
  isLoading.value = true
  
  try {
    // Здесь будет реальный запрос к API
    // Временная заглушка для демонстрации
    setTimeout(() => {
      console.log('Регистрация пользователя:', { 
        fullName: fullName.value,
        email: email.value
      })
      
      // После успешной регистрации перенаправляем на страницу входа
      router.push('/login')
      
      isLoading.value = false
    }, 1000)
  } catch (err) {
    error.value = 'Ошибка при регистрации'
    isLoading.value = false
  }
}
</script>

<template>
  <div class="register-page">
    <BaseCard title="Регистрация" class="register-card">
      <form @submit.prevent="handleRegister">
        <div v-if="error" class="alert alert-error">{{ error }}</div>

        <BaseInput
          v-model="fullName"
          label="ФИО"
          placeholder="Введите ваше полное имя"
          required
        />

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
          placeholder="Создайте пароль"
          required
        />

        <BaseInput
          v-model="confirmPassword"
          label="Подтверждение пароля"
          type="password"
          placeholder="Повторите пароль"
          required
        />

        <div class="form-actions">
          <BaseButton type="submit" :disabled="isLoading">
            {{ isLoading ? 'Регистрация...' : 'Зарегистрироваться' }}
          </BaseButton>
        </div>

        <div class="register-footer">
          <p>
            Уже есть аккаунт?
            <RouterLink to="/login">Войти</RouterLink>
          </p>
        </div>
      </form>
    </BaseCard>
  </div>
</template>

<style scoped>
.register-page {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 2rem 0;
}

.register-card {
  width: 100%;
  max-width: 420px;
}

.form-actions {
  margin-top: 1.5rem;
}

.register-footer {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.9rem;
}

.register-footer a {
  color: var(--color-primary);
  font-weight: 500;
}
</style>
