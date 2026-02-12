<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-50 via-pink-50 to-white">

    <!-- Hero Section -->
    <section class="py-12 md:py-20 px-4">
      <div class="max-w-5xl mx-auto text-center">
        <div class="inline-flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-purple-100 to-pink-100 rounded-full mb-6">
          <span class="text-2xl">⭐</span>
          <span class="text-purple-700 font-semibold text-sm">Премиум подписка</span>
        </div>
        
        <h1 class="text-4xl md:text-5xl lg:text-6xl font-bold text-gray-900 mb-6 leading-tight">
          Откройте все возможности <br class="hidden md:block">
          <span class="bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-transparent">
            UnicornStar Premium
          </span>
        </h1>
        
        <p class="text-lg md:text-xl text-gray-600 mb-12 max-w-3xl mx-auto">
          Выделяйтесь среди конкурентов, получайте больше откликов и находите идеальных кандидатов быстрее
        </p>

        <!-- Main Pricing Card -->
        <div class="bg-white rounded-3xl shadow-2xl p-6 md:p-12 mb-12 border-2 border-purple-200 relative overflow-hidden max-w-4xl mx-auto">
          <!-- Background decoration -->
          <div class="absolute top-0 right-0 w-64 h-64 bg-gradient-to-br from-purple-100 to-pink-100 rounded-full blur-3xl opacity-30 -mr-32 -mt-32"></div>
          <div class="absolute bottom-0 left-0 w-64 h-64 bg-gradient-to-tr from-pink-100 to-purple-100 rounded-full blur-3xl opacity-30 -ml-32 -mb-32"></div>
          
          <div class="relative">
            <!-- Price Header -->
            <div class="flex flex-col lg:flex-row items-center justify-between gap-8 mb-10">
              <div class="text-center lg:text-left">
                <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-3">
                  Premium подписка
                </h2>
                <p class="text-gray-600 text-lg">Всё для успешного найма и поиска работы</p>
              </div>
              <div class="text-center lg:text-right">
                <div class="flex items-baseline gap-2 justify-center lg:justify-end mb-2">
                  <span class="text-5xl md:text-6xl font-bold bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-transparent">
                    250₽
                  </span>
                  <span class="text-gray-500 text-xl">/мес</span>
                </div>
                <p class="text-sm text-green-600 font-semibold mt-1">
                  💰 Это как 1 шаурма в сырном лаваше
                </p>
              </div>
            </div>

            <!-- CTA Button -->
            <button 
              @click="subscribe"
              class="w-full bg-gradient-to-r from-purple-600 to-pink-600 text-white text-lg md:text-xl font-bold py-5 px-8 rounded-2xl hover:from-purple-700 hover:to-pink-700 transform hover:scale-[1.02] transition-all duration-200 shadow-lg hover:shadow-2xl mb-6 flex items-center justify-center gap-3"
            >
              <span class="text-2xl">🚀</span>
              <span>Подписаться на Premium</span>
            </button>

            <!-- Trust Badges -->
            <div class="flex flex-wrap items-center justify-center gap-4 md:gap-6 text-sm text-gray-600 mb-10">
              <div class="flex items-center gap-2">
                <svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
                <span class="font-medium">Не подписка, покупка каждый раз в ручном режиме, вы не получите оплату которую не ожидали</span>
              </div>
              <div class="flex items-center gap-2">
                <svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
                <span class="font-medium">Без комиссий</span>
              </div>
              <div class="flex items-center gap-2">
                <svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
                </svg>
                <span class="font-medium">Без скрытых условий</span>
              </div>
            </div>

            <!-- Quick Benefits -->
            <div class="grid md:grid-cols-2 gap-4">
              <div class="bg-gradient-to-br from-purple-50 to-pink-50 p-4 rounded-xl text-center">
                <div class="text-3xl mb-2">📝</div>
                <div class="font-bold text-gray-900 mb-1">До 16</div>
                <div class="text-sm text-gray-600">вакансий/резюме</div>
              </div>
              <div class="bg-gradient-to-br from-purple-50 to-pink-50 p-4 rounded-xl text-center">
                <div class="text-3xl mb-2">🎨</div>
                <div class="font-bold text-gray-900 mb-1">Цветной акцент</div>
                <div class="text-sm text-gray-600">и приоритет в выдаче</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    

    <!-- Comparison Section -->
    <section class="py-16 px-4 bg-gradient-to-b from-gray-50 to-white">
      <div class="max-w-5xl mx-auto">
        <h2 class="text-3xl md:text-4xl font-bold text-center text-gray-900 mb-4">
          Сравнение тарифов
        </h2>
        <p class="text-center text-gray-600 mb-12">
          Посмотрите, что доступно в бесплатной и премиум версиях
        </p>

        <div class="bg-white rounded-3xl shadow-2xl overflow-hidden border border-gray-100">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gradient-to-r from-purple-600 to-pink-600 text-white">
                <tr>
                  <th class="px-4 md:px-8 py-5 text-left font-bold text-base md:text-lg">Возможности</th>
                  <th class="px-4 md:px-8 py-5 text-center font-bold text-base md:text-lg">Бесплатно</th>
                  <th class="px-4 md:px-8 py-5 text-center font-bold text-base md:text-lg">Premium ⭐</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100">
                <tr v-for="(row, idx) in comparisonTable" :key="idx" class="hover:bg-gray-50 transition-colors">
                  <td class="px-4 md:px-8 py-5 text-gray-900 font-semibold text-sm md:text-base">{{ row.feature }}</td>
                  <td class="px-4 md:px-8 py-5 text-center">
                    <span v-if="row.free === true" class="text-green-500 text-3xl">✓</span>
                    <span v-else-if="row.free === false" class="text-gray-300 text-3xl">−</span>
                    <span v-else class="text-gray-600 text-sm md:text-base font-medium">{{ row.free }}</span>
                  </td>
                  <td class="px-4 md:px-8 py-5 text-center">
                    <span v-if="row.premium === true" class="text-purple-600 text-3xl font-bold">✓</span>
                    <span v-else class="text-purple-600 font-bold text-sm md:text-base">{{ row.premium }}</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup>
import { ref } from 'vue'

// SEO Meta tags
useHead({
  title: 'Premium подписка - UnicornStar | Найм и поиск работы',
  meta: [
    { name: 'description', content: 'Премиум подписка UnicornStar: до 16 вакансий и резюме, цветной акцент и приоритет в выдаче.' },
    { name: 'keywords', content: 'UnicornStar Premium, премиум подписка, вакансии, резюме, работа, найм сотрудников, HeadHunter альтернатива' },
    { property: 'og:title', content: 'Premium подписка - UnicornStar' },
    { property: 'og:description', content: 'Выделяйтесь среди конкурентов с Premium подпиской UnicornStar. 14 дней бесплатно!' },
    { property: 'og:type', content: 'website' },
    { property: 'og:image', content: '/og-premium.jpg' },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: 'Premium подписка - UnicornStar' },
    { name: 'twitter:description', content: 'Все возможности для успешного найма и поиска работы' },
    { name: 'robots', content: 'index, follow' },
    { name: 'author', content: 'UnicornStar' }
  ],
  htmlAttrs: {
    lang: 'ru'
  },
  link: [
    { rel: 'canonical', href: 'https://unicornstar.com/subscribe/offer' }
  ]
})

const comparisonTable = [
  { feature: 'Количество вакансий/резюме', free: 'До 2-х', premium: 'До 16' },
  { feature: 'Цветной акцент в выдаче', free: false, premium: true },
  { feature: 'Приоритет в результатах поиска', free: false, premium: true }
]

const subscribe = () => {
  alert('Подписку можно оформить в профиле')
}
</script>

<style scoped>
/* Smooth animations */
* {
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
}

/* Smooth scrolling */
html {
  scroll-behavior: smooth;
}
</style>