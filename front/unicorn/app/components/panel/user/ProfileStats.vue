<script setup lang="ts">
import { useI18n } from '~/composables/useI18n'
const { t } = useI18n()

const props = defineProps<{
  profile: any
}>()

const hasMfa = props.profile?.mfaEnabled
</script>

<template>
  <div class="stats">
    <div class="badge resume">{{ t('profileHeader.resumesCount') }}: {{ profile.resumesCount || 0 }}</div>

    <div
      class="badge mfa"
      :class="{ enabled: hasMfa }"
    >
      🔐 2FA
    </div>
  </div>
</template>

<style scoped>
.stats {
  display: flex;
  gap: 12px;
  margin-top: 6px;
}

.badge {
  padding: 6px 14px;
  border-radius: 999px;
  font-size: 13px;
  background: rgba(0,0,0,.05);
  color: #111827;
}

.mfa {
  background: rgba(0,0,0,.08);
}

.mfa.enabled {
  background: rgba(34,197,94,.25);
}
</style>
