<template>
  <div v-if="loaded && !subscription?.active" class="sub-banner">
    <div class="sub-banner__bg" aria-hidden="true" />
    <div class="sub-banner__overlay" aria-hidden="true" />

    <div class="sub-banner__content">
      <div class="sub-banner__left">
        <div class="sub-banner__kicker">
          <span class="sub-banner__pill">{{ t('subscription.premiumLabel') }}</span>
          <span class="sub-banner__spark">+</span>
        </div>

        <h3 class="sub-banner__title">{{ dict.subscription.banner.title }}</h3>
        <p class="sub-banner__price">{{ dict.subscription.banner.price }}</p>

        <ul class="sub-banner__list">
          <li v-for="(item, idx) in dict.subscription.banner.benefits" :key="idx">
            <span class="sub-banner__dot" />
            <span class="sub-banner__text">{{ item }}</span>
          </li>
        </ul>
      </div>

      <div class="sub-banner__right">
        <button
          class="sub-banner__cta"
          :disabled="loading"
          @click="onSubscribe"
        >
          <span class="sub-banner__cta-text">
            {{ loading ? t('subscription.banner.buttonLoading') : t('subscription.banner.buttonSubscribe') }}
          </span>
          <span class="sub-banner__cta-arrow">?</span>
        </button>
        <div class="sub-banner__hint">{{ t('subscription.banner.benefits')[0] }}</div>
      </div>
    </div>
  </div>

  <div v-else-if="loaded && subscription?.active" class="sub-banner sub-banner--active">
    <div class="sub-banner__bg" aria-hidden="true" />
    <div class="sub-banner__overlay" aria-hidden="true" />

    <div class="sub-banner__content">
      <div class="sub-banner__left">
        <div class="sub-banner__kicker">
          <span class="sub-banner__pill sub-banner__pill--active">{{ t('subscription.premiumLabel') }}</span>
          <span class="sub-banner__spark sub-banner__spark--active">*</span>
        </div>

        <h3 class="sub-banner__title">{{ t('subscription.banner.activeTitle') }}</h3>
        <p class="sub-banner__subtitle">{{ activeUntilLabel }}</p>
      </div>

      <div class="sub-banner__right sub-banner__right--active">
        <div class="sub-banner__badge">
          <span class="sub-banner__badge-dot" />
          <span class="sub-banner__badge-text">{{ t('subscription.premiumLabel') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useSubscription } from '~/composables/useSubscription'
import { useToast } from '~/composables/useToast'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'

const { getStatus, goToPayment } = useSubscription()
const toast = useToast()

const subscription = ref<any>(null)
const loading = ref(false)
const loaded = ref(false)
const { t, dict, lang } = useI18n()

const formatDate = (date: string) => {
  if (!date) return ''
  const locale = lang.value === 'ru' ? 'ru-RU' : 'en-US'
  return new Date(date).toLocaleDateString(locale, {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const activeUntilLabel = computed(() => {
  if (!subscription.value?.endDate || !subscription.value?.daysLeft) return ''
  return t('subscription.banner.activeUntil')
    .replace('{{date}}', formatDate(subscription.value.endDate))
    .replace('{{days}}', String(subscription.value.daysLeft))
    .replace('{{label}}', t('subscription.daysLeft'))
})

const onSubscribe = async () => {
  loading.value = true
  try {
    const res = await goToPayment()
    if (!res.ok) {
      toast.error(getErrorMessage(res.error || 'payment_failed', t, 'subscription.errors.paymentFailed'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'payment_failed', t, 'subscription.errors.generic'))
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const status = await getStatus()
    subscription.value = status
  } catch (e) {
    // In case of error, show subscription banner
    subscription.value = { active: false }
  } finally {
    loaded.value = true
  }
})
</script>

<style scoped>
.sub-banner {
  position: relative;
  overflow: hidden;
  border-radius: 22px;
  padding: 18px;
  border: 1px solid rgba(99,102,241,0.22);
  background: #ffffff;
  box-shadow:
    0 18px 36px rgba(17,24,39,0.12),
    inset 0 1px 0 rgba(255,255,255,0.9);
}

.sub-banner--active {
  border-color: rgba(16,185,129,0.32);
  background: #f0fdf4;
}

.sub-banner__bg {
  position: absolute;
  inset: 0;
  background: url('/images/sub-banner.webp') center / cover no-repeat;
  opacity: 0.28;
  filter: saturate(0.75) contrast(0.9);
}

.sub-banner--active .sub-banner__bg {
  opacity: 0.2;
  filter: saturate(0.65) contrast(0.9);
}

.sub-banner__overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(115deg, rgba(255,255,255,0.9) 0%, rgba(255,255,255,0.7) 55%, rgba(255,255,255,0.35) 100%);
}

.sub-banner--active .sub-banner__overlay {
  background: linear-gradient(115deg, rgba(236,253,245,0.9) 0%, rgba(236,253,245,0.7) 55%, rgba(236,253,245,0.35) 100%);
}

.sub-banner__content {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 220px;
  gap: 16px;
  align-items: center;
}

.sub-banner__left {
  background: rgba(255,255,255,0.9);
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 10px 24px rgba(15,23,42,0.08);
  backdrop-filter: blur(6px);
}

.sub-banner__kicker {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.sub-banner__pill {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  background: rgba(15,23,42,0.9);
  color: #f8fafc;
}

.sub-banner__pill--active {
  background: rgba(6,95,70,0.9);
  color: #d1fae5;
}

.sub-banner__spark {
  width: 22px;
  height: 22px;
  border-radius: 9px;
  display: grid;
  place-items: center;
  font-weight: 800;
  background: rgba(15,23,42,0.9);
  color: #f8fafc;
}

.sub-banner__spark--active {
  background: rgba(6,95,70,0.9);
  color: #d1fae5;
}

.sub-banner__title {
  margin: 0 0 8px;
  font-size: 18px;
  font-weight: 800;
  color: #0f172a;
}

.sub-banner__price {
  margin: 0 0 10px;
  font-size: 14px;
  font-weight: 700;
  color: #b45309;
}

.sub-banner__subtitle {
  margin: 0;
  font-size: 13px;
  color: #475569;
}

.sub-banner__list {
  margin: 0;
  padding: 0;
  list-style: none;
  display: grid;
  gap: 8px;
}

.sub-banner__list li {
  display: grid;
  grid-template-columns: 8px 1fr;
  gap: 10px;
  font-size: 12px;
  color: #334155;
  align-items: center;
}

.sub-banner__dot {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: #6366f1;
}

.sub-banner__right {
  display: grid;
  gap: 8px;
  align-content: center;
  padding: 14px;
  border-radius: 16px;
  border: 1px solid rgba(79,70,229,0.22);
  background: rgba(255,255,255,0.95);
  box-shadow: 0 10px 22px rgba(79,70,229,0.12);
}

.sub-banner__right--active {
  border-color: rgba(16,185,129,0.3);
  background: rgba(236,253,245,0.85);
}

.sub-banner__cta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 8px 12px;
  min-height: 36px;
  border-radius: 14px;
  border: 1px solid rgba(79,70,229,0.3);
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  color: #fff;
  font-weight: 600;
  cursor: pointer;
  transition: transform .15s ease, box-shadow .15s ease, opacity .15s ease;
  box-shadow: 0 12px 24px rgba(79,70,229,0.3);
}

.sub-banner__cta:hover {
  transform: translateY(-1px);
  box-shadow: 0 16px 30px rgba(79,70,229,0.35);
}

.sub-banner__cta:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  box-shadow: none;
}

.sub-banner__cta-text {
  font-size: 12px;
  font-weight: 600;
}

.sub-banner__cta-arrow {
  font-size: 14px;
  font-weight: 800;
}

.sub-banner__hint {
  font-size: 10px;
  color: #64748b;
  text-align: center;
}

.sub-banner__badge {
  display: grid;
  justify-items: center;
  gap: 6px;
  padding: 10px;
  border-radius: 12px;
  border: 1px solid rgba(16,185,129,0.35);
  color: #065f46;
  background: rgba(236,253,245,0.8);
}

.sub-banner__badge-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #10b981;
  box-shadow: 0 0 0 6px rgba(16,185,129,0.2);
}

.sub-banner__badge-text {
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

@media (max-width: 900px) {
  .sub-banner__content {
    grid-template-columns: 1fr;
  }

  .sub-banner__right {
    justify-items: stretch;
  }
}

@media (max-width: 640px) {
  .sub-banner {
    padding: 14px;
  }

  .sub-banner__left {
    padding: 14px;
  }

  .sub-banner__title {
    font-size: 16px;
  }

  .sub-banner__cta {
    width: 100%;
  }
}
</style>
