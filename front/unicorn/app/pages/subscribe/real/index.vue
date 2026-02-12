<template>
  <div class="min-h-screen bg-gray-50 flex items-center justify-center px-4">
    <div class="w-full max-w-md space-y-6">
      <div v-if="success" class="bg-green-50 border-l-4 border-green-500 p-6 rounded-lg shadow">
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

      <div v-else class="bg-red-50 border-l-4 border-red-500 p-6 rounded-lg shadow">
        <div class="flex items-start gap-4">
          <div class="text-4xl">❌</div>
          <div>
            <h1 class="text-2xl font-semibold text-red-900 mb-2">Оплата отменена</h1>
            <p class="text-red-700 mb-4">
              {{ message || 'Оплата была отменена или произошла ошибка.' }}
            </p>
            <p class="text-sm text-gray-600 mb-4">
              Вы можете повторить попытку, перейдя в профиль и кликнув "Подписаться".
            </p>
            <button
              @click="goBack"
              class="w-full px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
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
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const success = ref(false)
const message = ref('')

onMounted(() => {
  const query = route.query
  success.value = route.path.includes('/success')
  message.value = (query.message as string) || ''
})

const goBack = () => {
  router.push('/panel')
}
</script>
