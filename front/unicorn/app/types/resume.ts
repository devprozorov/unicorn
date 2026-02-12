// types/resume.ts
export interface Resume {
  id: string
  userId?: string

  title: string
  about: string

  skills: string[]
  links: string[]

  experience: any[]

  status: 'draft' | 'published' | 'active'

  createdAt?: string
  updatedAt?: string
}
