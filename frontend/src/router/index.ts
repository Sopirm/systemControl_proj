import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

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
      meta: { requiresAuth: true }
    },
  ],
})

// Navigation guard для проверки авторизации
// На данном этапе просто заглушка, потом будет реализована настоящая проверка
// router.beforeEach((to, from, next) => {
//   const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
//   const isAuthenticated = false // Заглушка, пока нет реальной авторизации

//   if (requiresAuth && !isAuthenticated) {
//     next('/login')
//   } else {
//     next()
//   }
// })

export default router