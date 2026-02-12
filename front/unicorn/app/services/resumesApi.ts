// services/resumesApi.ts
import { authFetch, AuthFetchError } from '~/composables/useAuthFetch'
import { useAuthStore } from '~/stores/auth'
import { normalizeResume } from '~/utils/normalizeResume'
import type { Resume } from '~/types/resume'
import { useRouter } from 'vue-router'

/* =========================
 * HELPERS
 * ========================= */

function extractItems(res: any): Resume[] {
  if (Array.isArray(res)) return res.map(normalizeResume)
  if (Array.isArray(res?.items)) return res.items.map(normalizeResume)
  return []
}

function extractItem(res: any): Resume | null {
  if (!res) return null
  if (res.item) return normalizeResume(res.item)
  if (res.resume) return normalizeResume(res.resume)
  if (res.id) return normalizeResume(res)
  return null
}

function isMfaRequired(e: any) {
  return (
    e instanceof AuthFetchError &&
    e.status === 403 &&
    e.data?.error === 'mfa_required'
  )
}

/**
 * Backend (по wiki) ожидает FULL BODY на PATCH.
 * Поэтому: любые update/publish сначала берут текущий Resume,
 * мёржат с payload и отправляют полное тело.
 */
function buildFullResumeBody(base: Resume, patch: Partial<Resume>) {
  const merged: Resume = {
    ...base,
    ...patch,
    // на всякий случай нормализуем массивы, чтобы не отправить null
    skills: (patch.skills ?? base.skills ?? []) as any,
    links: (patch.links ?? base.links ?? []) as any,
  } as Resume

  return {
    title: (merged as any).title ?? '',
    about: (merged as any).about ?? '',
    skills: (merged as any).skills ?? [],
    links: (merged as any).links ?? [],
    status: (merged as any).status ?? 'draft',
  }
}

/* =========================
 * API
 * ========================= */

export function useResumesApi() {
  const auth = useAuthStore()
  const router = useRouter()


  /* =========================
   * МОИ РЕЗЮМЕ (LIST)
   * ========================= */
  const getMyResumes = async (): Promise<Resume[]> => {
    try {
      const res: any = await authFetch('/api/resumes/my')
      if (auth.mfaEnabled === null) auth.mfaEnabled = true
      return extractItems(res)
    } catch (e: any) {
      if (isMfaRequired(e)) {
        auth.mfaEnabled = false
        return []
      }
      throw e
    }
  }

  /* =========================
   * МОЁ РЕЗЮМЕ (SINGLE, OWNER)
   * ========================= */
  const getMyResume = async (id: string): Promise<Resume | null> => {
    if (!id) return null

    try {
      const res: any = await authFetch(`/api/resumes/${id}`)
      if (auth.mfaEnabled === null) auth.mfaEnabled = true
      return extractItem(res)
    } catch (e: any) {
      if (isMfaRequired(e)) {
        auth.mfaEnabled = false
        return null
      }
      if (e instanceof AuthFetchError && e.status === 404) return null
      throw e
    }
  }

  /* =========================
   * COMPAT — старый UI
   * ========================= */
  const getResume = async (id: string): Promise<Resume> => {
    const r = await getMyResume(id)
    if (!r) throw new Error('Resume not found')
    return r
  }

  const getResumeById = getResume

  /* =========================
   * ПУБЛИЧНОЕ РЕЗЮМЕ
   * ========================= */
  const getPublicResume = async (id: string): Promise<Resume> => {
    if (!id) throw new Error('getPublicResume: id required')
    
    const config = useRuntimeConfig()
    const baseUrl = config.public.apiBase || ''
    const res = await fetch(`${baseUrl}/api/resumes/${id}`)
    const data = await res.json()
    return normalizeResume(data)
  }

  /* =========================
   * СОЗДАНИЕ (DRAFT) — FIX
   * =========================
   * Backend НЕ принимает пустое body.
   * Отправляем валидное full body сразу.
   */
  const createResume = async (): Promise<{ resumeId: string }> => {
  try {
    const res: any = await authFetch('/api/resumes', {
      method: 'POST',
      body: {
        title: "Новое резюме",
        about: "blank about",
        skills: ["no"],
        links: ["uniconrstar.online"]
        // status: 'draft',
      },
    })

    const resumeId = res?.resumeId ?? res?.id
    if (!resumeId) {
      throw new Error('resumeId not returned from backend')
    }

    if (auth.mfaEnabled === null) auth.mfaEnabled = true
    return { resumeId }
  } catch (e: any) {
    if (isMfaRequired(e)) auth.mfaEnabled = false
    throw e
  }
}


  /* =========================
   * ОБНОВЛЕНИЕ — FIX
   * =========================
   * Принимаем Partial<Resume>, но на backend отправляем FULL BODY.
   */
  const updateResume = async (
  id: string,
  payload: Partial<Resume>
): Promise<void> => {
  if (!id) throw new Error('updateResume: id required')

  const current = await getMyResume(id)
  if (!current) throw new Error('Resume not found')

  await authFetch(`/api/resumes/${id}`, {
    method: 'PATCH',
    body: {
      title: payload.title ?? current.title ?? '',
      about: payload.about ?? current.about ?? '',
      skills: payload.skills ?? current.skills ?? [],
      links: payload.links ?? current.links ?? [],
    },
  })

  if (auth.mfaEnabled === null) auth.mfaEnabled = true
}



  /* =========================
   * ПУБЛИКАЦИЯ — FIX
   * =========================
   * Раньше отправляли только {status:'published'} и ловили 400.
   * Теперь: берём текущее резюме и шлём FULL BODY + status.
   */
const publishResume = async (id: string): Promise<void> => {
  if (!id) throw new Error('publishResume: id required')

  // 1️⃣ берём актуальное резюме
  const current = await getMyResume(id)
  if (!current) throw new Error('Resume not found')

  // 2️⃣ отправляем FULL BODY + status
  await authFetch(`/api/resumes/${id}`, {
    method: 'PATCH',
    body: {
      title: current.title ?? '',
      about: current.about ?? '',
      skills: current.skills ?? [],
      links: current.links ?? [],
      status: 'published',
    },
  })

  if (auth.mfaEnabled === null) auth.mfaEnabled = true
}



  /* =========================
   * УДАЛЕНИЕ
   * ========================= */
  const deleteResume = async (id: string): Promise<void> => {
  if (!id) throw new Error('deleteResume: id required')

  try {
    await authFetch(`/api/resumes/${id}`, { method: 'DELETE' })

    if (auth.mfaEnabled === null) auth.mfaEnabled = true

    // ✅ редирект после успешного удаления
    await router.push(`/profile/${auth.userId}`, { replace: true })
  } catch (e: any) {
    if (isMfaRequired(e)) auth.mfaEnabled = false
    throw e
  }
}


  /* =========================
   * COMPAT ДЛЯ UI
   * ========================= */
  const askDeleteResume = async (id: string): Promise<void> => deleteResume(id)

  return {
    // list / single
    getMyResumes,
    getMyResume,

    // compat
    getResume,
    getResumeById,

    // public
    getPublicResume,

    // mutations
    createResume,
    updateResume,
    publishResume,
    deleteResume,

    // ui compat
    askDeleteResume,
  }
}
