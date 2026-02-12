<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ profile: any }>()

const website = computed(() => props.profile?.website || props.profile?.site)
</script>

<template>
  <div class="actions">
    <NuxtLink class="btn" to="/panel/company/profile">Редактировать профиль</NuxtLink>
    <NuxtLink class="btn" to="/panel/company/vacancies">Управлять вакансиями</NuxtLink>
    <a v-if="website" :href="website" target="_blank" rel="noopener" class="btn secondary">
      Открыть сайт
    </a>
  </div>
</template>

<style scoped>
.actions {
  display: grid;
  gap: 10px;
  margin-top: 16px;
}
.btn {
  display: inline-flex;
  justify-content: center;
  align-items: center;
  padding: 12px 14px;
  border-radius: 12px;
  background: #111827;
  color: #fff;
  font-weight: 700;
  text-decoration: none;
}
.btn.secondary {
  background: #f3f4f6;
  color: #111827;
}
</style>
