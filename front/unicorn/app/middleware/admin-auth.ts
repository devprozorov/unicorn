// app/middleware/admin-auth.ts
export default defineNuxtRouteMiddleware(() => {
  const token = useState<string | null>('adminToken')
  if (process.client) {
    if (!token.value) {
      // Try restore from localStorage
      try {
        const restored = localStorage.getItem('adminToken')
        if (restored) token.value = restored
      } catch {}
    }
  }

  if (!token.value) {
    return navigateTo('/admin/login')
  }
})
