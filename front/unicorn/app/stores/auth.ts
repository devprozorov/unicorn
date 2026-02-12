import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthApi } from '~/services/authApi'

/**
 * ЕДИНЫЙ контракт пользователя для всего фронта
 * userId используется ВЕЗДЕ
 */
export interface AuthUser {
  userId: string
  login: string
  displayName: string
  type: 'user' | 'company' | 'admin'
}

export const useAuthStore = defineStore('auth', () => {
  const api = useAuthApi()

  /* =====================
     STATE
     ===================== */

  /**
   * Access JWT
   * Только в памяти
   */
  const accessToken = ref<string | null>(null)

  /**
   * Текущий пользователь
   */
  const user = ref<AuthUser | null>(null)

  /**
   * Был ли выполнен init()
   * Сбрасывается ТОЛЬКО при смене accessToken
   */
  const initialized = ref(false)

  /**
   * Флаг загрузки init
   */
  const loading = ref(false)

  /**
   * MFA (временный токен)
   */
  const mfaToken = ref<string | null>(null)

  /* =====================
     GETTERS
     ===================== */

  /**
   * Есть ли accessToken
   * НЕ значит, что user загружен
   */
  const isAuthenticated = computed(() => !!accessToken.value)

  /**
   * Унифицированный userId
   */
  const userId = computed(() => user.value?.userId ?? null)

  /**
   * Auth полностью готов к использованию
   */
  const isReady = computed(() => initialized.value && !loading.value)

  /**
    * Косвенный флаг MFA
    * true  → 2FA включена
    * false → не включена
    * null  → ещё не определено
    */
  const mfaEnabled = ref<boolean | null>(null)

  /* =====================
     ACTIONS
     ===================== */

  /**
   * Установка accessToken
   * ВСЕГДА инвалидирует init
   */
  function setAccessToken(token: string) {
    accessToken.value = token
    initialized.value = false
    
    // Сохраняем в localStorage для восстановления после перезагрузки
    if (process.client) {
      try {
        localStorage.setItem('accessToken', token)
      } catch (e) {
        console.warn('Failed to save accessToken to localStorage', e)
      }
    }
  }

  /**
   * Установка пользователя
   */
  function setUser(data: AuthUser) {
    user.value = data
  }

  /**
   * MFA: требуется подтверждение
   */
  function requireMfa(token: string) {
    mfaToken.value = token
  }

  /**
   * MFA подтверждена
   */
  function confirmMfa() {
    mfaToken.value = null
  }

  /**
   * UI-restore hook for MFA state (currently noop).
   * Middleware expects this to exist; extend when MFA UI persists state.
   */
  function restoreMfaState() {
    // No persisted MFA UI state yet
  }

  /**
   * INIT
   *
   * Единственный источник истины для /me
   * Можно вызывать сколько угодно раз — выполнится один раз
   */
  async function init() {
    if (initialized.value || loading.value) return

    loading.value = true

    if (!accessToken.value) {
      initialized.value = true
      loading.value = false
      return
    }

    try {
      const res = await api.me()

      /**
       * Поддержка ОБОИХ форматов backend:
       * - userId
       * - id
       */
      if (res?.ok) {
        setUser({
          userId: res.userId ?? res.id,
          login: res.login,
          displayName: res.displayName,
          type: res.type
        })
        initialized.value = true
      } else {
        // Только если backend явно говорит, что токен невалиден
        clear()
      }
    } catch (err: any) {
      // Сбрасываем только при 401/403 (невалидный токен)
      // При других ошибках (сеть, 500) сохраняем токен
      const status = err?.response?.status || err?.status
      if (status === 401 || status === 403) {
        clear()
      } else {
        // При временных ошибках оставляем токен, но помечаем как initialized
        initialized.value = true
        console.warn('Failed to load user profile, but keeping token:', err)
      }
    } finally {
      loading.value = false
    }
  }

  /**
   * REFRESH ACCESS TOKEN
   * Используется interceptor'ами
   */
  async function refreshAccessToken(): Promise<string | null> {
    try {
      const res = await api.refresh()

      if (res?.ok && res.accessToken) {
        setAccessToken(res.accessToken)
        return res.accessToken
      }
    } catch {
      clear()
    }
    return null
  }

  /**
   * Полный logout / сброс
   */
  function clear(options?: { keepStorage?: boolean }) {
    accessToken.value = null
    user.value = null
    mfaToken.value = null
    initialized.value = false
    loading.value = false
    
    // Очищаем localStorage
    if (process.client && !options?.keepStorage) {
      try {
        localStorage.removeItem('accessToken')
      } catch (e) {
        console.warn('Failed to remove accessToken from localStorage', e)
      }
    }
  }

  /* =====================
     EXPORT
     ===================== */

  return {
    // state
    accessToken,
    user,
    mfaToken,
    initialized,
    loading,
    mfaEnabled,

    // getters
    isAuthenticated,
    userId,
    isReady,

    // actions
    setAccessToken,
    setUser,
    requireMfa,
    confirmMfa,
    restoreMfaState,
    init,
    refreshAccessToken,
    clear
  }
})
