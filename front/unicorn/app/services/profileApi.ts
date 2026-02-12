import { authFetch } from '~/composables/useAuthFetch'

export function useProfileApi() {
  /**
   * ТОЛЬКО для владельца
   * Использовать ТОЛЬКО на страницах:
   * - /profile/edit
   * - /panel/security
   */
  const getMyProfile = async () => {
    return await authFetch('/api/profile/me')
  }

  /**
   * Публичный профиль
   * ❌ НИКОГДА не использовать authFetch
   * ❌ НИКОГДА не вешать middleware
   */
  const getPublicProfile = async (loginOrId: string) => {
    if (!loginOrId) {
        throw new Error('login or userId is required')
    }

    const config = useRuntimeConfig()
    const baseUrl = config.public.apiBase || ''
    const res = await fetch(`${baseUrl}/api/profile/${loginOrId}`)
    return await res.json()
  }


  /**
   * Создание профиля (1 раз)
   */
  const createProfile = async (payload: any) => {
    if (!payload || typeof payload !== 'object') {
      throw new Error('Invalid payload for createProfile')
    }

    return await authFetch('/api/profile/me', {
      method: 'POST',
      body: payload
    })
  }

  /**
   * Частичное обновление профиля
   */
  const updateProfile = async (payload: any) => {
    if (!payload || typeof payload !== 'object') {
      throw new Error('Invalid payload for updateProfile')
    }

    return await authFetch('/api/profile/me', {
      method: 'PATCH',
      body: payload
    })
  }

  return {
    getMyProfile,
    getPublicProfile,
    createProfile,
    updateProfile
  }
}
