<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useI18n } from '~/composables/useI18n'
import { authFetch } from '~/composables/useAuthFetch'

const props = defineProps<{
  profile: any
  isOwner: boolean
}>()

const auth = useAuthStore()
const router = useRouter()
const { t } = useI18n()

/* =========================
 * BASIC VIEW MODEL
 * ========================= */
const isOwner = computed(() => Boolean(props.isOwner))

const companyId = computed(() => {
  return String(
    props.profile?.companyId ||
      props.profile?.userId ||
      props.profile?.id ||
      ''
  )
})

const displayName = computed(
  () =>
    props.profile?.displayName ||
    props.profile?.companyName ||
    t('companyPanel.company')
)
const location = computed(
  () => props.profile?.location || t('profilePage.locationNotSpecified')
)

const avatarUrl = computed(() => {
  return props.profile?.avatarUrl || props.profile?.logoUrl || '/images/com-base.webp'
})

const links = computed<string[]>(() => {
  const l = props.profile?.links
  if (Array.isArray(l)) return l.filter(Boolean)
  return []
})

const website = computed(() => String(props.profile?.website || '').trim())
const industry = computed(() => String(props.profile?.industry || '').trim())

const aboutText = computed(() => {
  const v =
    props.profile?.about ??
    props.profile?.description ??
    props.profile?.companyDescription ??
    props.profile?.bio ??
    ''
  return String(v || '').trim()
})

const hasLinks = computed(() => links.value.length > 0)

function normalizeUrl(u: string) {
  const s = String(u || '').trim()
  if (!s) return ''
  if (/^https?:\/\//i.test(s)) return s
  return `https://${s}`
}

function openWebsite(u: string) {
  if (!process.client) return
  const url = normalizeUrl(u)
  if (!url) return
  window.open(url, '_blank', 'noopener,noreferrer')
}

/* fallback для битого лого */
function onAvatarError(e: Event) {
  const img = e.target as HTMLImageElement | null
  if (!img) return
  // чтобы не уходить в бесконечный loop
  if (!img.src.includes('/images/com-base.webp')) {
    img.src = '/images/com-base.webp'
  }
}

/* =========================
 * BADGES (CHATS + APPLICATIONS INBOX)
 * ========================= */
const badgesPending = ref(false)
const chatsUnreadTotal = ref(0)
const unansweredAppsTotal = ref(0)
let badgeTimer: number | null = null

// чтобы можно было “обнулить” бейдж на Отклики при клике (и не пересчитать его обратно сразу)
const lastInboxItems = ref<any[]>([])

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

function chatLastSeenKey(appId: string) {
  const cid = companyId.value || 'anon'
  return `unicorn:company:chat:last_seen:${cid}:${appId}`
}

function loadChatLastSeen(appId: string) {
  if (!process.client) return 0
  try {
    const raw = localStorage.getItem(chatLastSeenKey(appId))
    const num = Number(raw)
    return Number.isFinite(num) ? num : 0
  } catch {
    return 0
  }
}

function saveChatLastSeen(appId: string, ts: number) {
  if (!process.client) return
  try {
    localStorage.setItem(chatLastSeenKey(appId), String(ts || 0))
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
  const res = await authFetch(`/api/chat/${appId}/messages`, { method: 'GET' })
  const items = Array.isArray(res?.items) ? res.items : []
  const last = items.length ? items[items.length - 1] : null
  const ts = pickMessageTimestamp(last)
  const senderType = String(last?.senderType || '')
  return { ts, senderType }
}


function isUnansweredByCompany(status?: string) {
  const v = String(status || '').toLowerCase()
  return (
    v.includes('pending') ||
    v.includes('review') ||
    v.includes('new') ||
    v.includes('created') ||
    v.includes('waiting') ||
    v.includes('sent')
  )
}

/**
 * localStorage map: applicationId -> lastStatusSeen
 */
function inboxSeenKey() {
  const cid = companyId.value || 'anon'
  return `unicorn:company:applications:inbox_seen_status:${cid}`
}

function loadInboxSeenMap(): Record<string, string> {
  if (!process.client) return {}
  try {
    return JSON.parse(localStorage.getItem(inboxSeenKey()) || '{}') || {}
  } catch {
    return {}
  }
}

function saveInboxSeenMap(map: Record<string, string>) {
  if (!process.client) return
  try {
    localStorage.setItem(inboxSeenKey(), JSON.stringify(map))
  } catch {
    // ignore
  }
}

function markInboxSeen(items?: any[]) {
  const arr = Array.isArray(items) ? items : lastInboxItems.value
  const map = loadInboxSeenMap()

  for (const a of arr) {
    const id = String(a?.id || a?.applicationId || a?._id || '')
    if (!id) continue
    const st = String(a?.status || a?.state || '')
    map[id] = st
  }

  saveInboxSeenMap(map)
  unansweredAppsTotal.value = 0
}

function onApplicationsClick() {
  markInboxSeen()
}

/* =========================
 * CHATS LIST PREVIEW (unread + last message)
 * ========================= */
type ChatPreview = {
  id: string
  unreadCount: number
  lastText: string
  lastAt?: string
  peerName?: string
}

const chatsPreview = ref<ChatPreview[]>([])

function pickChatId(c: any): string {
  return String(c?.id || c?.chatId || c?._id || '')
}

function pickLastText(c: any): string {
  const lists = [
    Array.isArray(c?.messages) ? c.messages : null,
    Array.isArray(c?.chat?.messages) ? c.chat.messages : null,
    Array.isArray(c?.data?.messages) ? c.data.messages : null,
    Array.isArray(c?.items) ? c.items : null,
  ].filter(Boolean) as any[][]

  const all = lists.flat()
  const latest = all.length ? pickLatestMessage(all) : null
  const latestText = pickMessageText(latest)
  if (latestText) return latestText

  const lm = c?.lastMessage || c?.last_message || null
  const t =
    lm?.text ??
    lm?.body ??
    lm?.message ??
    c?.lastMessageText ??
    c?.last_message_text ??
    c?.preview ??
    c?.snippet ??
    ''
  return String(t || '').trim()
}

function chatBadgeKey() {
  const cid = companyId.value || 'anon'
  return `unicorn:company:chats:unread:${cid}`
}

function loadChatBadge() {
  if (!process.client) return 0
  try {
    const raw = localStorage.getItem(chatBadgeKey())
    const num = Number(raw)
    return Number.isFinite(num) ? num : 0
  } catch {
    return 0
  }
}

function saveChatBadge(value: number) {
  if (!process.client) return
  try {
    localStorage.setItem(chatBadgeKey(), String(value || 0))
  } catch {}
}
function pickLastAt(c: any): string | undefined {
  const lists = [
    Array.isArray(c?.messages) ? c.messages : null,
    Array.isArray(c?.chat?.messages) ? c.chat.messages : null,
    Array.isArray(c?.data?.messages) ? c.data.messages : null,
    Array.isArray(c?.items) ? c.items : null,
  ].filter(Boolean) as any[][]

  const all = lists.flat()
  const latest = all.length ? pickLatestMessage(all) : null
  const ts = pickMessageTimestamp(latest)
  if (ts) return new Date(ts).toISOString()

  const lm = c?.lastMessage || c?.last_message || null
  const v = lm?.createdAt ?? lm?.created_at ?? c?.updatedAt ?? c?.updated_at ?? c?.lastAt ?? null
  return v ? String(v) : undefined
}

function pickMessageText(m: any): string {
  if (!m) return ''
  if (typeof m === 'string') return m.trim()
  return String(m?.text ?? m?.message ?? m?.content ?? m?.body ?? m?.preview ?? '').trim()
}

function pickLatestMessage(list: any[]): any {
  let best: any = null
  let bestTs = 0
  for (const m of list) {
    const ts = pickMessageTimestamp(m)
    if (ts > bestTs) {
      best = m
      bestTs = ts
    } else if (!best && pickMessageText(m)) {
      best = m
    }
  }
  return best
}

function pickPeerName(c: any): string | undefined {
  const u = c?.peer || c?.user || c?.otherUser || c?.companion || null
  const v = u?.displayName ?? u?.name ?? c?.title ?? c?.name ?? null
  return v ? String(v) : undefined
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

async function tryFetchList(urls: string[]) {
  let lastErr: any = null
  for (const url of urls) {
    try {
      const res: any = await authFetch(url, {
        method: 'GET'
      })

      const raw =
        (Array.isArray(res?.items) ? res.items : null) ||
        (Array.isArray(res?.chats) ? res.chats : null) ||
        (Array.isArray(res?.applications) ? res.applications : null) ||
        (Array.isArray(res?.data) ? res.data : null) ||
        (Array.isArray(res) ? res : null)

      if (Array.isArray(raw)) return raw
    } catch (e) {
      lastErr = e
    }
  }
  throw lastErr || new Error('No endpoint matched')
}

async function ensureAuthReady() {
  if (!process.client) return
  if (auth.isReady) return
  if (!auth.hasToken) return
  try {
    await auth.init()
  } catch {
    // ignore: бейджи просто будут 0
  }
}

async function loadCompanyBadges() {
  if (!process.client) return
  if (!isOwner.value) return
  if (badgesPending.value) return

  badgesPending.value = true
  try {
    await ensureAuthReady()

    // 1) INBOX applications
    let inbox: any[] = []
    try {
      inbox = await tryFetchList([
        '/api/applications/inbox',
        '/api/company/applications/inbox',
        '/api/applications/company/inbox',
      ])
    } catch {
      inbox = []
    }
    lastInboxItems.value = inbox

    const seen = loadInboxSeenMap()
    let unansweredNew = 0
    for (const a of inbox) {
      const id = String(a?.id || a?.applicationId || a?._id || '')
      if (!id) continue

      const st = String(a?.status || a?.state || '')
      if (!isUnansweredByCompany(st)) continue

      if (seen[id] !== st) unansweredNew += 1
    }

    // 2) CHATS list (unread + last message)
    let chatsRaw: any[] = []
    try {
      chatsRaw = await tryFetchList([
        '/api/company/chats',
        '/api/chats/my',
        '/api/chat/my',
        '/api/chats/company',
        '/api/company/chat',
      ])
    } catch {
      chatsRaw = []
    }

    const mapped = chatsRaw
      .map((c: any) => {
        const id = pickChatId(c)
        return {
          id,
          unreadCount: pickUnreadCount(c),
          lastText: pickLastText(c),
          lastAt: pickLastAt(c),
          peerName: pickPeerName(c),
        } as ChatPreview
      })
      .filter((x) => Boolean(x.id))

    let unreadByMessages = 0
    const list = chatsRaw.slice(0, 20)
    for (const c of list) {
      const appId = String(c?.applicationId || c?.id || c?._id || '')
      if (!appId) continue
      try {
        const info = await fetchLastMessageInfo(appId)
        if (!info.ts) continue
        if (info.senderType !== 'user') continue
        const lastSeen = loadChatLastSeen(appId)
        if (info.ts > lastSeen) unreadByMessages += 1
      } catch {}
    }

    chatsPreview.value = mapped
    chatsUnreadTotal.value = unreadByMessages
    saveChatBadge(chatsUnreadTotal.value)
    unansweredAppsTotal.value = unansweredNew
  } catch {
    chatsUnreadTotal.value = 0
    saveChatBadge(0)
    unansweredAppsTotal.value = 0
    chatsPreview.value = []
    lastInboxItems.value = []
  } finally {
    badgesPending.value = false
  }
}

/* =========================
 * MFA (2FA) STATE — company by backend signal
 * ========================= */
const mfaLoaded = ref(false)
const backendMfaEnabled = ref(false)
const showLogoutModal = ref(false)

type MfaState = 'loading' | 'enabled' | 'required'

const mfaState = computed<MfaState>(() => {
  if (!mfaLoaded.value) return 'loading'
  return backendMfaEnabled.value ? 'enabled' : 'required'
})

async function checkCompanyMfa() {
  if (!process.client) return
  try {
    await ensureAuthReady()
    await authFetch('/api/vacancies/my')
    backendMfaEnabled.value = true
  } catch (e: any) {
    const status = getHttpStatus(e)
    const code =
      e?.data?.error ||
      e?.data?.code ||
      e?.response?.data?.error ||
      ''

    if (status === 403 && String(code) === 'mfa_required') {
      backendMfaEnabled.value = false
    } else {
      backendMfaEnabled.value = true
    }
  } finally {
    mfaLoaded.value = true
  }
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

/* =========================
 * LIFECYCLE
 * ========================= */
const onFocus = () => {
  if (isOwner.value) loadCompanyBadges()
}
const onVis = () => {
  if (document.visibilityState === 'visible' && isOwner.value) loadCompanyBadges()
}

onMounted(async () => {
  if (!process.client) return

  if (isOwner.value) {
    chatsUnreadTotal.value = loadChatBadge()
    await loadCompanyBadges()
    await checkCompanyMfa()
  }

  window.addEventListener('focus', onFocus)
  document.addEventListener('visibilitychange', onVis)
  badgeTimer = window.setInterval(() => {
    if (document.visibilityState === 'visible' && isOwner.value) loadCompanyBadges()
  }, 10_000)
})

onBeforeUnmount(() => {
  if (!process.client) return
  window.removeEventListener('focus', onFocus)
  document.removeEventListener('visibilitychange', onVis)
  if (badgeTimer) window.clearInterval(badgeTimer)
})

// перезагрузка бейджей при смене профиля/владельца
watch(
  [() => isOwner.value, () => companyId.value],
  ([owner]) => {
    if (!process.client) return
    if (owner) {
      loadCompanyBadges()
      checkCompanyMfa()
    } else {
      chatsUnreadTotal.value = 0
      unansweredAppsTotal.value = 0
    }
  },
  { immediate: true }
)
</script>
<template>
  <section class="wrap">
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

    <!-- SUBSCRIPTION BANNER -->
    <div class="mb-6 px-4 md:px-6">
      <SubscriptionBanner />
    </div>

    <!-- HERO -->
    <div class="hero">
      <div class="hero-left">
        <div class="hero-kicker">{{ t('companyPanel.company') }}</div>
        <h1 class="hero-title">{{ displayName }}</h1>

        <div class="hero-subline">
          <span class="muted">{{ location }}</span>
          <span class="dot" />
          <span class="muted">{{ t('companyPanel.id') }} #{{ companyId }}</span>
        </div>

        <div class="hero-badges">
          <!-- 2FA status (owner only) -->
          <template v-if="isOwner">
            <span v-if="mfaState === 'loading'" class="badge">{{ t('companyPanel.mfaChecking') }}</span>
            <span v-else-if="mfaState === 'enabled'" class="badge badge-ok">{{ t('companyPanel.mfaActive') }}</span>
            <span v-else class="badge badge-warn">{{ t('companyPanel.mfaRequired') }}</span>
          </template>

          <span v-if="industry" class="badge badge-soft">{{ t('companyPanel.industry') }} {{ industry }}</span>
        </div>
      </div>

      <div v-if="isOwner" class="hero-actions">
        <button class="hero-btn danger" type="button" @click="askLogout">
          {{ t('panel.common.logout') }}
        </button>
      </div>
    </div>

    <!-- LAYOUT -->
    <div class="layout">
      <!-- LEFT -->
      <aside class="col left-col wide">
        <div class="card company-card">
          <div class="company-top">
            <div class="avatar-wrap">
              <img class="avatar" :src="avatarUrl" alt="logo" @error="onAvatarError" />
              <div v-if="isOwner && mfaState === 'enabled'" class="mfa-dot" :title="t('companyPanel.mfaActive')" />
            </div>

            <div class="company-meta">
              <div class="name-row">
                <div class="name">{{ displayName }}</div>
              </div>

              <div class="minor">
                <span class="muted">{{ location }}</span>
                <span class="dot" />
                <span class="muted">{{ t('companyPanel.company') }}</span>
              </div>
            </div>
          </div>

          <div class="section">
            <div class="section-title">{{ t('companyPanel.aboutCompany') }}</div>
            <div class="section-text">
              {{ aboutText || t('companyPanel.noDescription') }}
            </div>

            <div v-if="isOwner && !aboutText" class="section-cta">
              <NuxtLink class="btn small" to="/profile/edit">{{ t('companyPanel.addDescription') }}</NuxtLink>
            </div>
          </div>

          <div v-if="hasLinks || website" class="section">
            <div class="section-title">{{ t('companyPanel.links') }}</div>
            <div class="links">
              <button
                v-if="website"
                class="link-pill"
                type="button"
                @click="openWebsite(website)"
              >
                {{ website }}
              </button>

              <a
                v-for="(l, i) in links"
                :key="i"
                class="link-pill"
                :href="normalizeUrl(l)"
                target="_blank"
                rel="noopener noreferrer"
              >
                {{ l }}
              </a>
            </div>
          </div>
        </div>

        <!-- OWNER ONLY: 2FA NOTE -->
        <div v-if="isOwner && mfaState === 'required'" class="card note-card">
          <div class="note-title">{{ t('companyPanel.hintTitle') }}</div>
          <div class="note-text">
            {{ t('companyPanel.hintText') }}
          </div>
          <NuxtLink class="btn small" to="/profile/edit">{{ t('companyPanel.goToSecurity') }}</NuxtLink>
        </div>
      </aside>

      <!-- RIGHT (OWNER ONLY) -->
      <aside class="col right-col" v-if="isOwner">
        <div class="card actions-card">
          <NuxtLink class="action-btn primary" to="/company/chat">
            <span class="ic">💬</span><span>{{ t('companyPanel.chats') }}</span>
            <span v-if="chatsUnreadTotal > 0" class="count-badge">
              {{ formatBadge(chatsUnreadTotal) }}
            </span>
          </NuxtLink>

          <!-- ВОЗВРАЩЕНО: "Мои вакансии" (как на старом скрине).
               Маршрут поставил стандартный для Nuxt: /company/vacancies
               Если у тебя другой — просто поменяй to. -->
          <NuxtLink class="action-btn" to="/company/vacancies">
            <span class="ic">📌</span><span>{{ t('companyPanel.myVacancies') }}</span>
          </NuxtLink>

          <NuxtLink class="action-btn" to="/profile/edit">
            <span class="ic">✏️</span><span>{{ t('companyPanel.profile') }}</span>
          </NuxtLink>
        </div>

        <div class="card actions-card subtle">
          <div class="sub-title">{{ t('companyPanel.status') }}</div>
          <div class="sub-row">
            <span class="muted">{{ t('companyPanel.twoFa') }}</span>
            <span v-if="mfaState === 'loading'" class="pill">{{ t('companyPanel.checking') }}</span>
            <span v-else-if="mfaState === 'enabled'" class="pill pill-ok">{{ t('companyPanel.active') }}</span>
            <span v-else class="pill pill-warn">{{ t('companyPanel.required') }}</span>
          </div>

          <div class="sub-hint">
            {{ t('companyPanel.jobsHint') }}
          </div>
        </div>
      </aside>
    </div>

    <!-- MOBILE BOTTOM BAR (OWNER ONLY) -->
    <nav v-if="isOwner" class="bottom-bar">
      <NuxtLink to="/company/chat" class="bb-item">
        <span class="bb-ic">💬</span>
        <span class="bb-t">{{ t('companyPanel.chats') }}</span>
        <span v-if="chatsUnreadTotal > 0" class="count-badge">
          {{ formatBadge(chatsUnreadTotal) }}
        </span>
      </NuxtLink>

      <!-- ВОЗВРАЩЕНО: Вакансии вместо Откликов -->
      <NuxtLink to="/company/vacancies" class="bb-item">
        <span class="bb-ic">📌</span>
        <span class="bb-t">{{ t('companyPanel.vacancies') }}</span>
      </NuxtLink>

      <NuxtLink to="/profile/edit" class="bb-item">
        <span class="bb-ic">✏️</span>
        <span class="bb-t">{{ t('companyPanel.profile') }}</span>
      </NuxtLink>
    </nav>
  </section>
</template>

<style scoped>
.wrap {
  max-width: 1320px;
  margin: 0 auto;
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
    radial-gradient(600px 260px at 85% 30%, rgba(99,102,241,0.18), transparent 60%),
    linear-gradient(135deg, rgba(79,70,229,0.10), rgba(99,102,241,0.06));
  border: 1px solid rgba(99,102,241,0.15);
  box-shadow: 0 12px 40px rgba(17,24,39,.06);
}

.hero::after {
  content: '';
  position: absolute;

  right: -40px;
  top: 50%;
  transform: translateY(-52%);

  width: 980px;
  height: 980px;

  background-image: url('/images/company_panel.webp');
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;

  opacity: 0.40;
  pointer-events: none;
  z-index: 0;
}

.hero-left,
.hero-actions {
  position: relative;
  z-index: 1;
}

.hero-kicker {
  font-size: 12px;
  font-weight: 900;
  letter-spacing: .08em;
  text-transform: uppercase;
  color: rgba(55,48,163,.85);
}

.hero-title {
  margin: 8px 0 0;
  font-size: 28px;
  font-weight: 900;
  letter-spacing: -0.02em;
}

.hero-subline {
  margin-top: 8px;
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.hero-badges {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.hero-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 240px;
  align-items: flex-end;
}

.hero-hint {
  font-size: 12px;
  color: #6b7280;
  text-align: right;
}

.hero-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  padding: 12px 14px;
  font-weight: 900;
  background: rgba(17,24,39,0.03);
  color: #111827;
  border: 1px solid rgba(17,24,39,0.06);
  transition: transform .15s ease, background .15s ease;
  text-align: center;
  width: 100%;
}
.hero-btn:hover { transform: translateY(-1px); background: rgba(79,70,229,0.06); }
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

/* Badges */
.badge {
  font-size: 12px;
  padding: 6px 10px;
  border-radius: 999px;
  font-weight: 900;
  border: 1px solid rgba(17,24,39,0.06);
  background: rgba(17,24,39,0.03);
  color: #111827;
}
.badge-ok {
  background: rgba(34,197,94,0.12);
  border-color: rgba(34,197,94,0.20);
  color: #166534;
}
.badge-warn {
  background: rgba(245,158,11,0.14);
  border-color: rgba(245,158,11,0.24);
  color: #92400e;
}
.badge-soft {
  background: rgba(79,70,229,0.06);
  border-color: rgba(99,102,241,0.14);
  color: #3730a3;
}

.muted { color: #6b7280; }
.dot { width: 4px; height: 4px; border-radius: 999px; background: #cbd5e1; display: inline-block; }

/* Layout grid */
.layout {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr) 260px;
  gap: 18px;
  align-items: start;
}

.col { min-width: 0; }

.left-col.wide {
  grid-column: span 2;
}

/* Cards */
.card {
  background: #fff;
  border-radius: 18px;
  border: 1px solid rgba(17,24,39,0.06);
  box-shadow: 0 10px 30px rgba(0,0,0,.04);
}

.company-card { padding: 18px; }

.company-top {
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
  background: rgba(79,70,229,0.06);
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

.company-meta { min-width: 0; flex: 1 1 auto; }
.name-row { display: flex; gap: 8px; align-items: center; }
.name {
  font-weight: 900;
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

.section {
  margin-top: 14px;
  padding-top: 14px;
  border-top: 1px solid rgba(17,24,39,0.06);
}
.section-title {
  font-size: 12px;
  color: #6b7280;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: .08em;
  margin-bottom: 8px;
}
.section-text {
  font-size: 14px;
  color: #111827;
  line-height: 1.55;
  white-space: pre-line;
}
.section-cta { margin-top: 12px; }

.links {
  display: grid;
  gap: 8px;
}
.link-pill {
  display: inline-flex;
  align-items: center;
  justify-content: flex-start;
  padding: 10px 12px;
  border-radius: 12px;
  background: rgba(79,70,229,0.06);
  color: #3730a3;
  font-weight: 700;
  border: 1px solid rgba(99,102,241,0.14);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: transform .15s ease, background .15s ease;
}
.link-pill:hover {
  background: rgba(79,70,229,0.10);
  transform: translateY(-1px);
}

/* Right column */
.right-col {
  position: sticky;
  top: 96px;
  display: grid;
  gap: 14px;
}
.actions-card { padding: 14px; }
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
  font-weight: 900;
  color: #111827;
  transition: transform .15s ease, background .15s ease;
}
.action-btn:hover { transform: translateY(-1px); background: rgba(79,70,229,0.06); }
.action-btn.primary { background: linear-gradient(90deg, #4f46e5, #6366f1); color: #fff; }
.ic { width: 20px; text-align: center; }

.sub-title { font-weight: 900; margin-bottom: 10px; }
.sub-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; }

.pill {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 999px;
  font-weight: 900;
  background: rgba(17,24,39,0.03);
  border: 1px solid rgba(17,24,39,0.06);
}
.pill-ok { background: rgba(34,197,94,0.14); color: #166534; border-color: rgba(34,197,94,0.20); }
.pill-warn { background: rgba(245,158,11,0.16); color: #92400e; border-color: rgba(245,158,11,0.22); }

.sub-hint {
  margin-top: 10px;
  font-size: 12px;
  color: #6b7280;
  line-height: 1.5;
}

/* Note */
.note-card {
  margin-top: 14px;
  padding: 14px;
  border: 1px solid rgba(99,102,241,0.14);
  background: linear-gradient(135deg, rgba(79,70,229,0.05), rgba(99,102,241,0.02));
}
.note-title { font-weight: 900; margin-bottom: 8px; }
.note-text { font-size: 12px; color: #6b7280; line-height: 1.6; margin-bottom: 10px; }

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-radius: 14px;
  padding: 12px 14px;
  font-weight: 900;
  background: rgba(17,24,39,0.03);
  color: #111827;
  border: 1px solid rgba(17,24,39,0.06);
  transition: transform .15s ease, background .15s ease;
}
.btn:hover { transform: translateY(-1px); background: rgba(79,70,229,0.06); }
.btn.small { padding: 10px 12px; font-size: 13px; }
.btn.danger {
  background: #dc2626;
  color: #fff;
  border: none;
}
.btn.danger:hover {
  background: #b91c1c;
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
.bb-item:hover { background: rgba(79,70,229,0.10); }
.bb-ic { font-size: 16px; line-height: 1; }
.bb-t { font-size: 12px; }

/* Responsive */
@media (max-width: 1280px) {
  .layout { grid-template-columns: 360px minmax(0, 1fr); }
  .right-col { display: none; }
}

@media (max-width: 900px) {
  .hero { flex-direction: column; align-items: flex-start; }
  .hero-actions { width: 100%; align-items: stretch; }
  .hero-actions .hero-btn { width: 100%; }
}

@media (max-width: 1280px) {
  .hero::after {
    width: 720px;
    height: 720px;
    right: -120px;
    top: 45%;
    transform: translateY(-45%);
    opacity: 0.52;
  }

  .layout { grid-template-columns: 1fr; }
  .bottom-bar { display: grid; }
}

.hero-btn, .action-btn, .bb-item { position: relative; }

.count-badge{
  position:absolute;
  top:8px;
  right:10px;
  display:inline-flex;
  align-items:center;
  justify-content:center;
  height:20px;
  min-width:20px;
  padding:0 6px;
  border-radius:999px;
  background:#ef4444;
  color:#fff;
  font-weight:900;
  font-size:12px;
  line-height:1;
  box-shadow:0 10px 20px rgba(239,68,68,.22);
}

/* чтобы на маленьких кнопках bottom-bar не залезало на текст */
.bottom-bar .count-badge{
  top:6px;
  right:8px;
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

</style>
