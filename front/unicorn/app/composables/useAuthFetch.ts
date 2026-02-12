// composables/useAuthFetch.ts
import { useAuthStore } from '~/stores/auth'
import { useAuthApi } from '~/services/authApi'

type AuthFetchErrorData = any

export class AuthFetchError extends Error {
  status: number
  data: AuthFetchErrorData

  constructor(message: string, status: number, data: AuthFetchErrorData) {
    super(message)
    this.name = 'AuthFetchError'
    this.status = status
    this.data = data
  }
}

function isPlainObject(v: any) {
  return v !== null && typeof v === 'object' && (v.constructor === Object || Object.getPrototypeOf(v) === Object.prototype)
}

export async function authFetch<T = any>(
  url: string,
  options: any = {}
): Promise<T> {
  const auth = useAuthStore()
  const authApi = useAuthApi()
  const config = useRuntimeConfig()
  const baseApiUrl = config.public.apiBase || ''

  // Если URL относительный (начинается с /), добавляем baseURL
  const fullUrl = url.startsWith('/') ? `${baseApiUrl}${url}` : url

  const attempt = async (): Promise<Response> => {
    const headers = new Headers(options.headers || {})

    // Authorization
    if (auth.accessToken && !headers.has('Authorization')) {
      headers.set('Authorization', `Bearer ${auth.accessToken}`)
    }

    // credentials for refresh-cookie endpoints
    const credentials = options.credentials || 'include'

    // Body -> JSON (если передали объект)
    let body = options.body
    if (isPlainObject(body) || Array.isArray(body)) {
      if (!headers.has('Content-Type')) {
        headers.set('Content-Type', 'application/json')
      }
      body = JSON.stringify(body)
    }

    return await fetch(fullUrl, {
      ...options,
      headers,
      credentials,
      body
    })
  }

  // 1) первый запрос
  let res = await attempt()

  // 2) если access протух — пробуем refresh и повторяем 1 раз
  if (res.status === 401 && auth.accessToken) {
    try {
      const refreshed = await authApi.refresh()
      if (refreshed?.ok && refreshed?.accessToken) {
        auth.setAccessToken(refreshed.accessToken)
        await auth.init()
        res = await attempt()
      }
    } catch {
      // refresh не удался — чистим сессию
      auth.clear?.()
    }
  }

  // 3) обработка ошибок
  if (!res.ok) {
    let data: any = null
    try {
      data = await res.json()
    } catch {
      data = await res.text().catch(() => null)
    }

    const message =
      (data && typeof data === 'object' && (data.error || data.message)) ||
      `HTTP ${res.status}`

    throw new AuthFetchError(String(message), res.status, data)
  }

  // 4) успех
  const ct = res.headers.get('content-type') || ''
  if (ct.includes('application/json')) {
    return (await res.json()) as T
  }
  // на случай текстовых ответов
  return (await res.text()) as unknown as T
}
