<template>
  <div class="page">
    <Header />

    <section class="wrap">
      <div class="hero">
        <div class="hero-left">
          <h1 class="hero-title">{{ t('chatList.title') }}</h1>
          <p class="hero-sub">{{ t('chatList.description') }}</p>
        </div>

        <div class="hero-actions">
          <NuxtLink class="hero-btn" to="/user/applications">{{ t('chatList.backToApps') }}</NuxtLink>
          <button class="hero-btn primary" :disabled="pending" @click="load(true)">
            {{ pending ? t('chatList.refreshing') : t('chatList.refresh') }}
          </button>
        </div>
      </div>

      <div class="toolbar">
        <div class="search">
          <span class="s-ic">⌕</span>
          <input v-model="q" class="s-inp" :placeholder="t('chatList.searchPlaceholder')" />
        </div>

        <div class="meta">
          <span class="muted">{{ t('chatList.total') }}</span>
          <b>{{ filtered.length }}</b>
        </div>
      </div>

      <div v-if="err" class="banner danger">
        {{ err }}
      </div>

      <div v-if="hint" class="banner warn">
        {{ hint }}
      </div>

      <!-- Loading -->
      <div v-if="pending" class="state">
        <div class="skeleton-grid">
          <div v-for="i in 6" :key="i" class="sk card" />
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="!filtered.length" class="empty">
        <div class="empty-title">{{ t('chatList.noDialogsTitle') }}</div>
        <div class="empty-text">{{ t('chatList.noDialogsText') }}</div>
        <NuxtLink class="btn primary" to="/user/applications">{{ t('chatList.goToApps') }}</NuxtLink>
      </div>

      <!-- List -->
      <div v-else class="grid">
        <div
          v-for="a in filtered"
          :key="a.applicationId"
          class="chat-item"
        >
          <NuxtLink
            class="card chat-card"
            :to="`/user/chat/${a.applicationId}`"
          >
            <div class="cc-top">
              <div class="cc-left">
                <div class="logo">
                  <img
                    :src="a.companyAvatar"
                    alt="company"
                    @error="onAvatarError($event, '/images/com-base.webp')"
                  />
                </div>

                <div class="cc-meta">
                  <div class="cc-title-row">
                    <div class="cc-title">
                      {{ a.vacancyTitle || t('chatList.vacancy') }}
                    </div>

                    <span
                      v-if="(a.unreadCount || 0) > 0"
                      class="unread-badge"
                      :title="`${t('chatList.unreadCount')}: ${a.unreadCount}`"
                    >
                      {{ formatBadge(a.unreadCount || 0) }}
                    </span>
                  </div>

                  <div class="cc-sub">
                    <span class="cc-company">{{ a.companyName || t('chatList.company') }}</span>
                    <span class="dot" />
                    <span class="muted">{{ a.statusLabel }}</span>
                  </div>
                </div>
              </div>

              <div class="cc-right">
                <div class="cc-date">{{ formatDate(a.updatedAt) }}</div>
                <div class="cc-open">{{ t('chatList.open') }}</div>
              </div>
            </div>

            <div class="cc-last" :class="{ 'cc-last--unread': (a.unreadCount || 0) > 0 }" v-if="a.lastMessage">
              {{ a.lastMessage }}
            </div>
            <div class="cc-last muted" v-else>{{ t('chatList.noMessages') }}</div>
          </NuxtLink>
          
          <button 
            class="hide-btn" 
            @click.prevent="openDeleteConfirm(a.applicationId)"
            :title="t('chatList.hide')"
          >
            ✕
          </button>
        </div>
      </div>
    </section>



    <div
      v-if="confirmDeleteId"
      class="confirm-overlay"
      @click.self="closeDeleteConfirm"
    >
      <div class="confirm-modal">
        <div class="confirm-title">{{ t('chatList.confirmHide') }}</div>
        <div class="confirm-actions">
          <button class="btn ghost" @click="closeDeleteConfirm">{{ t('panel.common.logoutConfirmCancel') }}</button>
          <button class="btn danger" @click="confirmDelete">{{ t('apps.hide') }}</button>
        </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'
import { useAuthStore } from '~/stores/auth'
import { useApplications } from '~/composables/useApplications'
import { useToast } from '~/composables/useToast'
import { authFetch } from '~/composables/useAuthFetch'

definePageMeta({ middleware: ['user-auth'] })

const auth = useAuthStore()
const { t, lang } = useI18n()
const { hideApplication } = useApplications()
const toast = useToast()

const pending = ref(false)
const err = ref('')
const hint = ref('')
const q = ref('')

type AppChatItem = {
  applicationId: string
  vacancyTitle?: string
  companyName?: string
  companyAvatar?: string
  status?: string
  updatedAt?: string
  lastMessage?: string
  unreadCount?: number
}

type AppChatViewItem = AppChatItem & {
  statusLabel: string
}

const items = ref<AppChatItem[]>([])
const confirmDeleteId = ref<string | null>(null)

function formatBadge(n: number) {
  if (!n || n <= 0) return t('chatList.status.label')
  return n > 99 ? '99+' : String(n)
}

function normalizeStatus(s?: string) {
  const v = (s || '').toLowerCase()
  if (!v) return t('chatList.status.label')
  if (v.includes('pending')) return t('chatList.status.pending')
  if (v.includes('review')) return t('chatList.status.review')
  if (v.includes('accept') || v.includes('approve') || v.includes('offer') || v.includes('hire')) return t('chatList.status.accept')
  if (v.includes('reject') || v.includes('decline') || v.includes('deny')) return t('chatList.status.reject')
  return t('chatList.status.label')
}


function getHttpStatus(e: any): number {
  return Number(e?.statusCode) || Number(e?.status) || Number(e?.response?.status) || 0
}

function buildAuthHeaders(): Record<string, string> {
  const h: Record<string, string> = {}
  const token =
    (auth as any)?.accessToken ||
    (auth as any)?.access?.token ||
    (auth as any)?.tokens?.accessToken ||
    (auth as any)?.token
  if (token) h.Authorization = `Bearer ${token}`
  return h
}

async function tryFetchArray(urls: string[]) {
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
    } catch (e: any) {
      lastErr = e
      const st = getHttpStatus(e)
      if (st === 404) continue
      if (st === 401 || st === 403) continue
      // прочие ошибки — прерываем, чтобы не усугублять rate limit
      break
    }
  }

  throw lastErr || new Error('no_endpoint_matched')
}

async function hydrateMissingLastMessages() {
  const need = items.value.filter((x) => !x.lastMessage).slice(0, 20)
  if (!need.length) return

  const probeId = need[0]?.applicationId
  if (!probeId) return

  const candidates = [
    `/api/chat/${probeId}/messages`,
  ]

  let workingBase: string | null = null

  for (const url of candidates) {
    try {
      await authFetch(url, { method: 'GET' })
      workingBase = url.replace(probeId, '{id}')
      break
    } catch (e: any) {
      const st = getHttpStatus(e)
      if (st === 404) continue
      break
    }
  }

  if (!workingBase) return

  for (const it of need) {
    const url = workingBase.replace('{id}', it.applicationId)
    try {
      const res: any = await authFetch(url, { method: 'GET' })
      const list = Array.isArray(res?.items) ? res.items : []
      const lastMsg = list.length ? list[list.length - 1] : null
      const last = pickMessageText(lastMsg)
      if (last) it.lastMessage = last

      const u = String(lastMsg?.createdAt || lastMsg?.created_at || res?.updatedAt || res?.updated_at || '').trim()
      if (u) it.updatedAt = u
    } catch {
      // ignore
    }
  }

  items.value = [...items.value]
}

function pickApplicationId(x: any): string {
  return String(
    x?.applicationId ||
      x?.application?.id ||
      x?.application?._id ||
      x?.id ||
      x?._id ||
      x?.chat?.applicationId ||
      ''
  )
}

function pickVacancyTitle(x: any): string {
  return (
    x?.vacancyTitle ||
    x?.vacancy?.title ||
    x?.vacancy?.name ||
    x?.vacancy?.position ||
    x?.application?.vacancyTitle ||
    x?.application?.vacancy?.title ||
    x?.application?.vacancy?.name ||
    x?.job?.title ||
    x?.position ||
    x?.title ||
    x?.name ||
    t('chatList.vacancy')
  )
}


function pickCompanyName(x: any): string {
  return (
    x?.companyDisplayName ||
    x?.companyName ||
    x?.company?.displayName ||
    x?.company?.name ||
    x?.employer?.name ||
    x?.employer?.displayName ||
    x?.organization?.name ||
    x?.organization?.displayName ||
    x?.vacancy?.companyName ||
    x?.vacancy?.companyDisplayName ||
    x?.vacancy?.company?.name ||
    x?.vacancy?.company?.displayName ||
    x?.vacancy?.employer?.name ||
    x?.application?.companyName ||
    x?.application?.companyDisplayName ||
    x?.application?.company?.name ||
    x?.application?.company?.displayName ||
    x?.company?.login ||
    x?.employer?.login ||
    t('chatList.company')
  )
}


function pickCompanyAvatar(x: any): string {
  return normalizeAvatarUrl(
    x?.companyAvatar ||
    x?.companyAvatarUrl ||
    x?.company?.avatarUrl ||
    x?.company?.logoUrl ||
    x?.company?.photoUrl ||
    x?.company?.logo ||
    x?.company?.photo ||
    x?.vacancy?.company?.logoUrl ||
    x?.vacancy?.company?.avatarUrl ||
    x?.vacancy?.company?.photoUrl ||
    x?.vacancy?.company?.logo ||
    x?.vacancy?.company?.photo ||
    x?.application?.company?.avatarUrl ||
    x?.application?.company?.logoUrl ||
    '/images/com-base.webp',
    '/images/com-base.webp'
  )
}

function pickStatus(x: any): string {
  return String(x?.status || x?.state || x?.application?.status || x?.application?.state || '')
}

function pickUpdatedAt(x: any): string {
  return String(
    x?.updatedAt ||
      x?.updated_at ||
      x?.lastMessageAt ||
      x?.last_message_at ||
      x?.chat?.updatedAt ||
      x?.chat?.updated_at ||
      x?.createdAt ||
      x?.created_at ||
      ''
  )
}

function normalizeAvatarUrl(raw: any, fallback: string) {
  const v = String(raw || '').trim()
  if (!v || v === 'null' || v === 'undefined') return fallback
  if (v.startsWith('http') || v.startsWith('data:') || v.startsWith('/')) return v
  return `/${v}`
}

function onAvatarError(event: Event, fallback: string) {
  const img = event.target as HTMLImageElement | null
  if (img && img.src !== fallback) img.src = fallback
}

function pickMessageText(m: any): string {
  if (!m) return ''
  if (typeof m === 'string') return m.trim()
  return String(m?.text ?? m?.message ?? m?.content ?? m?.body ?? m?.preview ?? '').trim()
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

function pickLastMessage(x: any): string {
  const msgLists = [
    Array.isArray(x?.messages) ? x.messages : null,
    Array.isArray(x?.chat?.messages) ? x.chat.messages : null,
    Array.isArray(x?.chat?.items) ? x.chat.items : null,
    Array.isArray(x?.data?.messages) ? x.data.messages : null,
  ].filter(Boolean) as any[][]

  const allMessages = msgLists.flat()
  const latest = allMessages.length ? pickLatestMessage(allMessages) : null
  const latestText = pickMessageText(latest)
  if (latestText) return latestText

  const candidates = [
    x?.lastMessage,
    x?.last_message,
    x?.chat?.lastMessage,
    x?.chat?.last_message,
    x?.lastMessageText,
    x?.last_message_text,
    x?.chatLastMessage,
    x?.messagePreview,
    x?.preview,
    x?.snippet,
  ]

  for (const c of candidates) {
    const t = pickMessageText(c)
    if (t) return t
  }
  return ''
}


function pickUnreadCount(x: any): number {
  const n =
    x?.unreadCount ??
    x?.chat?.unreadCount ??
    x?.unread_messages ??
    x?.unread ??
    x?.hasUnread ??
    0

  if (typeof n === 'boolean') return n ? 1 : 0
  const num = Number(n)
  return Number.isFinite(num) ? num : 0
}

function normalizeApiError(e: any): { code: string; status?: number } {
  const status = getHttpStatus(e) || undefined
  const code =
    e?.data?.error ||
    e?.data?.code ||
    e?.response?._data?.error ||
    e?.response?._data?.code ||
    (status ? `http_${status}` : 'load_failed')
  return { code: String(code), status }
}

async function load(force = false) {
  if (pending.value && !force) return

  pending.value = true
  err.value = ''
  hint.value = ''

  try {
    if (process.client && !auth.isReady && (auth as any)?.hasToken) {
      try {
        await auth.init()
      } catch {
        // не валим страницу
      }
    }

    // Сначала пробуем ЧАТЫ (если такой endpoint есть),
    // затем fallback на APPLICATIONS (если чаты построены через отклики).
    const raw = await tryFetchArray([
      '/api/user/chats',
      '/api/chats/my',
      '/api/chat/my',
      '/api/user/applications',
      '/api/applications/my',
      '/api/applications/me',
      '/api/user/me/applications',
    ])

    items.value = raw
      .map((x: any): AppChatItem => {
        const applicationId = pickApplicationId(x)
        const unreadCount = pickUnreadCount(x)

        return {
          applicationId,
          vacancyTitle: pickVacancyTitle(x),
          companyName: pickCompanyName(x),
          companyAvatar: pickCompanyAvatar(x),
          status: pickStatus(x),
          updatedAt: pickUpdatedAt(x),
          lastMessage: pickLastMessage(x),
          unreadCount,
        }
      })
      .filter((x) => Boolean(x.applicationId))

    hydrateMissingLastMessages().catch(() => {})
  } catch (e: any) {
    const { code, status } = normalizeApiError(e)

    if (String(code).toLowerCase().includes('mfa')) {
      hint.value = t('chatList.enable2fa')
      items.value = []
    } else if (code === 'unauthorized' || status === 401) {
      err.value = getErrorMessage('unauthorized', t)
      items.value = []
    } else if (status === 429 || String(code).includes('rate')) {
      err.value = getErrorMessage('rate_limited', t)
      items.value = []
    } else {
      err.value = getErrorMessage(code || 'load_failed', t, 'errors.loadFailed')
      items.value = []
    }
  } finally {
    pending.value = false
  }
}

const filtered = computed<AppChatViewItem[]>(() => {
  const s = q.value.trim().toLowerCase()

  const enriched: AppChatViewItem[] = items.value.map((x) => ({
    ...x,
    statusLabel: normalizeStatus(x.status),
  }))

  if (!s) return enriched

  return enriched.filter((x) => {
    const a = (x.vacancyTitle || '').toLowerCase()
    const b = (x.companyName || '').toLowerCase()
    return a.includes(s) || b.includes(s)
  })
})

function formatDate(v?: string) {
  if (!v) return t('chatList.status.label')
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString(lang.value === 'ru' ? 'ru-RU' : 'en-US', { day: '2-digit', month: '2-digit' })
}

async function hideChat(applicationId: string) {
  try {
    const res = await hideApplication(applicationId)
    if (res.ok) {
      toast.success(t('chatList.hidden'))
      await load(true)
    } else {
      toast.error(getErrorMessage(res.error || 'hide_failed', t, 'errors.deleteFailed'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'hide_failed', t, 'errors.deleteFailed'))
  }
}

function openDeleteConfirm(id: string) {
  confirmDeleteId.value = id
}

function closeDeleteConfirm() {
  confirmDeleteId.value = null
}

async function confirmDelete() {
  if (!confirmDeleteId.value) return
  const id = confirmDeleteId.value
  confirmDeleteId.value = null
  await hideChat(id)
}

onMounted(() => {
  load()
})
</script>


<style scoped>
.page {
  min-height: 100vh;
  background:
    radial-gradient(1200px 520px at 20% -10%, rgba(79,70,229,0.16), transparent 60%),
    radial-gradient(900px 460px at 90% 10%, rgba(99,102,241,0.12), transparent 60%),
    #ffffff;
}

.wrap {
  max-width: 1100px;
  margin: 0 auto;
  padding: 28px 16px 80px;
}

.hero {
  position: relative;
  overflow: hidden;

  border-radius: 20px;
  padding: 22px;
  margin-bottom: 14px;

  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;

  background:
    radial-gradient(600px 240px at 85% 30%, rgba(99,102,241,0.16), transparent 60%),
    linear-gradient(135deg, rgba(79,70,229,0.10), rgba(99,102,241,0.06));
  border: 1px solid rgba(99,102,241,0.14);
  box-shadow: 0 12px 40px rgba(17,24,39,.06);
}

/* 👇 ИКОНКА ЧАТА */
.hero::after {
  content: '';
  position: absolute;
  pointer-events: none;

  left: -45px;
  top: 50%;
  transform: translateY(-52%);

  width: 1200px;
  height: 1200px;

  background-image: url('/images/chat_panel.webp');
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;

  opacity: 0.42;
  z-index: 0;
}

/* чтобы текст был поверх */
.hero > * {
  position: relative;
  z-index: 1;
}

@media (max-width: 900px) {
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
  margin: 0;
  font-size: 26px;
  font-weight: 900;
  letter-spacing: -0.02em;
}

.hero-sub {
  margin: 6px 0 0;
  font-size: 14px;
  color: #6b7280;
  max-width: 720px;
}

.hero-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.hero-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  padding: 10px 14px;
  font-weight: 900;
  background: rgba(17,24,39,0.03);
  color: #111827;
  border: 1px solid rgba(17,24,39,0.06);
  transition: transform .15s ease, background .15s ease;
}
.hero-btn:hover { transform: translateY(-1px); background: rgba(79,70,229,0.06); }
.hero-btn.primary {
  background: linear-gradient(90deg, #4f46e5, #6366f1);
  color: #fff;
  border-color: transparent;
}
.hero-btn:disabled { opacity: .6; cursor: not-allowed; }

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin: 10px 0 16px;
}

.search {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 12px;
  border-radius: 16px;
  background: rgba(255,255,255,0.9);
  border: 1px solid rgba(17,24,39,0.08);
  box-shadow: 0 10px 24px rgba(17,24,39,0.04);
  flex: 1 1 auto;
  max-width: 640px;
}
.s-ic { opacity: .6; font-weight: 900; }
.s-inp {
  width: 100%;
  border: none;
  outline: none;
  background: transparent;
  font-size: 14px;
}

.meta { color: #6b7280; font-size: 13px; display: flex; gap: 6px; align-items: baseline; }
.muted { color: #6b7280; }
.dot { width: 4px; height: 4px; border-radius: 999px; background: #cbd5e1; display: inline-block; }

.banner {
  padding: 12px 14px;
  border-radius: 16px;
  border: 1px solid rgba(17,24,39,0.08);
  margin-bottom: 12px;
  font-weight: 700;
}
.banner.danger { background: rgba(239,68,68,0.08); border-color: rgba(239,68,68,0.18); color: #991b1b; }
.banner.warn { background: rgba(245,158,11,0.10); border-color: rgba(245,158,11,0.22); color: #92400e; }

.state { padding: 12px 0; }
.skeleton-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}
.sk {
  height: 108px;
  border-radius: 18px;
  background: linear-gradient(90deg, rgba(17,24,39,0.04), rgba(17,24,39,0.08), rgba(17,24,39,0.04));
  background-size: 200% 100%;
  animation: shimmer 1.2s ease-in-out infinite;
}
@keyframes shimmer { 0% { background-position: 0% 0; } 100% { background-position: 200% 0; } }

.empty {
  padding: 22px;
  border-radius: 18px;
  border: 1px dashed rgba(17,24,39,0.14);
  background: rgba(17,24,39,0.02);
  text-align: center;
}
.empty-title { font-weight: 900; margin-bottom: 6px; }
.empty-text { color: #6b7280; margin-bottom: 14px; }
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  padding: 10px 14px;
  font-weight: 900;
  border: 1px solid rgba(17,24,39,0.06);
  background: rgba(17,24,39,0.03);
}
.btn.primary { background: linear-gradient(90deg, #4f46e5, #6366f1); color: #fff; border-color: transparent; }

.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.chat-item {
  position: relative;
  display: flex;
  flex-direction: column;
}

.hide-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  width: 32px;
  height: 32px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  background: #fff;
  color: #9ca3af;
  display: grid;
  place-items: center;
  cursor: pointer;
  opacity: 0;
  transition: opacity .2s ease, transform .2s ease, background .15s ease, border-color .15s ease;
  z-index: 10;
}

.chat-item:hover .hide-btn {
  opacity: 1;
}

.hide-btn:hover {
  border-color: #fecaca;
  color: #b91c1c;
  background: #fff1f2;
  transform: scale(1.05);
}


.confirm-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 80;
  padding: 16px;
}

.confirm-modal {
  background: #fff;
  border-radius: 16px;
  padding: 18px;
  width: 100%;
  max-width: 360px;
  box-shadow: 0 20px 40px rgba(15,23,42,.2);
}

.confirm-title {
  font-weight: 600;
  margin-bottom: 12px;
  color: #111827;
}

.confirm-actions {
  display: flex;
  gap: 10px;
}

.confirm-actions .btn {
  flex: 1;
}

.btn.danger {
  background: #b91c1c;
  color: #fff;
  border: none;
}

.hide-btn:hover {
  transform: scale(1.1);
  background: rgba(220, 38, 38, 1);
}

.card {
  background: #fff;
  border-radius: 18px;
  border: 1px solid rgba(17,24,39,0.06);
  box-shadow: 0 10px 30px rgba(0,0,0,.04);
}

.chat-card {
  padding: 14px 14px;
  transition: transform .16s ease, box-shadow .16s ease, border-color .16s ease;
}
.chat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 18px 40px rgba(17,24,39,0.08);
  border-color: rgba(99,102,241,0.18);
}

.cc-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.cc-left {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 0;
}

.logo {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  overflow: hidden;
  border: 1px solid rgba(99,102,241,0.18);
  background: rgba(79,70,229,0.06);
  display: grid;
  place-items: center;
  flex: 0 0 auto;
}
.logo img { width: 100%; height: 100%; object-fit: cover; }
.logo-fallback { font-weight: 900; color: #3730a3; }

.cc-meta { min-width: 0; }
.cc-title-row { display: flex; align-items: center; gap: 8px; }
.cc-title {
  font-weight: 900;
  color: #0f172a;
  letter-spacing: -0.01em;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 420px;
}
.unread-badge {
  min-width: 20px;
  height: 20px;
  padding: 0 6px;
  border-radius: 999px;
  background: #4f46e5;
  color: #fff;
  font-size: 11px;
  font-weight: 800;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.unread-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #4f46e5;
  box-shadow: 0 0 0 3px rgba(79,70,229,0.12);
  flex: 0 0 auto;
}

.cc-sub {
  margin-top: 4px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.cc-company { font-weight: 800; color: #111827; }

.cc-right { text-align: right; flex: 0 0 auto; }
.cc-date { font-size: 12px; color: #9ca3af; font-weight: 800; }
.cc-open { margin-top: 6px; font-size: 12px; font-weight: 900; color: #4f46e5; }

.cc-last {
  margin-top: 10px;
  font-size: 13px;
  color: #374151;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.cc-last--unread { color: #111827; font-weight: 700; }

@media (max-width: 920px) {
  .grid { grid-template-columns: 1fr; }
  .skeleton-grid { grid-template-columns: 1fr; }
  .hero { flex-direction: column; align-items: flex-start; }
  .hero-actions { justify-content: flex-start; }
}
</style>
