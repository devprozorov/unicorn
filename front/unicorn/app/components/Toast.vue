<template>
  <div v-if="toastVisible" class="toast">
    {{ toastMessage }}
  </div>
</template>

<script setup lang="ts">
import { toastMessage, toastVisible } from '~/composables/useToast'
</script>

<style scoped>
.toast {
  position: fixed;
  bottom: 24px;
  right: 24px;

  background: #111827;
  color: white;

  padding: 12px 16px;
  border-radius: 10px;
  font-size: 14px;

  box-shadow: 0 10px 30px rgba(0,0,0,.2);
}
</style>
