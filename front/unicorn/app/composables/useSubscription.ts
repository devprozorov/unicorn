// app/composables/useSubscription.ts
import { useRuntimeConfig } from '#imports'
import { useAuthStore } from '~/stores/auth'

export interface SubscriptionStatus {
  ok: boolean
  active: boolean
  endDate?: string
  daysLeft?: number
}

export interface PaymentLink {
  ok: boolean
  paymentUrl?: string
  invId?: number
  amount?: string
  error?: string
}

export const useSubscription = () => {
  const config = useRuntimeConfig()
  const baseApi = (config.public as any)?.apiBase || '/'
  const auth = useAuthStore()

  const authHeaders = () => {
    return auth.accessToken
      ? { Authorization: `Bearer ${auth.accessToken}` }
      : {}
  }

  const getStatus = async (): Promise<SubscriptionStatus> => {
    try {
      const res: any = await $fetch(`${baseApi}/api/subscription/status`, {
        method: 'GET',
        headers: { ...authHeaders() }
      })
      return res
    } catch {
      return {
        ok: false,
        active: false
      }
    }
  }

  const createPayment = async (): Promise<PaymentLink> => {
    try {
      const res: any = await $fetch(`${baseApi}/api/subscription/create-payment`, {
        method: 'POST',
        headers: { ...authHeaders() }
      })
      return res
    } catch (e: any) {
      return {
        ok: false,
        error: e?.data?.error || e?.message || 'payment_failed'
      }
    }
  }

  const goToPayment = async () => {
    const res = await createPayment()
    if (res.ok && res.paymentUrl) {
      if (process.client) {
        window.location.href = res.paymentUrl
      }
    }
    return res
  }

  return {
    getStatus,
    createPayment,
    goToPayment
  }
}
