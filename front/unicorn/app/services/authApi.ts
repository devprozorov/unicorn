// app/services/authApi.ts
import axios from 'axios'
import { useAuthStore } from '~/stores/auth'

export function useAuthApi() {
  const config = useRuntimeConfig()
  const baseApiUrl = config.public.apiBase || ''
  
  const api = axios.create({
    baseURL: `${baseApiUrl}/api/auth`,
    withCredentials: true
  })

  /**
   * =====================================
   * REQUEST INTERCEPTOR
   * =====================================
   * Гарантирует, что ВСЕ запросы:
   * - login
   * - refresh
   * - totpEnroll / totpEnable
   * идут с Authorization: Bearer
   */
  api.interceptors.request.use((config) => {
    const auth = useAuthStore()

    // НЕ добавляем Authorization для login/register
    if (
        auth.accessToken &&
        !config.url?.endsWith('/login') &&
        !config.url?.endsWith('/register')
    ) {
        config.headers = config.headers || {}
        config.headers.Authorization = `Bearer ${auth.accessToken}`
    }

    return config
  })


  /**
   * =====================================
   * RESPONSE INTERCEPTOR (опционально)
   * =====================================
   * Можно расширить позже под auto-logout / refresh
   */
  api.interceptors.response.use(
    r => r,
    error => {
        if (process.dev) {
        console.error('[authApi] error', error?.response?.data || error)
        }
        throw error
    }
  )


  return {
    /**
     * =====================================
     * LOGIN
     * =====================================
     * Возможные ответы:
     * - ok + accessToken
     * - ok + mfaRequired + mfaToken
     */
    async login(login: string, password: string) {
      const { data } = await api.post('/login', { login, password })
      console.log('[authApi] login', data)
      return data
    },

    /**
     * =====================================
     * REGISTER
     * =====================================
     */
    async register(payload: {
      login: string
      password: string
      displayName: string
      type: 'user' | 'company'
    }) {
      const { data } = await api.post('/register', payload)
      console.log('[authApi] register', data)
      return data
    },

    /**


    /**
     * =====================================
     * REFRESH TOKEN
     * =====================================
     * Использует cookie (HttpOnly refresh)
     */
    async refresh() {
      const { data } = await api.post('/refresh')
      console.log('[authApi] refresh', data)
      return data
    },

    /**
     * =====================================
     * LOGOUT
     * =====================================
     */
    async logout() {
      await api.post('/logout')
      console.log('[authApi] logout')
    },

    /**
     * =====================================
     * CHANGE PASSWORD
     * =====================================
     */
    async changePassword(payload: { oldPassword: string; newPassword: string }) {
      const { data } = await api.post('/change-password', payload)
      console.log('[authApi] changePassword', data)
      return data
    },

    /**
     * =====================================
     * MFA / TOTP
     * =====================================
     */

    /**
     * Шаг 1 — начать подключение TOTP
     * backend вернёт otpauthUrl
     */
    async totpEnroll() {
      const { data } = await api.post('/totp/enroll')
      console.log('[authApi] totpEnroll', data)
      return data
    },

    /**
     * Шаг 2 — подтвердить TOTP код
     * backend после этого:
     * - включает MFA
     * - отзывает refresh
     */
    async totpEnable(code: string) {
      const { data } = await api.post('/totp/enable', { code })
      console.log('[authApi] totpEnable', data)
      return data
    },

    /**
     * ВТОРОЙ ШАГ LOGIN при MFA
     */
    async totpVerify(payload: {
      mfaToken: string
      code: string
    }) {
      const { data } = await api.post('/totp/verify', payload)
      console.log('[authApi] totpVerify', data)
      return data
    },

    /**
     * (на будущее)
     * Отключение MFA — если backend разрешит
     */
    async totpDisable(code: string) {
      const { data } = await api.post('/totp/disable', { code })
      console.log('[authApi] totpDisable', data)
      return data
    },

    /**
     * =====================================
     * ME
     * =====================================
     * Запрашивает основную информацию о текущем пользователе.  Backend
     * возвращает displayName и type, но не включает id.  ID и тип
     * извлекаются из accessToken (JWT) без проверки подписи.  Login
     * не возвращается и оставляется пустым.  Возвращаемый объект
     * совместим с useAuthStore.init().
     */
    async me() {
      const auth = useAuthStore()
      let payload: any = null
      // Attempt to fetch displayName and type from backend
      let displayName: string | undefined
      let type: string | undefined
      try {
        const { data } = await api.get('/me')
        displayName = data?.displayName
        type = data?.type
      } catch (e) {
        // ignore; we'll still try to decode JWT for id and type
      }
      // Decode JWT accessToken to extract userId (sub) and type (typ)
      const token = auth.accessToken
      let id: string | undefined
      try {
        if (token) {
          const parts = token.split('.')
          if (parts.length === 3) {
            let base64 = parts[1].replace(/-/g, '+').replace(/_/g, '/')
            // pad base64 string
            while (base64.length % 4) base64 += '='
            const json = atob(base64)
            payload = JSON.parse(json)
            id = payload?.sub
            // fallback to typ from token if backend didn't provide
            if (!type && payload?.typ) type = payload.typ
          }
        }
      } catch (e) {
        // decoding failed; id will remain undefined
      }
      if (!displayName) displayName = ''
      return { ok: true, id, login: '', displayName, type }
    }
  }
}
