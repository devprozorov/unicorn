<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from '~/composables/useI18n'

const { t } = useI18n()

const props = defineProps<{
  profile: any
  isOwner: boolean
  resumesCount?: number
}>()

const showLogin = ref(false)

const loginValue = computed(() => props.profile?.login || t('panel.user.type'))

const displayName = computed(() =>
  props.profile?.displayName ||
  props.profile?.login || t('panel.user.type')
)

const avatarUrl = computed(() => {
  return props.profile?.avatarUrl || '/images/user-base.webp'
})

const has2FA = computed(() =>
  props.isOwner && typeof props.resumesCount === 'number'
)

function toggleLogin() {
  if (!props.isOwner || !loginValue.value) return
  showLogin.value = !showLogin.value
}
</script>

<template>
  <section class="profile-header">
    <!-- AVATAR -->
    <div class="avatar-wrap">
      <div class="avatar">
        <img
          :src="avatarUrl"
          alt="Avatar"
        />
      </div>

      <NuxtLink
        v-if="isOwner"
        to="/profile/edit"
        class="edit-btn"
      >{{ t('profileHeader.edit') }}</NuxtLink>
    </div>

    <!-- INFO -->
    <div class="info">
      <h2 class="name">
        {{ displayName }}
      </h2>

      <!-- LOGIN -->
      <div
        v-if="isOwner && loginValue"
        class="login"
      >
        <span>
          {{ showLogin ? loginValue : '*****' }}
        </span>

        <button
          class="eye-btn"
          @click="toggleLogin"
        >
          ??
        </button>
      </div>

      <!-- PILLS -->
      <div class="pills">
        <span class="pill resume">
          {{ resumesCount ?? 0 }} {{ t('profileHeader.resumesCount') }}</span>

        <span
          class="pill twofa"
          :class="{ active: has2FA }"
        >
          2FA
        </span>
      </div>

      <!-- ABOUT -->
      <div class="block">
        <div class="label">{{ t('profileHeader.about') }}</div>
        <div class="value">
          {{ profile?.about || t('profileHeader.notSpecified') }}
        </div>
      </div>

      <!-- LOCATION -->
      <div class="block">
        <div class="label">{{ t('profileHeader.location') }}</div>
        <div class="value">
          {{ profile?.location || t('profileHeader.notSpecified') }}
        </div>
      </div>

      <!-- LINKS -->
      <div class="block">
        <div class="label">{{ t('profileHeader.links') }}</div>
        <div v-if="profile?.links?.length" class="links">
          <a
            v-for="l in profile.links"
            :key="l"
            :href="l"
            target="_blank"
          >
            {{ l }}
          </a>
        </div>
        <div v-else class="muted">{{ t('profileHeader.noLinks') }}</div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.profile-header {
  display: grid;
  grid-template-columns: 220px 1fr;
  gap: 28px;
  background: #f9fafb;
  border-radius: 24px;
  padding: 24px;
}

/* AVATAR */
.avatar-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 256px;
  height: 256px;
  border-radius: 20px;
  overflow: hidden;
  background: #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 700;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: inherit;
  display: block;
}

.edit-btn {
  width: 100%;
  padding: 10px 14px;
  border-radius: 14px;
  background: #111827;
  color: #fff;
  text-align: center;
  text-decoration: none;
  font-size: 14px;
}

/* INFO */
.info {
  display: flex;
  flex-direction: column;
}

.name {
  font-size: 26px;
  font-weight: 700;
  margin: 0;
}

.login {
  margin-top: 6px;
  display: flex;
  gap: 10px;
  align-items: center;
  color: #6b7280;
}

.eye-btn {
  border: none;
  background: none;
  cursor: pointer;
  opacity: .7;
}

/* PILLS */
.pills {
  margin-top: 12px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.pill {
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
}

.pill.resume {
  background: rgba(239,68,68,.15);
  color: #ef4444;
}

.pill.twofa {
  background: #e5e7eb;
  color: #6b7280;
}

.pill.twofa.active {
  background: rgba(34,197,94,.2);
  color: #166534;
}

/* BLOCKS */
.block {
  margin-top: 14px;
}

.label {
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  margin-bottom: 4px;
}

.value {
  font-size: 14px;
}

.links {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.links a {
  font-size: 14px;
  color: #111827;
}

.muted {
  font-size: 14px;
  color: #9ca3af;
}

/* MOBILE */
@media (max-width: 768px) {
  .profile-header {
    grid-template-columns: 1fr;
    padding: 20px;
  }

  .avatar {
    width: 120px;
    height: 120px;
  }

  .name {
    text-align: center;
  }

  .info {
    align-items: center;
    text-align: center;
  }

  .block {
    width: 100%;
  }
}
</style>
