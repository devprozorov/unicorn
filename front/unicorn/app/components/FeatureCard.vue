<template>
  <article class="faq-card">
    <div class="faq-media">
      <img
        :src="imgSrc"
        :alt="title"
        loading="lazy"
        @error="onImgError"
      />
    </div>

    <div class="faq-body">
      <h3 class="faq-title">{{ title }}</h3>
      <p class="faq-subtitle">{{ subtitle }}</p>
    </div>

    <NuxtLink :to="link" class="faq-link">
      {{ t('employers.more') }}
    </NuxtLink>
  </article>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from '~/composables/useI18n'

/**
 * id — СТАБИЛЬНЫЙ КЛЮЧ
 * НЕ ЗАВИСИТ ОТ ЯЗЫКА
 */
const props = defineProps<{
  id: 'pochemu-my' | 'kak-nachat' | 'oferta' | 'o-nas'
  title: string
  subtitle: string
  link?: string
}>()

const { dict, t } = useI18n()


const imgSrc = computed(() => {
  return `/images/mainpage/${props.id}.png`
})

const link = computed(() => props.link || '/doc')

function onImgError(e: Event) {
  ;(e.target as HTMLImageElement).src = '/images/mainpage/default.png'
}
</script>

<style scoped>
/* === СТАРАЯ FAQ СТИЛИЗАЦИЯ === */

.faq-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  height: 100%;
  transition: box-shadow .2s ease, transform .2s ease;
}

.faq-card:hover {
  box-shadow: 0 10px 30px rgba(0,0,0,.06);
  transform: translateY(-2px);
}

.faq-media {
  width: 100%;
  height: 140px;
  background: #f3f4f6;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 12px;
}

.faq-media img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.faq-body {
  flex: 1;
}

.faq-title {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
}

.faq-subtitle {
  font-size: 12px;
  color: #6b7280;
}

.faq-link {
  margin-top: 10px;
  font-size: 13px;
  font-weight: 600;
  color: #111827;
  text-decoration: none;
}

.faq-link:hover {
  opacity: .7;
}
</style>
