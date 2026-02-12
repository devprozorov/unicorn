import axios from 'axios'
import { useAuthStore } from '~/stores/auth'

let isRefreshing = false
let queue: any[] = []

function resolveQueue(error: any, token: string | null) {
  queue.forEach(p => {
    if (error) p.reject(error)
    else p.resolve(token)
  })
  queue = []
}

const config = useRuntimeConfig()
const baseApiUrl = config.public.apiBase || ''

export const http = axios.create({
  baseURL: `${baseApiUrl}/api`,
  withCredentials: true
})

/**
 * Backward-compatible helper.
 * Some modules previously imported `useHttp()`; keep it to avoid build errors.
 */
export function useHttp() {
  return http
}

http.interceptors.request.use(config => {
  const auth = useAuthStore()
  if (auth.accessToken) {
    config.headers.Authorization = `Bearer ${auth.accessToken}`
  }
  return config
})

http.interceptors.response.use(
  r => r,
  async error => {
    const auth = useAuthStore()
    const original = error.config

    if (error.response?.status !== 401 || original._retry) {
      throw error
    }

    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        queue.push({ resolve, reject })
      }).then(token => {
        original.headers.Authorization = `Bearer ${token}`
        return http(original)
      })
    }

    original._retry = true
    isRefreshing = true

    try {
      const { data } = await http.post('/auth/refresh')
      auth.setAccessToken(data.accessToken)
      resolveQueue(null, data.accessToken)

      original.headers.Authorization = `Bearer ${data.accessToken}`
      return http(original)
    } catch (err) {
      resolveQueue(err, null)
      auth.clear()
      navigateTo('/auth')
      throw err
    } finally {
      isRefreshing = false
    }
  }
)
