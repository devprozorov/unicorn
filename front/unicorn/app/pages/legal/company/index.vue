<template>
  <Header />
  <main class="legal-page">
    <div class="legal-container">
      <h1 class="legal-title">Реквизиты компании</h1>
      
      <div class="legal-section">
        <h2>Юридическое лицо</h2>
        <p>ООО "Унистар"</p>
      </div>

      <div class="legal-section">
        <h2>ИНН</h2>
        <p>1103048119</p>
      </div>

      <div class="legal-section">
        <h2>ОГРН</h2>
        <p>1251100003934</p>
      </div>

      <div class="legal-section">
        <h2>Адрес</h2>
        <p>169915, Республика Коми, Воркута р-н, г Воркута, ул Суворова, д. 26а, кв. 10</p>
      </div>

      <div class="legal-section">
        <h2>E-mail поддержки</h2>
        <p><a href="mailto:tumyp2003map@mail.ru">tumyp2003map@mail.ru</a></p>
      </div>

      <div class="legal-section">
        <h2>Телефон</h2>
        <p>Выдается при обращении на почту</p>
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

a {
  color: #3b82f6;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
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