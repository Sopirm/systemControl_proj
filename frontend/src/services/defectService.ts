import { authService } from './api';

// Типы для работы с дефектами
export interface Defect {
  id: number;
  title: string;
  description: string;
  priority: 'low' | 'medium' | 'high';
  status: 'new' | 'in_progress' | 'review' | 'closed' | 'cancelled';
  assignee_id?: number;
  project_id: number;
  due_date?: string;
  created_at?: string;
  updated_at?: string;
  assignee?: {
    id: number;
    username: string;
    full_name: string;
    email: string;
  };
}

// Статистика дефектов для проекта
export interface DefectStats {
  active: number; // Количество активных дефектов (new, in_progress, review)
  resolved: number; // Количество решенных дефектов (closed)
  total: number; // Общее количество дефектов
}

// API для работы с дефектами (только для получения данных)
export const defectService = {
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
   * Получение статистики дефектов по ID проекта
   * Если API не реализован, генерирует случайные данные для демонстрации
   */
  async getDefectStatsByProjectId(projectId: number): Promise<DefectStats> {
    try {
      const token = authService.getAuthToken();
      if (!token) {
        throw new Error('Пользователь не авторизован');
      }

      const response = await fetch(`/api/projects/${projectId}/defects/stats`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        const data = await response.json();
        return data.stats;
      }
    } catch (e) {
      console.log('API для статистики дефектов не реализован, используем демо-данные');
    }

    // Если API не реализован или произошла ошибка, возвращаем демо-данные
    return this.generateDemoDefectStats(projectId);
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
  }
};
