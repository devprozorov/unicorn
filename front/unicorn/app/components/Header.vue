<script setup lang="ts">
import { computed, ref, onBeforeUnmount, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '~/composables/useI18n'
import { useAuthStore } from '~/stores/auth'
import { authFetch } from '~/composables/useAuthFetch'
import logo from '~/images/unicornstar.svg'

const router = useRouter()
const { lang, dict } = useI18n()
const auth = useAuthStore()

/* =====================
 * AUTH STATE
 * ===================== */
const isAuthenticated = computed(() => auth.isAuthenticated)

const displayName = computed(() => {
  if (!auth.user) return ''
  return auth.user.displayName || auth.user.login || ''
})

const profileAvatar = ref('')

const avatarUrl = computed(() => {
  if (profileAvatar.value) return profileAvatar.value
  if (!auth.user) return ''
  const userAvatar =
    (auth.user as any)?.avatarUrl ||
    (auth.user as any)?.photoUrl ||
    ''
  if (userAvatar) return userAvatar
  return auth.user.type === 'company'
    ? '/images/com-base.webp'
    : '/images/user-base.webp'
})
const isUserAccount = computed(() => {
  return auth.isAuthenticated && auth.user?.type === 'user'
})

async function loadProfileAvatar() {
  if (!process.client) return
  if (!auth.isAuthenticated) return
  try {
    const res: any = await authFetch('/api/profile/me')
    profileAvatar.value = res?.profile?.avatarUrl || res?.avatarUrl || ''
  } catch {
    profileAvatar.value = ''
  }
}

function onAvatarUpdated(event: Event) {
  const detail = (event as CustomEvent<string>)?.detail
  if (typeof detail === 'string') {
    profileAvatar.value = detail
  } else {
    loadProfileAvatar()
  }
}

watch(
  () => auth.userId,
  () => {
    loadProfileAvatar()
  },
  { immediate: true }
)

onMounted(() => {
  loadProfileAvatar()
  if (process.client) {
    window.addEventListener('avatar-updated', onAvatarUpdated as EventListener)
  }
})

onBeforeUnmount(() => {
  if (process.client) {
    window.removeEventListener('avatar-updated', onAvatarUpdated as EventListener)
  }
})

/* =====================
 * NAVIGATION
 * ===================== */
function goHome() {
  router.push('/')
}

function goToAuth() {
  router.push('/auth')
}

function gotoJobs() {
  router.push('/jobs')
}

function gotoDoc() {
  router.push('/doc')
}

function goToProfile() {
  if (!auth.userId || !auth.user) {
    auth.clear()
    router.push('/auth')
    return
  }

  if (auth.user.type === 'company') {
    router.push(`/company/${auth.userId}`)
  } else {
    router.push(`/profile/${auth.userId}`)
  }
}

/* =====================
 * MOBILE MENU
 * ===================== */
const mobileMenu = ref(false)

function toggleMenu() {
  mobileMenu.value = !mobileMenu.value
}

function closeMenu() {
  mobileMenu.value = false
}
</script>

<template>
  <header class="header">
    <div class="container">

      <!-- LEFT -->
      <div class="left">
        <button class="logo" @click="goHome">
          <img :src="logo" alt="Unicornstar" class="logo-img" />
          <span class="logo-text">Unicornstar</span>
        </button>
      </div>

      <!-- DESKTOP NAV -->
      <nav class="nav desktop">
        <button class="nav-btn" @click="goHome">
          {{ dict.nav.home }}
        </button>
        <button
            v-if="isUserAccount"
            class="nav-btn"
            @click="gotoJobs"
        >
            {{ dict.nav.candidates }}
        </button>
        <button class="nav-btn" @click="gotoDoc">
          {{ dict.nav.doc }}
        </button>
        <div class="nav-dropdown">
          <button class="nav-btn">
            {{ dict.nav.info }}
          </button>
          <div class="nav-dropdown-menu">
            <NuxtLink class="nav-dropdown-link" to="/contacts">
              {{ dict.nav.contacts }}
            </NuxtLink>
            <NuxtLink class="nav-dropdown-link" to="/career">
              {{ dict.nav.career }}
            </NuxtLink>
          </div>
        </div>
      </nav>

      <!-- RIGHT -->
      <div class="right">
        <!-- LANG -->
        <div class="lang-switch">
          <button
            class="lang-btn"
            :class="{ active: lang === 'ru' }"
            @click="lang = 'ru'"
          >
            RU
          </button>
          <button
            class="lang-btn"
            :class="{ active: lang === 'en' }"
            @click="lang = 'en'"
          >
            EN
          </button>
        </div>

        <!-- AUTH / PROFILE -->
        <template v-if="!isAuthenticated">
          <button class="auth-btn" @click="goToAuth">
            {{ dict.nav.auth }}
          </button>
        </template>

        <template v-else>
          <div class="profile">
            <span
              v-if="displayName"
              class="profile-name"
            >
              {{ displayName }}
            </span>

            <SubscriptionStatus />

            <button class="profile-avatar" @click="goToProfile">
              <img :src="avatarUrl" />
            </button>
          </div>
        </template>

        <!-- BURGER -->
        <button class="burger" @click="toggleMenu">
          ☰
        </button>
      </div>
    </div>

    <!-- MOBILE MENU -->
    <div v-if="mobileMenu" class="mobile-menu">
      <button class="mobile-link" @click="goHome(); closeMenu()">
        {{ dict.nav.home }}
      </button>
      <button
        v-if="isUserAccount"
        class="mobile-link"
        @click="gotoJobs(); closeMenu()"
      >
        {{ dict.nav.candidates }}
      </button>
      <button class="mobile-link" @click="gotoDoc(); closeMenu()">
        {{ dict.nav.doc }}
      </button>
      <button class="mobile-link" @click="router.push('/contacts'); closeMenu()">
        {{ dict.nav.contacts }}
      </button>
      <button class="mobile-link" @click="router.push('/career'); closeMenu()">
        {{ dict.nav.career }}
      </button>
      <div class="mobile-divider" />

      <button
        v-if="!isAuthenticated"
        class="mobile-link primary"
        @click="goToAuth"
      >
        {{ dict.nav.auth }}
      </button>

      <button
        v-else
        class="mobile-link"
        @click="goToProfile"
      >
        {{ dict.nav.profile || 'Profile' }}
      </button>
    </div>
  </header>
</template>

<style scoped>
/* =====================
 * HEADER BASE
 * ===================== */
.header {
  position: sticky;
  top: 0;
  z-index: 50;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
}

.container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 16px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* =====================
 * LOGO
 * ===================== */
.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  font-size: 18px;
}

.logo-img {
  width: 28px;
  height: 28px;
}

/* =====================
 * NAV
 * ===================== */
.nav {
  display: flex;
  gap: 24px;
}

.nav-btn {
  font-size: 14px;
  color: #4b5563;
}

.nav-btn:hover {
  color: #111827;
}

.nav-dropdown {
  position: relative;
}

.nav-dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 160px;
  padding: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.08);
  display: none;
  z-index: 20;
}

.nav-dropdown:hover .nav-dropdown-menu {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.nav-dropdown-link {
  padding: 8px 10px;
  border-radius: 8px;
  color: #374151;
  text-decoration: none;
  font-size: 14px;
}

.nav-dropdown-link:hover {
  background: #f3f4f6;
  color: #111827;
}

/* =====================
 * RIGHT
 * ===================== */
.right {
  display: flex;
  align-items: center;
  gap: 14px;
}

/* =====================
 * LANG
 * ===================== */
.lang-switch {
  display: flex;
  background: #f3f4f6;
  border-radius: 10px;
  padding: 2px;
}

.lang-btn {
  padding: 4px 10px;
  font-size: 12px;
  border-radius: 8px;
}

.lang-btn.active {
  background: #fff;
  font-weight: 600;
}

/* =====================
 * PROFILE
 * ===================== */
.profile {
  display: flex;
  align-items: center;
  gap: 8px;
}

.profile-name {
  max-width: 120px;
  font-size: 14px;
  color: #374151;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.profile-avatar {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.profile-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: inherit;
  display: block;
}

/* =====================
 * BURGER
 * ===================== */
.burger {
  display: none;
  font-size: 20px;
}

/* =====================
 * MOBILE MENU
 * ===================== */
.mobile-menu {
  border-top: 1px solid #e5e7eb;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.mobile-link {
  padding: 10px;
  border-radius: 10px;
  text-align: left;
}

.mobile-link.primary {
  background: #111827;
  color: #fff;
}

.mobile-divider {
  height: 1px;
  background: #e5e7eb;
  margin: 6px 0;
}

/* =====================
 * RESPONSIVE
 * ===================== */
@media (max-width: 768px) {
  .desktop {
    display: none;
  }

  .burger {
    display: block;
  }

  .profile-name {
    display: none;
  }
}
</style>
