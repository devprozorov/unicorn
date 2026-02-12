<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
    <div class="w-full max-w-md space-y-6">
      <div class="bg-red-50 border-l-4 border-red-500 p-6 rounded-lg shadow">
        <div class="flex items-start gap-4">
          <div class="text-4xl">❌</div>
          <div>
            <h1 class="text-2xl font-semibold text-red-900 mb-2">Оплата отменена</h1>
            <p class="text-red-700 mb-4">
              Оплата была отменена или произошла ошибка.
            </p>
            <p class="text-sm text-gray-600 mb-4">
              Вы можете повторить попытку, перейдя в профиль и кликнув "Подписаться".
            </p>
            <div class="space-y-2 mb-4 text-sm text-gray-700">
              <p>• Проверьте данные карты</p>
              <p>• Убедитесь в наличии средств</p>
              <p>• Попробуйте другой способ оплаты</p>
            </div>
            <button
              @click="goBack"
              class="w-full px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
            >
              Вернуться в профиль
            </button>
            <button
              @click="tryAgain"
              class="w-full mt-2 px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
            >
              Попробовать снова
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useSubscription } from '~/composables/useSubscription'

const router = useRouter()
const { goToPayment } = useSubscription()

const goBack = () => {
  router.push('/panel')
}

const tryAgain = async () => {
  await goToPayment()
}
</script>
