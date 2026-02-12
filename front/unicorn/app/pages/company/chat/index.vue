<template>
  <div class="page">
    <Header />

    <section class="wrap">
      <!-- HERO -->
      <div class="hero">
        <div class="hero-left">
          <h1 class="hero-title">Чаты компании</h1>
          <p class="hero-sub">
            Диалоги с кандидатами по вашим вакансиям. Открой чат, чтобы продолжить общение.
          </p>
        </div>

        <div class="hero-actions">
          <NuxtLink class="hero-btn" to="/company/vacancies">
            ← К вакансиям
          </NuxtLink>
          <button class="hero-btn primary" :disabled="pending" @click="load(true)">
            {{ pending ? 'Обновление…' : 'Обновить' }}
          </button>
        </div>
      </div>

      <!-- TOOLBAR -->
      <div class="toolbar">
        <div class="search">
          <span class="s-ic">⌕</span>
          <input
            v-model="q"
            class="s-inp"
            placeholder="Поиск по вакансии или кандидату…"
          />
        </div>

        <!-- ФИЛЬТР ПО ВАКАНСИИ (строго из чатов компании) -->
        <div class="filters">
          <div class="select-wrap" :class="{ disabled: pending }">
            <select v-model="vacancyFilter" class="select" :disabled="pending">
              <option value="">Все вакансии</option>
              <option v-for="v in vacancyOptions" :key="v.id" :value="v.id">
                {{ v.title }}
              </option>
            </select>
          </div>
        </div>

        <div class="meta">
          <span class="muted">Всего:</span>
          <b>{{ filtered.length }}</b>
        </div>
      </div>

      <!-- ERRORS -->
      <div v-if="err" class="banner danger">
        {{ err }}
      </div>

      <div v-if="hint" class="banner warn">
        {{ hint }}
      </div>

      <!-- LOADING -->
      <div v-if="pending" class="state">
        <div class="skeleton-grid">
          <div v-for="i in 6" :key="i" class="sk card" />
        </div>
      </div>

      <!-- EMPTY -->
      <div v-else-if="!filtered.length" class="empty">
        <div class="empty-title">Диалогов пока нет</div>
        <div class="empty-text">
          Чаты появятся, когда кандидаты откликнутся на ваши вакансии и начнётся переписка.
        </div>
        <NuxtLink class="btn primary" to="/company/vacancies">
          Перейти к вакансиям
        </NuxtLink>
      </div>

      <!-- LIST -->
      <div v-else class="grid">
        <div
          v-for="c in filtered"
          :key="c.applicationId"
          class="chat-item"
        >
          <NuxtLink
            class="card chat-card"
            :to="`/company/chat/${c.applicationId}`"
          >
            <div class="cc-top">
              <div class="cc-left">
                <div class="logo">
                  <img
                    :src="c.candidateAvatar"
                    alt="candidate"
                    @error="onAvatarError($event, '/images/user-base.webp')"
                  />
                </div>

                <div class="cc-meta">
                  <div class="cc-title-row">
                    <div class="cc-title">
                      {{ c.vacancyTitle || 'Вакансия' }}
                    </div>

                    <span
                      v-if="(c.unreadCount || 0) > 0"
                      class="unread-badge"
                      :title="`Непрочитанных: ${c.unreadCount}`"
                    >
                      {{ formatBadge(c.unreadCount || 0) }}
                    </span>
                  </div>

                  <div class="cc-sub">
                    <span class="cc-company">
                      {{ c.candidateName || 'Без имени' }}
                    </span>
                    <span class="dot" />
                    <span class="muted">{{ c.statusLabel }}</span>
                  </div>
                </div>
              </div>

              <div class="cc-right">
                <div class="cc-date">{{ formatDate(c.updatedAt) }}</div>
                <div class="cc-open">Открыть →</div>
              </div>
            </div>

            <!-- ПОСЛЕДНЕЕ СООБЩЕНИЕ + ОТПРАВИТЕЛЬ -->
          <div
            class="cc-last"
            :class="{ 'cc-last--unread': (c.unreadCount || 0) > 0 }"
            v-if="c.lastLine"
          >
            {{ c.lastLine }}
          </div>
          <div class="cc-last muted" v-else>
            <!-- если бэк не отдаёт lastMessage и догрузка ещё не успела -->
            Обновляем историю…
          </div>
        </NuxtLink>

        <button 
          class="hide-btn" 
          @click.prevent="openDeleteConfirm(c.applicationId)"
          title="Скрыть чат"
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
        <div class="confirm-title">Скрыть чат?</div>
        <div class="confirm-actions">
          <button class="btn ghost" @click="closeDeleteConfirm">Отмена</button>
          <button class="btn danger" @click="confirmDelete">Скрыть</button>
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
import { useAuthStore } from '~/stores/auth'
import { useApplications } from '~/composables/useApplications'
import { useToast } from '~/composables/useToast'
import { authFetch } from '~/composables/useAuthFetch'

definePageMeta({ middleware: ['company-auth'] })

const auth = useAuthStore()
const { hideApplication } = useApplications()
const toast = useToast()

const pending = ref(false)
const err = ref('')
const hint = ref('')
const q = ref('')
const vacancyFilter = ref<string>('')

type VacancyOption = { id: string; title: string }

type CompanyChatItem = {
  applicationId: string
  vacancyId?: string
  vacancyTitle?: string
  candidateName?: string
  candidateAvatar?: string
  status?: string
  updatedAt?: string
  unreadCount?: number

  // last message
  lastText?: string
  lastFrom?: string
  lastLine?: string
}

type CompanyChatViewItem = CompanyChatItem & {
  statusLabel: string
}

const items = ref<CompanyChatItem[]>([])
const confirmDeleteId = ref<string | null>(null)

/* ---------------- helpers ---------------- */

function formatBadge(n: number) {
  if (!n || n <= 0) return ''
  return n > 99 ? '99+' : String(n)
}

function normalizeStatus(s?: string) {
  const v = (s || '').toLowerCase()
  if (!v) return 'Статус'
  if (v.includes('pending')) return 'В ожидании'
  if (v.includes('review')) return 'На рассмотрении'
  if (v.includes('accept') || v.includes('approve') || v.includes('offer') || v.includes('hire')) return 'Принят'
  if (v.includes('reject') || v.includes('decline') || v.includes('deny')) return 'Отклонён'
  return 'Статус'
}

function getHttpStatus(e: any): number {
  return Number(e?.statusCode) || Number(e?.status) || Number(e?.response?.status) || 0
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

/* ---------------- data pickers ---------------- */

function pickUnreadCount(a: any): number {
  const n =
    a?.unreadCount ??
    a?.chat?.unreadCount ??
    a?.unread_messages ??
    a?.unread ??
    a?.hasUnread ??
    0

  if (typeof n === 'boolean') return n ? 1 : 0
  const num = Number(n)
  return Number.isFinite(num) ? num : 0
}

function pickCandidateDisplayName(a: any): string {
  // Максимально расширенный поиск displayName
  const v =
    a?.userDisplayName ||
    a?.userName ||
    a?.user?.displayName ||
    a?.user?.profile?.displayName ||
    a?.candidate?.displayName ||
    a?.candidate?.profile?.displayName ||
    a?.applicant?.displayName ||
    a?.applicant?.profile?.displayName ||
    a?.resume?.owner?.displayName ||
    a?.resume?.user?.displayName ||
    a?.from?.displayName ||
    a?.from?.profile?.displayName ||
    a?.user?.name ||
    a?.candidate?.name ||
    a?.candidateName ||
    a?.user?.login ||
    a?.candidate?.login ||
    a?.from?.name ||
    a?.from?.login ||
    a?.userLogin ||
    ''

  return String(v || '').trim()
}

function pickCandidateAvatar(a: any): string {
  return normalizeAvatarUrl(
    a?.userAvatar ||
    a?.userAvatarUrl ||
    a?.user?.avatarUrl ||
    a?.user?.photoUrl ||
    a?.user?.photo ||
    a?.candidate?.avatarUrl ||
    a?.candidate?.photoUrl ||
    a?.candidate?.photo ||
    a?.applicant?.avatarUrl ||
    a?.from?.avatarUrl ||
    a?.from?.photoUrl ||
    '/images/user-base.webp',
    '/images/user-base.webp'
  )
}


function pickVacancyId(a: any): string {
  return String(
    a?.vacancyId ||
      a?.vacancy?.id ||
      a?.vacancy?._id ||
      a?.application?.vacancyId ||
      a?.vacancy_id ||
      ''
  ).trim()
}

function pickVacancyTitle(a: any): string {
  return String(
    a?.vacancy?.title ||
      a?.vacancyTitle ||
      a?.vacancy?.name ||
      a?.vacancy?.position ||
      a?.vacancy_position ||
      'Вакансия'
  ).trim()
}

function pickUpdatedAt(a: any): string {
  return String(
    a?.updatedAt ||
      a?.updated_at ||
      a?.lastMessageAt ||
      a?.last_message_at ||
      a?.createdAt ||
      a?.created_at ||
      ''
  )
}

/* ---------------- last message parsing ---------------- */

function pickMsgText(m: any): string {
  if (!m) return ''
  if (typeof m === 'string') return m.trim()

  return String(m?.text ?? m?.message ?? m?.content ?? m?.body ?? m?.preview ?? '').trim()
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
    } else if (!best && pickMsgText(m)) {
      best = m
    }
  }
  return best
}


function pickUserName(u: any): string {
  if (!u) return ''
  return String(u?.displayName ?? u?.name ?? u?.login ?? u?.username ?? '').trim()
}

function pickUserId(u: any): string {
  if (!u) return ''
  return String(u?.id ?? u?._id ?? u?.userId ?? '').trim()
}

function getMyUserId(): string {
  return String((auth as any)?.user?.id || (auth as any)?.profile?.id || (auth as any)?.me?.id || '')
}

function pickLastMessagePair(a: any): { text: string; from: string } {
  const msgLists = [
    Array.isArray(a?.messages) ? a.messages : null,
    Array.isArray(a?.chat?.messages) ? a.chat.messages : null,
    Array.isArray(a?.chat?.items) ? a.chat.items : null,
    Array.isArray(a?.data?.messages) ? a.data.messages : null,
  ].filter(Boolean) as any[][]

  const allMessages = msgLists.flat()
  const latest = allMessages.length ? pickLatestMessage(allMessages) : null

  const candidates = [
    latest,
    a?.lastMessage,
    a?.last_message,
    a?.chat?.lastMessage,
    a?.chat?.last_message,
    a?.chatLastMessage,
    a?.lastMessageText,
    a?.last_message_text,
    a?.messagePreview,
    a?.preview,
  ]

  let msg: any = null
  let text = ''
  for (const c of candidates) {
    const t = pickMsgText(c)
    if (t) {
      msg = c
      text = t
      break
    }
  }

  if (!text) return { text: '', from: '' }

  const senderObj =
    (typeof msg === 'object'
      ? msg?.from || msg?.sender || msg?.author || msg?.user || msg?.createdBy || null
      : null) ||
    a?.lastMessageFrom ||
    a?.chat?.lastMessageFrom ||
    null

  let from = ''
  if (typeof senderObj === 'string') {
    from = senderObj.trim()
  } else {
    const senderName = pickUserName(senderObj)
    const senderId = pickUserId(senderObj)
    const myId = getMyUserId()
    if (myId && senderId && myId === senderId) from = 'Вы'
    else from = senderName
  }

  if (!from && typeof msg === 'object') {
    const role = String(msg?.role ?? msg?.senderRole ?? msg?.fromRole ?? msg?.senderType ?? '').toLowerCase()
    if (role.includes('company') || role.includes('employer')) from = 'Вы'
    if (role.includes('user') || role.includes('candidate')) from = pickCandidateDisplayName(a)
  }

  return { text, from }
}





function buildLastLine(a: any): { lastText: string; lastFrom: string; lastLine: string } {
  const { text, from } = pickLastMessagePair(a)
  if (!text) return { lastText: '', lastFrom: '', lastLine: '' }
  return { lastText: text, lastFrom: from, lastLine: from ? `${from}: ${text}` : text }
}

/* ---------------- fetching ---------------- */

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
      break
    }
  }
  throw lastErr || new Error('no_endpoint_matched')
}

async function tryFetchObject(urls: string[]) {
  let lastErr: any = null
  for (const url of urls) {
    try {
      const res: any = await authFetch(url, {
        method: 'GET'
      })
      if (res && typeof res === 'object') return res
    } catch (e: any) {
      lastErr = e
      const st = getHttpStatus(e)
      if (st === 404) continue
      if (st === 401 || st === 403) continue
      break
    }
  }
  throw lastErr || new Error('no_endpoint_matched')
}

// best-effort догрузка последнего сообщения, если inbox его не отдаёт
async function hydrateMissingLastMessages() {
  const need = items.value.filter((x) => !x.lastLine).slice(0, 20)
  if (!need.length) return

  // Try to resolve a working messages endpoint.
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
      // если не 404 — тоже выходим, чтобы не DDOS-ить
      break
    }
  }

  // НЕТ рабочего endpoint → выходим (иначе будет вечный спам 404)
  if (!workingBase) return

  for (const it of need) {
    const url = workingBase.replace('{id}', it.applicationId)
    try {
      const res: any = await authFetch(url, {
        method: 'GET'
      })

      const list = Array.isArray(res?.items) ? res.items : []
      const lastMsg = list.length ? list[list.length - 1] : null
      const lastText = pickMsgText(lastMsg)
      if (lastText) {
        const from = String(lastMsg?.senderType || '') === 'company' ? '??' : pickCandidateDisplayName(it)
        it.lastText = lastText
        it.lastFrom = from
        it.lastLine = from ? `${from}: ${lastText}` : lastText
      }

      const u = String(lastMsg?.createdAt || lastMsg?.created_at || res?.updatedAt || res?.updated_at || '').trim()
      if (u) it.updatedAt = u
    } catch {
      // игнор
    }
  }

  items.value = [...items.value]
}


/* ---------------- main load ---------------- */

async function load(force = false) {
  if (pending.value && !force) return

  pending.value = true
  err.value = ''
  hint.value = ''

  try {
    if (process.client && !(auth as any).isReady && (auth as any)?.hasToken) {
      try { await (auth as any).init?.() } catch {}
    }

    // ВАЖНО: НЕ начинаем с company/*, пока не знаем что они существуют.
    // Ставь сюда те URL, которые у тебя реально не 404 (по Network).
    const raw = await tryFetchArray([
      '/api/applications/inbox',
      '/api/applications/company/inbox',
      '/api/company/applications/inbox',
      '/api/company/chats',
      '/api/chats/company',
    ])

    items.value = raw
      .map((a: any): CompanyChatItem => {
        const unreadCount = pickUnreadCount(a)
        const last = buildLastLine(a)

        return {
          applicationId: String(a?.id || a?.applicationId || a?._id || '').trim(),
          vacancyId: pickVacancyId(a),
          vacancyTitle: pickVacancyTitle(a),
          candidateName: pickCandidateDisplayName(a) || 'Без имени',
          candidateAvatar: pickCandidateAvatar(a),
          status: String(a?.status || a?.state || ''),
          updatedAt: pickUpdatedAt(a),
          unreadCount,
          lastText: last.lastText,
          lastFrom: last.lastFrom,
          lastLine: last.lastLine,
        }
      })
      .filter((x) => Boolean(x.applicationId))

    // pending снимаем СРАЗУ — UI больше не “вечный”
    pending.value = false

    // догрузка lastMessage — только если нужно, и НЕ блокируем страницу
    hydrateMissingLastMessages().catch(() => {})
  } catch (e: any) {
    const { code, status } = normalizeApiError(e)

    if (status === 404) {
      err.value = 'inbox_endpoint_404'
    } else if (String(code).toLowerCase().includes('mfa')) {
      hint.value = 'Включи 2FA, чтобы использовать чаты'
    } else if (code === 'unauthorized' || status === 401) {
      err.value = getErrorMessage('unauthorized', t)
    } else if (status === 429 || String(code).includes('rate')) {
      err.value = getErrorMessage('rate_limited', t)
    } else {
      err.value = getErrorMessage(code || 'load_failed', t, 'errors.loadFailed')
    }

    items.value = []
    pending.value = false
  }
}


/* ---------------- vacancy options (строго из чатов компании) ---------------- */

const vacancyOptions = computed<VacancyOption[]>(() => {
  const map = new Map<string, string>()
  for (const it of items.value) {
    const id = String(it.vacancyId || '').trim()
    if (!id) continue
    const title = String(it.vacancyTitle || 'Вакансия').trim()
    if (!map.has(id)) map.set(id, title)
  }
  return Array.from(map.entries())
    .map(([id, title]) => ({ id, title }))
    .sort((a, b) => a.title.localeCompare(b.title, 'ru'))
})

/* ---------------- filtered view ---------------- */

const filtered = computed<CompanyChatViewItem[]>(() => {
  const s = q.value.trim().toLowerCase()
  const vf = vacancyFilter.value

  const enriched: CompanyChatViewItem[] = items.value.map((x) => ({
    ...x,
    statusLabel: normalizeStatus(x.status),
  }))

  const byVacancy = vf ? enriched.filter((x) => String(x.vacancyId || '') === vf) : enriched
  if (!s) return byVacancy

  return byVacancy.filter((x) => {
    const a = (x.vacancyTitle || '').toLowerCase()
    const b = (x.candidateName || '').toLowerCase()
    return a.includes(s) || b.includes(s)
  })
})

function formatDate(v?: string) {
  if (!v) return ''
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return ''
  return d.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' })
}

async function hideChat(applicationId: string) {
  try {
    const res = await hideApplication(applicationId)
    if (res.ok) {
      toast.success('Чат скрыт')
      await load(true)
    } else {
      toast.error(res.error || 'Ошибка скрытия')
    }
  } catch (e: any) {
    toast.error(e?.message || 'Ошибка скрытия')
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

onMounted(async () => {
  await load()
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

@media (max-width: 920px) {
  .grid { grid-template-columns: 1fr; }
  .skeleton-grid { grid-template-columns: 1fr; }
  .hero { flex-direction: column; align-items: flex-start; }
  .hero-actions { justify-content: flex-start; }
}

.cc-title-row { position: relative; }

.unread-badge{
  display:inline-flex;
  align-items:center;
  justify-content:center;
  height:20px;
  min-width:20px;
  padding:0 6px;
  border-radius:999px;
  background:#4f46e5;
  color:#fff;
  font-weight:900;
  font-size:12px;
  line-height:1;
  box-shadow:0 8px 18px rgba(79,70,229,.25);
}

.cc-last--unread{
  color:#111827;
  font-weight:800;
}
/* ---------------- FILTERS (vacancy select) ---------------- */

.filters{
  flex: 0 0 auto;
  display: flex;
  align-items: center;
  gap: 10px;
}

.select{
  height: 44px; /* вровень с search */
  padding: 0 12px;
  border-radius: 16px;

  background: rgba(255,255,255,0.9);
  border: 1px solid rgba(17,24,39,0.08);
  box-shadow: 0 10px 24px rgba(17,24,39,0.04);

  font-size: 14px;
  font-weight: 800;
  color: #111827;

  outline: none;
  cursor: pointer;

  min-width: 260px;
  max-width: 380px;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;

  transition: transform .15s ease, border-color .15s ease, background .15s ease;
}

.select:hover{
  transform: translateY(-1px);
  border-color: rgba(99,102,241,0.18);
  background: rgba(255,255,255,1);
}

.select:focus{
  border-color: rgba(79,70,229,0.35);
  box-shadow:
    0 10px 24px rgba(17,24,39,0.04),
    0 0 0 4px rgba(79,70,229,0.10);
}

.select:disabled{
  opacity: .6;
  cursor: not-allowed;
}

/* если в option очень длинные названия — не даём разъехать */
.select option{
  font-weight: 800;
}

/* ---------------- RESPONSIVE toolbar ---------------- */

@media (max-width: 920px) {
  .toolbar{
    flex-wrap: wrap;
    align-items: stretch;
  }

  .search{
    max-width: 100%;
    flex: 1 1 100%;
  }

  .filters{
    flex: 1 1 100%;
  }

  .select{
    width: 100%;
    min-width: 0;
    max-width: none;
  }

  .meta{
    flex: 1 1 100%;
    justify-content: flex-end;
  }
}
/* ===== Select: стабильный вид + стрелка (Firefox/Chrome) ===== */
.select-wrap{
  position: relative;
  display: inline-flex;
  width: 100%;
  max-width: 380px;
}
.select-wrap::after{
  content: '';
  position: absolute;
  right: 14px;
  top: 50%;
  width: 10px;
  height: 10px;
  transform: translateY(-50%);
  pointer-events: none;
  opacity: .65;
  background-repeat: no-repeat;
  background-size: 10px 10px;
  /* маленькая стрелка вниз */
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24'%3E%3Cpath fill='%23111827' d='M7 10l5 5 5-5z'/%3E%3C/svg%3E");
}
.select{
  -webkit-appearance: none;
  -moz-appearance: none;
  appearance: none;
  padding-right: 38px; /* место под стрелку */
}
.select-wrap.disabled{
  opacity: .7;
}

/* ===== Cards: адаптив под мобилу ===== */
@media (max-width: 560px) {
  .wrap{
    padding: 18px 12px 70px;
  }

  .hero{
    padding: 16px;
    border-radius: 18px;
  }

  .hero-title{
    font-size: 22px;
  }

  .hero-actions{
    width: 100%;
  }

  .hero-btn{
    flex: 1 1 auto;
  }

  .grid{
    gap: 10px;
  }

  .chat-card{
    padding: 12px;
  }

  .cc-top{
    align-items: flex-start;
    gap: 10px;
  }

  .cc-right{
    text-align: left;
    display: flex;
    align-items: baseline;
    gap: 10px;
    margin-left: 56px; /* визуально под метаданные */
  }

  .cc-open{
    margin-top: 0;
    font-weight: 900;
    background: rgba(79,70,229,0.06);
    border: 1px solid rgba(99,102,241,0.18);
    padding: 6px 10px;
    border-radius: 999px;
    display: inline-flex;
    width: fit-content;
  }

  .logo{
    width: 40px;
    height: 40px;
    border-radius: 12px;
  }

  .cc-title{
    max-width: 100%;
  }

  .cc-last{
    margin-top: 8px;
  }
}

</style>
