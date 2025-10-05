<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import BaseCard from '../components/BaseCard.vue'
import BaseInput from '../components/BaseInput.vue'
import BaseButton from '../components/BaseButton.vue'
import BaseSelect from '../components/BaseSelect.vue'
import { authService } from '../services/api'

const router = useRouter()
const fullName = ref('')
const username = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const role = ref('')
const error = ref('')
const isLoading = ref(false)

const roleOptions = [
  { value: 'manager', label: 'Менеджер' },
  { value: 'engineer', label: 'Инженер' },
  { value: 'observer', label: 'Наблюдатель' }
]

const handleRegister = async () => {
  error.value = ''
  
  if (!fullName.value || !username.value || !email.value || !password.value || !confirmPassword.value || !role.value) {
    error.value = 'Пожалуйста, заполните все поля'
    return
  }
  
  if (password.value !== confirmPassword.value) {
    error.value = 'Пароли не совпадают'
    return
  }
  
  isLoading.value = true
  
  try {
    // Регистрация пользователя через API сервис
    const response = await authService.register({
      username: username.value,
      full_name: fullName.value,
      email: email.value,
      password: password.value,
      role: role.value
    })
    
    console.log('Пользователь успешно зарегистрирован:', response.user)
    
    // После успешной регистрации перенаправляем на страницу входа
    router.push('/login')
  } catch (err) {
    console.error('Ошибка регистрации:', err)
    error.value = err instanceof Error ? err.message : 'Ошибка при регистрации'
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
          v-model="username"
          label="Имя пользователя"
          placeholder="Введите имя пользователя"
          required
        />
        
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
        
        <BaseSelect
          v-model="role"
          label="Роль в системе"
          :options="roleOptions"
          placeholder="Выберите вашу роль"
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
