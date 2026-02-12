<script setup lang="ts">
import { useRoute } from 'vue-router'
import { computed, ref, onMounted } from 'vue'
import { useResumesApi } from '~/services/resumesApi'
import { useAuthStore } from '~/stores/auth'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import { useI18n } from '~/composables/useI18n'
import type { Resume } from '~/types/resume'

const route = useRoute()
const api = useResumesApi()
const auth = useAuthStore()
const { t } = useI18n()

const id = route.params.id as string
if (!id) throw createError({ statusCode: 404 })

// на всякий случай, если страница SSR/переходами
const resume = ref<Resume | null>(null)
const error = ref(false)
const loading = ref(true)

async function loadResume() {
  error.value = false
  loading.value = true

  if (process.client && !auth.isReady) {
    try {
      await auth.init()
    } catch {}
  }

  const myId = auth.userId
  if (!myId) {
    error.value = true
    loading.value = false
    return
  }

  try {
    const r = await api.getResumeById(id)

    if (r.userId !== myId) {
      error.value = true
    } else {
      resume.value = r
    }
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

onMounted(loadResume)

const statusLabel = computed(() => {
  if (!resume.value) return ''
  if (resume.value.status === 'active') return t('resumeView.active')
  if (resume.value.status === 'draft') return t('resumeView.draft')
  if (resume.value.status === 'hidden') return t('resumeView.hidden')
  return resume.value.status
})
</script>

<template>
  <Header />

  <main class="resume-page">
    <section v-if="error" class="card error">{{ t('resumeView.notFound') }}</section>

    <section v-else-if="resume" class="card">
      <div class="top">
        <h1 class="title">{{ resume.title }}</h1>
        <span class="badge">{{ statusLabel }}</span>
      </div>

      <p v-if="resume.about" class="about whitespace-pre-wrap">
        {{ resume.about }}
      </p>

      <section v-if="resume.skills?.length">
        <h3>{{ t('resumeView.linksLabel') }}</h3>
        <ul>
          <li v-for="s in resume.skills" :key="s">
            {{ s }}
          </li>
        </ul>
      </section>

      <section v-if="resume.links?.length">
        <h3>{{ t('resumeView.linksLabel') }}</h3>
        <ul>
          <li v-for="l in resume.links" :key="l">
            <a :href="l" target="_blank" rel="noopener">
              {{ l }}
            </a>
          </li>
        </ul>
      </section>
    </section>
  </main>

  <Footer />
</template>

<style scoped>
.resume-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 96px 16px;
}

.card {
  background: #ffffff;
  padding: 32px;
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.06);
}

.top {
  display: flex;
  gap: 12px;
  align-items: center;
  justify-content: space-between;
}

.title {
  font-size: 28px;
  font-weight: 700;
  margin: 0;
}

.badge {
  font-size: 12px;
  padding: 6px 10px;
  border-radius: 999px;
  background: rgba(99, 102, 241, 0.12);
  color: #4f46e5;
  border: 1px solid rgba(99, 102, 241, 0.25);
  white-space: nowrap;
}

.error {
  text-align: center;
  color: #b91c1c;
}

.about {
  margin: 16px 0;
  color: #374151;
  line-height: 1.6;
}

h3 {
  margin-top: 24px;
  font-size: 18px;
}

ul {
  margin-top: 8px;
  padding-left: 20px;
}
</style>
