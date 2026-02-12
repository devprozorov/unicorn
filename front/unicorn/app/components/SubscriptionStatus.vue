<template>
  <div v-if="loaded && subscription" class="subscription-status">
    <div v-if="subscription.active" class="status-badge active">
      <span class="icon">*</span>
      <span class="text">{{ t('subscription.premiumLabel') }}</span>
      <span v-if="subscription.daysLeft" class="days">
        {{ subscription.daysLeft }} {{ t('subscription.daysLeft') }}
      </span>
    </div>
    <NuxtLink v-else to="/subscribe/offer" class="status-badge inactive">
      <span class="icon">+</span>
      <span class="text">{{ t('subscription.subscribe') }}</span>
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSubscription } from '~/composables/useSubscription'
import { useI18n } from '~/composables/useI18n'

const { getStatus } = useSubscription()
const { t } = useI18n()

const subscription = ref<any>(null)
const loaded = ref(false)

onMounted(async () => {
  try {
    const status = await getStatus()
    subscription.value = status
  } catch (e) {
    subscription.value = { active: false }
  } finally {
    loaded.value = true
  }
})
</script>

<style scoped>
.subscription-status {
  display: flex;
  align-items: center;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  transition: all 0.2s;
  white-space: nowrap;
}

.status-badge.active {
  background: linear-gradient(135deg, #fef3c7 0%, #fde68a 100%);
  color: #92400e;
  border: 1px solid #fbbf24;
}

.status-badge.inactive {
  background: linear-gradient(135deg, #f8fafc 0%, #eef2f7 100%);
  color: #0f172a;
  border: 1px solid #cbd5e1;
  box-shadow: 0 6px 16px rgba(15,23,42,0.06);
  cursor: pointer;
  text-decoration: none;
}

.status-badge.inactive:hover {
  background: linear-gradient(135deg, #f1f5f9 0%, #e2e8f0 100%);
  border-color: #94a3b8;
  transform: translateY(-1px);
}

.status-badge .icon {
  font-size: 14px;
  line-height: 1;
}

.status-badge .text {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.02em;
}

.status-badge .days {
  font-size: 10px;
  opacity: 0.8;
}

@media (max-width: 768px) {
  .status-badge .text {
    display: none;
  }
  
  .status-badge {
    padding: 4px 8px;
  }
}
</style>

