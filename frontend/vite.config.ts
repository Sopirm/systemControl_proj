import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    proxy: {
      // Настройка прокси для API-запросов
      '/api': {
        target: 'http://localhost:8080', // Порт, на котором работает Go-сервер
        changeOrigin: true,
        secure: false,
      },
      // Прокси для аутентификации
      '/auth': {
        target: 'http://localhost:8080', 
        changeOrigin: true,
        secure: false,
      },
      // Прокси для отладочных запросов
      '/debug': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false,
      }
    }
  }
})
