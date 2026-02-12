<script setup lang="ts">
import { useI18n } from '~/composables/useI18n'
defineProps<{ links: string[] }>()
const { t } = useI18n()
</script>

<template>
  <div class="section">
    <h3>{{ t('profileHeader.links') }}</h3>

    <ul v-if="links?.length">
      <li v-for="l in links" :key="l">
        <a :href="l" target="_blank">{{ l }}</a>
      </li>
    </ul>

    <p v-else class="muted">{{ t('profileHeader.noLinks') }}</p>
  </div>
</template>
