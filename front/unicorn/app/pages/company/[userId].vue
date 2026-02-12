<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'

import CompanyPanel from '~/components/panel/company/CompanyPanel.vue'
import ProfileSkeleton from '~/components/panel/user/ProfileSkeleton.vue'

import { useAuthStore } from '~/stores/auth'
import { useProfileApi } from '~/services/profileApi'
import { useI18n } from '~/composables/useI18n'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const profileApi = useProfileApi()
const { t } = useI18n()

const loading = ref(true)
const error = ref<string | null>(null)
const profile = ref<any | null>(null)

const userId = computed(() => String(route.params.userId || ''))

/**
 * ✅ ВЛАДЕЛЕЦ = только company + совпадение id
 */
const isOwner = computed(() => {
  if (!auth.isReady) return false
  if (!auth.user) return false
  if (auth.user.type !== 'company') return false
  return auth.user.userId === userId.value
})


function getProfileType(p: any): string | null {
  return p?.type ?? p?.accountType ?? null
}

onMounted(async () => {
  loading.value = true
  error.value = null

  try {
    await auth.init()

    let res: any = null

    /**
     * ✅ ТОЛЬКО владелец-компания
     */
    if (
      auth.user &&
      auth.user.type === 'company' &&
      auth.user.userId === userId.value
    ) {
      res = await profileApi.getMyProfile()
    } else {
      /**
       * ✅ ВСЕ ОСТАЛЬНЫЕ
       * user + guest
       */
      res = await profileApi.getPublicCompanyProfile(userId.value)
    }

    const p = res?.profile ?? null
    profile.value = p

    if (!p) {
      error.value = 'not_found'
      return
    }

    /**
     * ❗ НИКАКИХ редиректов
     */
    const type = p?.type ?? p?.accountType ?? null
    if (type !== 'company') {
      error.value = 'not_found'
      return
    }
  } catch (e) {
    error.value = 'load_error'
  } finally {
    loading.value = false
  }
})

</script>


<template>
  <div class="page">
    <Header />

    <main class="shell">
      <ProfileSkeleton v-if="loading || redirecting" />

      <div v-else-if="error" class="error">
        <div class="error-card">
                    <div class="error-title">
            <span v-if="error === 'not_found'">
              {{ t('companyProfilePage.errors.notFoundTitle') }}
            </span>
            <span v-else>{{ t('companyProfilePage.errors.loadTitle') }}</span>
          </div>
          <div class="error-text">
            {{ t('companyProfilePage.errors.text') }}
          </div>
          <div class="error-actions">
            <NuxtLink class="btn" to="/">
              {{ t('companyProfilePage.errors.toHome') }}
            </NuxtLink>
          </div>
        </div>
      </div>

      <CompanyPanel
        v-else
        :profile="profile"
        :is-owner="isOwner"
      />
    </main>

    <Footer />
  </div>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background:
    radial-gradient(1200px 600px at 20% -10%, rgba(79,70,229,0.18), transparent 60%),
    radial-gradient(900px 500px at 90% 10%, rgba(99,102,241,0.14), transparent 60%),
    #ffffff;
}

.shell {
  max-width: 1320px;
  margin: 0 auto;
  padding: 28px 16px 110px;
}

.error {
  padding: 18px 0;
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

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  padding: 12px 14px;
  font-weight: 800;
  background: rgba(17,24,39,0.03);
  color: #111827;
  border: 1px solid rgba(17,24,39,0.06);
  transition: transform .15s ease, background .15s ease;
}
.btn:hover { transform: translateY(-1px); background: rgba(79,70,229,0.06); }

@media (max-width: 768px) {
  .shell { padding: 16px 12px 120px; }
}
</style>

