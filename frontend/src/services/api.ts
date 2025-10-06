/**
 * API service for all backend requests
 */

// API base URL - using relative path to work with any deployment setup
const API_BASE_URL = '/api';

// Types for authentication
export interface UserRegistration {
  username: string;
  full_name: string;
  email: string;
  password: string;
  role: string;
}

export interface UserLogin {
  username: string;
  password: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  full_name?: string;
  role: string;
}

export interface AuthResponse {
  message: string;
  token: string;
  user: User;
}

/**
 * Authentication service
 */
export const authService = {
  /**
   * Register a new user
   */
  async register(userData: UserRegistration): Promise<AuthResponse> {
    const response = await fetch(`/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(userData)
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка регистрации:', data);
      throw new Error(data?.error || `Ошибка при регистрации пользователя (${response.status})`);
    }

    return data;
  },

  /**
   * Login a user
   */
  async login(credentials: UserLogin): Promise<AuthResponse> {
    const response = await fetch(`/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(credentials)
    });

    let data;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка авторизации:', data);
      throw new Error(data?.error || `Неверное имя пользователя или пароль (${response.status})`);
    }

    // Store token and user data in localStorage
    localStorage.setItem('authToken', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));

    return data;
  },

  /**
   * Logout the current user
   */
  logout(): void {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
  },

  /**
   * Check if user is logged in
   */
  isLoggedIn(): boolean {
    return !!localStorage.getItem('authToken');
  },

  /**
   * Get the current user
   */
  getCurrentUser(): User | null {
    const userJson = localStorage.getItem('user');
    if (!userJson) return null;
    
    try {
      return JSON.parse(userJson) as User;
    } catch (error) {
      console.error('Failed to parse user data:', error);
      return null;
    }
  },

  /**
   * Get the auth token
   */
  getAuthToken(): string | null {
    return localStorage.getItem('authToken');
  },

  /**
   * Fetch current user profile from backend and update local cache
   */
  async getProfile(): Promise<User> {
    const token = this.getAuthToken();
    if (!token) throw new Error('Пользователь не авторизован');

    const response = await fetch('/api/profile', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    });

    let data: any;
    try {
      data = await response.json();
    } catch (e) {
      console.error('Ошибка при парсинге ответа профиля:', e);
      throw new Error('Ошибка соединения с сервером');
    }

    if (!response.ok) {
      console.error('Ошибка получения профиля:', data);
      throw new Error(data?.error || `Ошибка получения профиля (${response.status})`);
    }

    const user = data.user as User;
    // Обновляем локальный кэш, чтобы во всех местах появились актуальные поля (в т.ч. full_name)
    localStorage.setItem('user', JSON.stringify(user));
    return user;
  }
};
