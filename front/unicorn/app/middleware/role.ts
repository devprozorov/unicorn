import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuthStore()

  // /auth не перехватываем
  if (to.path.startsWith('/auth')) return

  // без токена — на /auth
  if (!auth.accessToken) {
    return navigateTo('/auth')
  }

  // защита “company-only”
  if (to.path.startsWith('/panel/company') && auth.accountType !== 'company') {
    return navigateTo('/panel')
  }

  // защита “user-only”
  if (to.path.startsWith('/panel/user') && auth.accountType !== 'user') {
    return navigateTo('/panel')
  }
})
