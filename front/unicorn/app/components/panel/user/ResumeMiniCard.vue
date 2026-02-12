<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from '~/composables/useI18n'

const { t } = useI18n()

const props = defineProps<{
  resume: {
    id: string
    title: string
    about?: string
    skills?: string[]
    links?: string[]
  }
}>()

const filledPercent = computed(() => {
  let filled = 0
  if (props.resume.title) filled++
  if (props.resume.about) filled++
  if (props.resume.skills?.length) filled++
  if (props.resume.links?.length) filled++
  return Math.round((filled / 4) * 100)
})

const isPublished = computed(() => filledPercent.value === 100)
</script>

<template>
  <div class="resume-card">
    <div class="header">
      <h3>{{ resume.title || t('resumeMini.noTitle') }}</h3>
      <span
        class="badge"
        :class="isPublished ? 'published' : 'draft'"
      >
        {{ isPublished ? t('resumeMini.published') : t('resumeMini.draft') }}
      </span>
    </div>

    <div class="progress">
      <div class="bar">
        <div
          class="fill"
          :style="{ width: filledPercent + '%' }"
        />
      </div>
      <span>{{ filledPercent }}%</span>
    </div>

    <div class="actions">
      <NuxtLink :to="`/resumes/${resume.id}/edit`">{{ t('resumeMini.edit') }}</NuxtLink>
      <NuxtLink :to="`/resumes/${resume.id}`">{{ t('resumeMini.view') }}</NuxtLink>
    </div>
  </div>
</template>

<style scoped>
.resume-card {
  background: #fff;
  border-radius: 20px;
  padding: 20px;
  border: 1px solid #e5e7eb;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.badge {
  font-size: 11px;
  padding: 4px 8px;
  border-radius: 999px;
  font-weight: 500;
}

.badge.draft {
  background: #fef3c7;
  color: #92400e;
}

.badge.published {
  background: #dcfce7;
  color: #166534;
}

.progress {
  margin: 12px 0;
}

.bar {
  height: 6px;
  background: #e5e7eb;
  border-radius: 999px;
}

.fill {
  height: 100%;
  background: #6366f1;
}

.actions {
  display: flex;
  gap: 12px;
}
</style>
