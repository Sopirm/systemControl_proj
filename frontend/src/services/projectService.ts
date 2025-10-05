import { authService } from './api';

// Типы для работы с проектами
export interface Project {
  id: number;
  name: string;
  description: string;
  location: string;
  start_date: string;
  end_date: string;
  status: string;
  manager_id: number;
  manager?: {
    id: number;
    username: string;
    full_name: string;
    email: string;
  };
  created_at?: string;
  updated_at?: string;
}

export interface ProjectCreate {
  name: string;
  description: string;
  location: string;
  start_date: string;
  end_date: string;
  status?: string;
  manager_id: number;
}

export interface ProjectUpdate {
  name?: string;
  description?: string;
  location?: string;
  start_date?: string;
  end_date?: string;
  status?: string;
  manager_id?: number;
}

// API для работы с проектами
export const projectService = {
  /**
   * Получение списка всех проектов
   */
  async getAllProjects(): Promise<Project[]> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch('/api/projects', {
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
      console.error('Ошибка получения проектов:', data);
      throw new Error(data?.error || `Ошибка получения проектов (${response.status})`);
    }

    return data.projects;
  },

  /**
   * Получение проекта по ID
   */
  async getProjectById(id: number): Promise<Project> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/projects/${id}`, {
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
      console.error('Ошибка получения проекта:', data);
      throw new Error(data?.error || `Проект не найден (${response.status})`);
    }

    return data.project;
  },

  /**
   * Создание нового проекта
   */
  async createProject(project: ProjectCreate): Promise<Project> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch('/api/projects', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(project)
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка создания проекта:', data);
      throw new Error(data?.error || `Ошибка создания проекта (${response.status})`);
    }

    return data.project;
  },

  /**
   * Обновление существующего проекта
   */
  async updateProject(id: number, project: ProjectUpdate): Promise<Project> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/projects/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(project)
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка обновления проекта:', data);
      throw new Error(data?.error || `Ошибка обновления проекта (${response.status})`);
    }

    return data.project;
  },

  /**
   * Удаление проекта
   */
  async deleteProject(id: number): Promise<void> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/projects/${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      let errorData;
      try {
        errorData = await response.json();
      } catch (e) {
        throw new Error(`Ошибка удаления проекта (${response.status})`);
      }
      throw new Error(errorData?.error || `Ошибка удаления проекта (${response.status})`);
    }
  }
};
