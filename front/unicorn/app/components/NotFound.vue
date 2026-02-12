<script setup lang="ts">
import { useI18n } from '~/composables/useI18n'

const { t } = useI18n()
</script>

<template>
  <section class="nf">
    <h1 class="title">{{ t('notFound.title') }}</h1>
    <p class="text">{{ t('notFound.text') }}</p>

    <NuxtLink to="/" class="btn">{{ t('notFound.home') }}</NuxtLink>
  </section>
</template>

<style scoped>
.nf {
  padding: 48px 16px;
  text-align: center;
}
.title {
  font-size: 26px;
  font-weight: 700;
  margin-bottom: 10px;
}
.text {
  color: #6b7280;
  margin-bottom: 18px;
}
.btn {
  display: inline-block;
  padding: 10px 14px;
  border-radius: 10px;
  background: #6366f1;
  color: #fff;
  text-decoration: none;
  font-weight: 500;
}
</style>
