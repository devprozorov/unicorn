// app/plugins/axios.client.ts
import axios from 'axios'
import { useAuthStore } from '~/stores/auth'
import { useAuthApi } from '~/services/authApi'

// Защита от race condition
let isRefreshing = false
let refreshSubscribers: Array<(token: string | null) => void> = []

function subscribeTokenRefresh(cb: (token: string | null) => void) {
  refreshSubscribers.push(cb)
}

function onRefreshed(token: string | null) {
  refreshSubscribers.forEach(cb => cb(token))
  refreshSubscribers = []
}

export default defineNuxtPlugin(() => {
  const auth = useAuthStore()
  const api = useAuthApi()

  axios.interceptors.request.use((config) => {
    if (auth.accessToken) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${auth.accessToken}`
    }
    return config
  })

  axios.interceptors.response.use(
    (response) => response,
    async (error) => {
      const originalRequest = error.config

      // Не 401 или уже повторный запрос
      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(error)
      }

      // Отмечаем, что уже пытались обновить
      originalRequest._retry = true

      // Если уже идет refresh, ждем его завершения
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          subscribeTokenRefresh((token) => {
            if (token) {
              originalRequest.headers.Authorization = `Bearer ${token}`
              resolve(axios(originalRequest))
            } else {
              reject(error)
            }
          })
        })
      }

      isRefreshing = true
      console.warn('[axios] 401, trying refresh')

      try {
        const res = await api.refresh()

        if (res?.accessToken) {
          auth.setAccessToken(res.accessToken)
          originalRequest.headers.Authorization = `Bearer ${res.accessToken}`
          
          // Уведомляем всех ожидающих
          onRefreshed(res.accessToken)
          isRefreshing = false
          
          return axios(originalRequest)
        } else {
          throw new Error('No access token in refresh response')
        }
      } catch (refreshError) {
        console.error('[axios] refresh failed', refreshError)
        onRefreshed(null)
        isRefreshing = false
        
        if (process.client) {
          const stored = localStorage.getItem('accessToken')
          if (stored && stored !== auth.accessToken) {
            auth.setAccessToken(stored)
            originalRequest.headers.Authorization = `Bearer ${stored}`
            onRefreshed(stored)
            return axios(originalRequest)
          }
        }

        auth.clear()
        
        // Перенаправляем на страницу авторизации
        if (process.client && !window.location.pathname.startsWith('/auth')) {
          navigateTo('/auth')
        }
        
        return Promise.reject(refreshError)
      }
    }
  )
})
