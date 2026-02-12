<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useResumesApi } from '~/services/resumesApi'
import { AuthFetchError } from '~/composables/useAuthFetch'
import { showToast } from '~/composables/useToast'
import { useI18n } from '~/composables/useI18n'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
definePageMeta({ middleware: ['auth-redirect'] })

const router = useRouter()
const api = useResumesApi()
const { t } = useI18n()

const creating = ref(false)

const create = async () => {
  if (creating.value) return
  creating.value = true

  try {
    /**
     * ВАЖНО:
     * backend НЕ принимает пустой POST
     * минимум — title
     */
    const res = await api.createResume({
      title: t('resumeCreate.defaultTitle')
    })

    /**
     * Backend возвращает:
     * { ok: true, resumeId: "..." }
     */
    if (!res?.resumeId) {
      throw new Error('resumeId not returned')
    }

    await router.push(`/resumes/${res.resumeId}/edit`)
  } catch (e: any) {
    /**
     * MFA REQUIRED
     */
    if (
      e?.name === 'AuthFetchError' &&
      e?.status === 403 &&
      e?.data?.error === 'mfa_required'
    ) {
      showToast(t('resumeCreate.mfaRequired'))
      await router.push('/profile/edit')
      return
    }

    console.error('[resumes/create] create failed', e)
    showToast(t('resumeCreate.creationFailed'))
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <Header />
  <main class="center">
    <div class="card">
      <h1>{{ t('resumeCreate.title') }}</h1>
      <p>{{ t('resumeCreate.draftInfo') }}</p>

      <button class="primary" :disabled="creating" @click="create">
        {{ creating ? t('resumeCreate.creating') : t('resumeCreate.createResume') }}
      </button>
    </div>
  </main>
  <Footer />
</template>

<style scoped>
.center {
  min-height: 80vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
.card {
  background: #fff;
  padding: 32px;
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0,0,0,.06);
  max-width: 420px;
  width: 100%;
  text-align: center;
}
.primary {
  width: 100%;
  padding: 14px;
  border-radius: 14px;
  background: #4f46e5;
  color: #fff;
  font-weight: 600;
  margin-top: 16px;
}
.primary:disabled {
  opacity: .65;
  cursor: not-allowed;
}
</style>
