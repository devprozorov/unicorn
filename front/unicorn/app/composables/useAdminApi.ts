// app/composables/useAdminApi.ts
import { useRuntimeConfig } from '#imports'

export type AdminUserType = 'user' | 'company'

export interface AdminUserItem {
  userId: string
  login: string
  displayName: string
  type: AdminUserType
}

export interface AdminUserDetails {
  ok: boolean
  userId: string
  login: string
  displayName: string
  type: AdminUserType
  status: {
    deleted: boolean
    blocked: boolean
  }
  subscription: {
    active: boolean
    until: string
  }
  mfa: {
    totpEnabled: boolean
  }
  createdAt: string
  updatedAt: string
}

export interface UserFilters {
  type?: AdminUserType
  search?: string
  blocked?: 'true' | 'false'
  deleted?: 'true' | 'false'
  premium?: 'true' | 'false'
}

export const useAdminApi = () => {
  const config = useRuntimeConfig()
  const baseApi = (config.public as any)?.apiBase || '/'
  const baseUrl = `${baseApi}/api/admin`
  const token = useState<string | null>('adminToken', () => null)

  const authHeaders = () => {
    return token.value
      ? { Authorization: `Bearer ${token.value}` }
      : {}
  }

  const login = async (login: string, password: string) => {
    try {
      const res: any = await $fetch(`${baseUrl}/login`, {
        method: 'POST',
        body: { login, password }
      })

      if (res?.ok && res.accessToken) {
        token.value = res.accessToken
        try {
          localStorage.setItem('adminToken', res.accessToken)
        } catch {}
      }
      return res
    } catch (e: any) {
      return { ok: false, error: e?.data?.error || 'login_failed' }
    }
  }

  const getUsers = async (filters?: UserFilters) => {
    try {
      const res: any = await $fetch(`${baseUrl}/users`, {
        method: 'GET',
        headers: { ...authHeaders() },
        query: filters || {}
      })
      return (res?.items ?? []) as AdminUserItem[]
    } catch (e: any) {
      throw new Error(e?.data?.error || 'load_failed')
    }
  }

  const getUserDetails = async (userId: string): Promise<AdminUserDetails> => {
    try {
      const res: any = await $fetch(`${baseUrl}/users/${userId}`, {
        method: 'GET',
        headers: { ...authHeaders() }
      })
      return res
    } catch (e: any) {
      throw new Error(e?.data?.error || 'user_not_found')
    }
  }

  const updateUser = async (userId: string, updates: { displayName?: string, login?: string }) => {
    try {
      const res: any = await $fetch(`${baseUrl}/users/${userId}`, {
        method: 'PATCH',
        headers: { ...authHeaders() },
        body: updates
      })
      return res
    } catch (e: any) {
      return { ok: false, error: e?.data?.error || 'update_failed' }
    }
  }

  const blockUser = async (userId: string) => {
    try {
      const res: any = await $fetch(`${baseUrl}/users/${userId}/block`, {
        method: 'POST',
        headers: { ...authHeaders() }
      })
      return res
    } catch (e: any) {
      return { ok: false, error: e?.data?.error || 'block_failed' }
    }
  }

  const unblockUser = async (userId: string) => {
    try {
      const res: any = await $fetch(`${baseUrl}/users/${userId}/unblock`, {
        method: 'POST',
        headers: { ...authHeaders() }
      })
      return res
    } catch (e: any) {
      return { ok: false, error: e?.data?.error || 'unblock_failed' }
    }
  }

  const deleteUser = async (userId: string) => {
    try {
      const res: any = await $fetch(`${baseUrl}/users/${userId}`, {
        method: 'DELETE',
        headers: { ...authHeaders() }
      })
      return res
    } catch (e: any) {
      return { ok: false, error: e?.data?.error || 'delete_failed' }
    }
  }

  const activateSubscription = async (userId: string, days: number) => {
    try {
      const res: any = await $fetch(`${baseUrl}/users/${userId}/subscription/activate`, {
        method: 'POST',
        headers: { ...authHeaders() },
        body: { days }
      })
      return res
    } catch (e: any) {
      return { ok: false, error: e?.data?.error || 'activation_failed' }
    }
  }

  const deactivateSubscription = async (userId: string) => {
    try {
      const res: any = await $fetch(`${baseUrl}/users/${userId}/subscription/deactivate`, {
        method: 'POST',
        headers: { ...authHeaders() }
      })
      return res
    } catch (e: any) {
      return { ok: false, error: e?.data?.error || 'deactivation_failed' }
    }
  }

  return {
    token,
    login,
    getUsers,
    getUserDetails,
    updateUser,
    blockUser,
    unblockUser,
    deleteUser,
    activateSubscription,
    deactivateSubscription
  }
}
