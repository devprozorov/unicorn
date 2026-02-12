<template>
  <div class="p-6 space-y-6">
    <h1 class="text-2xl font-semibold">Company Profile</h1>

    <div v-if="err" class="p-3 rounded bg-red-100 text-red-800">{{ err }}</div>
    <div v-if="hint" class="p-3 rounded bg-yellow-100 text-yellow-900">{{ hint }}</div>

    <div class="grid gap-6 md:grid-cols-2">
      <div class="p-4 rounded border space-y-3">
        <h2 class="font-semibold">Public info</h2>

        <label class="block">
          <div class="text-sm opacity-70">Display name</div>
          <input v-model="form.displayName" class="w-full border rounded p-2" />
        </label>

        <label class="block">
          <div class="text-sm opacity-70">About</div>
          <textarea v-model="form.about" class="w-full border rounded p-2" rows="5" />
        </label>

        <label class="block">
          <div class="text-sm opacity-70">Location</div>
          <input v-model="form.location" class="w-full border rounded p-2" />
        </label>

        <label class="block">
          <div class="text-sm opacity-70">Industry</div>
          <input v-model="form.industry" class="w-full border rounded p-2" />
        </label>

        <label class="block">
          <div class="text-sm opacity-70">Website</div>
          <input v-model="form.website" class="w-full border rounded p-2" />
        </label>

        <label class="block">
          <div class="text-sm opacity-70">Links (one per line)</div>
          <textarea v-model="linksText" class="w-full border rounded p-2" rows="4" />
        </label>

        <button class="px-4 py-2 rounded bg-black text-white" @click="saveProfile" :disabled="loading">
          {{ loading ? 'Saving...' : 'Save' }}
        </button>
      </div>

      <div class="p-4 rounded border space-y-3">
        <h2 class="font-semibold">Avatar</h2>

        <div class="flex items-center gap-3">
          <img :src="avatarUrl" class="w-16 h-16 rounded object-cover border" />

          <input type="file" accept="image/png,image/jpeg,image/webp" @change="onAvatar" />
        </div>

        <div class="text-sm opacity-70">
          Upload PNG/JPG/WEBP. URL будет вида <code>/uploads/avatars/...</code>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import axios from 'axios'
import { getErrorMessage } from '~/utils/errorMessages'
import { useI18n } from '~/composables/useI18n'
definePageMeta({ middleware: ['company-auth'] })

const auth = useAuthStore()
const { t } = useI18n()

const loading = ref(false)
const err = ref('')
const hint = ref('')

type Profile = {
  userId: string
  type: 'company' | 'user'
  displayName: string
  about?: string
  location?: string
  industry?: string
  website?: string
  links?: string[]
  avatarUrl?: string
}

const profile = ref<Profile | null>(null)
const avatarUrl = computed(() => profile.value?.avatarUrl || '/images/com-base.webp')

const form = reactive({
  displayName: '',
  about: '',
  location: '',
  industry: '',
  website: '',
})

const linksText = ref('')

function normalizeLinks(txt: string) {
  return txt
    .split('\n')
    .map(s => s.trim())
    .filter(Boolean)
    .slice(0, 20)
}

async function load() {
  err.value = ''
  hint.value = ''

  const me = await axios.get('/auth/me')
  auth.displayName = me.data.displayName
  auth.type = me.data.type
  if (me.data.userId) auth.userId = me.data.userId

  const res = await axios.get('/profile/me')
  profile.value = res.data.profile

  if (profile.value) {
    form.displayName = profile.value.displayName || auth.displayName || ''
    form.about = profile.value.about || ''
    form.location = profile.value.location || ''
    form.industry = profile.value.industry || ''
    form.website = profile.value.website || ''
    linksText.value = (profile.value.links || []).join('\n')
  } else {
    form.displayName = auth.displayName || ''
  }
}

async function saveProfile() {
  loading.value = true
  err.value = ''
  hint.value = ''

  const payload: any = {
    displayName: form.displayName,
    about: form.about,
    location: form.location,
    industry: form.industry,
    website: form.website,
    links: normalizeLinks(linksText.value),
  }

  try {
    if (!profile.value) {
      await axios.post('/profile/me', payload)
    } else {
      await axios.patch('/profile/me', payload)
    }
    await load()
    hint.value = 'Saved'
  } catch (e: any) {
    err.value = getErrorMessage(e?.response?.data?.error || 'save_failed', t, 'errors.saveFailed')
  } finally {
    loading.value = false
  }
}

function compressImage(file: File, maxSizeMB = 2): Promise<File> {
  return new Promise((resolve, reject) => {
    const maxSizeBytes = maxSizeMB * 1024 * 1024
    
    // Если файл уже меньше лимита, возвращаем как есть
    if (file.size <= maxSizeBytes) {
      resolve(file)
      return
    }

    const reader = new FileReader()
    reader.readAsDataURL(file)
    
    reader.onload = (e) => {
      const img = new Image()
      img.src = e.target?.result as string
      
      img.onload = () => {
        const canvas = document.createElement('canvas')
        let width = img.width
        let height = img.height
        
        // Максимальные размеры
        const maxDimension = 1024
        
        if (width > height && width > maxDimension) {
          height = (height * maxDimension) / width
          width = maxDimension
        } else if (height > maxDimension) {
          width = (width * maxDimension) / height
          height = maxDimension
        }
        
        canvas.width = width
        canvas.height = height
        
        const ctx = canvas.getContext('2d')
        if (!ctx) {
          reject(new Error('Failed to get canvas context'))
          return
        }
        
        ctx.drawImage(img, 0, 0, width, height)
        
        // Пробуем разное качество сжатия
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              reject(new Error('Failed to compress image'))
              return
            }
            
            const compressedFile = new File([blob], file.name, {
              type: 'image/jpeg',
              lastModified: Date.now(),
            })
            
            resolve(compressedFile)
          },
          'image/jpeg',
          0.85 // качество сжатия
        )
      }
      
      img.onerror = () => reject(new Error('Failed to load image'))
    }
    
    reader.onerror = () => reject(new Error('Failed to read file'))
  })
}

async function onAvatar(ev: Event) {
  const input = ev.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  err.value = ''
  hint.value = ''

  // Проверка типа файла
  if (!file.type.startsWith('image/')) {
    err.value = 'Please select an image file'
    input.value = ''
    return
  }

  // Проверка размера (максимум 10MB для исходного файла)
  const maxOriginalSize = 10 * 1024 * 1024
  if (file.size > maxOriginalSize) {
    err.value = 'File is too large. Maximum size is 10MB'
    input.value = ''
    return
  }

  try {
    // Сжимаем изображение до 2MB
    const compressedFile = await compressImage(file, 2)
    
    const fd = new FormData()
    fd.append('avatar', compressedFile)

    const res = await axios.post('/profile/me/avatar', fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    
    // Проверяем успешный ответ от backend
    if (res.data && res.data.avatarUrl) {
      hint.value = 'Avatar successfully updated!'
      if (profile.value) profile.value.avatarUrl = res.data.avatarUrl
      else await load()
    } else {
      err.value = 'Failed to update avatar: invalid response from server'
    }
  } catch (e: any) {
    if (e.message && (e.message.includes('compress') || e.message.includes('canvas') || e.message.includes('Failed to'))) {
      err.value = 'Failed to process image. Please try a different file.'
    } else {
      err.value = getErrorMessage(e?.response?.data?.error || 'upload_failed', t, 'errors.uploadFailed')
    }
  } finally {
    input.value = ''
  }
}

onMounted(load)
</script>
