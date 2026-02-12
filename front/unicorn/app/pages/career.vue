<template>
  <Header />

  <main class="career-page">
    <section class="career-hero">
      <div class="career-hero-content">
        <p class="career-eyebrow">Unicornstar</p>
        <h1>{{ t('career.title') }}</h1>
        <p class="career-lead">{{ t('career.intro') }}</p>
        <div class="career-cta-row">
          <NuxtLink class="career-cta" to="/jobs">{{ t('jobs.title') }}</NuxtLink>
          <a class="career-cta ghost" href="mailto:tumyp2003map@mail.ru">
            {{ t('career.sendResume') }}
          </a>
        </div>
        <div class="career-meta">
          <span class="career-pill">{{ t('career.openRoles') }}</span>
          <span class="career-pill">0</span>
        </div>
      </div>
      <div class="career-hero-panel">
        <div class="career-panel-card">
          <p class="career-panel-title">{{ t('career.openRoles') }}</p>
          <p class="career-panel-value">0</p>
          <p class="career-panel-sub">{{ t('career.openRolesEmpty') }}</p>
        </div>
        <div class="career-panel-card light">
          <p class="career-panel-title">{{ t('career.sendResume') }}</p>
          <p class="career-panel-value">{{ t('career.email') }}</p>
        </div>
      </div>
    </section>

    <section class="career-board">
      <div class="career-board-header">
        <h2>{{ t('career.openRoles') }}</h2>
        <span class="career-tag">Unicornstar</span>
      </div>

      <div class="career-empty">
        <div class="career-empty-count">0</div>
        <div>
          <p class="career-empty-title">{{ t('career.openRolesEmpty') }}</p>
          <p class="career-empty-text">{{ t('career.intro') }}</p>
        </div>
      </div>
    </section>

    <section class="career-resume">
      <div>
        <h2>{{ t('career.sendResume') }}</h2>
        <p>{{ t('career.intro') }}</p>
      </div>
      <a class="career-resume-link" href="mailto:tumyp2003map@mail.ru">
        {{ t('career.email') }}
      </a>
    </section>
  </main>

  <Footer />
</template>

<script setup lang="ts">
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import { useI18n } from '~/composables/useI18n'

const { t, lang } = useI18n()

useHead(() => ({
  title: t('career.title'),
  htmlAttrs: { lang: lang.value },
  meta: [
    { name: 'description', content: t('career.description') },
    { name: 'robots', content: 'index, follow' }
  ],
  link: [
    { rel: 'canonical', href: 'https://unicornstar.online/career' }
  ]
}))
</script>

<style scoped>
.career-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 16px 80px;
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.career-hero {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(0, 0.8fr);
  gap: 24px;
  padding: 44px;
  border-radius: 30px;
  background: linear-gradient(120deg, #0b1020, #111827 55%, #1f2937);
  color: #fff;
  position: relative;
  overflow: hidden;
}

.career-hero::before {
  content: '';
  position: absolute;
  right: -120px;
  top: -80px;
  width: 280px;
  height: 280px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(124, 58, 237, 0.35), transparent 70%);
  z-index: 0;
}

.career-hero::after {
  content: '';
  position: absolute;
  inset: 0;
  background: url('/images/career-bg.webp') right center / cover no-repeat;
  opacity: 0.28;
  z-index: 0;
}

.career-hero-content,
.career-hero-panel {
  position: relative;
  z-index: 1;
}

.career-eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.24em;
  font-size: 12px;
  color: #c7d2fe;
}

.career-hero-content h1 {
  font-size: 42px;
  margin: 12px 0;
}

.career-lead {
  font-size: 18px;
  color: #e5e7eb;
  max-width: 600px;
}

.career-cta-row {
  margin-top: 24px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.career-cta {
  padding: 12px 18px;
  border-radius: 12px;
  background: #f8fafc;
  color: #111827;
  font-weight: 600;
  text-decoration: none;
}

.career-cta.ghost {
  background: transparent;
  border: 1px solid #e5e7eb;
  color: #f8fafc;
}

.career-meta {
  margin-top: 18px;
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.career-pill {
  padding: 6px 12px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
}

.career-hero-panel {
  display: grid;
  gap: 16px;
  align-content: center;
}

.career-panel-card {
  padding: 18px;
  border-radius: 18px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.08);
  display: grid;
  gap: 6px;
}

.career-panel-title {
  font-size: 12px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #cbd5f5;
}

.career-panel-value {
  font-size: 22px;
  font-weight: 700;
}

.career-panel-sub {
  color: #cbd5f5;
  font-size: 13px;
}

.career-panel-card.light {
  background: #ffffff;
  color: #111827;
}

.career-panel-card.light .career-panel-title,
.career-panel-card.light .career-panel-sub {
  color: #6b7280;
}

.career-board {
  padding: 28px;
  border-radius: 24px;
  border: 1px solid #e5e7eb;
  background: #fff;
  box-shadow: 0 14px 32px rgba(15, 23, 42, 0.08);
}

.career-board-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.career-tag {
  padding: 6px 12px;
  border-radius: 999px;
  background: #111827;
  color: #fff;
  font-size: 12px;
}

.career-empty {
  display: grid;
  grid-template-columns: 80px 1fr;
  gap: 18px;
  padding: 20px;
  border-radius: 18px;
  background: #f8fafc;
  border: 1px dashed #cbd5f5;
}

.career-empty-count {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: #111827;
  color: #fff;
  display: grid;
  place-items: center;
  font-size: 26px;
  font-weight: 700;
}

.career-empty-title {
  font-weight: 600;
  margin-bottom: 6px;
}

.career-empty-text {
  color: #6b7280;
  margin: 0;
}

.career-resume {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 24px 28px;
  border-radius: 24px;
  border: 1px solid #e5e7eb;
  background: linear-gradient(120deg, #eef2ff, #ffffff);
}

.career-resume-link {
  padding: 10px 18px;
  border-radius: 999px;
  background: #4f46e5;
  color: #fff;
  text-decoration: none;
  font-weight: 600;
}

@media (max-width: 900px) {
  .career-hero {
    grid-template-columns: 1fr;
    padding: 28px;
  }

  .career-hero-content h1 {
    font-size: 32px;
  }

  .career-board-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }

  .career-resume {
    flex-direction: column;
    align-items: flex-start;
  }

  .career-empty {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 360px) {
  .career-page {
    padding: 20px 12px 60px;
  }

  .career-hero {
    padding: 22px;
    border-radius: 22px;
  }

  .career-hero-content h1 {
    font-size: 26px;
  }

  .career-lead {
    font-size: 15px;
  }

  .career-cta {
    width: 100%;
    justify-content: center;
    text-align: center;
  }

  .career-panel-card {
    padding: 14px;
  }

  .career-empty-count {
    width: 64px;
    height: 64px;
    font-size: 22px;
  }
}

@media (max-width: 320px) {
  .career-hero {
    padding: 18px;
  }

  .career-hero-content h1 {
    font-size: 24px;
  }

  .career-cta-row {
    gap: 10px;
  }
}
</style>
