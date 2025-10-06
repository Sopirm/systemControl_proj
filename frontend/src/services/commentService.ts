import { authService } from './api';
import type { User } from './userService';

// Типы для работы с комментариями
export interface Comment {
  id: number;
  defect_id: number;
  user_id: number;
  user?: User;
  content: string;
  created_at: string;
}

// Интерфейс для создания комментария
export interface CommentCreate {
  defect_id: number;
  content: string;
}

// API для работы с комментариями
export const commentService = {
  /**
   * Получение всех комментариев для конкретного дефекта
   */
  async getDefectComments(defectId: number): Promise<Comment[]> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/defects/${defectId}/comments`, {
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
      console.error('Ошибка получения комментариев:', data);
      throw new Error(data?.error || `Ошибка получения комментариев (${response.status})`);
    }

    return data.comments || [];
  },

  /**
   * Создание нового комментария к дефекту
   */
  async createComment(comment: CommentCreate): Promise<Comment> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch('/api/defects/comments', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(comment)
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка создания комментария:', data);
      throw new Error(data?.error || `Ошибка создания комментария (${response.status})`);
    }

    return data.comment;
  },

  /**
   * Удаление комментария
   */
  async deleteComment(commentId: number): Promise<void> {
    const token = authService.getAuthToken();
    if (!token) {
      throw new Error('Пользователь не авторизован');
    }

    const response = await fetch(`/api/defects/comments/${commentId}`, {
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
        throw new Error(`Ошибка удаления комментария (${response.status})`);
      }
      throw new Error(errorData?.error || `Ошибка удаления комментария (${response.status})`);
    }
  },

  /**
   * Проверка, может ли текущий пользователь удалить комментарий
   */
  canDeleteComment(comment: Comment): boolean {
    const currentUser = authService.getCurrentUser();
    
    if (!currentUser) return false;
    
    // Менеджер может удалить любой комментарий
    if (currentUser.role === 'manager') return true;
    
    // Автор может удалить свой комментарий
    return comment.user_id === currentUser.id;
  }
};
