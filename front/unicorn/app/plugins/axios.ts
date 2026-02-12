// app/plugins/axios.ts
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
  const config = useRuntimeConfig()
  
  // Устанавливаем baseURL для всех axios запросов
  axios.defaults.baseURL = `${config.public.apiBase}/api`
  axios.defaults.withCredentials = true

  axios.interceptors.request.use((config) => {
    if (auth.accessToken) {
      config.headers.Authorization = `Bearer ${auth.accessToken}`
    }
    return config
  })

  axios.interceptors.response.use(
    (r) => r,
    async (error) => {
      const originalRequest = error.config

      if (error.response?.status !== 401 || originalRequest._retry) {
        return Promise.reject(error)
      }

      originalRequest._retry = true

      // Если уже идет refresh, ждем
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

      try {
        const res = await api.refresh()

        if (res?.accessToken) {
          auth.setAccessToken(res.accessToken)
          originalRequest.headers.Authorization = `Bearer ${res.accessToken}`
          
          onRefreshed(res.accessToken)
          isRefreshing = false
          
          return axios(originalRequest)
        } else {
          throw new Error('No access token')
        }
      } catch (refreshError) {
        onRefreshed(null)
        isRefreshing = false
        auth.clear()
        
        if (process.client && !window.location.pathname.startsWith('/auth')) {
          navigateTo('/auth')
        }
        
        return Promise.reject(refreshError)
      }
    }
  )
})
