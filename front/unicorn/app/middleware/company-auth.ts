export default defineNuxtRouteMiddleware(async () => {
  // ✅ на сервере ничего не делаем
  if (import.meta.server) return

  const auth = useAuthStore()

  if (!auth.accessToken) return navigateTo('/auth')

  // Ждем инициализации пользователя
  if (!auth.initialized) {
    await auth.init()
  }

  if (!auth.user) {
    auth.clear()
    return navigateTo('/auth')
  }

  if (auth.user.type !== 'company') return navigateTo('/')
})
