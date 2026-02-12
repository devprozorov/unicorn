import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin(() => {
  if (!process.client) return

  const auth = useAuthStore()

  function onStorage(e: StorageEvent) {
    if (e.key !== 'accessToken') return

    const next = e.newValue
    if (next === auth.accessToken) return

    if (next) {
      auth.setAccessToken(next)
      auth.init()
      return
    }

    auth.clear({ keepStorage: true })
  }

  window.addEventListener('storage', onStorage)
})
