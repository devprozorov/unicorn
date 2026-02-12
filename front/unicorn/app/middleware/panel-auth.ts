import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()

  /**
   * 🔑 КРИТИЧНО:
   * восстановление UI-состояния ДО рендера
   */
  auth.restoreMfaState()

  if (!auth.isAuthenticated) {
    return navigateTo('/auth')
  }
})
