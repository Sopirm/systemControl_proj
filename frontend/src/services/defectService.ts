import { authService } from './api';
import type { User } from './userService';

// Типы для работы с дефектами
export interface Defect {
  id: number;
  title: string;
  description: string;
  priority: 'low' | 'medium' | 'high';
  status: 'new' | 'in_progress' | 'review' | 'closed' | 'cancelled';
  assignee_id?: number;
  reporter_id?: number;
  project_id: number;
  due_date?: string;
  created_at?: string;
  updated_at?: string;
  assignee?: User;
  reporter?: User;
  project?: {
    id: number;
    name: string;
  };
}

// Интерфейс для создания дефекта
export interface DefectCreate {
  title: string;
  description: string;
  project_id: number;
  priority?: 'low' | 'medium' | 'high';
  assignee_id?: number;
  due_date?: string;
}

// Интерфейс для обновления дефекта
export interface DefectUpdate {
  title?: string;
  description?: string;
  status?: 'new' | 'in_progress' | 'review' | 'closed' | 'cancelled';
  priority?: 'low' | 'medium' | 'high';
  assignee_id?: number;
  due_date?: string;
}

// Статистика дефектов для проекта
export interface DefectStats {
  active: number; // Количество активных дефектов (new, in_progress, review)
  resolved: number; // Количество решенных дефектов (closed)
  total: number; // Общее количество дефектов
}

// API для работы с дефектами
export const defectService = {
  /**
   * Получение всех дефектов
   */
  async getAllDefects(): Promise<Defect[]> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch('/api/defects', {
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
      console.error('Ошибка получения дефектов:', data);
      throw new Error(data?.error || `Ошибка получения дефектов (${response.status})`);
    }

    return data.defects;
  },

  /**
   * Получение дефектов по ID проекта
   */
  async getDefectsByProjectId(projectId: number): Promise<Defect[]> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/projects/${projectId}/defects`, {
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
      console.error('Ошибка получения дефектов:', data);
      throw new Error(data?.error || `Ошибка получения дефектов (${response.status})`);
    }

    return data.defects;
  },

  /**
   * Получение дефекта по ID
   */
  async getDefectById(defectId: number): Promise<Defect> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/defects/${defectId}`, {
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
      console.error('Ошибка получения дефекта:', data);
      throw new Error(data?.error || `Дефект не найден (${response.status})`);
    }

    return data.defect;
  },

  /**
   * Создание нового дефекта
   */
  async createDefect(defect: DefectCreate): Promise<Defect> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    // Форматирование даты, если она есть
    const formattedDefect = { ...defect };
    if (formattedDefect.due_date) {
      const date = new Date(formattedDefect.due_date);
      formattedDefect.due_date = date.toISOString();
    }

    const response = await fetch('/api/defects', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(formattedDefect)
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка создания дефекта:', data);
      throw new Error(data?.error || `Ошибка создания дефекта (${response.status})`);
    }

    return data.defect;
  },

  /**
   * Обновление существующего дефекта
   */
  async updateDefect(defectId: number, defect: DefectUpdate): Promise<Defect> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    // Форматирование даты, если она есть
    const formattedDefect = { ...defect };
    if (formattedDefect.due_date) {
      const date = new Date(formattedDefect.due_date);
      formattedDefect.due_date = date.toISOString();
    }

    const response = await fetch(`/api/defects/${defectId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(formattedDefect)
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка обновления дефекта:', data);
      throw new Error(data?.error || `Ошибка обновления дефекта (${response.status})`);
    }

    return data.defect;
  },

  /**
   * Удаление дефекта
   */
  async deleteDefect(defectId: number): Promise<void> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/defects/${defectId}`, {
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
        throw new Error(`Ошибка удаления дефекта (${response.status})`);
      }
      throw new Error(errorData?.error || `Ошибка удаления дефекта (${response.status})`);
    }
  },

  /**
   * Получение статистики дефектов по ID проекта
   */
  async getDefectStatsByProjectId(projectId: number): Promise<DefectStats> {
    // Сначала пробуем получить статистику через специальный API-эндпоинт
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    try {
      const response = await fetch(`/api/projects/${projectId}/defects/stats`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        const data = await response.json();
        console.log(`Получена статистика дефектов для проекта ${projectId} через API:`, data.stats);
        return data.stats;
      }
    } catch (e) {
      console.log(`API для статистики дефектов проекта ${projectId} недоступен, будем считать вручную`);
    }
    
    // Если специальный эндпоинт не реализован, считаем статистику вручную
    const defects = await this.getDefectsByProjectId(projectId);
    console.log(`Получено ${defects.length} дефектов для проекта ${projectId}, считаем статистику`);
    
    // Подсчет активных и решенных дефектов
    const active = defects.filter(d => 
      d.status === 'new' || d.status === 'in_progress' || d.status === 'review'
    ).length;
    
    const resolved = defects.filter(d => 
      d.status === 'closed'
    ).length;
    
    const total = defects.length;
    
    return { active, resolved, total };
  },

  /**
   * Генерация демонстрационных данных статистики дефектов
   * Используется только если API не реализован
   */
  generateDemoDefectStats(projectId: number): DefectStats {
    // Генерируем псевдослучайные числа на основе projectId для стабильности
    const seed = projectId * 13 % 100;
    const active = Math.floor(seed / 10) + 1;  // 1-10
    const resolved = Math.floor(seed / 5) % 10 + 1;  // 1-10
    const total = active + resolved;

    return {
      active,
      resolved,
      total
    };
  },

  /**
   * Проверка прав доступа для различных операций с дефектами
   */
  checkPermissions() {
    const currentUser = authService.getCurrentUser();
    
    if (!currentUser) {
      return {
        canCreate: false,
        canEdit: false,
        canDelete: false,
        canChangeStatus: false,
        canCloseOrCancel: false
      };
    }

    const isManager = currentUser.role === 'manager';
    const isEngineer = currentUser.role === 'engineer';
    
    return {
      canCreate: isManager || isEngineer,
      canEdit: isManager || isEngineer,
      canDelete: isManager,
      canChangeStatus: isManager || isEngineer,
      canCloseOrCancel: isManager
    };
  },

  /**
   * Проверка, может ли пользователь изменить статус дефекта на указанный
   */
  canChangeStatusTo(status: string): boolean {
    const currentUser = authService.getCurrentUser();
    
    if (!currentUser) return false;
    
    const isManager = currentUser.role === 'manager';
    
    // Менеджер может изменить на любой статус
    if (isManager) return true;
    
    // Инженер не может изменить на 'closed' или 'cancelled'
    if (status === 'closed' || status === 'cancelled') {
      return false;
    }
    
    // Для остальных статусов инженер может менять
    return currentUser.role === 'engineer';
  }
};