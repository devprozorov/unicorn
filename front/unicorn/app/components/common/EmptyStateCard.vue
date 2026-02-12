<script setup lang="ts">
defineProps<{
  title: string
  text: string
  actionLabel: string
  actionTo: string
  hint?: string
}>()
</script>

<template>
  <div class="empty-card">
    <div class="icon">
      <slot name="icon">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none">
          <path
            d="M7 3h7l5 5v13a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1V4a1 1 0 0 1 1-1z"
            stroke="currentColor"
            stroke-width="1.6"
          />
          <path d="M14 3v5h5" stroke="currentColor" stroke-width="1.6" />
        </svg>
      </slot>
    </div>

    <h3 class="title">{{ title }}</h3>
    <p class="text">{{ text }}</p>

    <NuxtLink :to="actionTo" class="btn">
      {{ actionLabel }}
    </NuxtLink>

    <p v-if="hint" class="hint">{{ hint }}</p>
  </div>
</template>

<style scoped>
.empty-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 20px;
  padding: 24px;
  text-align: center;
}

.icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  background: linear-gradient(135deg, #6366f1, #3b82f6);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 16px;
}

.title {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.text {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 16px;
  line-height: 1.5;
}

.btn {
  display: block;
  width: 100%;
  padding: 14px 18px;
  border-radius: 14px;
  background: #111827;
  color: #fff;
  font-weight: 500;
  text-decoration: none;
  margin-bottom: 10px;
}

.btn:hover {
  background: #1f2937;
}

.hint {
  font-size: 12px;
  color: #9ca3af;
}
</style>
