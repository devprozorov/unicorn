<script setup lang="ts">
import { ref, computed } from 'vue'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'

import { useAuthApi } from '~/services/authApi'
import { useAuthStore } from '~/stores/auth'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'
import { useProfileApi } from '~/services/profileApi'
import { AuthFetchError } from '~/composables/useAuthFetch'

definePageMeta({
  middleware: [] // без auth middleware
})

const { t } = useI18n()
const errorText = computed(() => (errorCode.value ? getErrorMessage(errorCode.value, t, 'auth.errors.generic') : ''))
const api = useAuthApi()
const auth = useAuthStore()
const profileApi = useProfileApi()

/* MODE */
const isLogin = ref(true)

/* FORM */
const login = ref('')
const password = ref('')
const type = ref<'user' | 'company'>('user')

/* MFA */
const showMfa = ref(false)
const mfaCode = ref('')
const mfaToken = ref<string | null>(null)

/* UI */
const loading = ref(false)
const errorCode = ref('')

/* =========================
   REDIRECT
   ========================= */
async function redirectAfterAuth() {
  await auth.init()

  if (!auth.userId) {
    errorCode.value = 'invalid_session'
    return
  }

  // у тебя в store нет accountType, поэтому ориентируемся по auth.user.type
  const accType = auth.user?.type

  if (accType === 'company') {
    await navigateTo(`/company/${auth.userId}`, { replace: true })
    return
  }

  // ВАЖНО: редирект на реально существующий роут профиля
  await navigateTo(`/profile/${auth.userId}`, { replace: true })
}

/* =========================
   HELPERS
   ========================= */
function normalizeError(e: any) {
  // axios ошибки
  const axiosErr = e?.response?.data?.error
  if (axiosErr) return String(axiosErr)

  // AuthFetchError из authFetch
  if (e instanceof AuthFetchError) {
    // backend у тебя отвечает { error: "bad_request", ok:false }
    if (e?.data?.error) return String(e.data.error)
    return `${e.message}`
  }

  // fallback
  if (e?.message) return String(e.message)
  return 'auth_error'
}

/* =========================
   SUBMIT
   ========================= */
async function submit() {
  errorCode.value = ''
  loading.value = true

  try {
    /* =========================
       LOGIN (STEP 1)
       ========================= */
    if (isLogin.value && !showMfa.value) {
      const res = await api.login(login.value, password.value)

      if (res?.mfaRequired) {
        mfaToken.value = res.mfaToken
        showMfa.value = true
        return
      }

      if (!res?.accessToken) {
        errorCode.value = 'auth_error'
        return
      }

      auth.setAccessToken(res.accessToken)
      await redirectAfterAuth()
      return
    }

    /* =========================
       MFA VERIFY (STEP 2)
       ========================= */
    if (showMfa.value) {
      if (!mfaToken.value) {
        errorCode.value = 'invalid_mfa_state'
        return
      }

      const res = await api.totpVerify({
        mfaToken: mfaToken.value,
        code: mfaCode.value
      })

      if (!res?.accessToken) {
        errorCode.value = 'invalid_code'
        return
      }

      auth.setAccessToken(res.accessToken)
      await redirectAfterAuth()
      return
    }

    /* =========================
       REGISTER
       ========================= */
    const reg = await api.register({
      login: login.value,
      password: password.value,
      displayName: login.value,
      type: type.value
    })

    if (!reg?.ok) {
      errorCode.value = reg?.error || 'auth_error'
      return
    }

    // backend должен вернуть accessToken
    if (!reg?.accessToken) {
      errorCode.value = 'auth_error'
      return
    }

    auth.setAccessToken(reg.accessToken)

    // ВАЖНО: после регистрации делаем ЕДИНОРАЗОВЫЙ POST /api/profile/me
    try {
      await profileApi.createProfile({
        about: '',
        location: '',
        links: [],
        industry: '',
        website: ''
      })
    } catch (e: any) {
      // Если уже есть профиль — 409 нормально
      if (e instanceof AuthFetchError && e.status === 409) {
        // ignore
      } else {
        throw e
      }
    }

    await redirectAfterAuth()
  } catch (e: any) {
    errorCode.value = normalizeError(e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Header />

  <div class="auth-page">
    <div
      class="auth-card"
      :class="{
        'bg-user': !isLogin && type === 'user',
        'bg-company': !isLogin && type === 'company'
      }"
    >
      <h1 class="auth-title">
        {{
          showMfa
            ? t('auth.mfaTitle')
            : isLogin
              ? t('auth.login')
              : t('auth.register')
        }}
      </h1>

      <form @submit.prevent="submit">
        <template v-if="!showMfa">
          <input
            v-model="login"
            class="auth-input"
            :placeholder="t('auth.loginPlaceholder')"
            required
          />
          <input
            v-model="password"
            type="password"
            class="auth-input"
            :placeholder="t('auth.passwordPlaceholder')"
            required
          />

          <div v-if="!isLogin" class="auth-type">
            <div class="type-highlight" :class="type" />

            <button type="button" class="type-btn" :class="{ active: type === 'user' }" @click="type = 'user'">
              {{ t('auth.userType.user') }}
            </button>

            <button type="button" class="type-btn" :class="{ active: type === 'company' }" @click="type = 'company'">
              {{ t('auth.userType.company') }}
            </button>
          </div>
        </template>

        <template v-else>
          <div class="mfa-box">
            <div class="mfa-icon">🔐</div>

            <p class="mfa-text">
              {{ t('auth.mfaText') }}
            </p>

            <input
              v-model="mfaCode"
              class="auth-input mfa-input"
              placeholder="123456"
              maxlength="6"
              inputmode="numeric"
              autofocus
              @input="mfaCode = mfaCode.replace(/\D/g, '')"
              required
            />

            <button
              type="button"
              class="mfa-back"
              @click="showMfa = false; mfaCode = ''; mfaToken = null"
            >
              {{ t('auth.mfaBack') }}
            </button>
          </div>
        </template>

        <button class="auth-btn" :disabled="loading">
          {{
            loading
              ? '...'
              : showMfa
                ? t('auth.confirm')
                : isLogin
                  ? t('auth.login')
                  : t('auth.register')
          }}
        </button>
      </form>

      <button v-if="!showMfa" class="auth-switch" @click="isLogin = !isLogin">
        {{ isLogin ? t('auth.toRegister') : t('auth.toLogin') }}
      </button>

      <p v-if="errorText" class="auth-error">
        {{ errorText }}
      </p>
    </div>
  </div>

  <Footer />
</template>

<style scoped>
/* стили без изменений */
.auth-page {
  min-height: 80dvh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.auth-card {
  position: relative;
  width: 100%;
  max-width: 360px;
  padding: 28px;
  border-radius: 18px;
  background: #fff;
  box-shadow:
    0 10px 30px rgba(0,0,0,.08),
    0 0 0 rgba(99,102,241,0);
  overflow: hidden;
  transition:
    box-shadow 0.35s ease,
    transform 0.35s ease;
}

@media (hover: hover) {
  .auth-card:hover {
    box-shadow:
      0 16px 60px rgba(0,0,0,.12),
      0 0 0 15px rgba(99,102,241,0.02);
    transform: translateY(-2px);
  }
}

.auth-title {
  font-size: 22px;
  font-weight: 600;
  margin-bottom: 16px;
  text-align: center;
}

.auth-input {
  width: 100%;
  padding: 12px;
  margin-bottom: 12px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  text-align: center;
  font-size: 15px;
}

.auth-btn {
  width: 100%;
  padding: 12px;
  border-radius: 10px;
  background: #6366f1;
  color: white;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.auth-btn:active {
  transform: scale(0.98);
}

.auth-switch {
  margin-top: 14px;
  font-size: 13px;
  color: #6366f1;
  background: none;
  border: none;
  cursor: pointer;
}

.auth-error {
  margin-top: 12px;
  font-size: 13px;
  color: #ef4444;
  text-align: center;
}

.auth-type {
  position: relative;
  display: flex;
  gap: 6px;
  margin-bottom: 14px;
  padding: 4px;
  background: #f3f4f6;
  border-radius: 12px;
  overflow: hidden;
}

.type-highlight {
  position: absolute;
  top: 4px;
  bottom: 4px;
  width: calc(50% - 4px);
  background: #6366f1;
  border-radius: 10px;
  z-index: 0;
  transition:
    transform 0.35s cubic-bezier(.34,1.56,.64,1),
    box-shadow 0.35s ease;
}

.type-highlight.user { transform: translateX(0%); }
.type-highlight.company { transform: translateX(100%); }

.type-btn {
  flex: 1;
  padding: 10px 0;
  border-radius: 10px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-weight: 500;
  font-size: 14px;
  color: #374151;
  position: relative;
  z-index: 1;
  transition:
    color 0.25s ease,
    transform 0.2s ease;
}

.type-btn.active { color: white; }
.type-btn:active { transform: scale(0.96); }

.mfa-box { text-align: center; }
.mfa-icon { font-size: 36px; margin-bottom: 12px; }
.mfa-text { font-size: 14px; color: #6b7280; margin-bottom: 16px; }
.mfa-input { font-size: 20px; letter-spacing: 6px; }
.mfa-back {
  margin-top: 12px;
  font-size: 13px;
  color: #6366f1;
  background: none;
  border: none;
  cursor: pointer;
}

.auth-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background-repeat: no-repeat;
  background-position: center 55%;
  background-size: 350px;
  left: -40%;
  bottom: -20%;
  opacity: 0;
  transition: opacity 0.45s ease;
  pointer-events: none;
  z-index: 0;
}

.auth-card.bg-user::before { background-image: url('/images/user.png'); opacity: 0.08; }
.auth-card.bg-company::before { background-image: url('/images/office-building.png'); opacity: 0.08; }

.auth-card > * { position: relative; z-index: 1; }

@media (max-width: 420px) {
  .auth-card { padding: 22px; border-radius: 16px; }
  .auth-title { font-size: 20px; }
  .type-btn { font-size: 13px; }
}
</style>
