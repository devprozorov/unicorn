<script setup lang="ts">
import axios from 'axios'
import { ref, computed, reactive, onMounted } from 'vue'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'
import { useAuthStore } from '~/stores/auth'

const auth = useAuthStore()
const { t } = useI18n()

type Vacancy = {
  vacancyId: string
  companyId: string
  title: string
  description: string
  location?: string
  tags?: string[]
}

const all = ref<Vacancy[]>([])
const err = ref('')
const hint = ref('')
const loading = ref(false)

const myVacancies = computed(() =>
  all.value.filter(v => v.companyId === auth.userId)
)

const modal = reactive({
  open: false,
  mode: 'create' as 'create' | 'edit',
  editingId: '',
  form: { title: '', description: '', location: '' },
  tagsText: '',
})

function parseTags(s: string) {
  return s.split(',').map(x => x.trim()).filter(Boolean).slice(0, 20)
}

function openCreate() {
  modal.open = true
  modal.mode = 'create'
  modal.editingId = ''
  modal.form = { title: '', description: '', location: '' }
  modal.tagsText = ''
}

function openEdit(v: Vacancy) {
  modal.open = true
  modal.mode = 'edit'
  modal.editingId = v.vacancyId
  modal.form.title = v.title
  modal.form.description = v.description
  modal.form.location = v.location || ''
  modal.tagsText = (v.tags || []).join(', ')
}

function closeModal() {
  modal.open = false
}

async function load() {
  err.value = ''
  try {
    const res = await axios.get('/vacancies')
    all.value = res.data.items || []
  } catch (e: any) {
    err.value = getErrorMessage(e?.response?.data?.error || 'load_failed', t, 'errors.loadFailed')
  }
}

async function saveVacancy() {
  loading.value = true
  try {
    const payload = {
      title: modal.form.title,
      description: modal.form.description,
      location: modal.form.location,
      tags: parseTags(modal.tagsText),
    }

    if (modal.mode === 'create') {
      await axios.post('/vacancies', payload)
    } else {
      await axios.patch(`/vacancies/${modal.editingId}`, payload)
    }

    closeModal()
    await load()
    hint.value = t('companyVacancies.saved')
  } catch (e: any) {
    const code = e?.response?.data?.error
    hint.value =
      code === 'mfa_required'
        ? t('companyVacancies.enableTotp')
        : code === 'limit_reached'
        ? t('companyVacancies.limitReached')
        : getErrorMessage(code || 'save_failed', t, 'errors.saveFailed')
  } finally {
    loading.value = false
  }
}

async function remove(v: Vacancy) {
  if (!confirm(t('companyVacancies.confirmDelete'))) return
  try {
    await axios.delete(`/api/vacancies/${v.vacancyId}`)
    await load()
  } catch (e: any) {
    hint.value = getErrorMessage(e?.response?.data?.error || 'delete_failed', t, 'errors.deleteFailed')
  }
}

onMounted(load)
</script>

<template>
  <div class="card">
    <div class="header">
      <h2>{{ t('companyVacancies.title') }}</h2>
      <button class="btn-black" @click="openCreate">{{ t('companyVacancies.newButton') }}</button>
    </div>

    <div v-if="err" class="error">{{ err }}</div>
    <div v-if="hint" class="hint">{{ hint }}</div>

    <div v-if="myVacancies.length" class="list">
      <div v-for="v in myVacancies" :key="v.vacancyId" class="vacancy">
        <strong>{{ v.title }}</strong>
        <div class="muted">{{ v.location || t('profileHeader.notSpecified') }}</div>
       <div class="text-sm mt-2 whitespace-pre-wrap max-h-24 overflow-hidden"> {{ v.description }} </div>

        <div class="actions">
          <button @click="openEdit(v)">{{ t('companyVacancies.edit') }}</button>
          <button class="danger" @click="remove(v)">{{ t('companyVacancies.deleteBtn') }}</button>
        </div>
      </div>
    </div>

    <div v-else class="empty">
      <h3>{{ t('companyVacancies.emptyTitle') }}</h3>
      <p>{{ t('companyVacancies.emptyText') }}</p>
    </div>

    <!-- MODAL -->
    <div v-if="modal.open" class="modal">
      <div class="modal-card">
        <input v-model="modal.form.title" :placeholder="t('companyVacancies.modal.titlePlaceholder')" />
        <input v-model="modal.form.location" :placeholder="t('companyVacancies.modal.cityPlaceholder')" />
        <input v-model="modal.tagsText" :placeholder="t('companyVacancies.modal.tagsPlaceholder')" />
        <textarea v-model="modal.form.description" rows="5" :placeholder="t('companyVacancies.modal.descriptionPlaceholder')" />

        <div class="modal-actions">
          <button @click="closeModal">{{ t('companyVacancies.modal.cancel') }}</button>
          <button class="btn-black" @click="saveVacancy">{{ t('companyVacancies.modal.saveChanges') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: #fff;
  border-radius: 20px;
  padding: 24px;
  border: 1px solid #e5e7eb;
}
.header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}
.vacancy {
  border-bottom: 1px solid #e5e7eb;
  padding: 12px 0;
}
.actions {
  display: flex;
  gap: 8px;
}
.danger {
  color: #dc2626;
}
.modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.4);
  display: flex;
  align-items: center;
  justify-content: center;
}
.modal-card {
  background: #fff;
  padding: 24px;
  border-radius: 16px;
  width: 100%;
  max-width: 520px;
}
.btn-black {
  background: #111827;
  color: #fff;
  padding: 10px 16px;
  border-radius: 12px;
}
</style>
