import { useAuthStore } from '~/stores/auth'
import { useAuthApi } from '~/services/authApi'

export default defineNuxtPlugin(async () => {
  const auth = useAuthStore()
  const api = useAuthApi()

  // ⚠️ важно: плагин может выполниться несколько раз
  if (auth.initialized) return

  try {
    // 1. Пробуем восстановить accessToken из localStorage
    if (!auth.accessToken && process.client) {
      const stored = localStorage.getItem('accessToken')
      if (stored) {
        // Проверяем, не истек ли токен
        try {
          const parts = stored.split('.')
          if (parts.length === 3) {
            let base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/')
            while (base64.length % 4) base64 += '='
            const payload = JSON.parse(atob(base64))
            const exp = payload?.exp
            
            // Если токен истек, удаляем его
            if (exp && exp * 1000 < Date.now()) {
              localStorage.removeItem('accessToken')
              console.warn('[auth-init] Token expired, removed from storage')
            } else {
              auth.setAccessToken(stored)
            }
          }
        } catch (e) {
          // Если не удалось разобрать токен, удаляем его
          localStorage.removeItem('accessToken')
          console.warn('[auth-init] Invalid token format, removed from storage')
        }
      }
    }

    // 2. Есть accessToken → просто грузим user
    if (auth.accessToken) {
      await auth.init()
      return
    }

    // 3. Нет accessToken → пробуем refresh-cookie
    const res = await api.refresh()

    if (res?.ok && res.accessToken) {
      auth.setAccessToken(res.accessToken)
      await auth.init()
      return
    }
  } catch {
    // тихо — гость
  } finally {
    auth.initialized = true
  }
})
