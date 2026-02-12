<script setup lang="ts">
import ProfileLinks from './ProfileLinks.vue'
import { useI18n } from '~/composables/useI18n'
const { t } = useI18n()
</script>

<template>
  <section class="card">
        <div class="actions">
      <NuxtLink to="/profile/edit" class="btn primary">
        {{ t('profileContent.editProfile') }}
      </NuxtLink>

      <NuxtLink to="/resumes/create" class="btn secondary">
        {{ t('profileContent.createResume') }}
      </NuxtLink>
    </div>


        <div class="section">
      <h3>{{ t('profileContent.about') }}</h3>
      <p>{{ profile.about || t('profileHeader.notSpecified') }}</p>
    </div>


        <div class="section">
      <h3>{{ t('profileContent.location') }}</h3>
      <p>{{ profile.location || t('profileHeader.notSpecified') }}</p>
    </div>


    <ProfileLinks :links="profile.links" />
  </section>
</template>

<style scoped>
.actions {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.btn {
  padding: 10px 16px;
  border-radius: 12px;
  font-weight: 600;
  text-decoration: none;
}

.primary {
  background: #6366f1;
  color: #fff;
}

.secondary {
  background: #f3f4f6;
  color: #111827;
}
</style>
