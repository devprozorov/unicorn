<script setup lang="ts">
import { useI18n } from '~/composables/useI18n'
const { t } = useI18n()
</script>

<template>
  <div class="card">
    <h2>{{ t('profileActions.title') }}</h2>

    <div class="actions">
      <NuxtLink to="/resumes" class="action">
        <span>📄</span> {{ t('profileActions.resumes') }}</NuxtLink>

      <NuxtLink to="/vacancies" class="action">
        <span>🔍</span> {{ t('profileActions.findJobs') }}</NuxtLink>

      <NuxtLink to="/applications" class="action">
        <span>💬</span> {{ t('profileActions.applications') }}</NuxtLink>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0,0,0,.06);
}

h2 {
  font-size: 18px;
  margin-bottom: 16px;
}

.actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.action {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-radius: 10px;
  background: #f3f4f6;
  text-decoration: none;
  color: #111827;
  font-weight: 500;
  transition: transform .15s ease, background .15s ease;
}

.action:hover {
  background: #e5e7eb;
  transform: translateY(-2px);
}
</style>
