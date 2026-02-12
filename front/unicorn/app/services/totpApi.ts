import axios from 'axios'
import { useAuthStore } from '~/stores/auth'

export function useTotpApi() {
  const auth = useAuthStore()
  const config = useRuntimeConfig()
  const baseApiUrl = config.public.apiBase || ''

  const api = axios.create({
    baseURL: `${baseApiUrl}/api/auth`,
    withCredentials: true
  })

  api.interceptors.request.use((config) => {
    if (auth.accessToken) {
      config.headers.Authorization = `Bearer ${auth.accessToken}`
    }
    return config
  })

  return {
    async enroll() {
      const { data } = await api.post('/totp/enroll')
      return data as { ok: true; otpauthUrl: string }
    },

    async enable(code: string) {
      const { data } = await api.post('/totp/enable', { code })
      return data as { ok: true }
    },

    async verify(mfaToken: string, code: string) {
      const { data } = await api.post('/totp/verify', {
        mfaToken,
        code
      })

      return data as {
        ok: true
        accessToken: string
      }
    }
  }
}
