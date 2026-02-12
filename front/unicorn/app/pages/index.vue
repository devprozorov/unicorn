<template>
  <!-- Вся страница = viewport -->
  <div class="page">

    <!-- HEADER -->
    <Header />

    <!-- FIRST SCREEN -->
    <section class="hero-screen">
      <Hero />
    </section>

    <!-- SECOND SCREEN -->
    <section class="content-screen feature padding-fast-fix">
      <Features />
      <Footer />
    </section>

  </div>
</template>
<style scoped>
/*быстрый фикс отступа */
.padding-fast-fix{
  margin-top: 30px;
}

@media (max-width: 768px) {
  .padding-fast-fix {
    margin-top: 48px;
  }
}
/* базовые easing */

:root {
  --ease-out: cubic-bezier(.22,.61,.36,1);
}

/* скрыто до появления */
.reveal {
  opacity: 0;
  transform: translateY(24px);
  filter: blur(6px);
  transition:
    opacity .8s var(--ease-out),
    transform .8s var(--ease-out),
    filter .8s var(--ease-out);
}

/* показ */
.reveal.is-visible {
  opacity: 1;
  transform: translateY(0);
  filter: blur(0);
}

/* задержки */
.reveal.delay-1 { transition-delay: .1s }
.reveal.delay-2 { transition-delay: .2s }
.reveal.delay-3 { transition-delay: .3s }

</style>
<script setup lang="ts">
import Header from '~/components/Header.vue'
import Hero from '~/components/Hero.vue'
import Features from '~/components/Features.vue'
import Footer from '~/components/Footer.vue'
import { useI18n } from '~/composables/useI18n'

const { t, lang } = useI18n()

useHead(() => ({
  title: t('seo.home.title'),
  htmlAttrs: {
    lang: lang.value
  },
  meta: [
    { name: 'description', content: t('seo.home.description') },
    { name: 'keywords', content: t('seo.keywords.home') },
    { name: 'robots', content: 'index, follow' },
    { property: 'og:title', content: t('seo.home.title') },
    { property: 'og:description', content: t('seo.home.description') },
    { property: 'og:type', content: 'website' },
    { property: 'og:url', content: 'https://unicornstar.online/' },
    { property: 'og:site_name', content: 'Unicornstar' },
    { property: 'og:locale', content: lang.value === 'ru' ? 'ru_RU' : 'en_US' },
    { property: 'og:image', content: 'https://unicornstar.online/images/mainpage/main.png' },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: t('seo.home.title') },
    { name: 'twitter:description', content: t('seo.home.description') },
    { name: 'twitter:image', content: 'https://unicornstar.online/images/mainpage/main.png' }
  ],
  link: [
    { rel: 'canonical', href: 'https://unicornstar.online/' }
  ],
  script: [
    {
      type: 'application/ld+json',
      children: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'Organization',
        name: 'Unicornstar',
        url: 'https://unicornstar.online',
        logo: 'https://unicornstar.online/favicon.png'
      })
    },
    {
      type: 'application/ld+json',
      children: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'WebSite',
        name: 'Unicornstar',
        url: 'https://unicornstar.online'
      })
    }
  ]
}))
</script>
