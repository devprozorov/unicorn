<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
import { useI18n } from '~/composables/useI18n'

const auth = useAuthStore()
const { dict } = useI18n()

function goProfile() {
  if (auth.accountType === 'company') {
    navigateTo('/panel/company')
  } else if (auth.userId) {
    // перенаправляем на страницу профиля текущего пользователя
    navigateTo(`/profile/${auth.userId}`)
  } else {
    // если идентификатор неизвестен, отправляем на страницу авторизации
    navigateTo('/auth')
  }
}
</script>

<template>
  <header class="header">
    <NuxtLink to="/" class="logo">Unicornstar</NuxtLink>

    <nav class="nav">
      <NuxtLink to="/">{{ dict.nav.home }}</NuxtLink>
      <NuxtLink to="/jobs">{{ dict.nav.employers }}</NuxtLink>
      <NuxtLink to="/resumes">{{ dict.nav.candidates }}</NuxtLink>
      <NuxtLink to="/doc">{{ dict.nav.doc }}</NuxtLink>
    </nav>

    <div class="actions">
      <NuxtLink v-if="!auth.isAuthenticated" to="/auth" class="login-btn">
        {{ dict.nav.auth }}
      </NuxtLink>

      <button
        v-else
        class="profile-icon"
        @click="goProfile"
        :title="dict.nav.profile"
      />
    </div>
  </header>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 24px;
}

.actions {
  display: flex;
  align-items: center;
}

.profile-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #2dd4bf, #38bdf8);
  border: none;
  cursor: pointer;
}
</style>
