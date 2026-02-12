import type { PublicResume } from '~/types/publicResume'

function normalizePublicResume(r: any): PublicResume {
  return {
    id: r.id || r._id,
    title: r.title ?? '',
    about: r.about ?? '',
    skills: Array.isArray(r.skills) ? r.skills : [],
    links: Array.isArray(r.links) ? r.links : [],
    experience: Array.isArray(r.experience) ? r.experience : [],
    updatedAt: r.updatedAt ?? ''
  }
}

export function usePublicResumesApi() {
  const getById = async (id: string): Promise<PublicResume> => {
    const config = useRuntimeConfig()
    const baseUrl = config.public.apiBase || ''
    const res = await fetch(`${baseUrl}/api/public/resumes/${id}`)
    const data = await res.json()
    return normalizePublicResume(data)
  }

  return { getById }
}
