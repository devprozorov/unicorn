<script setup lang="ts">
import axios from 'axios'
import { computed, ref } from 'vue'
import { useI18n } from '~/composables/useI18n'
import { useToast } from '~/composables/useToast'
import { useAuthStore } from '~/stores/auth'
import { useRouter } from 'vue-router'
const { t } = useI18n()
const toast = useToast()
const auth = useAuthStore()
const router = useRouter()

const props = defineProps<{
  avatarUrl?: string
}>()

const avatarPH = '/images/user-base.webp'
const previewUrl = ref('')

const emit = defineEmits<{
  (e: 'uploaded', url: string): void
}>()

const currentAvatar = computed(() => {
  return previewUrl.value || props.avatarUrl || avatarPH
})

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

  // Проверка типа файла
  if (!file.type.startsWith('image/')) {
    toast.error(
      t('profileEdit.avatar.invalidFileType') === 'profileEdit.avatar.invalidFileType'
        ? 'Please select an image file'
        : t('profileEdit.avatar.invalidFileType')
    )
    input.value = ''
    return
  }

  // Проверка размера (максимум 10MB для исходного файла)
  const maxOriginalSize = 10 * 1024 * 1024
  if (file.size > maxOriginalSize) {
    toast.error(
      t('profileEdit.avatar.fileTooLarge') === 'profileEdit.avatar.fileTooLarge'
        ? 'File is too large. Maximum size is 10MB'
        : t('profileEdit.avatar.fileTooLarge')
    )
    input.value = ''
    return
  }

  try {
    // Сжимаем изображение до 2MB
    const compressedFile = await compressImage(file, 2)
    
    const fd = new FormData()
    fd.append('avatar', compressedFile)

    const res = await axios.post('/profile/me/avatar', fd)
    
    // Проверяем успешный ответ от backend
    if (res.data && res.data.avatarUrl) {
      const url = res.data.avatarUrl
      previewUrl.value = url
      emit('uploaded', url)
      if (process.client) {
        window.dispatchEvent(new CustomEvent('avatar-updated', { detail: url }))
      }
      
      // Устанавливаем что TOTP включен после успешной загрузки
      if (auth.mfaEnabled === null) auth.mfaEnabled = true
      
      const successText = t('profileEdit.avatar.uploadSuccess')
      toast.success(
        successText === 'profileEdit.avatar.uploadSuccess'
          ? 'Avatar successfully updated!'
          : successText
      )
    } else {
      toast.error('Failed to update avatar: invalid response from server')
    }
  } catch (e: any) {
    // Проверяем ошибку MFA
    if (e?.response?.status === 403 && e?.response?.data?.error === 'mfa_required') {
      auth.mfaEnabled = false
      toast.error('Для изменения аватара необходимо подключить двухэтапную аутентификацию в профиле')
      await router.push('/profile/edit')
      input.value = ''
      return
    }
    
    // Проверяем, это ошибка сжатия или загрузки
    if (e.message && (e.message.includes('compress') || e.message.includes('canvas') || e.message.includes('Failed to'))) {
      toast.error(
        t('profileEdit.avatar.compressionFailed') === 'profileEdit.avatar.compressionFailed'
          ? 'Failed to process image'
          : t('profileEdit.avatar.compressionFailed')
      )
    } else {
      const failText = t('profileEdit.avatar.uploadFailed')
      toast.error(
        failText === 'profileEdit.avatar.uploadFailed'
          ? 'Failed to upload avatar'
          : failText
      )
    }
  } finally {
    input.value = ''
  }
}
</script>

<template>
  <div class="avatar-box">
    <img
      :src="currentAvatar"
      class="avatar"
      alt="avatar"
    />

    <label class="upload">
      {{ t('profileContent.changePhoto') }}
      <input
        type="file"
        hidden
        accept="image/png,image/jpeg,image/webp"
        @change="onAvatar"
      />
    </label>
  </div>
</template>

<style scoped>
.avatar {
  width: 96px;
  height: 96px;
  border-radius: 18px;
  object-fit: cover;
  display: block;
}
.upload {
  margin-top: 8px;
  display: inline-block;
  cursor: pointer;
}
</style>
