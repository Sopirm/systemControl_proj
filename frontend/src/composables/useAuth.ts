import { ref, computed } from 'vue';
import { authService } from '../services/api';
import type { User } from '../services/api';

/**
 * Composable for handling authentication state and user roles
 */
export function useAuth() {
  const currentUser = ref<User | null>(authService.getCurrentUser());
  
  // Is the user authenticated?
  const isAuthenticated = computed(() => !!currentUser.value);
  
  // Role-based access control helpers
  const isManager = computed(() => currentUser.value?.role === 'manager');
  const isEngineer = computed(() => currentUser.value?.role === 'engineer');
  const isObserver = computed(() => currentUser.value?.role === 'observer');
  
  // Check if user has a specific role
  const hasRole = (role: string) => currentUser.value?.role === role;
  
  // Logout function
  const logout = () => {
    authService.logout();
    currentUser.value = null;
  };
  
  // Update user data in local state
  const updateUser = (user: User) => {
    currentUser.value = user;
  };
  
  return {
    currentUser,
    isAuthenticated,
    isManager,
    isEngineer,
    isObserver,
    hasRole,
    logout,
    updateUser,
    authService // Экспортируем сервис для прямого доступа
  };
}
