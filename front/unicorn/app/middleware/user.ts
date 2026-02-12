export default defineNuxtRouteMiddleware(() => {
  const auth = useAuthStore()

  if (!auth.isAuthenticated) {
    return navigateTo('/auth')
  }

  if (auth.accountType !== 'user') {
    return navigateTo('/')
  }
})
