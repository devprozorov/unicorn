<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from '~/composables/useI18n'

const { t } = useI18n()

const props = defineProps<{ profile: any }>()

const vacanciesCount = computed(() => props.profile?.vacanciesCount ?? 0)
const followersCount = computed(() => props.profile?.followersCount ?? 0)
</script>

<template>
  <div class="stats">
    <div class="stat">
      <div class="value">{{ vacanciesCount }}</div>
      <div class="label">{{ t('companyPanel.vacancies') }}</div>
    </div>
    <div class="stat">
      <div class="value">{{ followersCount }}</div>
      <div class="label">{{ t('companyPanel.followers') }}</div>
    </div>
  </div>
</template>

<style scoped>
.stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 14px;
}
.stat {
  background: #f9fafb;
  border-radius: 14px;
  padding: 14px;
}
.value {
  font-weight: 900;
  font-size: 18px;
}
.label {
  color: #6b7280;
  margin-top: 4px;
  font-weight: 700;
  font-size: 12px;
}
</style>
