<template>
  <section class="max-w-7xl mx-auto px-6 pt-[24px]">
    <h2 class="emp-h2">
      {{ dict.employers.title }}
    </h2>

    <div class="emp-row">
     <FeatureCard
        v-for="c in companies"
        :key="c.id"
        :id="c.id"
        :title="c.title"
        :subtitle="c.subtitle"
     />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import FeatureCard from '~/components/FeatureCard.vue'
import { useI18n } from '~/composables/useI18n'

const { dict, t } = useI18n()


const companies = computed(() => [
  {
    id: 'pochemu-my',
    title: t('features.items.whyUs'),
    subtitle: t('features.industries.itCloud')
  },
  {
    id: 'kak-nachat',
    title: t('features.items.howStart'),
    subtitle: t('features.industries.gamedevMobile')
  },
  {
    id: 'oferta',
    title: t('features.items.offer'),
    subtitle: t('features.industries.fintechPayments')
  },
  {
    id: 'o-nas',
    title: t('features.items.aboutUs'),
    subtitle: t('features.industries.rndAi')
  }
])

</script>

<style scoped>
.emp-h2{
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 18px;
}
.emp-row{
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.features {
  margin-top: 12px;
}

@media (max-width: 768px) {
  .max-w-7xl {
    padding-top: 32px;
  }
}

</style>

