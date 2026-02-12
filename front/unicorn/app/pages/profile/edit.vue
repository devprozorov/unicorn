<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import AvatarUploader from '~/components/panel/user/AvatarUploader.vue'
import { useProfileApi } from '~/services/profileApi'
import { useAuthApi } from '~/services/authApi'
import { useAuthStore } from '~/stores/auth'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'
import { authFetch } from '~/composables/useAuthFetch'
import { locationRegions, findRegionByCountry } from '~/utils/locations'

definePageMeta({
  middleware: ['auth-redirect']
})

/* =======================
 * CORE
 * ======================= */
const router = useRouter()
const profileApi = useProfileApi()
const api = useAuthApi()
const auth = useAuthStore()
const { t } = useI18n()
const requireReauth = ref(false)
type MfaState = 'loading' | 'enabled' | 'required'
const mfaState = computed<MfaState>(() => {
  if (!mfaLoaded.value) return 'loading'
  return backendMfaEnabled.value ? 'enabled' : 'required'
})
const qrSrc = computed(() => {
  if (!otpauthUrl.value) return ''
  return `https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=${encodeURIComponent(otpauthUrl.value)}`
})

const mfaConfigured = computed(() => backendMfaEnabled.value)


/* =======================
 * UI STATE
 * ======================= */
const activeTab = ref<'profile' | 'avatar' | 'security' | '2fa'>('profile')
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const tabOrder = ['profile', 'avatar', 'security', '2fa'] as const
const prevTab = ref<typeof activeTab.value>(activeTab.value)
const slideDir = ref<'left' | 'right'>('right')

watch(activeTab, (next, prev) => {
  const nextIndex = tabOrder.indexOf(next)
  const prevIndex = tabOrder.indexOf(prev)

  slideDir.value = nextIndex > prevIndex ? 'right' : 'left'
  prevTab.value = prev
})

watch(activeTab, () => {
  requestAnimationFrame(() => {
    document.querySelector('.settings-page')?.scrollTo({
      top: 0,
      behavior: 'smooth'
    })
  })
})

const mobileMenuOpen = ref(false)

const tabLabels = computed(
  () =>
    ({
      profile: t('profileEdit.tabs.profile'),
      avatar: t('profileEdit.tabs.avatar'),
      security: t('profileEdit.tabs.security'),
      '2fa': t('profileEdit.tabs.twofa')
    }) as Record<typeof activeTab.value, string>
)


/* =======================
 * PROFILE FORM
 * ======================= */
const notCreated = ref(false)

const displayName = ref('')
const about = ref('')
const locationRegion = ref('')
const locationCountry = ref('')
const locationCustom = ref('')
const links = ref('')
const industry = ref('')
const website = ref('')

const backendMfaEnabled = ref(false)
const mfaLoaded = ref(false)

const locationCountryOptions = computed(() => {
  const region = locationRegions.find(r => r.key === locationRegion.value)
  return region?.countries || []
})

function hydrateLocationFields(loc: string) {
  const reg = findRegionByCountry(loc || '')
  if (reg) {
    locationRegion.value = reg.key
    locationCountry.value =
      reg.countries.find(c => c.toLowerCase() === (loc || '').toLowerCase()) || ''
    locationCustom.value = ''
  } else {
    locationRegion.value = ''
    locationCountry.value = ''
    locationCustom.value = loc || ''
  }
}

async function checkUserMfa() {
  await authFetch('/api/resumes/my')
}

async function checkCompanyMfa() {
  await authFetch('/api/vacancies/my')
}

async function loadMfaState() {
  mfaLoaded.value = false

  try {
    if (!auth.user?.type) {
      backendMfaEnabled.value = false
      return
    }

    if (auth.user.type === 'user') {
      await checkUserMfa()
    }

    if (auth.user.type === 'company') {
      await checkCompanyMfa()
    }

    // если запрос прошёл — MFA включена
    backendMfaEnabled.value = true

  } catch (e: any) {
  if (e?.status === 403 && e?.data?.error === 'mfa_required') {
    backendMfaEnabled.value = false
  } else {
    // ⚠️ ВАЖНО: НЕ СЧИТАЕМ "ВЫКЛЮЧЕНА"
    backendMfaEnabled.value = true
  }
}
finally {
    mfaLoaded.value = true
  }
}





onMounted(async () => {
  loading.value = true
  requireReauth.value = false

  try {
    // PROFILE
    const profileRes: any = await profileApi.getMyProfile()
    if (profileRes?.ok && profileRes.profile) {
      const p = profileRes.profile
      displayName.value = p.displayName || ''
      about.value = p.about || ''
      hydrateLocationFields(p.location || '')
      links.value = Array.isArray(p.links) ? p.links.join('\n') : ''
      industry.value = p.industry || ''
      website.value = p.website || ''
    }

    // MFA (ОДИН РАЗ)
    await loadMfaState()
  } finally {
    loading.value = false
  }
})






async function saveProfile() {
  if (saving.value) return

  saving.value = true
  error.value = ''

  try {
    const payload: Record<string, any> = {}

    // === BASIC ===
    const dn = displayName.value.trim()
    if (dn) payload.displayName = dn
    if (about.value.trim()) payload.about = about.value.trim()
    const locValue = locationCountry.value || locationCustom.value.trim()
    if (locValue) payload.location = locValue
    if (industry.value.trim()) payload.industry = industry.value.trim()
    if (website.value.trim()) payload.website = website.value.trim()

    // === LINKS ===
    const linkLines = links.value
      .split(/\r?\n/)
      .map(l => l.trim())
      .filter(Boolean)

    if (linkLines.length) payload.links = linkLines

    // ❗ Нечего сохранять — выходим тихо
    if (Object.keys(payload).length === 0) {
      saving.value = false
      return
    }

    // === API CALL ===
    if (notCreated.value) {
      await profileApi.createProfile(payload)
      notCreated.value = false
    } else {
      await profileApi.updateProfile(payload)
    }

    // Устанавливаем что TOTP включен после успешного сохранения
    if (auth.mfaEnabled === null) auth.mfaEnabled = true

    // === REDIRECT ===
    await router.push(`/profile/${auth.userId}`)

  } catch (e: any) {
    // Проверяем ошибку MFA
    if (e?.response?.status === 403 && e?.response?.data?.error === 'mfa_required') {
      auth.mfaEnabled = false
      error.value = 'Для сохранения изменений необходимо подключить двухэтапную аутентификацию'
      activeTab.value = '2fa'
      saving.value = false
      return
    }
    
    error.value = getErrorMessage(
      e?.response?.data?.error || e?.data?.error || 'save_failed',
      t,
      'profileEdit.errors.saveProfile'
    )
  } finally {
    saving.value = false
  }
}


/* =======================
 * CHANGE PASSWORD
 * ======================= */
const changePwdOpen = ref(false)
const currentPassword = ref('')
const newPassword = ref('')
const repeatPassword = ref('')
const pwdLoading = ref(false)
const pwdError = ref('')
const pwdSuccess = ref(false)

async function changePassword() {
  pwdError.value = ''
  pwdSuccess.value = false

  if (!currentPassword.value.trim()) {
    pwdError.value = t('profileEdit.errors.passwordChangeFailed')
    return
  }

  if (newPassword.value.trim().length < 8) {
    pwdError.value = t('profileEdit.errors.passwordChangeFailed')
    return
  }

  if (newPassword.value !== repeatPassword.value) {
    pwdError.value = t('profileEdit.errors.passwordMismatch')
    return
  }

  pwdLoading.value = true
  try {
    await api.changePassword({
      oldPassword: currentPassword.value,
      newPassword: newPassword.value
    })

    pwdSuccess.value = true
    changePwdOpen.value = false

    currentPassword.value = ''
    newPassword.value = ''
    repeatPassword.value = ''
  } catch (e: any) {
    pwdError.value = getErrorMessage(
      e?.response?.data?.error || 'password_change_failed',
      t,
      'profileEdit.errors.passwordChangeFailed'
    )
  } finally {
    pwdLoading.value = false
  }
}

/* =======================
 * 2FA / TOTP — 1:1 КАК security.vue
 * ======================= */

const enrolling = ref(false)
const otpauthUrl = ref('')
const code = ref('')
const twoFaError = ref('')

/* =======================
 * START ENROLL
 * ======================= */
async function start2FA() {
  if (loading.value || mfaConfigured.value) return

  twoFaError.value = ''
  loading.value = true

  try {
    const res = await api.totpEnroll()

    /**
     * backend говорит: уже включено
     */
    if (res?.ok === false && res?.error === 'already_enabled') {
      backendMfaEnabled.value = true
      enrolling.value = false
      return
    }

    const otpUrl: unknown = res?.otpauthUrl ?? res?.otpauth

    if (res?.ok !== true || typeof otpUrl !== 'string' || otpUrl.length === 0) {
      twoFaError.value = 'MFA_ENROLL_FAILED'
      return
    }

    otpauthUrl.value = otpUrl
    enrolling.value = true
  } catch (e: any) {
    twoFaError.value = getErrorMessage(e?.response?.data?.error || 'MFA_ENROLL_FAILED', t, 'errors.mfaEnrollFailed')
  } finally {
    loading.value = false
  }
}


/* =======================
 * CONFIRM ENROLL
 * ======================= */

async function confirm2FA() {
  if (loading.value) return

  twoFaError.value = ''
  loading.value = true

  try {
    const res = await api.totpEnable(code.value)

    /**
     * backend может вернуть already_enabled — считаем успехом
     */
    if (res?.ok !== true && res?.error !== 'already_enabled') {
      twoFaError.value = getErrorMessage('invalid_code', t, 'auth.errors.invalidCode')
      return
    }

    /**
     * ✅ backend подтвердил 2FA
     * 1. закрываем enroll
     * 2. обновляем ЕДИНСТВЕННЫЙ источник правды
     */
    enrolling.value = false
    backendMfaEnabled.value = true

  } catch (e: any) {
    twoFaError.value = getErrorMessage(e?.response?.data?.error || 'invalid_code', t, 'auth.errors.invalidCode')
  } finally {
    loading.value = false
  }
}


async function forceReauth() {
  requireReauth.value = false
  auth.clear()
  await router.push('/auth')
}

const mfaInitialized = ref(false)

watch(backendMfaEnabled, (enabled, prev) => {
  if (!mfaInitialized.value) {
    mfaInitialized.value = true
    return
  }

  if (enabled && prev === false) {
    requireReauth.value = true
  }
})




/* =======================
 * 2FA DISABLE
 * ======================= */
const disableOpen = ref(false)
const disablePassword = ref('')
const disableCode = ref('')
const disableLoading = ref(false)
const disableError = ref('')

async function disable2FA() {
  if (disableLoading.value) return

  disableError.value = ''
  disableLoading.value = true

  try {
    const res: any = await authFetch('/api/auth/change-password', {
      method: 'POST',
      body: {
        oldPassword: disablePassword.value,
        newPassword: disablePassword.value, // пароль не меняем
        totpCode: disableCode.value
      }
    })

    if (res?.ok !== true) {
      disableError.value = 'INVALID_CREDENTIALS'
      return
    }

    /**
     * ✅ MFA отключена
     */
    backendMfaEnabled.value = false
    requireReauth.value = true


    // reset UI
    disableOpen.value = false
    disablePassword.value = ''
    disableCode.value = ''
  } catch (e: any) {
    disableError.value = getErrorMessage(e?.data?.error || 'INVALID_CODE', t, 'auth.errors.invalidCode')
  } finally {
    disableLoading.value = false
  }
}


</script>

<template>
  <Header />

  <main class="settings-page">
    <div class="settings-layout">
      <!-- SIDEBAR -->
            <aside class="sidebar">
        <h3 class="sidebar-title">{{ t('profileEdit.sidebar.title') }}</h3>

        <div class="sidebar-section">
          <div class="sidebar-section-title">{{ t('profileEdit.sidebar.account') }}</div>

          <button
            class="menu-item"
            :class="{ active: activeTab === 'profile' }"
            @click="activeTab = 'profile'"
          >
            <span class="icon">👨‍🦲</span>
            {{ t('profileEdit.tabs.profile') }}
          </button>

          <button
            class="menu-item"
            :class="{ active: activeTab === 'avatar' }"
            @click="activeTab = 'avatar'"
          >
            <span class="icon">👤</span>
            {{ t('profileEdit.tabs.avatar') }}
          </button>
        </div>

        <div class="sidebar-section">
          <div class="sidebar-section-title">{{ t('profileEdit.sidebar.security') }}</div>

          <button
            class="menu-item"
            :class="{ active: activeTab === 'security' }"
            @click="activeTab = 'security'"
          >
            <span class="icon">📝</span>
            {{ t('profileEdit.tabs.security') }}
          </button>

          <button
            class="menu-item"
            :class="{ active: activeTab === '2fa' }"
            @click="activeTab = '2fa'"
          >
            <span class="icon">🔐</span>
            {{ t('profileEdit.tabs.twofa') }}
          </button>
        </div>
      </aside>
<!-- MOBILE SECTION SELECTOR -->
<div class="mobile-section">
  <button
    class="mobile-select"
    @click="mobileMenuOpen = !mobileMenuOpen"
  >
    <span>{{ tabLabels[activeTab] }}</span>
    <span class="chevron" :class="{ open: mobileMenuOpen }">▾</span>
  </button>

  <transition name="fade-slide">
    <div v-if="mobileMenuOpen" class="mobile-dropdown">
      <button
        v-for="key in tabOrder"
        :key="key"
        class="mobile-option"
        :class="{ active: activeTab === key }"
        @click="
          activeTab = key;
          mobileMenuOpen = false
        "
      >
        {{ tabLabels[key] }}
      </button>
    </div>
  </transition>
</div>


      <!-- CONTENT -->
      <section class="content tab-container">
        <transition
  :name="`tab-${slideDir}`"
  mode="out-in"
>

                    <!-- ================= PROFILE ================= -->
          <div v-if="activeTab === 'profile'" class="card">
            <h2 class="section-title">{{ t('profileEdit.profile.title') }}</h2>

            <form @submit.prevent="saveProfile">
              <div class="form-block">
                <h3 class="block-title">{{ t('profileEdit.profile.basicTitle') }}</h3>

                <label>{{ t('profileEdit.profile.displayName') }}</label>
                <input v-model="displayName" />

                <label>{{ t('profileEdit.profile.about') }}</label>
                <textarea v-model="about" rows="4" />
              </div>

              <div class="form-block">
                <h3 class="block-title">{{ t('profileEdit.profile.contactsTitle') }}</h3>

                <label>{{ t('profileEdit.profile.location') }}</label>
                <div class="location-row">
                  <select v-model="locationRegion">
                    <option value="">{{ t('profileEdit.profile.region') }}</option>
                    <option
                      v-for="r in locationRegions"
                      :key="r.key"
                      :value="r.key"
                    >
                      {{ r.label }}
                    </option>
                    <option value="custom">{{ t('profileEdit.profile.other') }}</option>
                  </select>

                  <select
                    v-model="locationCountry"
                    :disabled="!locationRegion || locationRegion === 'custom'"
                  >
                    <option value="">{{ t('profileEdit.profile.country') }}</option>
                    <option
                      v-for="c in locationCountryOptions"
                      :key="c"
                      :value="c"
                    >
                      {{ c }}
                    </option>
                  </select>
                </div>

                <input
                  v-if="locationRegion === 'custom'"
                  v-model="locationCustom"
                  class="mt-2"
                  :placeholder="t('profileEdit.profile.cityPlaceholder')"
                />

                <label>{{ t('profileEdit.profile.website') }}</label>
                <input v-model="website" />

                <label>{{ t('profileEdit.profile.links') }}</label>
                <textarea v-model="links" rows="3" />
              </div>

              <div class="form-block">
                <h3 class="block-title">{{ t('profileEdit.profile.additionalTitle') }}</h3>

                <label>{{ t('profileEdit.profile.industry') }}</label>
                <input v-model="industry" />
              </div>

              <div class="form-actions">
                <button type="submit" class="primary-btn" :disabled="saving">
                  {{ t('profileEdit.profile.save') }}
                </button>
              </div>
            </form>
          </div>

                    <!-- ================= AVATAR ================= -->
          <div v-else-if="activeTab === 'avatar'" key="avatar" class="card">
            <h2 class="section-title">{{ t('profileEdit.avatar.title') }}</h2>

            <AvatarUploader
              :avatar-url="profile?.avatarUrl"
              @uploaded="(url) => { if (profile) profile.avatarUrl = url }"
            />

            <p class="hint">{{ t('profileEdit.avatar.hint') }}</p>
          </div>


                    <!-- ================= SECURITY ================= -->
          <div v-else-if="activeTab === 'security'" key="security" class="card">
            <h1 class="mobile-title">
              {{ tabLabels[activeTab] }}
            </h1>

            <button class="link-btn" @click="changePwdOpen = !changePwdOpen">
              {{ t('profileEdit.security.changePassword') }}
            </button>

            <transition name="fade-slide">
              <div v-if="changePwdOpen" class="password-box">
                <input
                  type="password"
                  v-model="currentPassword"
                  :placeholder="t('profileEdit.security.currentPassword')"
                />

                <input
                  type="password"
                  v-model="newPassword"
                  :placeholder="t('profileEdit.security.newPassword')"
                />

                <input
                  type="password"
                  v-model="repeatPassword"
                  :placeholder="t('profileEdit.security.repeatPassword')"
                />

                <button
                  class="primary-btn"
                  :disabled="pwdLoading"
                  @click="changePassword"
                >
                  {{ t('profileEdit.security.confirm') }}
                </button>

                <p v-if="pwdError" class="error">{{ pwdError }}</p>
                <p v-if="pwdSuccess" class="success">
                  {{ t('profileEdit.security.passwordChanged') }}
                </p>
              </div>
            </transition>
          </div>

                    <!-- ================= 2FA ================= -->
          <div v-else key="2fa" class="card">
            <h2 class="section-title">{{ t('profileEdit.twofa.title') }}</h2>

            <ClientOnly>
              <template v-if="mfaState === 'loading'">
              </template>

              <template v-else-if="enrolling">
                <div class="qr-box">
                  <img :src="qrSrc" alt="QR" />
                </div>

                <input
                  v-model="code"
                  class="input"
                  maxlength="6"
                  inputmode="numeric"
                  placeholder="123456"
                />

                <button
                  class="primary-btn"
                  :disabled="loading || code.length !== 6"
                  @click="confirm2FA"
                >
                  {{ t('profileEdit.twofa.confirm') }}
                </button>
              </template>

              <template v-else-if="mfaState === 'enabled'">
                <div class="twofa-wrapper">
                  <div class="twofa-card enabled">
                    <div class="twofa-line" />
                    <div class="twofa-bg" />

                    <div class="twofa-content">
                      <div class="twofa-title">
                        {{ t('profileEdit.twofa.enabledTitle') }}
                      </div>
                      <div class="twofa-subtitle">
                        {{ t('profileEdit.twofa.enabledSubtitle') }}
                      </div>
                    </div>

                    <div class="twofa-action active">
                      <img src="/images/security.png" />
                    </div>
                  </div>
                </div>
              </template>

              <template v-else>
                <div class="twofa-wrapper">
                  <div class="twofa-card">
                    <div class="twofa-line" />
                    <div class="twofa-bg" />

                    <div class="twofa-content">
                      <div class="twofa-title">
                        {{ t('profileEdit.twofa.disabledTitle') }}
                      </div>
                      <div class="twofa-subtitle">
                        {{ t('profileEdit.twofa.disabledSubtitle') }}
                      </div>
                    </div>

                    <button class="twofa-action" @click="start2FA">
                      <img src="/images/key.png" />
                    </button>
                  </div>
                </div>
              </template>
            </ClientOnly>

          </div>
        </transition>
      </section>
    </div>
  </main>

  <Footer />
</template>

<style scoped>
.settings-page {
  background: #f5f6f8;
  padding: 32px 16px;
  min-height: 100vh;
}
.settings-layout {
  max-width: 1100px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 24px;
}
.sidebar {
  background: #fff;
  padding: 20px;
  border-radius: 14px;
}
.menu-item {
  padding: 10px 12px;
  border-radius: 8px;
  border: none;
  background: transparent;
  cursor: pointer;
}
.menu-item.active {
  background: #2563eb;
  color: #fff;
}
.card {
  background: #fff;
  border-radius: 16px;
  padding: 28px;
}
.avatar-preview {
  width: 96px;
  height: 96px;
  border-radius: 18px;
}
.twofa-card {
  position: relative;
  display: flex;
  align-items: center;
  padding: 24px;
  border-radius: 20px;
  overflow: hidden;
}

.twofa-card.enabled .line {
  background: #0FE44C;
}
.twofa-bg {
  position: absolute;
  inset: 0;
  background: url('/images/2fa.png') center / 160% no-repeat;
  opacity: 0.06;
}
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all .25s ease;
}
.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

/* wrapper для центровки */
.twofa-wrapper {
  display: flex;
  justify-content: flex-start;
}

/* ОСНОВНОЙ БЛОК */
.twofa-card {
  position: relative;
  width: 90%;
  height: 80px;

  display: flex;
  align-items: center;
  gap: 16px;

  padding: 0 20px;
  border-radius: 24px;

  background: #ffffff;
  overflow: hidden;

  box-shadow: 0 6px 20px rgba(0,0,0,0.06);
}

/* ЛЕВАЯ ЛИНИЯ */
.twofa-line {
  width: 4px;
  height: 48px;
  border-radius: 4px;
  background: #d1d5db;
  z-index: 2;
}

.twofa-card.enabled .twofa-line {
  background: #0FE44C;
}

/* ФОН */
.twofa-bg {
  position: absolute;
  inset: 0;
  background: url('/images/2fa.png') center right / 45% no-repeat;
  opacity: 0.08;
  pointer-events: none;
}

/* ТЕКСТ */
.twofa-content {
  position: relative;
  flex: 1;
  z-index: 2;
}

.twofa-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  line-height: 1.2;
}

.twofa-subtitle {
  font-size: 12px;
  color: #6b7280;
  margin-top: 2px;
}

/* КНОПКА */
.twofa-action {
  z-index: 2;

  width: 44px;
  height: 44px;
  border-radius: 50%;

  display: flex;
  align-items: center;
  justify-content: center;

  background: #f3f4f6;
  border: none;
  cursor: pointer;

  transition:
    background 0.15s ease,
    transform 0.1s ease,
    box-shadow 0.15s ease;
}

.twofa-action img {
  width: 20px;
  height: 20px;
}

/* hover / focus */
.twofa-action:hover:not(:disabled) {
  background: #e5e7eb;
}

.twofa-action:active:not(:disabled) {
  transform: scale(0.96);
}

.twofa-action:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

/* АКТИВНОЕ СОСТОЯНИЕ */
.twofa-action.active {
  background: #e6fbea;
  cursor: default;
}

.sidebar {
  background: #fff;
  padding: 20px;
  border-radius: 16px;
}

.sidebar-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 16px;
}

.sidebar-section {
  margin-bottom: 20px;
}

.sidebar-section-title {
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  margin-bottom: 8px;
  text-transform: uppercase;
}

.menu-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;

  padding: 10px 12px;
  border-radius: 10px;
  border: none;
  background: transparent;
  cursor: pointer;

  font-size: 14px;
  color: #111827;
}

.menu-item:hover {
  background: #f3f4f6;
}

.menu-item.active {
  background: #2563eb;
  color: #fff;
}

.menu-item .icon {
  width: 18px;
  text-align: center;
}

.form-block {
  margin-bottom: 28px;
}

.block-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 12px;
}

label {
  display: block;
  font-size: 13px;
  color: #374151;
  margin-bottom: 4px;
}

input,
textarea {
  width: 100%;
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  margin-bottom: 12px;
}

.location-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  margin-bottom: 8px;
}

.location-row select {
  width: 100%;
}

input:focus,
textarea:focus {
  outline: none;
  border-color: #2563eb;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}
.settings-page {
  background: #f5f6f8;
  padding: 32px 16px;
  min-height: 100vh;
}

.settings-layout {
  max-width: 1100px;
  margin: 0 auto;

  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 24px;
}

.card {
  background: #fff;
  border-radius: 20px;
  padding: 32px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 24px;
}

.danger-btn {
  margin-top: 12px;
  padding: 10px 16px;
  border-radius: 10px;
  border: none;
  background: #ef4444;
  color: #fff;
  font-weight: 500;
  cursor: pointer;
}

.danger-btn:hover {
  background: #dc2626;
}

.disable-box {
  margin-top: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.reauth-box {
  margin-top: 16px;
  padding: 16px;
  border-radius: 12px;
  background: #f8fafc;
  border: 1px solid #e5e7eb;
}

.reauth-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.reauth-text {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 12px;
}

/* =========================
 * MOBILE SETTINGS PAGE
 * ========================= */
@media (max-width: 900px) {
  .settings-page {
    padding: 16px 12px 32px;
  }

  .settings-layout {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  /* SIDEBAR → TOP TABS */
  .sidebar {
    display: flex;
    flex-direction: column;
    padding: 12px;
  }

  .sidebar-title {
    font-size: 16px;
    margin-bottom: 10px;
  }

  .sidebar-section {
    margin-bottom: 10px;
  }

  .sidebar-section-title {
    font-size: 11px;
    margin-bottom: 6px;
  }

  /* HORIZONTAL MENU */
  .sidebar-section {
    display: flex;
    gap: 8px;
    overflow-x: auto;
    padding-bottom: 6px;
  }

  .menu-item {
    flex-shrink: 0;
    padding: 10px 14px;
    border-radius: 999px;
    background: #f3f4f6;
    font-size: 13px;
  }

  .menu-item.active {
    background: #2563eb;
    color: #fff;
  }

  /* CONTENT */
  .content {
    padding: 0;
  }

  .card {
    padding: 18px;
    border-radius: 16px;
  }

  .section-title {
    font-size: 18px;
    margin-bottom: 16px;
  }

  /* FORM */
  .form-block {
    margin-bottom: 20px;
  }

  input,
  textarea {
    font-size: 14px;
    padding: 12px;
  }

  .form-actions {
    justify-content: stretch;
  }

  .primary-btn {
    width: 100%;
    padding: 14px;
    border-radius: 14px;
  }
}
@media (max-width: 420px) {
  .card {
    padding: 14px;
  }

  .menu-item {
    font-size: 12px;
    padding: 8px 12px;
  }

  .section-title {
    font-size: 16px;
  }
}
/* =========================
 * TAB SWITCH — NATIVE FEEL
 * ========================= */

/* RIGHT (вперёд) */
.tab-right-enter-active,
.tab-right-leave-active {
  transition:
    transform 0.28s cubic-bezier(.4,0,.2,1),
    opacity 0.22s ease;
}

.tab-right-enter-from {
  opacity: 0;
  transform: translateX(24px);
}

.tab-right-leave-to {
  opacity: 0;
  transform: translateX(-16px);
}

/* LEFT (назад) */
.tab-left-enter-active,
.tab-left-leave-active {
  transition:
    transform 0.28s cubic-bezier(.4,0,.2,1),
    opacity 0.22s ease;
}

.tab-left-enter-from {
  opacity: 0;
  transform: translateX(-24px);
}

.tab-left-leave-to {
  opacity: 0;
  transform: translateX(16px);
}
.tab-container {
  position: relative;
  overflow: hidden;
}
@media (max-width: 768px) {
  .tab-right-enter-from,
  .tab-left-enter-from {
    transform: translateX(18px);
  }
}
.menu-item:active {
  transform: scale(0.97);
}
.sidebar {
  position: sticky;
  top: 56px; /* высота Header */
  z-index: 20;

  background: #f5f6f8;
  padding: 8px 0 6px;
  border-radius: 0;
}
@media (max-width: 900px) {
  .sidebar-section {
    display: flex;
    gap: 8px;
    overflow-x: auto;

    padding: 0 12px 6px;
    margin: 0;

    scrollbar-width: none;
  }

  .sidebar-section::-webkit-scrollbar {
    display: none;
  }

  .menu-item {
    flex-shrink: 0;

    padding: 10px 14px;
    border-radius: 999px;

    background: #e5e7eb;
    font-size: 13px;
    font-weight: 500;

    transition:
      background .2s ease,
      color .2s ease,
      transform .1s ease;
  }

  .menu-item.active {
    background: #2563eb;
    color: #fff;
  }
}
@media (max-width: 900px) {
  .card {
    border-radius: 18px;
    box-shadow: 0 8px 24px rgba(0,0,0,.06);
  }
}
@media (max-width: 900px) {
  .settings-page {
    padding-top: 12px;
  }

  .section-title {
    margin-top: 4px;
  }
}
@media (max-width: 900px) {
  label {
    font-size: 12px;
  }

  input,
  textarea {
    font-size: 15px;
    padding: 14px 14px;
  }

  textarea {
    min-height: 96px;
  }
}
@media (max-width: 768px) {
  .tab-right-enter-active,
  .tab-left-enter-active {
    transition:
      transform 0.24s cubic-bezier(.33,1,.68,1),
      opacity 0.18s ease;
  }
}
.menu-item:active {
  transform: scale(0.95);
}
@media (max-width: 900px) {
  .sidebar {
    display: none;
  }

  .location-row {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .settings-page {
    padding: 18px 14px;
  }

  .settings-layout {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .card {
    padding: 18px;
  }

  .twofa-card {
    width: 100%;
  }

  .twofa-wrapper {
    justify-content: center;
  }
}
/* =========================
 * MOBILE SECTION DROPDOWN
 * ========================= */
.mobile-section {
  display: none;
  position: sticky;
  top: 56px; /* header */
  z-index: 30;
  background: #f5f6f8;
  padding: 8px 0 12px;
}

.mobile-select {
  width: 100%;
  height: 48px;

  display: flex;
  align-items: center;
  justify-content: space-between;

  padding: 0 16px;
  border-radius: 14px;

  background: #ffffff;
  border: 1px solid #e5e7eb;

  font-size: 15px;
  font-weight: 600;
}

.chevron {
  transition: transform .2s ease;
}

.chevron.open {
  transform: rotate(180deg);
}

.mobile-dropdown {
  margin-top: 8px;

  background: #ffffff;
  border-radius: 14px;
  overflow: hidden;

  box-shadow: 0 12px 30px rgba(0,0,0,.12);
}

.mobile-option {
  width: 100%;
  text-align: left;

  padding: 14px 16px;
  border: none;
  background: transparent;

  font-size: 14px;
  font-weight: 500;
}

.mobile-option:not(:last-child) {
  border-bottom: 1px solid #f1f5f9;
}

.mobile-option.active {
  background: #2563eb;
  color: #ffffff;
}

@media (max-width: 900px) {
  .mobile-section {
    display: block;
  }
}
.mobile-title {
  font-size: 22px;
  font-weight: 700;
  margin-bottom: 20px;
}

</style>






