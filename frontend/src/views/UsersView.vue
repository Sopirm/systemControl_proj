<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseCard from '../components/BaseCard.vue'
import BaseInput from '../components/BaseInput.vue'
import BaseSelect from '../components/BaseSelect.vue'
import BaseButton from '../components/BaseButton.vue'

interface User {
  id: number
  username: string
  email: string
  fullName: string
  role: string
}

const users = ref<User[]>([])
const isLoading = ref(true)
const error = ref('')

const roleOptions = [
  { value: 'manager', label: 'Менеджер' },
  { value: 'engineer', label: 'Инженер' },
  { value: 'observer', label: 'Наблюдатель' }
]

// Имитация загрузки пользователей с сервера
onMounted(async () => {
  try {
    // В дальнейшем заменить на реальный API запрос
    setTimeout(() => {
      users.value = [
        { 
          id: 1, 
          username: 'manager1', 
          email: 'manager@example.com', 
          fullName: 'Иванов Иван Иванович', 
          role: 'manager' 
        },
        { 
          id: 2, 
          username: 'engineer1', 
          email: 'engineer@example.com', 
          fullName: 'Петров Петр Петрович', 
          role: 'engineer' 
        },
        { 
          id: 3, 
          username: 'observer1', 
          email: 'observer@example.com', 
          fullName: 'Сидоров Сидор Сидорович', 
          role: 'observer' 
        }
      ]
      isLoading.value = false
    }, 1000)
  } catch (err) {
    console.error('Ошибка при загрузке пользователей:', err)
    error.value = 'Не удалось загрузить список пользователей'
    isLoading.value = false
  }
})

const updateUserRole = (userId: number, newRole: string) => {
  // В дальнейшем заменить на реальный API запрос
  const user = users.value.find(u => u.id === userId)
  if (user) {
    user.role = newRole
    console.log(`Роль пользователя ${user.username} изменена на ${newRole}`)
  }
}
</script>

<template>
  <div class="users-page container">
    <h1>Управление пользователями</h1>
    
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
              <td>{{ user.fullName }}</td>
              <td>
                <BaseSelect
                  v-model="user.role"
                  :options="roleOptions"
                  @change="updateUserRole(user.id, user.role)"
                />
              </td>
              <td class="actions-cell">
                <BaseButton size="small">Редактировать</BaseButton>
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
