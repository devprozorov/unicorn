<script setup lang="ts">
import { ref, onMounted, defineExpose } from 'vue'
import { useI18n } from '~/composables/useI18n'

const visible = ref(false)
const { t } = useI18n()

function checkConsent() {
  visible.value = !localStorage.getItem('cookie-consent')
}

onMounted(checkConsent)

function accept() {
  localStorage.setItem('cookie-consent', 'accepted')
  visible.value = false
}

function decline() {
  localStorage.setItem('cookie-consent', 'declined')
  visible.value = false
}

/**
 * 🔁 ПУБЛИЧНЫЙ МЕТОД
 * Позволяет вернуть окно обратно
 */
function resetConsent() {
  localStorage.removeItem('cookie-consent')
  visible.value = true
}

defineExpose({ resetConsent })
</script>

<template>
  <div v-if="visible" class="cookie-consent">
    <img src="/images/cookies.png" class="cookie-icon" alt="Cookies" />

    <div class="cookie-text">
      {{ t('cookie.text') }}
    </div>

    <div class="cookie-actions">
      <button class="btn ghost" @click="decline">
        {{ t('cookie.decline') }}
      </button>
      <button class="btn primary" @click="accept">
        {{ t('cookie.accept') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.cookie-consent {
  position: fixed;
  left: 16px;
  bottom: 16px;
  z-index: 1000;

  max-width: 320px;
  padding: 12px 14px;

  display: grid;
  grid-template-columns: 36px 1fr;
  gap: 6px 10px;

  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0,0,0,.12);

  font-size: 12px;
}

.cookie-icon {
  width: 36px;
  height: 36px;
}

.cookie-actions {
  display: flex;
  gap: 6px;
}

.btn {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
}

.btn.primary {
  background: #6366f1;
  color: #fff;
}

.btn.ghost {
  background: transparent;
  color: #6b7280;
}
</style>

