<template>
  <div class="page">
    <Header />

    <main class="content">
      <div class="edit-card">
        <h1>Редактирование резюме</h1>

        <form v-if="!loading" @submit.prevent="save">
          <label>
            Название
            <input
              v-model="form.title"
              type="text"
              placeholder="Введите название"
            />
          </label>

          <label>
            О себе
            <textarea
              v-model="form.about"
              rows="5"
              placeholder="Расскажите о себе"
            />
          </label>

          <label>
            Навыки
            <HierarchicalTagSelector
              :selected-tags="selectedTags"
              @toggle-tag="toggleTag"
              @clear-all="clearAllTags"
              :enable-search="true"
            />
          </label>


          <label>
            Ссылки
            <input
              v-model="linksString"
              type="text"
              placeholder="Введите ссылки через запятую"
            />
          </label>

          <div class="actions">
            <button
              type="submit"
              class="btn primary"
              :disabled="saving"
            >
              {{ saving ? 'Сохранение…' : 'Сохранить изменения' }}
            </button>
          </div>
        </form>

        <p v-else class="loading">Загрузка…</p>
      </div>
    </main>

    <Footer />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useResumesApi } from '~/services/resumesApi'
import { showToast } from '~/composables/useToast'
import { useI18n } from '~/composables/useI18n'
import { useAuthStore } from '~/stores/auth'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import HierarchicalTagSelector from '~/components/tags/HierarchicalTagSelector.vue'
const { t } = useI18n()
const auth = useAuthStore()

/* =====================
 * ROUTE & API
 * ===================== */
const route = useRoute()
const router = useRouter()
const api = useResumesApi()

const resumeId = route.params.id as string
if (!resumeId) {
  throw new Error('[ResumeEdit] resumeId is required')
}

// выбранные теги = skills
const selectedTags = computed<string[]>({
  get: () => form.value.skills,
  set: v => (form.value.skills = v),
})

function toggleTag(tag: string) {
  if (selectedTags.value.includes(tag)) {
    selectedTags.value = selectedTags.value.filter(t => t !== tag)
  } else {
    selectedTags.value = [...selectedTags.value, tag]
  }
}

function clearAllTags() {
  selectedTags.value = []
}


/* =====================
 * STATE
 * ===================== */
const loading = ref(true)
const saving = ref(false)

const form = ref({
  title: '',
  about: '',
  skills: [] as string[],
  links: [] as string[],
})

/* =====================
 * COMPUTED
 * ===================== */


const linksString = computed({
  get: () => form.value.links.join(', '),
  set: v =>
    (form.value.links = v
      .split(',')
      .map(s => s.trim())
      .filter(Boolean)),
})

/* =====================
 * LOAD (OWNER)
 * ===================== */
onMounted(async () => {
  loading.value = true
  try {
    const resume = await api.getMyResume(resumeId)
    if (!resume) {
      showToast('Резюме не найдено или недоступно')
      await router.replace('/profile')
      return
    }

    form.value = {
      title: resume.title ?? '',
      about: resume.about ?? '',
      skills: resume.skills ?? [],
      links: resume.links ?? [],
    }
  } catch (e) {
    console.error('[ResumeEdit] load failed', e)
    showToast('Ошибка загрузки резюме')
    await router.replace('/profile')
  } finally {
    loading.value = false
  }
})

/* =====================
 * SAVE (FULL PATCH)
 * ===================== */
const save = async () => {
  if (saving.value) return
  saving.value = true

  try {
    await api.updateResume(resumeId, {
      title: form.value.title,
      about: form.value.about,
      skills: form.value.skills,
      links: form.value.links,
    })

    // Проверяем статус TOTP после сохранения
    if (auth.mfaEnabled === false) {
      showToast('Для сохранения изменений необходимо подключить двухэтапную аутентификацию в профиле')
      await router.push('/profile/edit')
      return
    }

    showToast('Изменения сохранены')
  } catch (e) {
    console.error('[ResumeEdit] save failed', e)
    
    // Проверяем, не является ли ошибка требованием MFA
    if (auth.mfaEnabled === false) {
      showToast('Для сохранения изменений необходимо подключить двухэтапную аутентификацию в профиле')
      await router.push('/profile/edit')
    } else {
      showToast('Ошибка сохранения')
    }
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background: #f5f7fb;
}

.content {
  flex: 1;
  display: flex;
  justify-content: center;
  padding: 96px 16px 64px;
}

.edit-card {
  width: 100%;
  max-width: 680px;
  background: #ffffff;
  border-radius: 20px;
  padding: 36px;
  box-shadow:
    0 10px 30px rgba(0, 0, 0, 0.05),
    0 1px 3px rgba(0, 0, 0, 0.04);
}

h1 {
  font-size: 24px;
  font-weight: 700;
  text-align: center;
  margin-bottom: 28px;
}

label {
  display: block;
  margin-bottom: 18px;
  font-weight: 600;
  color: #111827;
}

input,
textarea {
  width: 100%;
  margin-top: 8px;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  font-size: 15px;
  background: #fff;
  color: #111;
}

select {
  width: 100%;
  margin-top: 8px;
  padding: 12px 14px;
  border-radius: 10px;
  border: 1px solid #e5e7eb;
  font-size: 15px;
  background: #fff;
  color: #111;
}

input:focus,
textarea:focus {
  border-color: #4f46e5;
  outline: none;
}

textarea {
  resize: vertical;
}

.actions {
  margin-top: 32px;
  display: flex;
  justify-content: center;
}

.btn {
  min-width: 240px;
  padding: 14px 20px;
  border-radius: 12px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s ease, transform 0.05s ease;
}

.btn.primary {
  background: #4f46e5;
  color: #ffffff;
}

.btn.primary:hover {
  background: #4338ca;
}

.btn:active {
  transform: translateY(1px);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading {
  text-align: center;
  padding: 48px;
  color: #6b7280;
}



</style>
