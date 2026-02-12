<script setup lang="ts">
import { ref, computed, watch, watchEffect, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useRuntimeConfig } from '#imports'

import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import { useI18n } from '~/composables/useI18n'

import { useAuthStore } from '~/stores/auth'
import { useProfileApi } from '~/services/profileApi'
import { useResumesApi } from '~/services/resumesApi'

/* =========================
 * INIT
 * ========================= */
const route = useRoute()
const router = useRouter()
const config = useRuntimeConfig()

const auth = useAuthStore()
const profileApi = useProfileApi()
const resumesApi = useResumesApi()
const { t } = useI18n()

/**
 * ВАЖНО:
 * Не инициируем auth.init() принудительно с этой страницы.
 * Иначе при любой сетевой ошибке/429/временной проблеме можно получить "вылет".
 * Бейджи/owner-логика просто ждут, пока auth станет готовым штатно.
 */

/* =========================
 * HELPERS: API BASE + SAFE FETCH
 * ========================= */
function getApiBase(): string {
  const pub: any = (config as any)?.public || {}
  return (
    pub.apiBase ||
    pub.apiURL ||
    pub.apiUrl ||
    pub.backendBase ||
    pub.backendURL ||
    pub.backendUrl ||
    ''
  )
}

const apiBase = computed(() => String(getApiBase() || '').trim())

function authHeaders() {
  // если accessToken пуст — просто не шлём Authorization
  const token = (auth as any)?.accessToken
  return token ? { Authorization: `Bearer ${token}` } : undefined
}

function extractArray(res: any): any[] | null {
  if (!res) return null
  if (Array.isArray(res)) return res
  if (Array.isArray(res.items)) return res.items
  if (Array.isArray(res.data)) return res.data
  if (Array.isArray(res.chats)) return res.chats
  if (Array.isArray(res.applications)) return res.applications
  if (Array.isArray(res.results)) return res.results
  return null
}

function getHttpStatus(e: any): number {
  return (
    Number(e?.statusCode) ||
    Number(e?.status) ||
    Number(e?.response?.status) ||
    Number(e?.data?.statusCode) ||
    0
  )
}

async function tryFetchList(paths: string[]): Promise<any[]> {
  let lastErr: any = null
  const baseURL = apiBase.value || ''

  for (const path of paths) {
    try {
      const fullUrl = baseURL ? `${baseURL}${path}` : path
      const res = await fetch(fullUrl, {
        method: 'GET',
        credentials: 'include',
        headers: authHeaders()
      })
      const data: any = await res.json()
      const arr = extractArray(data)
      if (arr) return arr
    } catch (e) {
      lastErr = e
    }
  }
  throw lastErr || new Error('No endpoint matched')
}

/* =========================
 * PROFILE (PUBLIC)
 * ========================= */
const {
  data: profileRes,
  pending: profilePending,
  error: profileError,
  refresh: reloadProfile,
} = await useAsyncData(
  () => `profile:${route.params.userId}`,
  () => profileApi.getPublicProfile(route.params.userId as string),
  { server: false }
)

const profile = computed(() => profileRes.value?.profile ?? profileRes.value ?? null)
const loadError = computed(() => Boolean(profileError.value))

const profileTypeLabel = computed(() => {
  const value = String(profile.value?.type || '').toLowerCase()
  if (value === 'company') return t('panel.company.type')
  if (value) return value
  return t('panel.user.type')
})
/* =========================
 * REDIRECT COMPANY PROFILE
 * ========================= */
watchEffect(() => {
  if (!profile.value) return
  if (String(profile.value.type || '').toLowerCase() === 'company') {
    navigateTo(`/company/${route.params.userId}`)
  }
})

/* =========================
 * OWNER CHECK
 * ========================= */
const isOwner = computed(() => {
  const u = (auth as any)?.user
  if (!u || !profile.value) return false
  return String(u.userId) === String(profile.value.userId)
})

/* =========================
 * AVATAR / LINKS
 * ========================= */
const profileAvatar = computed(() => {
  if (profile.value?.avatarUrl) return profile.value.avatarUrl
  return String(profile.value?.type || '').toLowerCase() === 'company'
    ? '/images/com-base.webp'
    : '/images/user-base.webp'
})

const normalizedLinks = computed(() => {
  const links = profile.value?.links
  return Array.isArray(links) ? links.filter(Boolean) : []
})

const hasLinks = computed(() => normalizedLinks.value.length > 0)

/* =========================
 * BADGES: CHATS + APPLICATIONS (OWNER ONLY)
 * ========================= */
const chatsUnreadTotal = ref(0)
const appUpdatesTotal = ref(0)
const badgesPending = ref(false)
const lastApplications = ref<any[]>([])

/** анти-шторм: не чаще чем раз в 12 сек */
const BADGES_MIN_INTERVAL_MS = 5_000
const lastBadgesAt = ref(0)

/** если словили 429 — не трогаем бейджи N секунд */
const cooldownUntil = ref(0)
let badgeTimer: number | null = null

function formatBadge(n: number) {
  if (!n || n <= 0) return ''
  return n > 99 ? '99+' : String(n)
}

function pickUnreadCount(x: any): number {
  const n =
    x?.unreadCount ??
    x?.unread_count ??
    x?.unread ??
    x?.hasUnread ??
    x?.chat?.unreadCount ??
    x?.chat?.unread_count ??
    0

  if (typeof n === 'boolean') return n ? 1 : 0
  const num = Number(n)
  return Number.isFinite(num) ? num : 0
}

function chatLastSeenKey(uid: string, appId: string) {
  return `unicorn:user:chat:last_seen:${uid || 'anon'}:${appId}`
}

function loadChatLastSeen(uid: string, appId: string) {
  if (!process.client) return 0
  try {
    const raw = localStorage.getItem(chatLastSeenKey(uid, appId))
    const num = Number(raw)
    return Number.isFinite(num) ? num : 0
  } catch {
    return 0
  }
}

function saveChatLastSeen(uid: string, appId: string, ts: number) {
  if (!process.client) return
  try {
    localStorage.setItem(chatLastSeenKey(uid, appId), String(ts || 0))
  } catch {}
}

function pickMessageTimestamp(m: any): number {
  if (!m || typeof m === 'string') return 0
  const v =
    m?.createdAt ||
    m?.created_at ||
    m?.sentAt ||
    m?.sent_at ||
    m?.updatedAt ||
    m?.updated_at ||
    m?.timestamp ||
    ''
  const t = Date.parse(String(v))
  return Number.isNaN(t) ? 0 : t
}

async function fetchLastMessageInfo(appId: string) {
  const baseURL = apiBase.value || ''
  const url = baseURL ? `${baseURL}/api/chat/${appId}/messages` : `/api/chat/${appId}/messages`
  const res = await fetch(url, {
    method: 'GET',
    credentials: 'include',
    headers: authHeaders()
  })
  const data: any = await res.json()
  const items = Array.isArray(data?.items) ? data.items : []
  const last = items.length ? items[items.length - 1] : null
  const ts = pickMessageTimestamp(last)
  const senderType = String(last?.senderType || '')
  return { ts, senderType }
}


function chatBadgeKey(uid: string) {
  return `unicorn:user:chats:unread:${uid || 'anon'}`
}

function loadChatBadge(uid: string) {
  if (!process.client) return 0
  try {
    const raw = localStorage.getItem(chatBadgeKey(uid))
    const num = Number(raw)
    return Number.isFinite(num) ? num : 0
  } catch {
    return 0
  }
}

function saveChatBadge(uid: string, value: number) {
  if (!process.client) return
  try {
    localStorage.setItem(chatBadgeKey(uid), String(value || 0))
  } catch {}
}

/* =========================
 * APPLICATIONS "SEEN" MAP
 * ========================= */
function isFinalCompanyDecision(status?: string) {
  const v = String(status || '').toLowerCase()
  return (
    v.includes('accept') ||
    v.includes('approved') ||
    v.includes('hire') ||
    v.includes('offer') ||
    v.includes('reject') ||
    v.includes('decline') ||
    v.includes('denied')
  )
}

function seenKey(uid: string) {
  return `unicorn:user:applications:seen_status:${uid || 'anon'}`
}

function loadSeenMap(uid: string): Record<string, string> {
  if (!process.client) return {}
  try {
    return JSON.parse(localStorage.getItem(seenKey(uid)) || '{}') || {}
  } catch {
    return {}
  }
}

function saveSeenMap(uid: string, map: Record<string, string>) {
  if (!process.client) return
  try {
    localStorage.setItem(seenKey(uid), JSON.stringify(map))
  } catch {}
}

function markAllFinalDecisionsAsSeen() {
  if (!process.client) return
  const uid = String((auth as any)?.user?.userId || 'anon')
  const seen = loadSeenMap(uid)

  for (const a of lastApplications.value) {
    const id = String(a?.id || a?.applicationId || a?._id || '')
    const st = String(a?.status || a?.state || '')
    if (!id) continue
    if (isFinalCompanyDecision(st)) seen[id] = st
  }

  saveSeenMap(uid, seen)
  appUpdatesTotal.value = 0
}

function onApplicationsClick() {
  // важно: не мешаем навигации — только фиксируем "seen"
  markAllFinalDecisionsAsSeen()
}

/* =========================
 * LOAD BADGES (SAFE)
 * ========================= */
async function loadUserBadges(opts?: { force?: boolean }) {
  if (!isOwner.value) return
  if (!process.client) return
  if (badgesPending.value) return

  const now = Date.now()
  if (cooldownUntil.value && now < cooldownUntil.value) return

  const force = Boolean(opts?.force)
  if (!force && now - lastBadgesAt.value < BADGES_MIN_INTERVAL_MS) return

  badgesPending.value = true
  lastBadgesAt.value = now

  const uid = String((auth as any)?.user?.userId || 'anon')

  try {
    // 1) CHATS unread (by last message)
    let chatsUnreadSum = 0
    try {
      const chats = await tryFetchList([
        '/api/user/chats',
        '/api/chats/my',
        '/api/chat/my',
        '/api/user/chat',
        '/api/chats/user',
      ])

      const uid = String((auth as any)?.user?.userId || 'anon')
      const list = Array.isArray(chats) ? chats.slice(0, 20) : []

      for (const c of list) {
        const appId = String(c?.applicationId || c?.id || c?._id || '')
        if (!appId) continue
        try {
          const info = await fetchLastMessageInfo(appId)
          if (!info.ts) continue
          if (info.senderType !== 'company') continue
          const lastSeen = loadChatLastSeen(uid, appId)
          if (info.ts > lastSeen) chatsUnreadSum += 1
        } catch {}
      }
    } catch (e: any) {
      const st = getHttpStatus(e)
      if (st === 429) {
        cooldownUntil.value = Date.now() + 20_000
      }
      chatsUnreadSum = 0
    }

    // 2) APPLICATIONS: final decision not seen yet
    let updates = 0
    try {
      const apps = await tryFetchList([
        '/api/user/applications',
        '/api/applications/my',
        '/api/applications/me',
        '/api/user/me/applications',
      ])

      lastApplications.value = Array.isArray(apps) ? apps : []

      const uid = String((auth as any)?.user?.userId || 'anon')
      const seen = loadSeenMap(uid)

      for (const a of lastApplications.value) {
        const id = String(a?.id || a?.applicationId || a?._id || '')
        const st = String(a?.status || a?.state || '')
        if (!id) continue
        if (isFinalCompanyDecision(st) && seen[id] !== st) updates += 1
      }
    } catch (e: any) {
      const st = getHttpStatus(e)
      if (st === 429) {
        cooldownUntil.value = Date.now() + 20_000
      }
      updates = 0
      lastApplications.value = []
    }

    chatsUnreadTotal.value = chatsUnreadSum
    saveChatBadge(uid, chatsUnreadSum)
    appUpdatesTotal.value = updates
  } catch (e: any) {
    const st = getHttpStatus(e)
    if (st === 429) cooldownUntil.value = Date.now() + 20_000
    chatsUnreadTotal.value = 0
    saveChatBadge(uid, 0)
    appUpdatesTotal.value = 0
    lastApplications.value = []
  } finally {
    badgesPending.value = false
  }
}

/**
 * Запускаем бейджи:
 * - когда isOwner становится true
 * - на focus/visibility (но с throttle/cooldown)
 */
watch(
  () => isOwner.value,
  (v) => {
    if (v) loadUserBadges({ force: true })
  },
  { immediate: true }
)

const onFocus = () => loadUserBadges()
const onVis = () => {
  if (document.visibilityState === 'visible') loadUserBadges()
}

onMounted(() => {
  if (!process.client) return
  const uid = String((auth as any)?.user?.userId || 'anon')
  chatsUnreadTotal.value = loadChatBadge(uid)
  window.addEventListener('focus', onFocus)
  document.addEventListener('visibilitychange', onVis)
  badgeTimer = window.setInterval(() => {
    if (document.visibilityState === 'visible') loadUserBadges()
  }, 10_000)
})

onBeforeUnmount(() => {
  if (!process.client) return
  window.removeEventListener('focus', onFocus)
  document.removeEventListener('visibilitychange', onVis)
  if (badgeTimer) window.clearInterval(badgeTimer)
})

/* =========================
 * RESUMES (OWNER ONLY)
 * ========================= */
const {
  data: resumesRes,
  pending: resumesPending,
  error: resumesError,
  refresh: reloadResumes,
} = await useAsyncData(
  () => `my-resumes:${String((auth as any)?.user?.userId || 'anon')}:${String(route.params.userId || '')}`,
  async () => {
    if (!isOwner.value) return []
    return await resumesApi.getMyResumes()
  },
  { server: false, watch: [isOwner] }
)

const resumes = computed(() => resumesRes.value ?? [])

/* =========================
 * DELETE MODAL STATE
 * ========================= */
const showDeleteModal = ref(false)
const resumeToDelete = ref<string | null>(null)
const deletingId = ref<string | null>(null)
const showLogoutModal = ref(false)

function askDeleteResume(id: string) {
  if (!isOwner.value) return
  resumeToDelete.value = id
  showDeleteModal.value = true
}

function cancelDeleteResume() {
  showDeleteModal.value = false
  resumeToDelete.value = null
}

function askLogout() {
  if (!isOwner.value) return
  showLogoutModal.value = true
}

function cancelLogout() {
  showLogoutModal.value = false
}

async function confirmLogout() {
  showLogoutModal.value = false
  auth.clear()
  await router.push('/auth')
}

async function confirmDeleteResume() {
  if (!resumeToDelete.value) return

  try {
    deletingId.value = resumeToDelete.value
    await resumesApi.deleteResume(resumeToDelete.value)
    await reloadResumes()
  } catch (e) {
    console.error('[DELETE RESUME FAILED]', e)
  } finally {
    deletingId.value = null
    resumeToDelete.value = null
    showDeleteModal.value = false
  }
}

/* =========================
 * COMPLETION %
 * ========================= */
function computeCompletion(r: any): number {
  const fields = ['about', 'skills', 'links']
  const filled = fields.filter((f) => {
    const v = r?.[f]
    if (Array.isArray(v)) return v.length > 0
    return Boolean(v && String(v).trim())
  }).length
  return Math.round((filled / fields.length) * 100)
}
</script>

<template>
  <!-- DELETE RESUME MODAL -->
  <div v-if="showDeleteModal" class="modal-backdrop">
    <div class="modal-card">
      <h3 class="modal-title">{{ t('profilePage.deleteResumeTitle') }}</h3>

      <p class="modal-text">{{ t('profilePage.deleteResumeText') }}</p>

      <div class="modal-actions">
        <button class="btn" @click="cancelDeleteResume">{{ t('profilePage.cancel') }}</button>

        <button
          class="btn danger"
          :disabled="deletingId !== null"
          @click="confirmDeleteResume"
        >
          {{ deletingId ? t('profilePage.deleting') : t('profilePage.delete') }}
        </button>
      </div>
    </div>
  </div>

  <!-- LOGOUT MODAL -->
  <div v-if="showLogoutModal" class="modal-backdrop">
    <div class="modal-card">
      <h3 class="modal-title">{{ t('panel.common.logoutConfirmTitle') }}</h3>

      <p class="modal-text">{{ t('panel.common.logoutConfirmText') }}</p>

      <div class="modal-actions">
        <button class="btn" @click="cancelLogout">{{ t('panel.common.logoutConfirmCancel') }}</button>

        <button class="btn danger" @click="confirmLogout">
          {{ t('panel.common.logoutConfirmOk') }}
        </button>
      </div>
    </div>
  </div>

  <div class="page">
    <Header />

    <section class="wrap">
      <!-- SUBSCRIPTION BANNER -->
      <div class="mb-6 px-4 md:px-6">
        <SubscriptionBanner />
      </div>
      <!-- HERO / заголовок страницы -->
      <div class="hero">
        <div class="hero-left">
          <h1 class="hero-title">{{ t('profilePage.userPanelTitle') }}</h1>
          <p class="hero-sub">{{ t('profilePage.userPanelSub') }}</p>
        </div>

        <!-- Owner only quick actions -->
        <div v-if="isOwner" class="hero-actions">
          <NuxtLink class="hero-btn primary" to="/profile/edit">{{ t('profilePage.editProfile') }}</NuxtLink>
          <button class="hero-btn danger" type="button" @click="askLogout">
            {{ t('panel.common.logout') }}
          </button>
        </div>

        <!-- Guest CTA -->
        <div v-else class="hero-actions">
          <NuxtLink class="hero-btn primary" to="/auth">{{ t('profilePage.loginToInteract') }}</NuxtLink>
          <div class="hero-hint">{{ t('profilePage.ownerOnly') }}</div>
        </div>
      </div>

      <!-- Состояния: loading / error -->
      <div v-if="profilePending" class="state">
        <div class="skeleton-grid">
          <div class="sk card sk-profile" />
          <div class="sk card sk-main" />
          <div class="sk card sk-actions" />
        </div>
      </div>

      <div v-else-if="loadError" class="state">
        <div class="error-card">
          <div class="error-title">{{ t('profilePage.loadErrorTitle') }}</div>
          <div class="error-text">{{ t('profilePage.loadErrorText') }}</div>
          <div class="error-actions">
            <button class="hero-btn primary" @click="reloadProfile()">{{ t('profilePage.retry') }}</button>
            <NuxtLink class="hero-btn" to="/">{{ t('profilePage.toHome') }}</NuxtLink>
          </div>
        </div>
      </div>

      <!-- Основной контент -->
      <div v-else class="layout">
        <!-- LEFT: profile card -->
        <aside class="col profile-col">
          <div class="card profile-card">
            <div class="profile-top">
              <div class="avatar-wrap">
                <img class="avatar" :src="profileAvatar" alt="avatar" />
                <div
                  v-if="(auth as any)?.mfaEnabled === true"
                  class="mfa-dot"
                  :title="t('profilePage.twoFaEnabled')"
                />
              </div>

              <div class="profile-meta">
                <div class="name-row">
                  <div class="display-name">
                    {{ profile?.displayName || t('profilePage.user') }}
                  </div>

                  <span class="badge badge-mfa" v-if="(auth as any)?.mfaEnabled === true">
                    {{ t('profilePage.twoFaEnabled') }}
                  </span>
                  <span class="badge badge-warn" v-else-if="(auth as any)?.mfaEnabled === false">{{ t('profilePage.twoFaDisabled') }}</span>
                </div>

                <div class="minor">
                  <span class="muted">
                    {{ profile?.location || t('profilePage.locationNotSpecified') }}
                  </span>
                  <span class="dot" />
                  <span class="muted">
                    {{ profileTypeLabel }}
                  </span>
                </div>
              </div>
            </div>

            <div class="profile-about">
              <div class="section-title">{{ t('profilePage.links') }}</div>
              <div class="about-text">
                {{ profile?.about || t('profilePage.notSpecified') }}
              </div>
            </div>

            <div class="profile-links" v-if="hasLinks">
              <div class="section-title">{{ t('profilePage.links') }}</div>
              <ul class="links">
                <li v-for="(l, i) in normalizedLinks" :key="i">
                  <a :href="l" target="_blank" rel="noopener noreferrer">{{ l }}</a>
                </li>
              </ul>
            </div>

            <div class="profile-footer" v-if="!isOwner">
              <NuxtLink class="btn wide" to="/auth">{{ t('profilePage.loginToWrite') }}</NuxtLink>
            </div>
          </div>

          <!-- Owner-only note card -->
          <div v-if="isOwner && (auth as any)?.mfaEnabled === false" class="card note-card">
            <div class="note-title">{{ t('profilePage.hint') }}</div>
            <div class="note-text">{{ t('profilePage.hintText') }}</div>
            <NuxtLink class="btn small" to="/profile/edit">
                  {{ t('profilePage.enable2fa') }}
                </NuxtLink>
          </div>
        </aside>

        <!-- CENTER: resumes -->
                <main class="col main-col">
          <div class="card main-card">
            <div class="main-head">
              <div class="main-title">{{ t('profilePage.resumes') }}</div>
              <div class="main-tools">
                <NuxtLink v-if="isOwner" class="btn small primary" to="/resumes/create">
                  {{ t('profilePage.createResume') }}
                </NuxtLink>
              </div>
            </div>

            <div v-if="resumesPending" class="mini-skeleton">
              <div class="sk-line" />
              <div class="sk-line" />
              <div class="sk-line" />
            </div>

            <div v-else-if="resumesError">
              <div class="empty">
                <div class="empty-title">{{ t('profilePage.resumesUnavailableTitle') }}</div>
                <div class="empty-text">{{ t('profilePage.resumesUnavailableText') }}</div>
                <NuxtLink class="btn small" to="/profile/edit">
                  {{ t('profilePage.enable2fa') }}
                </NuxtLink>
              </div>
            </div>

            <div v-else>
              <div v-if="!resumes?.length" class="empty">
                <div class="empty-title">{{ t('profilePage.noResumesTitle') }}</div>
                <div class="empty-text">{{ t('profilePage.noResumesText') }}</div>
                <NuxtLink class="btn primary" to="/resumes/create">
                  {{ t('profilePage.createResumeButton') }}
                </NuxtLink>
              </div>

              <div v-else class="resume-list">
                <div v-for="r in resumes" :key="r.id || r._id" class="resume-card">
                  <div class="resume-top">
                    <div class="resume-title">{{ r.title || t('profilePage.untitled') }}</div>
                    <span
                      class="badge"
                      :class="r.status === 'active' ? 'badge-pub' : 'badge-draft'"
                    >
                      {{ r.status === 'active' ? t('resumeView.active') : t('resumeView.draft') }}
                    </span>
                  </div>

                  <div class="progress">
                    <div class="progress-bg">
                      <div class="progress-bar" :style="{ width: computeCompletion(r) + '%' }" />
                    </div>
                    <div class="progress-label">{{ computeCompletion(r) }}%</div>
                  </div>

                  <div class="resume-actions">
                    <NuxtLink class="link view" :to="`/resumes/${r.id || r._id}`">
                      {{ t('profilePage.view') }}
                    </NuxtLink>
                    <NuxtLink
                      v-if="isOwner"
                      class="link edit"
                      :to="`/resumes/${r.id || r._id}/edit`"
                    >
                      {{ t('profilePage.editResume') }}
                    </NuxtLink>
                    <button
                      v-if="isOwner"
                      class="link danger"
                      @click="askDeleteResume(r.id || r._id)"
                    >
                      {{ t('profilePage.deleteResume') }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </main>


        <!-- RIGHT: actions (desktop only) -->
        <aside class="col actions-col" v-if="isOwner">
          <div class="card actions-card">
            <NuxtLink class="action-btn primary" to="/user/chat">
              <span class="ic">💬</span>
              <span>{{ t('profilePage.chats') }}</span>
              <span v-if="chatsUnreadTotal > 0" class="count-badge">{{ formatBadge(chatsUnreadTotal) }}</span>
            </NuxtLink>

            <NuxtLink class="action-btn" to="/user/applications" @click="onApplicationsClick">
              <span class="ic">📄</span>
              <span>{{ t('profilePage.applications') }}</span>
              <span v-if="appUpdatesTotal > 0" class="count-badge">{{ formatBadge(appUpdatesTotal) }}</span>
            </NuxtLink>

          
          </div>

          <div class="card actions-card subtle">
            <div class="sub-title">{{ t('profilePage.status') }}</div>
            <div class="sub-row">
              <span class="muted">2FA</span>
              <span class="pill pill-ok" v-if="(auth as any)?.mfaEnabled === true">{{ t('profilePage.twoFaEnabled') }}</span>
              <span class="pill pill-warn" v-else-if="(auth as any)?.mfaEnabled === false">{{ t('profilePage.twoFaDisabled') }}</span>
              <span class="pill" v-else>{{ t('profilePage.twoFaChecking') }}</span>
            </div>
            <div class="sub-hint">{{ t('profilePage.twoFaHint') }}</div>
          </div>
        </aside>
      </div>
    </section>

    <!-- Mobile bottom bar (owner only) -->
    <nav v-if="isOwner" class="bottom-bar">
      <NuxtLink to="/user/chat" class="bb-item">
        <span class="bb-ic">💬</span>
        <span class="bb-t">{{ t('profilePage.mobileChats') }}</span>
        <span v-if="chatsUnreadTotal > 0" class="count-badge">{{ formatBadge(chatsUnreadTotal) }}</span>
      </NuxtLink>

      <NuxtLink to="/user/applications" class="bb-item" @click="onApplicationsClick">
        <span class="bb-ic">📄</span>
        <span class="bb-t">{{ t('profilePage.mobileApplications') }}</span>
        <span v-if="appUpdatesTotal > 0" class="count-badge">{{ formatBadge(appUpdatesTotal) }}</span>
      </NuxtLink>

      <NuxtLink to="/profile/edit" class="bb-item">
        <span class="bb-ic">✏️</span>
        <span class="bb-t">{{ t('profilePage.mobileProfile') }}</span>
      </NuxtLink>
    </nav>

    <Footer />
  </div>
</template>





<style scoped>
/* ====== Unicornstar vibe: clean + gradients + depth ====== */
.page {
  min-height: 100vh;
  background:
    radial-gradient(1200px 600px at 20% -10%, rgba(79,70,229,0.18), transparent 60%),
    radial-gradient(900px 500px at 90% 10%, rgba(99,102,241,0.14), transparent 60%),
    #ffffff;
}

.wrap {
  max-width: 1320px;
  margin: 0 auto;
  padding: 28px 16px 110px;
}

/* HERO */
.hero {
  position: relative;
  overflow: hidden;

  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  padding: 28px 28px;
  border-radius: 20px;
  margin-bottom: 15px;

  background:
    radial-gradient(
      600px 260px at 85% 30%,
      rgba(99,102,241,0.18),
      transparent 60%
    ),
    linear-gradient(
      135deg,
      rgba(79,70,229,0.10),
      rgba(99,102,241,0.06)
    );

  border: 1px solid rgba(99,102,241,0.15);
  box-shadow: 0 12px 40px rgba(17,24,39,.06);
}

.hero::after {
  content: '';
  position: absolute;

  left: -45px;
  top: 50%;
  transform: translateY(-52%);

  width: 1200px;
  height: 1200px;

  background-image: url('/images/user_panel.webp');
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;

  opacity: 0.42;

  pointer-events: none; /* не мешает кликам */
  z-index: 0;
}

.hero-left,
.hero-actions {
  position: relative;
  z-index: 1;
}

@media (max-width: 768px) {
  .hero::after {
    width: 768px;
    height: 768px;
    right: -80px;
    top: 40%;
    transform: translateY(-40%);
    bottom: -80px;
    opacity: 0.5;
  }
}


.hero-title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.02em;
  margin: 0;
}

.hero-sub {
  margin: 6px 0 0;
  color: #6b7280;
  font-size: 14px;
}

.hero-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 0;
  width: auto;
}

.hero-hint {
  font-size: 12px;
  color: #6b7280;
  text-align: right;
}

/* Layout grid */
.layout {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr) 240px;
  gap: 18px;
  align-items: start;
}

.col {
  min-width: 0;
}

/* cards */
.card {
  background: #fff;
  border-radius: 18px;
  border: 1px solid rgba(17,24,39,0.06);
  box-shadow: 0 10px 30px rgba(0,0,0,.04);
}

.profile-card {
  padding: 18px;
}

.profile-top {
  display: flex;
  gap: 14px;
  align-items: center;
}

.avatar-wrap {
  position: relative;
  width: 74px;
  height: 74px;
  flex: 0 0 auto;
}

.avatar {
  width: 74px;
  height: 74px;
  border-radius: 18px;
  object-fit: cover;
  border: 1px solid rgba(99,102,241,0.18);
}

.mfa-dot {
  position: absolute;
  right: -6px;
  bottom: -6px;
  width: 16px;
  height: 16px;
  background: #22c55e;
  border-radius: 999px;
  border: 3px solid #fff;
  box-shadow: 0 6px 16px rgba(34,197,94,0.25);
}

.profile-meta {
  min-width: 0;
  flex: 1 1 auto;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.display-name {
  font-weight: 700;
  font-size: 18px;
  letter-spacing: -0.01em;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.minor {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  font-size: 12px;
}

.muted {
  color: #6b7280;
}

.dot {
  width: 4px;
  height: 4px;
  background: #cbd5e1;
  border-radius: 999px;
}

.section-title {
  font-size: 12px;
  color: #6b7280;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: .08em;
  margin-bottom: 8px;
}

.profile-about {
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid rgba(17,24,39,0.06);
}

.about-text {
  font-size: 14px;
  color: #111827;
  line-height: 1.55;
  white-space: pre-line;
}

.profile-links {
  margin-top: 14px;
}

.links {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 8px;
}

.links a {
  display: block;
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(79,70,229,0.06);
  color: #3730a3;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: transform .15s ease, background .15s ease;
}

.links a:hover {
  background: rgba(79,70,229,0.10);
  transform: translateY(-1px);
}

.profile-footer {
  margin-top: 16px;
  display: grid;
  gap: 10px;
}

/* Main card */
.main-card {
  padding: 18px;
}

.main-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 14px;
}

.main-title {
  font-size: 16px;
  font-weight: 700;
}

.resume-list {
  display: grid;
  gap: 12px;
}

.resume-card {
  border-radius: 16px;
  padding: 14px;
  border: 1px solid rgba(17,24,39,0.06);
  background: linear-gradient(180deg, rgba(255,255,255,1), rgba(249,250,255,1));
  transition: transform .15s ease, box-shadow .15s ease, border-color .15s ease;
}

.resume-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 14px 30px rgba(17,24,39,.07);
  border-color: rgba(99,102,241,0.18);
}

.resume-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
}

.resume-title {
  font-weight: 700;
  font-size: 14px;
  color: #111827;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Progress */
.progress {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 10px;
}

.progress-bg {
  flex: 1 1 auto;
  height: 8px;
  background: rgba(17,24,39,0.06);
  border-radius: 999px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #4f46e5, #6366f1);
  border-radius: 999px;
}

.progress-label {
  width: 42px;
  text-align: right;
  font-size: 12px;
  color: #6b7280;
}

/* Actions */
.resume-actions {
  display: flex;
  gap: 14px;
  margin-top: 10px;
}

.link {
  font-weight: 600;
  font-size: 13px;
  color: #4f46e5;
}

/* Right actions column */
.actions-col {
  position: sticky;
  top: 96px;
  display: grid;
  gap: 14px;
}

.actions-card {
  padding: 14px;
}

.actions-card.subtle {
  background: linear-gradient(135deg, rgba(79,70,229,0.05), rgba(99,102,241,0.03));
  border: 1px solid rgba(99,102,241,0.14);
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 12px;
  border-radius: 14px;
  background: rgba(17,24,39,0.03);
  font-weight: 700;
  color: #111827;
  transition: transform .15s ease, background .15s ease;
}

.action-btn:hover {
  transform: translateY(-1px);
  background: rgba(79,70,229,0.06);
}

.action-btn.primary {
  background: linear-gradient(90deg, #4f46e5, #6366f1);
  color: #fff;
}

.ic {
  width: 20px;
  text-align: center;
}

/* Pills & badges */
.badge {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 999px;
  font-weight: 700;
  border: 1px solid rgba(17,24,39,0.06);
}

.badge-mfa {
  background: rgba(34,197,94,0.12);
  color: #166534;
  border-color: rgba(34,197,94,0.20);
}

.badge-pub {
  background: rgba(34,197,94,0.12);
  color: #166534;
}

.badge-draft {
  background: rgba(245,158,11,0.14);
  color: #92400e;
}

.sub-title {
  font-weight: 800;
  margin-bottom: 10px;
}

.sub-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.pill {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 999px;
  font-weight: 800;
}

.pill-ok {
  background: rgba(34,197,94,0.14);
  color: #166534;
}

.pill-warn {
  background: rgba(245,158,11,0.16);
  color: #92400e;
}

.sub-hint {
  margin-top: 10px;
  font-size: 12px;
  color: #6b7280;
  line-height: 1.5;
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-radius: 14px;
  padding: 12px 14px;
  font-weight: 800;
  background: rgba(17,24,39,0.03);
  color: #111827;
  border: 1px solid rgba(17,24,39,0.06);
  transition: transform .15s ease, background .15s ease;
}

.btn:hover {
  transform: translateY(-1px);
  background: rgba(79,70,229,0.06);
}

.btn.primary {
  background: linear-gradient(90deg, #4f46e5, #6366f1);
  color: #fff;
  border-color: transparent;
}

.btn.wide {
  width: 100%;
}

.btn.small {
  padding: 10px 12px;
  font-size: 13px;
}

/* Hero buttons re-use */
.hero-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-radius: 14px;
  padding: 12px 14px;
  font-weight: 800;
  background: rgba(17,24,39,0.03);
  color: #111827;
  border: 1px solid rgba(17,24,39,0.06);
  transition: transform .15s ease, background .15s ease;
  text-align: center;
}

.hero-btn:hover {
  transform: translateY(-1px);
  background: rgba(79,70,229,0.06);
}

.hero-btn.primary {
  background: linear-gradient(90deg, #4f46e5, #6366f1);
  color: #fff;
  border-color: transparent;
}
.hero-btn.danger {
  background: linear-gradient(90deg, #ef4444, #f97316);
  color: #fff;
  border-color: transparent;
}

/* Notes */
.note-card {
  margin-top: 14px;
  padding: 14px;
  border: 1px solid rgba(99,102,241,0.14);
  background: linear-gradient(135deg, rgba(79,70,229,0.05), rgba(99,102,241,0.02));
}

.note-title {
  font-weight: 900;
  margin-bottom: 8px;
}

.note-text {
  font-size: 12px;
  color: #6b7280;
  line-height: 1.6;
  margin-bottom: 10px;
}

/* States */
.state {
  padding: 18px 0;
}

.skeleton-grid {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr) 240px;
  gap: 18px;
}

.sk {
  background: linear-gradient(90deg, rgba(17,24,39,0.04), rgba(17,24,39,0.08), rgba(17,24,39,0.04));
  background-size: 200% 100%;
  animation: shimmer 1.3s ease-in-out infinite;
}

@keyframes shimmer {
  0% { background-position: 0% 0; }
  100% { background-position: 200% 0; }
}

.sk.card {
  border-radius: 18px;
  height: 220px;
}

.sk-profile { height: 360px; }
.sk-main { height: 520px; }
.sk-actions { height: 240px; }

.mini-skeleton {
  display: grid;
  gap: 10px;
}
.sk-line {
  height: 16px;
  border-radius: 10px;
  background: rgba(17,24,39,0.06);
}

.error-card {
  border-radius: 18px;
  padding: 22px;
  border: 1px solid rgba(239,68,68,0.18);
  background: rgba(239,68,68,0.05);
}

.error-title {
  font-weight: 900;
  font-size: 18px;
  margin-bottom: 8px;
}

.error-text {
  color: #6b7280;
  font-size: 14px;
}

.error-actions {
  margin-top: 14px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.empty {
  padding: 18px;
  border-radius: 16px;
  border: 1px dashed rgba(17,24,39,0.12);
  background: rgba(17,24,39,0.02);
}

.empty-title {
  font-weight: 900;
  margin-bottom: 6px;
}

.empty-text {
  color: #6b7280;
  font-size: 13px;
  margin-bottom: 10px;
}

/* Mobile bottom bar */
.bottom-bar {
  position: fixed;
  left: 12px;
  right: 12px;
  bottom: 12px;
  display: none;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 10px;
  padding: 10px;
  border-radius: 18px;
  background: rgba(255,255,255,0.92);
  border: 1px solid rgba(17,24,39,0.08);
  box-shadow: 0 18px 40px rgba(17,24,39,.12);
  backdrop-filter: blur(10px);
  z-index: 999;
}

.bb-item {
  display: grid;
  place-items: center;
  gap: 4px;
  padding: 10px 8px;
  border-radius: 14px;
  font-weight: 900;
  background: rgba(79,70,229,0.06);
  color: #3730a3;
}

.bb-item:hover {
  background: rgba(79,70,229,0.10);
}

.bb-ic {
  font-size: 16px;
  line-height: 1;
}

.bb-t {
  font-size: 12px;
}

/* Responsiveness */
@media (max-width: 1280px) {
  .layout { grid-template-columns: 320px minmax(0, 1fr); }
  .actions-col { display: none; }
  .skeleton-grid { grid-template-columns: 320px minmax(0, 1fr); }
}

@media (max-width: 900px) {
  .hero { flex-direction: column; align-items: flex-start; }
  .hero-actions {
    width: 100%;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 10px;
  }
  .hero-actions .hero-btn {
    flex: 1 1 calc(50% - 8px);
    text-align: center;
  }
  .hero-hint { text-align: left; width: 100%; }
}

@media (max-width: 768px) {
  .layout { grid-template-columns: 1fr; }
  .profile-col { order: 1; }
  .main-col { order: 2; }
  .bottom-bar { display: grid; }
  .wrap { padding-bottom: 150px; }
  .skeleton-grid { grid-template-columns: 1fr; }
}
.link.danger {
  background: none;
  border: none;
  padding: 0;
  font: inherit;
  color: #dc2626;
  cursor: pointer;
  font-weight: 700;
}

.link.danger:hover {
  text-decoration: underline;
}

.link.danger:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
/* ===== MODAL ===== */

.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.55);
  display: grid;
  place-items: center;
  z-index: 200;
  backdrop-filter: blur(6px);
}

.modal-card {
  width: 100%;
  max-width: 420px;
  background: #fff;
  border-radius: 20px;
  padding: 22px;
  box-shadow: 0 30px 80px rgba(0,0,0,.25);
  animation: modalIn .18s ease-out;
}

@keyframes modalIn {
  from {
    opacity: 0;
    transform: translateY(6px) scale(.98);
  }
  to {
    opacity: 1;
    transform: none;
  }
}

.modal-title {
  font-size: 18px;
  font-weight: 900;
  margin-bottom: 10px;
}

.modal-text {
  font-size: 14px;
  color: #4b5563;
  line-height: 1.6;
  margin-bottom: 18px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn.danger {
  background: #dc2626;
  color: #fff;
  border: none;
}

.btn.danger:hover {
  background: #b91c1c;
}
.resume-card {
  position: relative;
  overflow: hidden;

  border-radius: 20px;
  padding: 16px 16px 14px;

  background:
    linear-gradient(180deg, #ffffff, #f9faff);

  border: 1px solid rgba(99,102,241,0.14);

  box-shadow:
    0 10px 24px rgba(17,24,39,0.06),
    inset 0 1px 0 rgba(255,255,255,.9);

  transition:
    transform .18s ease,
    box-shadow .18s ease,
    border-color .18s ease;
}

.resume-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background:
    radial-gradient(
      400px 120px at 10% -20%,
      rgba(99,102,241,0.10),
      transparent 60%
    );
  opacity: .9;
  pointer-events: none;
}

.resume-card:hover {
  transform: translateY(-3px);
  border-color: rgba(99,102,241,0.35);
  box-shadow:
    0 18px 40px rgba(17,24,39,0.12),
    inset 0 1px 0 rgba(255,255,255,.9);
}
.resume-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.resume-title {
  font-size: 15px;
  font-weight: 800;
  letter-spacing: -0.01em;
  color: #0f172a;

  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.progress {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
}

.progress-bg {
  position: relative;
  flex: 1;
  height: 10px;
  background: rgba(99,102,241,0.10);
  border-radius: 999px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background:
    linear-gradient(
      90deg,
      #4f46e5,
      #6366f1,
      #818cf8
    );
  border-radius: 999px;
  box-shadow:
    0 0 0 1px rgba(255,255,255,.4) inset,
    0 6px 14px rgba(99,102,241,.35);
}

.progress-label {
  min-width: 42px;
  text-align: right;
  font-size: 12px;
  font-weight: 700;
  color: #4f46e5;
}
.resume-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-top: 12px;
  padding-top: 10px;
  border-top: 1px dashed rgba(99,102,241,0.18);
}
.link {
  font-size: 13px;
  font-weight: 700;
  color: #4f46e5;
  position: relative;
  transition: color .15s ease;
}

.link::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: -2px;
  height: 2px;
  background: currentColor;
  transform: scaleX(0);
  transform-origin: left;
  transition: transform .15s ease;
}

.link:hover::after {
  transform: scaleX(1);
}

.link.edit {
  color: #0f766e;
}

.link.view {
  color: #4338ca;
}
.link.danger {
  color: #dc2626;
  font-weight: 800;
}

.link.danger::after {
  background: #dc2626;
}
@media (max-width: 640px) {
  .resume-card {
    padding: 14px;
  }

  .resume-actions {
    gap: 14px;
    flex-wrap: wrap;
  }

  .progress-label {
    font-size: 11px;
  }
}

/* ===== unread counters (chats / applications) ===== */
.count-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;

  min-width: 22px;
  height: 22px;
  padding: 0 8px;

  border-radius: 999px;
  font-size: 12px;
  font-weight: 900;
  line-height: 1;

  color: #fff;
  background: linear-gradient(90deg, #ef4444, #f97316);
  box-shadow:
    0 10px 22px rgba(239,68,68,0.22),
    inset 0 1px 0 rgba(255,255,255,0.25);
  border: 1px solid rgba(255,255,255,0.35);
}

/* когда бейдж стоит на светлых кнопках — чуть контрастнее */
.action-btn .count-badge,
.hero-btn .count-badge,
.bb-item .count-badge {
  margin-left: 8px;
}

/* для primary-кнопок можно сделать “стеклянный” вариант */
.action-btn.primary .count-badge {
  background: rgba(255,255,255,0.22);
  border-color: rgba(255,255,255,0.35);
  box-shadow: inset 0 1px 0 rgba(255,255,255,0.25);
}

</style>

