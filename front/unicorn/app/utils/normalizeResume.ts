// utils/normalizeResume.ts
import type { Resume } from '~/types/resume'

export function normalizeResume(r: any): Resume {
  if (!r) {
    throw new Error('normalizeResume: empty payload')
  }

  return {
    id: r.id || r.resumeId || r._id, // ← КЛЮЧЕВО
    userId: r.userId,
    title: r.title ?? '',
    about: r.about ?? '',
    skills: Array.isArray(r.skills) ? r.skills : [],
    links: Array.isArray(r.links) ? r.links : [],
    experience: Array.isArray(r.experience) ? r.experience : [],
    status: r.status ?? 'draft',
    createdAt: r.createdAt,
    updatedAt: r.updatedAt
  }
}
