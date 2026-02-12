<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthApi } from '~/services/authApi'
import { useAuthStore } from '~/stores/auth'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'

definePageMeta({ middleware: ['panel-auth'] })

const { t } = useI18n()
const errorText = computed(() => (errorCode.value ? getErrorMessage(errorCode.value, t) : ''))
const api = useAuthApi()
const auth = useAuthStore()

const loading = ref(false)
const errorCode = ref('')

const enrolling = ref(false)
const otpauthUrl = ref('')
const code = ref('')

/**
 * 🔑 Единственный источник правды для UI:
 * MFA включена или нет (session scoped, из Pinia)
 */
const mfaConfigured = computed(() => auth.isMfaConfigured)

/* ===============================
 * START TOTP ENROLL
 * =============================== */
async function startEnroll() {
  if (loading.value || mfaConfigured.value) return

  errorCode.value = ''
  loading.value = true

  try {
    const res = await api.totpEnroll()

    /**
     * 🔑 КЛЮЧЕВОЙ КЕЙС:
     * backend вернул already_enabled →
     * значит 2FA УЖЕ ВКЛЮЧЕНА
     */
    if (res?.ok === false && res?.error === 'already_enabled') {
      auth.markMfaConfigured()
      enrolling.value = false
      errorCode.value = ''
      return
    }

    const otpUrl: unknown = res?.otpauthUrl ?? res?.otpauth

    if (res?.ok !== true || typeof otpUrl !== 'string' || otpUrl.length === 0) {
      errorCode.value = 'MFA_ENROLL_FAILED'
      return
    }

    otpauthUrl.value = otpUrl
    enrolling.value = true
  } catch (e: any) {
    errorCode.value = e?.response?.data?.error || 'MFA_ENROLL_FAILED'
  } finally {
    loading.value = false
  }
}

/* ===============================
 * CONFIRM TOTP
 * =============================== */
async function confirmEnroll() {
  if (loading.value) return

  errorCode.value = ''
  loading.value = true

  try {
    const res = await api.totpEnable(code.value)

    if (res?.ok !== true) {
      await navigateTo(`/profile/${auth.userId}`, { replace: true })
      window.location.reload()
    }

    /**
     * ✅ MFA успешно включена
     */
    auth.markMfaConfigured()

    /**
     * backend отзывает refresh → принудительный logout
     */
    auth.clear()
    await navigateTo('/auth')
  } catch (e: any) {
    await navigateTo(`/profile/${auth.userId}`, { replace: true })
    window.location.reload()
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Header />

  <main class="panel">
    <h1 class="title">{{ t('panel.security.title') }}</h1>

    <div class="card">
      <!-- MFA ENABLED -->
      <template v-if="mfaConfigured && !enrolling">
        <p class="status success">
          {{ t('panel.security.enabled') }}
        </p>
        <p class="muted">
          {{ t('panel.security.enabledHint') }}
        </p>
      </template>

      <!-- ENROLL FLOW -->
      <template v-else-if="enrolling">
        <p class="muted">{{ t('panel.security.scanQr') }}</p>

        <div class="qr-box">
          <img
            :src="`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(otpauthUrl)}`"
            alt="TOTP QR Code"
          />
        </div>

        <input
          v-model="code"
          class="input"
          placeholder="123456"
          maxlength="6"
          inputmode="numeric"
          @input="code = code.replace(/\\D/g, '')"
        />

        <button
          class="btn primary"
          :disabled="loading || code.length !== 6"
          @click="confirmEnroll"
        >
          {{ t('panel.security.confirm') }}
        </button>
      </template>

      <!-- MFA NOT CONFIGURED -->
      <template v-else>
        <p class="muted">
          {{ t('panel.security.disabled') }}
        </p>

        <button
          class="btn primary"
          :disabled="loading"
          @click="startEnroll"
        >
          {{ t('panel.security.enable') }}
        </button>
      </template>

      <p v-if="errorText" class="error">{{ errorText }}</p>
    </div>
  </main>

  <Footer />
</template>

<style scoped>
.panel {
  max-width: 720px;
  margin: 40px auto;
  padding: 0 16px;
}

.title {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 24px;
}

.card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0,0,0,.06);
}

.muted {
  color: #6b7280;
  margin-bottom: 16px;
}

.status.success {
  color: #16a34a;
  font-weight: 600;
  margin-bottom: 8px;
}

.qr-box {
  display: flex;
  justify-content: center;
  margin: 16px 0;
}

.input {
  width: 100%;
  padding: 12px;
  font-size: 18px;
  text-align: center;
  letter-spacing: 6px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  margin-bottom: 12px;
}

.btn {
  padding: 10px 16px;
  border-radius: 8px;
  border: none;
  cursor: pointer;
  font-weight: 500;
}

.primary {
  background: #6366f1;
  color: white;
}

.error {
  margin-top: 12px;
  font-size: 13px;
  color: #ef4444;
}
</style>
