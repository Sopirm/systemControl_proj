import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { authService } from '../services/api'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
    },
    {
      path: '/projects',
      name: 'projects',
      component: () => import('../views/ProjectsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/projects/create',
      name: 'project-create',
      component: () => import('../views/ProjectCreateView.vue'),
      meta: { requiresAuth: true, requiredRole: 'manager' }
    },
    {
      path: '/projects/:id/edit',
      name: 'project-edit',
      component: () => import('../views/ProjectEditView.vue'),
      meta: { requiresAuth: true, requiredRole: 'manager' }
    },
    {
      path: '/projects/:id',
      name: 'project-detail',
      component: () => import('../views/ProjectDetailView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/defects',
      name: 'defects',
      component: () => import('../views/DefectsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/defects/:id',
      name: 'defect-detail',
      component: () => import('../views/DefectDetailView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/reports',
      name: 'reports',
      component: () => import('../views/ReportsView.vue'),
      meta: { requiresAuth: true, requiredRoles: ['manager', 'observer'] }
    },
    {
      path: '/users',
      name: 'users',
      component: () => import('../views/UsersView.vue'),
      meta: { requiresAuth: true, requiredRole: 'manager' }
    },
  ],
})

// Navigation guard для проверки авторизации и ролей
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  // Поддержка как requiredRole (string), так и requiredRoles (string[])
  let allowedRoles: string[] | null = null
  const metaWithRoles = to.matched.find(r => (r.meta as any).requiredRoles || (r.meta as any).requiredRole)?.meta as any
  if (metaWithRoles?.requiredRoles) {
    allowedRoles = metaWithRoles.requiredRoles as string[]
  } else if (metaWithRoles?.requiredRole) {
    allowedRoles = [metaWithRoles.requiredRole as string]
  }

  const isAuthenticated = authService.isLoggedIn()
  const currentUser = authService.getCurrentUser()

  // Если пользователь уже авторизован и пытается перейти на страницы входа или регистрации
  if (isAuthenticated && (to.path === '/login' || to.path === '/register' || to.path === '/')) {
    console.log('Пользователь авторизован, перенаправляем на /projects')
    next('/projects')
    return
  }

  // Проверка авторизации
  if (requiresAuth && !isAuthenticated) {
    next('/login')
    return
  }
  
  // Проверка роли пользователя
  if (allowedRoles && currentUser && !allowedRoles.includes(currentUser.role)) {
    // Если роль не разрешена
    next('/projects')
    return
  }
  
  next()
})

export default router