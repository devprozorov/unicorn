// app/composables/useApplications.ts
import axios from 'axios'

export const useApplications = () => {
  /**
   * Скрыть отклик/чат для текущего пользователя
   * @param applicationId - ID отклика
   */
  const hideApplication = async (applicationId: string): Promise<{ ok: boolean; error?: string }> => {
    try {
      const response = await axios.post(`/api/applications/${applicationId}/hide`)
      return response.data
    } catch (e: any) {
      return {
        ok: false,
        error: e?.response?.data?.error || e?.message || 'hide_failed'
      }
    }
  }

  return {
    hideApplication
  }
}
