<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4">
    <div class="w-full max-w-md space-y-8">
      <h1 class="text-2xl font-semibold text-gray-900 text-center">
        {{ t('admin.login.title') }}
      </h1>
      <form @submit.prevent="onSubmit" class="space-y-4 bg-white p-6 rounded-lg shadow">
        <div>
          <label class="block text-sm font-medium text-gray-700">{{ t('admin.login.login') }}</label>
          <input
            v-model="login"
            type="text"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            required
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">{{ t('admin.login.password') }}</label>
          <input
            v-model="password"
            type="password"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            required
          />
        </div>
        <div>
          <button
            type="submit"
            :disabled="loading"
            class="inline-flex justify-center rounded-md bg-indigo-600 px-4 py-2 text-white hover:bg-indigo-700 disabled:opacity-60 w-full"
          >
            {{ loading ? t('admin.login.loading') : t('admin.login.submit') }}
          </button>
        </div>
        <p v-if="error" class="text-sm text-red-600">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminApi } from '~/composables/useAdminApi'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'
import { useToast } from '~/composables/useToast'

definePageMeta({
  layout: false
})

const router = useRouter()
const toast = useToast()
const { t } = useI18n()
const { login: apiLogin } = useAdminApi()

const login = ref('')
const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)

const onSubmit = async () => {
  error.value = null
  loading.value = true
  try {
    const res = await apiLogin(login.value, password.value)
    if (res?.ok) {
      toast.success(t('admin.login.success'))
      router.push('/admin/users')
    } else {
      error.value = getErrorMessage(res?.error || 'login_failed', t, 'admin.errors.loginFailed')
      toast.error(error.value)
    }
  } catch (e: any) {
    error.value = getErrorMessage(e?.data?.error || e?.message || 'login_failed', t, 'admin.errors.loginFailed')
    toast.error(error.value)
  } finally {
    loading.value = false
  }
}
</script>
