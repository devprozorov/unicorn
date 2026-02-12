<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
    <div class="w-full max-w-md space-y-6">
      <div class="bg-green-50 border-l-4 border-green-500 p-6 rounded-lg shadow">
        <div class="flex items-start gap-4">
          <div class="text-4xl">✅</div>
          <div>
            <h1 class="text-2xl font-semibold text-green-900 mb-2">Оплата успешна!</h1>
            <p class="text-green-700 mb-4">Спасибо за подписку! Ваша премиум подписка активирована.</p>
            <p class="text-sm text-gray-600 mb-4">
              Теперь вы можете создавать до 16 вакансий/резюме, ваш контент выделяется цветом и показывается с приоритетом в выдаче.
            </p>
            <button
              @click="goBack"
              class="w-full px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
            >
              Вернуться в профиль
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.push('/')
}
</script>
