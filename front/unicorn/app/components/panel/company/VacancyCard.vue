<script setup lang="ts">
import { useI18n } from '~/composables/useI18n'
const { t } = useI18n()

defineProps<{
  vacancy: {
    vacancyId: string
    title: string
    description: string
    location?: string
    tags?: string[]
    isPremium?: boolean
    colorCode?: string
  }
  editable?: boolean
}>()

defineEmits<{
  (e: 'edit'): void
  (e: 'delete'): void
}>()
</script>

<template>
  <PremiumCard :is-premium="vacancy.isPremium" :color-code="vacancy.colorCode">
    <div class="vac-card">
      <div class="vac-head">
        <h3 class="vac-title">{{ vacancy.title }}</h3>
        <span v-if="vacancy.location" class="vac-location">
          {{ vacancy.location }}
        </span>
      </div>

      <p class="vac-desc">
        {{ vacancy.description || t('companyVacancies.noDescription') }}
      </p>

      <div v-if="vacancy.tags?.length" class="vac-tags">
        <span
          v-for="(t, i) in vacancy.tags"
          :key="i"
          class="tag"
        >
          {{ t }}
        </span>
      </div>

      <div v-if="editable" class="vac-actions">
        <button class="btn small" @click="$emit('edit')">{{ t('companyVacancies.edit') }}</button>
        <button class="btn small danger" @click="$emit('delete')">{{ t('companyVacancies.deleteBtn') }}</button>
      </div>
    </div>
  </PremiumCard>
</template>

<style scoped>
.vac-card {
  padding: 16px;
  border-radius: 16px;
  border: 1px solid rgba(17,24,39,.08);
  background: #fff;
}

.vac-head {
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.vac-title {
  font-weight: 900;
  font-size: 16px;
}

.vac-location {
  font-size: 12px;
  color: #6b7280;
}

.vac-desc {
  margin-top: 8px;
  font-size: 14px;
  color: #111827;
}

.vac-tags {
  margin-top: 10px;
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.tag {
  font-size: 11px;
  padding: 4px 8px;
  border-radius: 999px;
  background: rgba(79,70,229,.08);
  color: #3730a3;
  font-weight: 700;
}

.vac-actions {
  margin-top: 12px;
  display: flex;
  gap: 8px;
}

.btn {
  padding: 8px 10px;
  border-radius: 10px;
  font-weight: 800;
  border: none;
  cursor: pointer;
}

.btn.danger {
  background: rgba(239,68,68,.12);
  color: #991b1b;
}
</style>
