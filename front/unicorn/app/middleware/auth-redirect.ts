// app/middleware/auth-redirect.ts
import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(() => {
  if (import.meta.server) return

  const auth = useAuthStore()

  if (!auth.accessToken && process.client) {
    try {
      const stored = localStorage.getItem('accessToken')
      if (stored) auth.setAccessToken(stored)
    } catch {
      // ignore storage errors
    }
  }

  if (!auth.accessToken) {
    return navigateTo('/auth')
  }

  if (!auth.initialized) {
    return auth.init().then(() => {
      if (!auth.user) return navigateTo('/auth')
    })
  }

  if (!auth.user) return navigateTo('/auth')
})

