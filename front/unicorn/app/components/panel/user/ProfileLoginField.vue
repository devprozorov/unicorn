<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  login: string
}>()

const visible = ref(false)
</script>

<template>
  <div class="login">
    <span class="value">
      {{ visible ? login : '*****' }}
    </span>

    <img
      class="icon"
      :src="visible ? '/icons/eye.png' : '/icons/hidden.png'"
      @click="visible = !visible"
    />
  </div>
</template>

<style scoped>
.login {
  display: flex;
  align-items: center;
  gap: 8px;
}
.value {
  font-weight: 600;
}
.icon {
  width: 18px;
  cursor: pointer;
  opacity: .7;
  transition: transform .15s ease, opacity .15s ease;
}
.icon:hover {
  transform: scale(1.1);
  opacity: 1;
}
</style>
