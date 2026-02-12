<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useProfileApi } from '~/services/profileApi'
import { useResumesApi } from '~/services/resumesApi'
import { useI18n } from '~/composables/useI18n'
import { useAuthStore } from '~/stores/auth'
import { authFetch } from '~/composables/useAuthFetch'

import ProfileHeader from './ProfileHeader.vue'
import ProfileSkeleton from './ProfileSkeleton.vue'
import ResumeMiniCard from './ResumeMiniCard.vue'

/**
 * API
 */
const profileApi = useProfileApi()
const resumesApi = useResumesApi()
const { t } = useI18n()
const auth = useAuthStore()

/**
 * SSR-safe загрузка профиля
 */
const {
  data: profileRes,
  pending: profilePending,
  error: profileError
} = await useAsyncData(
  'my-profile',
  () => profileApi.getMyProfile(),
  { server: true }
)

/**
 * SSR-safe загрузка резюме
 */
const {
  data: resumesRes,
  pending: resumesPending
} = await useAsyncData(
  'my-resumes',
  () => resumesApi.getMyResumes(),
  { server: true }
)

/**
 * Normalized data
 */
const profile = computed(() =>
  profileRes.value?.ok ? profileRes.value.profile : null
)

const resumes = computed<any[]>(() =>
  resumesRes.value?.ok && Array.isArray(resumesRes.value.items)
    ? resumesRes.value.items
    : []
)

/**
 * UI states
 */
const loading = computed(() => profilePending.value || resumesPending.value)
const hasError = computed(() => !!profileError.value)

const chatsUnreadTotal = ref(0)

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

async function tryFetchList(urls: string[]) {
  let lastErr: any = null
  for (const url of urls) {
    try {
      const res: any = await authFetch(url, { method: 'GET' })
      const raw =
        (Array.isArray(res?.items) ? res.items : null) ||
        (Array.isArray(res?.chats) ? res.chats : null) ||
        (Array.isArray(res?.data) ? res.data : null) ||
        (Array.isArray(res) ? res : null)
      if (Array.isArray(raw)) return raw
    } catch (e: any) {
      lastErr = e
      const st = Number(e?.statusCode || e?.status || e?.response?.status || 0)
      if (st === 404 || st === 401 || st === 403) continue
      break
    }
  }
  throw lastErr || new Error('no_endpoint_matched')
}

async function loadChatsUnread() {
  if (!process.client) return
  try {
    if (!(auth as any)?.isReady && (auth as any)?.hasToken) {
      try { await (auth as any).init?.() } catch {}
    }
    const chats = await tryFetchList([
      '/api/user/chats',
      '/api/chats/my',
      '/api/chat/my',
      '/api/user/chat',
    ])
    chatsUnreadTotal.value = chats.reduce((sum: number, c: any) => sum + pickUnreadCount(c), 0)
  } catch {
    chatsUnreadTotal.value = 0
  }
}

onMounted(() => {
  loadChatsUnread()
})
</script>

<template>
  <section class="panel">
    <h1 class="title">{{ t('panel.user.title') }}</h1>

    <!-- SUBSCRIPTION BANNER -->
    <div class="mb-6">
      <SubscriptionBanner />
    </div>

    <!-- LOADING -->
    <ProfileSkeleton v-if="loading" />

    <!-- ERROR -->
    <div v-else-if="hasError" class="error">{{ t('panel.user.loadError') }}</div>

    <!-- CONTENT -->
    <div v-else class="layout">
      <!-- MAIN -->
      <div class="main">
        <ProfileHeader
          v-if="profile"
          :profile="profile"
          :is-owner="true"
        />
      </div>

      <!-- SIDEBAR -->
      <aside class="sidebar">
        <NuxtLink
          v-if="isOwner"
          to="/user/applications"
          class="chat-btn"
        >
          <span>{{ t('panel.user.myChats') }}</span>
          <span v-if="chatsUnreadTotal > 0" class="count-badge">{{ formatBadge(chatsUnreadTotal) }}</span>
        </NuxtLink>
        <!-- RESUMES LIST -->
        <template v-if="resumes.length">
          <ResumeMiniCard
            v-for="resume in resumes"
            :key="resume.id"
            :resume="resume"
          />
        </template>

        <!-- EMPTY STATE -->
        <div v-else class="resume-empty">
          <div class="resume-empty-icon">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
              <path
                d="M7 3h7l5 5v13a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"
                stroke="currentColor"
                stroke-width="1.6"
              />
              <path d="M14 3v5h5" stroke="currentColor" stroke-width="1.6" />
            </svg>
          </div>

          <h3 class="resume-empty-title">{{ t('panel.user.noResumesTitle') }}</h3>

          <p class="resume-empty-text">{{ t('panel.user.noResumesText') }}</p>

          <NuxtLink
            to="/resumes/create"
            class="resume-create-btn"
          >{{ t('panel.user.createResume') }}</NuxtLink>

          <p class="resume-empty-hint">{{ t('panel.user.emptyHint') }}</p>
        </div>
      </aside>
    </div>
  </section>
</template>

<style scoped>
.panel {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 16px;
}

.title {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 24px;
}

.layout {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 24px;
}

.main,
.sidebar {
  background: #f9fafb;
  border-radius: 20px;
  padding: 24px;
}

.sidebar {
  height: fit-content;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* ERROR */

.error {
  color: #ef4444;
  font-weight: 600;
}

/* EMPTY RESUME */

.resume-empty {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  border: 1px solid #e5e7eb;
  text-align: center;
}

.resume-empty-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: linear-gradient(135deg, #6366f1, #3b82f6);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.resume-empty-title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.resume-empty-text {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 16px;
  line-height: 1.5;
}

.resume-create-btn {
  display: block;
  width: 100%;
  padding: 14px 18px;
  border-radius: 14px;
  background: #111827;
  color: #ffffff;
  font-weight: 500;
  text-decoration: none;
  margin-bottom: 10px;
  transition: background 0.2s ease, transform 0.15s ease;
}

.resume-create-btn:hover {
  background: #1f2937;
  transform: translateY(-1px);
}

.resume-empty-hint {
  font-size: 12px;
  color: #9ca3af;
}

.chat-btn {
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 14px 18px;
  border-radius: 14px;
  background: #111827;
  color: #fff;
  text-decoration: none;

  transition:
    background 0.15s ease,
    transform 0.12s ease,
    box-shadow 0.15s ease;
}

.chat-btn:hover {
  background: #1f2937;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(0,0,0,0.18);
}

.count-badge {
  min-width: 22px;
  height: 22px;
  padding: 0 6px;
  border-radius: 999px;
  background: #ef4444;
  color: #fff;
  font-size: 12px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

</style>
