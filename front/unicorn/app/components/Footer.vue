<script setup lang="ts">
import { useI18n } from '~/composables/useI18n'
const { dict } = useI18n()
</script>

<template>
  <footer class="footer">
    <div class="container">

      <div class="footer-section">
        <div class="legal-links">
          <NuxtLink to="/legal/privacy" class="legal-link">{{ dict.footer.links.privacy }}</NuxtLink>
          <NuxtLink to="/legal/offer" class="legal-link">{{ dict.footer.links.offer }}</NuxtLink>
          <NuxtLink to="/legal/consent" class="legal-link">{{ dict.footer.links.consent }}</NuxtLink>
          <NuxtLink to="/legal/company" class="legal-link">{{ dict.footer.links.company }}</NuxtLink>
        </div>
      </div>


      <div class="footer-section">
        <div class="right">
          {{ dict.footer.rights }}
        </div>
      </div>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  border-top: 1px solid #e5e7eb;
  background: #fff;
  margin-top: 24px;
}

.container {
  max-width: 1280px;
  margin: 0 auto;
  padding: 20px 16px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  font-size: 13px;
  color: #6b7280;
  gap: 24px;
}

.footer-section {
  display: flex;
  flex-direction: column;
}

.legal-links {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 16px;
  align-items: center;
}

.legal-link {
  color: #6b7280;
  text-decoration: none;
  transition: color 0.2s;
  font-size: 13px;
}

.legal-link:hover {
  color: #374151;
  text-decoration: underline;
}

/* =====================
 * MOBILE
 * ===================== */
@media (max-width: 768px) {
  .container {
    flex-direction: column;
    gap: 16px;
    text-align: center;
    align-items: center;
  }

  .footer-section {
    width: 100%;
    align-items: center;
  }

  .legal-links {
    flex-direction: column;
    gap: 8px;
    align-items: center;
  }

  .legal-link {
    font-size: 12px;
  }
}
</style>
