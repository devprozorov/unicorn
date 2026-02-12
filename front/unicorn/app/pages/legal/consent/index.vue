<template>
  <Header />
  <main class="legal-page">
    <div class="legal-container">
      <h1 class="legal-title">Согласие на обработку персональных данных</h1>
      
      <div class="legal-section">
        <p>Используется совместно с политикой обработки персональных данных.</p>
      </div>

      <div class="legal-section">
        <p>
          Настоящим я, пользователь сервиса, выражаю согласие на обработку моих персональных данных в целях 
          предоставления функциональности сервиса: регистрации, создания и публикации поздравлений, поддержки 
          и связи со мной по вопросам работы сервиса.
        </p>
      </div>

      <div class="legal-section">
        <h2>Перечень данных</h2>
        <p>
          Логин, e-mail, контент поздравлений, загруженные фото, технические данные (IP, cookies).
        </p>
      </div>

      <div class="legal-section">
        <h2>Срок действия</h2>
        <p>
          До достижения целей обработки или до отзыва согласия, если иное не требуется законом.
        </p>
      </div>

      <div class="legal-section">
        <h2>Отзыв согласия</h2>
        <p>
          Согласие может быть отозвано путём обращения в поддержку, указанную в реквизитах.
        </p>
      </div>
    </div>
  </main>
  <Footer />
</template>

<script setup lang="ts">
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
</script>

<style scoped>
.legal-page {
  min-height: 70vh;
  padding: 32px 16px;
  background: #f9fafb;
}

.legal-container {
  max-width: 800px;
  margin: 0 auto;
  background: white;
  padding: 32px;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.legal-title {
  font-size: 28px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 24px;
  line-height: 1.3;
}

.legal-section {
  margin-bottom: 24px;
}

.legal-section:last-child {
  margin-bottom: 0;
}

h2 {
  font-size: 20px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
}

p {
  font-size: 15px;
  line-height: 1.6;
  color: #4b5563;
  margin-bottom: 8px;
}

p:last-child {
  margin-bottom: 0;
}

@media (max-width: 768px) {
  .legal-page {
    padding: 24px 12px;
  }

  .legal-container {
    padding: 24px 20px;
  }

  .legal-title {
    font-size: 24px;
  }

  h2 {
    font-size: 18px;
  }

  p {
    font-size: 14px;
  }
}
</style>