<script setup lang="ts">
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'

import { useAuthStore } from '~/stores/auth'
import UserPanel from '~/components/panel/user/UserPanel.vue'
import CompanyPanel from '~/components/panel/company/CompanyPanel.vue'

definePageMeta({
  middleware: ['panel-auth']
})

const auth = useAuthStore()
</script>

<template>
  <Header />

  <main class="panel-page">
    <UserPanel v-if="auth.accountType === 'user'" />
    <CompanyPanel v-else-if="auth.accountType === 'company'" />

    <div v-else class="error">
      Unknown account type
    </div>
  </main>

  <Footer />
</template>

<style scoped>
.panel-page {
  min-height: calc(100vh - 120px);
}
.error {
  padding: 40px;
  text-align: center;
}
</style>
