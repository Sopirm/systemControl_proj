<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseCard from '../components/BaseCard.vue'
import BaseInput from '../components/BaseInput.vue'
import BaseSelect from '../components/BaseSelect.vue'
import BaseButton from '../components/BaseButton.vue'
import { userService, User } from '../services/userService'
import { useAuth } from '../composables/useAuth'

// Используем интерфейс User из userService.ts

const users = ref<User[]>([])
const isLoading = ref(true)
const error = ref('')
const { currentUser } = useAuth()

const roleOptions = [
  { value: 'manager', label: 'Менеджер' },
  { value: 'engineer', label: 'Инженер' },
  { value: 'observer', label: 'Наблюдатель' }
]

// Загрузка списка пользователей с сервера
onMounted(async () => {
  try {
    isLoading.value = true
    users.value = await userService.getAllUsers()
    isLoading.value = false
  } catch (err) {
    console.error('Ошибка при загрузке пользователей:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить список пользователей'
    isLoading.value = false
  }
})

// Обновление роли пользователя
const isUpdating = ref(false)
const updateError = ref('')

const updateUserRole = async (userId: number, newRole: string) => {
  try {
    // Запрещаем менять роль у самого себя на клиенте
    if (currentUser.value && currentUser.value.id === userId) {
      updateError.value = 'Нельзя менять роль у самого себя'
      return
    }
    isUpdating.value = true
    updateError.value = ''

    // Отправляем запрос на обновление роли
    const updatedUser = await userService.updateUserRole(userId, newRole)
    
    // Обновляем пользователя в списке
    const index = users.value.findIndex(u => u.id === userId)
    if (index !== -1) {
      users.value[index] = {
        ...users.value[index],
        role: updatedUser.role
      }
    }
    
    console.log(`Роль пользователя ${updatedUser.username} изменена на ${updatedUser.role}`)
  } catch (err) {
    console.error('Ошибка при обновлении роли:', err)
    updateError.value = err instanceof Error ? err.message : 'Ошибка при обновлении роли пользователя'
    
    // Отменяем изменение роли в UI при ошибке
    const user = users.value.find(u => u.id === userId)
    if (user) {
      // Здесь мы должны вернуть предыдущее значение роли
      // Но поскольку мы его не сохраняли, придется перезагрузить список
      loadUsers()
    }
  } finally {
    isUpdating.value = false
  }
}

// Функция для загрузки пользователей
const loadUsers = async () => {
  try {
    isLoading.value = true
    users.value = await userService.getAllUsers()
    error.value = ''
  } catch (err) {
    console.error('Ошибка при загрузке пользователей:', err)
    error.value = err instanceof Error ? err.message : 'Не удалось загрузить список пользователей'
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <div class="users-page container">
    <h1>Управление пользователями</h1>

    <!-- Личный аккаунт -->
    <BaseCard v-if="currentUser" title="Личный аккаунт" class="self-card">
      <div class="self-grid">
        <div><strong>Имя пользователя:</strong> {{ currentUser.username }}</div>
        <div><strong>ФИО:</strong> {{ currentUser.full_name || '—' }}</div>
        <div><strong>Email:</strong> {{ currentUser.email }}</div>
        <div><strong>Роль:</strong> {{ currentUser.role }}</div>
      </div>
      <small class="hint">Собственную роль изменить нельзя</small>
    </BaseCard>
    
    <div v-if="error" class="alert alert-error">{{ error }}</div>
    
    <BaseCard v-if="isLoading" class="loading-card">
      <p>Загрузка списка пользователей...</p>
    </BaseCard>
    
    <BaseCard v-else title="Список пользователей" class="users-table-card">
      <div class="users-table-container">
        <table class="users-table">
          <thead>
            <tr>
              <th>Имя пользователя</th>
              <th>Email</th>
              <th>ФИО</th>
              <th>Роль</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>{{ user.username }}</td>
              <td>{{ user.email }}</td>
              <td>{{ user.full_name }}</td>
              <td>
                <BaseSelect
                  v-model="user.role"
                  :options="roleOptions"
                  :disabled="isUpdating || (currentUser && user.id === currentUser.id)"
                  @change="updateUserRole(user.id, user.role)"
                />
              </td>
              <td class="actions-cell">
                <div v-if="updateError && updateError.includes(String(user.id))" class="error-message">
                  {{ updateError }}
                </div>
                <BaseButton 
                  size="small" 
                  :disabled="isUpdating" 
                  @click="loadUsers"
                >
                  Обновить
                </BaseButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </BaseCard>
  </div>
</template>

<style scoped>
.users-page {
  padding: 2rem 0;
}

.loading-card {
  text-align: center;
  padding: 2rem;
}

.users-table-container {
  overflow-x: auto;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th, 
.users-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid var(--color-border);
}

.users-table th {
  font-weight: 600;
  background-color: var(--color-background-soft);
}

.actions-cell {
  white-space: nowrap;
  text-align: right;
}

@media (max-width: 768px) {
  .users-table th, 
  .users-table td {
    padding: 0.5rem;
  }
}
</style>
