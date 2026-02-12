<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useResumesApi } from '~/services/resumesApi'
import Header from '~/components/Header.vue'
import { useI18n } from '~/composables/useI18n'
import Footer from '~/components/Footer.vue'
import type { Resume } from '~/types/resume'

const route = useRoute()
const api = useResumesApi()
const { t } = useI18n()

const id = route.params.id as string
if (!id) throw createError({ statusCode: 404 })

const resume: Resume = await api.getPublicResume(id)

// Backend сейчас отдаёт status: "active" (а не "published")
const isPublic = resume?.status === 'active' || resume?.status === 'published'
if (!resume || !isPublic) {
  throw createError({ statusCode: 404 })
}

useHead({
  title: `${resume.title} ? ${t('resumeView.titleSuffix')}`,
  meta: [
    {
      name: 'description',
      content: (resume.about ?? '').slice(0, 160),
    },
  ],
})
</script>

<template>
  <Header />

  <main class="resume-page">
    <div class="card">
      <h1 class="title">{{ resume.title }}</h1>

      <p v-if="resume.about" class="about">
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
    </div>
  </main>

  <Footer />
</template>

<style scoped>
.resume-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 96px 16px 64px;
}

.card {
  background: #ffffff;
  padding: 32px;
  border-radius: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.06);
}

.title {
  font-size: 28px;
  font-weight: 700;
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
