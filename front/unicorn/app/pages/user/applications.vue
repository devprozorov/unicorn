<template>
  <Header />

  <div class="page p-6 space-y-6">
    <!-- HEADER -->
    <div class="header">
      <h1 class="text-2xl font-semibold">{{ t('apps.title') }}</h1>
    </div>

    <!-- ERROR -->
    <div v-if="err" class="alert error">
      {{ err }}
    </div>

    <!-- EMPTY -->
    <div v-if="!items.length && !err" class="empty">
      <div class="empty-title">{{ t('apps.noAppsTitle') }}</div>
      <div class="empty-sub">
        {{ t('apps.noAppsSub') }}
      </div>
    </div>

    <!-- LIST -->
    <div class="list">
      <div
        v-for="a in items"
        :key="a.applicationId"
        class="card"
        :class="[`status-${a.status}`, { active: activeId === a.applicationId }]"
        @click="select(a)"
      >
        <div class="card-body">
          <!-- INFO -->
          <div class="info">
            <!-- STATUS -->
            <div class="status-line">
              <span class="status-badge">
                {{ statusLabel(a.status) }}
              </span>

              <span
                v-if="a.createdAt"
                class="time"
              >
                · {{ timeAgo(a.createdAt) }}
              </span>
            </div>

            <!-- TITLE -->
            <div class="title">
              {{ a.vacancyTitle }}
            </div>

            <!-- META -->
            <div class="meta">
              {{ t('apps.companyLabel') }}
              <b>{{ a.companyDisplayName }}</b>
            </div>

            <!-- TIMELINE -->
            <div class="timeline">
              <div class="step done">{{ t('apps.timeline.applied') }}</div>
              <div
                class="step"
                :class="{
                  done: a.status !== 'pending',
                  current: a.status === 'pending'
                }"
              >
                {{ t('apps.timeline.reviewed') }}
              </div>
              <div
                class="step"
                :class="{
                  done: a.status === 'accepted',
                  rejected: a.status === 'rejected'
                }"
              >
                {{ t('apps.timeline.result') }}
              </div>
            </div>
          </div>

          <!-- ACTIONS (DESKTOP) -->
          <div class="actions desktop">
            <NuxtLink
              class="btn"
              :to="`/jobs/${a.vacancyId}`"
            >
              {{ t('apps.viewVacancy') }}
            </NuxtLink>

            <NuxtLink
              class="btn primary chat-btn"
              :to="`/user/chat/${a.applicationId}`"
            >
              {{ t('apps.chat') }}
              <span v-if="a.unreadCount" class="count-badge">{{ a.unreadCount }}</span>
            </NuxtLink>

            <button
              class="btn"
              @click.stop="hide(a)"
              :title="t('apps.hide')"
            >
              {{ t('apps.hide') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- STICKY ACTIONS (MOBILE) -->
  <div
    v-if="activeItem"
    class="sticky-actions"
  >
    <NuxtLink
      class="btn"
      :to="`/jobs/${activeItem.vacancyId}`"
    >
      {{ t('apps.vacancy') }}
    </NuxtLink>

    <NuxtLink
      class="btn primary chat-btn"
      :to="`/user/chat/${activeItem.applicationId}`"
    >
      {{ t('apps.chat') }}
      <span v-if="activeItem.unreadCount" class="count-badge">{{ activeItem.unreadCount }}</span>
    </NuxtLink>

    <button
      class="btn"
      @click="hide(activeItem)"
    >
      {{ t('apps.hide') }}
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

definePageMeta({ middleware: ['user-auth'] })

const { hideApplication } = useApplications()
const toast = useToast()
const { t } = useI18n()

const err = ref('')
const activeId = ref<string | null>(null)

type Item = {
  applicationId: string
  vacancyId: string
  companyId: string
  status: 'pending' | 'accepted' | 'rejected'
  vacancyTitle?: string
  createdAt?: string
  companyDisplayName: string
  unreadCount?: number
}

const items = ref<Item[]>([])

const activeItem = computed(() =>
  items.value.find(i => i.applicationId === activeId.value) || null
)

function select(a: Item) {
  if (window.innerWidth <= 768) {
    activeId.value = a.applicationId
  }
}

function timeAgo(date: string) {
  const diff = Date.now() - new Date(date).getTime()
  const days = Math.floor(diff / 86400000)
  if (days === 0) {
    const label = t('apps.time.today')
    return label === 'apps.time.today' ? 'today' : label
  }
  if (days === 1) {
    const label = t('apps.time.dayAgo')
    return label === 'apps.time.dayAgo' ? '1 day ago' : label
  }
  const label = t('apps.time.daysAgo')
  if (label === 'apps.time.daysAgo') return `${days} days ago`
  return label.replace('{{count}}', String(days))
}

function statusLabel(status: Item['status']) {
  const key = `apps.status.${status}`
  const label = t(key)
  return label === key ? status : label
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

async function load() {
  err.value = ''
  try {
    const res = await axios.get('/applications/my')
    const raw = Array.isArray(res.data.items) ? res.data.items : []
    items.value = raw.map((x: any) => ({
      ...x,
      unreadCount: pickUnreadCount(x),
    }))
  } catch (e: any) {
    const code = e?.response?.data?.error
    if (code === 'mfa_required') {
      err.value = getErrorMessage('mfa_required', t, 'apps.errors.mfaRequired')
    } else {
      err.value = getErrorMessage(code || 'load_failed', t, 'apps.errors.loadFailed')
    }
  }
}

async function hide(a: Item) {
  const confirmText = t('apps.hideConfirm')
  if (!confirm(confirmText)) return

  try {
    const res = await hideApplication(a.applicationId)
    if (res.ok) {
      toast.success(t('apps.hideSuccess'))
      await load()
      activeId.value = null
    } else {
      toast.error(getErrorMessage(res.error || 'hide_failed', t, 'apps.errors.hideFailed'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'hide_failed', t, 'apps.errors.hideFailed'))
  }
}

onMounted(load)
</script>

<style scoped>
/* PAGE */
.page {
  max-width: 900px;
  margin: 0 auto;
  padding-bottom: 0;
}

/* ALERT */
.alert.error {
  background: #fee2e2;
  color: #991b1b;
  padding: 12px;
  border-radius: 14px;
}

/* EMPTY */
.empty {
  text-align: center;
  padding: 48px 16px;
  border: 1px dashed #e5e7eb;
  border-radius: 18px;
}

.empty-title {
  font-weight: 600;
  font-size: 16px;
}

.empty-sub {
  margin-top: 6px;
  font-size: 14px;
  color: #6b7280;
}

/* LIST */
.list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* CARD */
.card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  transition: transform .15s ease, box-shadow .15s ease;
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 28px rgba(0,0,0,.06);
}

.card.active {
  border-color: #111827;
}

/* STATUS COLORS */
.card.status-pending { border-left: 4px solid #facc15; }
.card.status-accepted { border-left: 4px solid #22c55e; }
.card.status-rejected { border-left: 4px solid #ef4444; }

/* BODY */
.card-body {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
}

/* INFO */
.status-line {
  font-size: 12px;
  display: flex;
  gap: 6px;
  align-items: center;
}

.time {
  color: #6b7280;
}

.status-badge {
  padding: 4px 10px;
  border-radius: 999px;
  background: #f3f4f6;
  text-transform: capitalize;
}

.title {
  margin-top: 6px;
  font-weight: 600;
}

.meta {
  margin-top: 4px;
  font-size: 13px;
  color: #6b7280;
}

/* TIMELINE */
.timeline {
  margin-top: 10px;
  display: flex;
  gap: 12px;
  font-size: 12px;
}

.step {
  opacity: .4;
}

.step.done {
  opacity: 1;
  color: #16a34a;
}

.step.current {
  opacity: 1;
  color: #92400e;
}

.step.rejected {
  opacity: 1;
  color: #991b1b;
}

/* ACTIONS */
.actions.desktop {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 150px;
}

.btn {
  padding: 8px 12px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  background: #fff;
  font-size: 14px;
  text-align: center;
}

.btn.primary {
  background: #111827;
  color: #fff;
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

/* STICKY MOBILE ACTIONS */
.sticky-actions {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  gap: 10px;
  padding: 12px;
  background: #fff;
  border-top: 1px solid #e5e7eb;
  z-index: 50;
  padding-bottom: calc(12px + env(safe-area-inset-bottom));
}

/* MOBILE */
@media (max-width: 768px) {
  .page {
    padding: 16px;
    padding-bottom: 140px; /* space for sticky actions */
  }

  .card-body {
    flex-direction: column;
  }

  .actions.desktop {
    display: none;
  }

  .title {
    font-size: 16px;
  }

  .timeline {
    flex-wrap: wrap;
    gap: 8px;
  }
}
</style>
