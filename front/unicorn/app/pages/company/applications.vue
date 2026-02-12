<template>
  <Header />

  <div class="page p-6 space-y-4">
    <!-- HEADER -->
    <div class="header">
      <h1 class="text-2xl font-semibold">{{ t('companyApps.title') }}</h1>
    </div>

    <!-- FILTERS -->
        <div class="filters">
      <button
        v-for="f in filters"
        :key="f.value"
        class="filter-btn"
        :class="{ active: status === f.value }"
        @click="setStatus(f.value)"
      >
        {{ f.label }}
      </button>
    </div>

    <!-- ALERTS -->
    <div v-if="err" class="alert error">{{ err }}</div>
    <div v-if="hint" class="alert hint">{{ hint }}</div>

    <!-- LIST -->
    <div class="list">
      <div
        v-for="a in items"
        :key="a.applicationId"
        class="card rounded p-4 space-y-2 transition"
        :class="{ unread: !a.viewed, 'premium-user': a.userIsPremium }"
        :data-premium-label="t('companyApps.premiumUser')"
        @click="isMobile && (selectedApplication = a)"
      >
        <div class="card-body">
          <!-- INFO -->
          <div class="info">
            <div class="title">
              <span v-if="a.status === 'pending'" class="badge-new">{{ t('companyApps.newResumeBadge') }}</span>
            </div>

            <div class="text-sm opacity-70">
              {{ t('companyApps.from') }} <b>{{ a.userDisplayName }}</b>
            </div>

            <div class="text-sm opacity-70">
              {{ t('companyApps.vacancy') }} <b>{{ a.vacancyTitle }}</b>
            </div>

            <div class="text-sm opacity-70">
              {{ t('companyApps.resume') }} <b>{{ a.resumeTitle }}</b>
            </div>
            <div class="status">
              {{ t('companyApps.status') }}
              <b :class="`status-${a.status}`">
                {{ statusLabel(a.status) }}
              </b>
            </div>
            <button
              class="text-xs px-2 py-1 rounded border opacity-70 hover:opacity-100"
              @click="toggleIds(a.applicationId)"
              type="button"
            >
              {{ t('companyApps.showId') }}
            </button>
            <div v-if="isIdsOpen(a.applicationId)" class="mt-2 p-2 rounded bg-black/5 text-xs font-mono space-y-1">
              <div>{{ t('companyApps.showIdApplication') }} {{ a.applicationId }}</div>
              <div>{{ t('companyApps.showIdUser') }} {{ a.userId }}</div>
              <div>{{ t('companyApps.showIdVacancy') }} {{ a.vacancyId }}</div>
              <div>{{ t('companyApps.showIdResume') }} {{ a.resumeId }}</div>
            </div>
            <div
              v-if="a.message"
              class="message-preview whitespace-pre-wrap"
            >
              {{ a.message }}
            </div>
          </div>

          <!-- ACTIONS (desktop) -->
          <div class="actions" v-if="!isMobile">
            <NuxtLink
              class="btn chat-btn"
              :to="`/company/chat/${a.applicationId}`"
              @click.stop="markViewed(a)"
            >
              {{ t('companyApps.chat') }}
              <span v-if="a.unreadCount" class="count-badge">{{ a.unreadCount }}</span>
            </NuxtLink>

            <button class="btn" @click.stop="openResume(a.resumeId, a)">
              {{ t('companyApps.viewResume') }}
            </button>

            <button
              v-if="a.status === 'pending'"
              class="btn primary"
              @click.stop="accept(a)"
            >
              {{ t('companyApps.accept') }}
            </button>

            <button
              v-if="a.status === 'pending'"
              class="btn danger"
              @click.stop="reject(a)"
            >
              {{ t('companyApps.reject') }}
            </button>

            <button
              class="btn"
              @click.stop="hide(a)"
              :title="t('companyApps.hide')"
            >
              {{ t('companyApps.hide') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- RESUME MODAL / BOTTOM SHEET -->
    <div
      v-if="resumeModal.open"
      class="overlay"
      :class="{ mobile: isMobile }"
        @click.self="closeResume"
    >
      <div class="modal">
       
        <div class="modal-header">
          <div class="font-semibold">{{ t('companyApps.resumeModal.title') }}</div>
          <button class="close" @click="closeResume">x</button>
        </div>
      
        <div v-if="resumeModal.loading" class="loader">
          {{ t('companyApps.resumeModal.loading') }}
        </div>

        <div v-else-if="resumeModal.error" class="alert error">
          {{ resumeModal.error }}
        </div>

        <div v-else-if="resumeModal.data" class="modal-body">
          <div class="field">
            <div class="label font-bold pt-6">{{ t('companyApps.resumeModal.resumeId') }}</div>
            <div class="mono">{{ resumeModal.resumeId }}</div>
          </div>

          <div class="field">
            <div class="label font-bold pt-6">{{ t('companyApps.resumeModal.titleLabel') }}</div>
            <div>{{ resumeModal.data.title || t('profilePage.notSpecified') }}</div>
          </div>

          <div class="field">
            <div class="label font-bold pt-6">{{ t('companyApps.resumeModal.info') }}</div>
            <div class="whitespace-pre-wrap">
              {{ resumeModal.data.about || t('profilePage.notSpecified') }}
            </div>
          </div>

          <div class="grid-2">
            <div class="field">
              <div class="label font-bold pt-6">{{ t('companyApps.resumeModal.location') }}</div>
              <div>{{ resumeModal.data.location || t('profilePage.notSpecified') }}</div>
            </div>

            <div class="field">
              <div class="label font-bold pt-6">{{ t('companyApps.resumeModal.industry') }}</div>
              <div>{{ resumeModal.data.industry || t('profilePage.notSpecified') }}</div>
            </div>
          </div>

          <div
            v-if="Array.isArray(resumeModal.data.skills)"
            class="field"
          >
            <div class="label font-bold pt-6">{{ t('companyApps.resumeModal.skills') }}</div>
            <div>{{ resumeModal.data.skills.join(', ') }}</div>
          </div>

          <div
            v-if="Array.isArray(resumeModal.data.links)"
            class="field"
          >
            <div class="label font-bold pt-6">{{ t('companyApps.resumeModal.links') }}</div>
            <ul class="links">
              <li v-for="l in resumeModal.data.links" :key="l">
                <a :href="l" target="_blank">{{ l }}</a>
              </li>
            </ul>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn" @click="closeResume">
            {{ t('companyApps.resumeModal.close') }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- STICKY ACTIONS (mobile only) -->
  <div
    v-if="isMobile && selectedApplication"
    class="mobile-sticky-actions"
  >
    <NuxtLink
      class="btn chat-btn"
      :to="`/company/chat/${selectedApplication.applicationId}`"
    >
      {{ t('companyApps.mobileChat') }}
      <span v-if="selectedApplication.unreadCount" class="count-badge">{{ selectedApplication.unreadCount }}</span>
    </NuxtLink>

    <button
      class="btn"
      @click="openResume(selectedApplication.resumeId, selectedApplication)"
    >
      {{ t('companyApps.viewResume') }}
    </button>

    <button
      v-if="selectedApplication.status === 'pending'"
      class="btn primary"
      @click="accept(selectedApplication)"
    >
      {{ t('companyApps.mobileAccept') }}
    </button>

    <button
      v-if="selectedApplication.status === 'pending'"
      class="btn danger"
      @click="reject(selectedApplication)"
    >
      {{ t('companyApps.mobileReject') }}
    </button>

    <button
      class="btn"
      @click="hide(selectedApplication)"
    >
      {{ t('companyApps.hide') }}
    </button>
  </div>

  <Footer />
</template>

<script setup lang="ts">
import axios from 'axios'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import { useApplications } from '~/composables/useApplications'
import { useToast } from '~/composables/useToast'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'

definePageMeta({ middleware: ['company-auth'] })

const { hideApplication } = useApplications()
const toast = useToast()
const { t } = useI18n()

const err = ref('')
const hint = ref('')
const status = ref('')

const filters = computed(() => ([
  { label: t('companyApps.filters.all'), value: '' },
  { label: t('companyApps.filters.pending'), value: 'pending' },
  { label: t('companyApps.filters.accepted'), value: 'accepted' },
  { label: t('companyApps.filters.rejected'), value: 'rejected' },
]))

const selectedApplication = ref<App | null>(null)

const isMobile = computed(() => {
  if (process.server) return false
  return window.innerWidth <= 768
})

type App = {
  applicationId: string
  vacancyId: string
  vacancyTitle: string
  resumeId: string
  resumeTitle: string
  userId: string
  userDisplayName: string
  status: 'pending' | 'accepted' | 'rejected'
  message?: string
  viewed?: boolean
  companyId: string
  userIsPremium?: boolean
  unreadCount?: number
}
/* функция скрытия и показа ID */
const items = ref<App[]>([])

const showIds = ref<Record<string, boolean>>({})

function toggleIds(appId: string) {
  showIds.value[appId] = !showIds.value[appId]
}
function isIdsOpen(appId: string) {
  return !!showIds.value[appId]
}

function statusLabel(statusValue: App['status']) {
  const key = `companyApps.filters.${statusValue}`
  const label = t(key)
  return label === key ? statusValue : label
}

function pickUnreadCount(x: any): number {
  const n =
    x?.unreadCount ??
    x?.chat?.unreadCount ??
    x?.chat?.unread_count ??
    x?.unread_messages ??
    x?.unread ??
    x?.hasUnread ??
    0

  if (typeof n === 'boolean') return n ? 1 : 0
  const num = Number(n)
  return Number.isFinite(num) ? num : 0
}

/* RESUME MODAL */
const resumeModal = reactive({
  open: false,
  loading: false,
  error: '',
  resumeId: '',
  data: null as any,
})

async function load() {
  err.value = ''
  hint.value = ''
  try {
    const res = await axios.get('/applications/inbox', {
      params: status.value ? { status: status.value } : {},
    })
    const raw = Array.isArray(res.data.items) ? res.data.items : []
    items.value = raw.map((x: any) => ({
      ...x,
      unreadCount: pickUnreadCount(x),
    }))
  } catch (e: any) {
    const code = e?.response?.data?.error
    if (code === 'mfa_required') {
      hint.value = getErrorMessage('mfa_required', t, 'companyApps.errors.mfaRequired')
    } else {
      err.value = getErrorMessage(code || 'load_failed', t, 'companyApps.errors.loadFailed')
    }
  }
}

function setStatus(s: string) {
  status.value = s
  load()
}

async function markViewed(a: App) {
  if (a.viewed) return
  a.viewed = true
  try {
    await axios.post(`/api/applications/${a.applicationId}/viewed`)
  } catch {}
}

async function accept(a: App) {
  await axios.post(`/api/applications/${a.applicationId}/accept`)
  await load()
}

async function reject(a: App) {
  await axios.post(`/api/applications/${a.applicationId}/reject`)
  await load()
}

async function openResume(resumeId: string, a: App) {
  markViewed(a)
  resumeModal.open = true
  resumeModal.loading = true
  resumeModal.error = ''
  resumeModal.data = null
  resumeModal.resumeId = resumeId

  try {
    const res = await axios.get(`/resumes/${resumeId}`)
    resumeModal.data = res.data.resume
  } catch (e: any) {
    resumeModal.error = getErrorMessage(e?.response?.data?.error || 'resume_load_failed', t, 'companyApps.errors.resumeLoadFailed')
  } finally {
    resumeModal.loading = false
  }
}

function closeResume() {
  resumeModal.open = false
}

async function hide(a: App) {
  if (!confirm(t('companyApps.hideConfirm'))) return

  try {
    const res = await hideApplication(a.applicationId)
    if (res.ok) {
      toast.success(t('companyApps.hideSuccess'))
      await load()
    } else {
      toast.error(getErrorMessage(res.error || 'hide_failed', t, 'companyApps.errors.hideFailed'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'hide_failed', t, 'companyApps.errors.hideFailed'))
  }
}

onMounted(load)
</script>

<style scoped>
.page {
  max-width: 1100px;
  margin: 0 auto;
  padding-bottom: 0;
}

/* FILTERS */
.filters {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 6px 12px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  background: #fff;
}

.filter-btn.active {
  background: #111827;
  color: #fff;
}

/* LIST */
.list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* CARD */
.card {
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  background: #fff;
}

.card.unread {
  background: #f0fdf4;
  border-color: #86efac;
}

.card.premium-user {
  background: linear-gradient(135deg, rgba(79,70,229,0.08) 0%, #ffffff 60%);
  border: 1px solid rgba(99,102,241,0.35);
  box-shadow: 0 12px 28px rgba(79,70,229,0.12);
  position: relative;
}

.card.premium-user::before {
  content: attr(data-premium-label);
  position: absolute;
  top: 8px;
  right: 8px;
  background: linear-gradient(90deg, #4f46e5, #6366f1);
  color: #fff;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  z-index: 10;
  pointer-events: none;
}

.card.premium-user .card-body {
  padding-top: 32px;
}

.card.premium-user:hover {
  border-color: rgba(79,70,229,0.55);
  box-shadow: 0 18px 36px rgba(79,70,229,0.18);
}

.card-body {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
}

/* BADGE */
.badge-new {
  background: #22c55e;
  color: #fff;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 999px;
}

/* ACTIONS */
.actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* BUTTON */
.btn {
  padding: 8px 12px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  background: #fff;
}

.chat-btn {
  position: relative;
}

.count-badge {
  position: absolute;
  top: -6px;
  right: -6px;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 999px;
  background: #ef4444;
  color: #fff;
  font-size: 11px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #fff;
}

.btn.primary {
  background: #111827;
  color: #fff;
}

.btn.danger {
  border-color: #fecaca;
  color: #b91c1c;
}

/* MODAL / BOTTOM SHEET */
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.4);
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 16px;
  z-index: 50;
}

.overlay.mobile {
  align-items: flex-end;
}

.overlay.mobile .modal {
  border-radius: 18px 18px 0 0;
  animation: slideUp .25s ease;
}

@keyframes slideUp {
  from { transform: translateY(100%) }
  to { transform: translateY(0) }
}

.modal {
  background: #fff;
  max-width: 640px;
  width: 100%;
  border-radius: 18px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}
.modal-header .close {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1px solid #e5e7eb;
  background: #fff;
  cursor: pointer;
}
.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}
/* Custom scrollbar for modal */
.modal-body::-webkit-scrollbar {
  width: 8px;
}

.modal-body::-webkit-scrollbar-track {
  background: #f3f4f6;
  border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 4px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: #9ca3af;
}

.modal-body .loader {
  padding: 40px;
  text-align: center;
  color: #6b7280;
}

.modal-body .alert {
  margin: 16px 0;
}

.modal-footer {
  padding: 16px;
  border-top: 1px solid #e5e7eb;
}
/* STICKY ACTIONS */
.mobile-sticky-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;

  display: flex;
  gap: 8px;
  padding: 12px;
  flex-wrap: wrap;

  background: rgba(255,255,255,.96);
  backdrop-filter: blur(10px);
  border-top: 1px solid #e5e7eb;
  z-index: 60;
  padding-bottom: calc(12px + env(safe-area-inset-bottom));
}

.mobile-sticky-actions .btn {
  flex: 1 1 calc(50% - 8px);
  min-width: 140px;
}

/* MOBILE */
@media (max-width: 768px) {
  .page {
    padding: 16px;
    padding-bottom: 120px; /* leave room for sticky action bar */
  }

  .card-body {
    flex-direction: column;
  }

  .page {
    padding-bottom: 90px;
  }
}
</style>
