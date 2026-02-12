import { ref } from 'vue'

// Shared toast state used by the global Toast component
export const toastMessage = ref('')
export const toastVisible = ref(false)

export function showToast(msg: string) {
  toastMessage.value = msg
  toastVisible.value = true

  setTimeout(() => {
    toastVisible.value = false
  }, 3000)
}

// Composable used across pages/components
export function useToast() {
  return {
    show: (msg: string) => showToast(msg),
    success: (msg: string) => showToast(msg),
    error: (msg: string) => showToast(msg),
    info: (msg: string) => showToast(msg)
  }
}
