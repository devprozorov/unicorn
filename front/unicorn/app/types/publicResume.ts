// app/types/publicResume.ts

export interface PublicResume {
  id: string

  title: string
  about: string

  skills: string[]
  links: string[]
  experience: {
    title: string
    company?: string
    from?: string
    to?: string | null
    description?: string
  }[]

  updatedAt: string
}
