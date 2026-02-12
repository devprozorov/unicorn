<template>
  <Header />

  <main class="contact-page">
    <section class="contact-hero">
      <div class="contact-hero-inner">
        <p class="contact-kicker">Unicornstar</p>
        <h1>{{ t('contacts.title') }}</h1>
        <p class="contact-lead">{{ t('contacts.intro') }}</p>
      </div>
      <div class="contact-hero-panel">
        <div class="contact-chip">Telegram</div>
        <div class="contact-chip">Email</div>
        <div class="contact-chip">Community</div>
      </div>
    </section>

    <section class="contact-grid">
      <article class="contact-card">
        <div class="contact-icon">TG</div>
        <div>
          <h2>{{ t('contacts.telegram') }}</h2>
          <a
            class="contact-link"
            href="https://t.me/miraclescapedev"
            target="_blank"
            rel="noopener"
          >
            @miraclescapedev
          </a>
        </div>
      </article>

      <article class="contact-card">
        <div class="contact-icon">CH</div>
        <div>
          <h2>{{ t('contacts.telegramChannel') }}</h2>
          <a
            class="contact-link"
            href="https://t.me/MiracleScape"
            target="_blank"
            rel="noopener"
          >
            t.me/MiracleScape
          </a>
        </div>
      </article>

      <article class="contact-card">
        <div class="contact-icon">EM</div>
        <div>
          <h2>{{ t('contacts.email') }}</h2>
          <a class="contact-link" href="mailto:tumyp2003map@mail.ru">
            tumyp2003map@mail.ru
          </a>
        </div>
      </article>
    </section>

    <section class="contact-note">
      <p>{{ t('contacts.description') }}</p>
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
  title: t('contacts.title'),
  htmlAttrs: { lang: lang.value },
  meta: [
    { name: 'description', content: t('contacts.description') },
    { name: 'robots', content: 'index, follow' }
  ],
  link: [
    { rel: 'canonical', href: 'https://unicornstar.online/contacts' }
  ]
}))
</script>

<style scoped>
.contact-page {
  max-width: 1100px;
  margin: 0 auto;
  padding: 36px 16px 80px;
}

.contact-hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 260px;
  gap: 24px;
  padding: 32px;
  border-radius: 28px;
  background: radial-gradient(circle at top left, #f5f3ff, #ffffff 55%);
  border: 1px solid #e5e7eb;
  position: relative;
  overflow: hidden;
}

.contact-hero::after {
  content: '';
  position: absolute;
  inset: 0;
  background: url('/images/contact-bg.webp') center / cover no-repeat;
  opacity: 0.22;
  z-index: 0;
}

.contact-hero-inner,
.contact-hero-panel {
  position: relative;
  z-index: 1;
}

.contact-hero-inner h1 {
  font-size: 36px;
  margin-bottom: 10px;
}

.contact-kicker {
  text-transform: uppercase;
  letter-spacing: 0.2em;
  font-size: 12px;
  color: #7c3aed;
  font-weight: 600;
}

.contact-lead {
  color: #4b5563;
  font-size: 17px;
  max-width: 560px;
}

.contact-hero-panel {
  display: grid;
  align-content: center;
  gap: 10px;
}

.contact-chip {
  padding: 10px 14px;
  border-radius: 999px;
  background: #111827;
  color: #fff;
  font-size: 12px;
  text-align: center;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.contact-grid {
  margin-top: 28px;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 16px;
}

.contact-card {
  display: grid;
  grid-template-columns: 52px 1fr;
  gap: 16px;
  padding: 18px;
  border-radius: 20px;
  background: #fff;
  border: 1px solid #e5e7eb;
  box-shadow: 0 10px 28px rgba(17, 24, 39, 0.06);
}

.contact-icon {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  background: #ede9fe;
  color: #4c1d95;
  display: grid;
  place-items: center;
  font-weight: 700;
}

.contact-card h2 {
  font-size: 16px;
  margin-bottom: 6px;
}

.contact-link {
  color: #4f46e5;
  font-weight: 600;
  text-decoration: none;
}

.contact-note {
  margin-top: 24px;
  padding: 16px 20px;
  border-left: 3px solid #4f46e5;
  background: #f8fafc;
  border-radius: 14px;
  color: #374151;
}

@media (max-width: 900px) {
  .contact-hero {
    grid-template-columns: 1fr;
  }

  .contact-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 360px) {
  .contact-page {
    padding: 24px 12px 60px;
  }

  .contact-hero {
    padding: 22px;
    border-radius: 22px;
  }

  .contact-hero-inner h1 {
    font-size: 28px;
  }

  .contact-lead {
    font-size: 15px;
  }

  .contact-card {
    grid-template-columns: 44px 1fr;
    padding: 14px;
  }

  .contact-icon {
    width: 44px;
    height: 44px;
    border-radius: 12px;
    font-size: 12px;
  }
}

@media (max-width: 320px) {
  .contact-hero {
    padding: 18px;
  }

  .contact-hero-inner h1 {
    font-size: 24px;
  }
}
</style>
