<template>
  <Header />

  <main class="seo-landing">
    <div class="seo-hero">
      <p class="seo-kicker">{{ content.linkTitle }}</p>
      <h1>{{ content.heroTitle }}</h1>
      <p class="seo-lead">{{ content.heroText }}</p>
      <NuxtLink class="seo-cta" to="/jobs">{{ content.cta }}</NuxtLink>
    </div>

    <section class="seo-section">
      <h2>{{ content.linkTitle }}</h2>
      <ul class="seo-list">
        <li v-for="item in content.bullets" :key="item">{{ item }}</li>
      </ul>
    </section>

    <section class="seo-section">
      <h2>{{ content.faqTitle }}</h2>
      <div class="seo-faq">
        <div v-for="item in content.faq" :key="item.q" class="seo-faq-item">
          <h3>{{ item.q }}</h3>
          <p>{{ item.a }}</p>
        </div>
      </div>
    </section>

    <section class="seo-links">
      <h2>{{ t('seoLanding.linksTitle') }}</h2>
      <p>{{ t('seoLanding.linksText') }}</p>
      <div class="seo-links-grid">
        <NuxtLink to="/seo/startup-jobs">{{ t('seoLanding.startupJobs.linkTitle') }}</NuxtLink>
        <NuxtLink to="/jobs">{{ t('jobs.title') }}</NuxtLink>
      </div>
    </section>
  </main>

  <Footer />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import { useI18n } from '~/composables/useI18n'

const { t, dict, lang } = useI18n()

const content = computed(() => dict.value.seoLanding.gamedevJobs)

useHead(() => ({
  title: content.value.title,
  htmlAttrs: {
    lang: lang.value
  },
  meta: [
    { name: 'description', content: content.value.description },
    { name: 'keywords', content: content.value.keywords },
    { name: 'robots', content: 'index, follow' },
    { property: 'og:title', content: content.value.title },
    { property: 'og:description', content: content.value.description },
    { property: 'og:type', content: 'article' },
    { property: 'og:url', content: 'https://unicornstar.online/seo/gamedev-jobs' },
    { property: 'og:site_name', content: 'Unicornstar' },
    { property: 'og:locale', content: lang.value === 'ru' ? 'ru_RU' : 'en_US' },
    { property: 'og:image', content: 'https://unicornstar.online/images/mainpage/main.png' },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: content.value.title },
    { name: 'twitter:description', content: content.value.description },
    { name: 'twitter:image', content: 'https://unicornstar.online/images/mainpage/main.png' }
  ],
  link: [
    { rel: 'canonical', href: 'https://unicornstar.online/seo/gamedev-jobs' }
  ],
  script: [
    {
      type: 'application/ld+json',
      children: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'FAQPage',
        mainEntity: content.value.faq.map((item: any) => ({
          '@type': 'Question',
          name: item.q,
          acceptedAnswer: {
            '@type': 'Answer',
            text: item.a
          }
        }))
      })
    },
    {
      type: 'application/ld+json',
      children: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'BreadcrumbList',
        itemListElement: [
          {
            '@type': 'ListItem',
            position: 1,
            name: 'Unicornstar',
            item: 'https://unicornstar.online/'
          },
          {
            '@type': 'ListItem',
            position: 2,
            name: content.value.linkTitle,
            item: 'https://unicornstar.online/seo/gamedev-jobs'
          }
        ]
      })
    }
  ]
}))
</script>

<style scoped>
.seo-landing {
  max-width: 960px;
  margin: 0 auto;
  padding: 48px 16px 72px;
}

.seo-hero {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 36px;
}

.seo-kicker {
  color: #6b7280;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.seo-lead {
  color: #374151;
  font-size: 18px;
  line-height: 1.6;
}

.seo-cta {
  width: fit-content;
  padding: 10px 18px;
  border-radius: 10px;
  background: #3b35c7;
  color: #fff;
  font-weight: 600;
  text-decoration: none;
}

.seo-section {
  margin-top: 32px;
}

.seo-list {
  margin: 16px 0 0;
  padding-left: 18px;
  color: #374151;
}

.seo-faq {
  display: grid;
  gap: 16px;
  margin-top: 16px;
}

.seo-faq-item {
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #fff;
}

.seo-faq-item h3 {
  font-size: 16px;
  margin-bottom: 6px;
}

.seo-faq-item p {
  margin: 0;
  color: #4b5563;
}

.seo-links {
  margin-top: 36px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.seo-links-grid {
  margin-top: 12px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.seo-links-grid a {
  color: #3b35c7;
  font-weight: 600;
  text-decoration: none;
}
</style>
