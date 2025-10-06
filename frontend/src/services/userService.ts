import { authService } from './api';

// Тип для пользователя
export interface User {
  id: number;
  username: string;
  email: string;
  full_name: string;
  role: string;
  created_at?: string;
}

// Тип для обновления роли пользователя
export interface UserRoleUpdate {
  role: string;
}

// API для работы с пользователями
export const userService = {
  /**
   * Получение списка всех пользователей
   * Доступно только для менеджеров
   */
  async getAllUsers(): Promise<User[]> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch('/api/users', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка получения пользователей:', data);
      throw new Error(data?.error || `Ошибка получения пользователей (${response.status})`);
    }

    return data.users;
  },

  /**
   * Получение списка инженеров
   * Доступно менеджерам и инженерам
   */
  async getEngineers(): Promise<User[]> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch('/api/users/engineers', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка получения инженеров:', data);
      throw new Error(data?.error || `Ошибка получения инженеров (${response.status})`);
    }

    return data.users as User[];
  },

  /**
   * Обновление роли пользователя
   * Доступно только для менеджеров
   */
  async updateUserRole(userId: number, role: string): Promise<User> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/users/${userId}/role`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ role })
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка обновления роли:', data);
      throw new Error(data?.error || `Ошибка обновления роли пользователя (${response.status})`);
    }

    return data.user;
  }
};
