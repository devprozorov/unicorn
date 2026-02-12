export default defineNuxtRouteMiddleware(async () => {
  // SSR не трогаем (иначе будет 302 на /auth при прямом заходе)
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

  if (auth.user.type !== 'user') return navigateTo('/')
})
