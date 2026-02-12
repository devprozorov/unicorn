<template>
  <Header />

  <div class="job-page">
    <!-- BACK -->
    <NuxtLink class="back-link" to="/jobs">{{ t('job.backToJobs') }}</NuxtLink>

    <!-- ERRORS -->
    <div v-if="err" class="alert error">{{ err }}</div>
    <div v-if="hint" class="alert hint">{{ hint }}</div>

    <!-- =====================
     | VACANCY HEADER
     ===================== -->
    <div v-if="vacancy" class="vacancy-card">
      <div class="vacancy-head">
        <div class="vacancy-main">
          <div class="title-row">
            <h1 class="title">{{ vacancy.title }}</h1>


          </div>

          <div class="meta">
            {{ vacancy.location || t('profileHeader.notSpecified') }}
          </div>

          <!-- TAGS -->
          <div class="tags">
            <span v-for="t in vacancy.tags || []" :key="t" class="tag">
              {{ tagLabel(t) }}
            </span>
          </div>

          <!-- ANALYTICS -->
          <div class="analytics-grid">
  <!-- APPLIES -->
  <div class="metric">
    <img src="/images/job-application.png" alt="" />
    <div>
      <div class="value">{{ vacancy?.responsesCount ?? 0 }}</div>
      <div class="label">{{ t('job.applies') }}</div>
    </div>
  </div>

  <!-- DAYS AGO -->
  <div class="metric">
    <img src="/images/pending.png" alt="" />
    <div>
      <div class="value">{{ analytics.daysAgo }}</div>
      <div class="label">{{ t('job.daysAgo') }}</div>
    </div>
  </div>
</div>




          <!-- MATCH -->
          <div
            v-if="match.score !== null"
            class="match"
            :class="scoreClass(match.score)"
          >
            {{ t('job.resumeMatch') }} {{ match.score }}%
          </div>
        </div>

        <!-- DESKTOP CTA -->
        <div class="cta desktop">
          <button
            v-if="!myApp"
            class="btn primary"
            @click="openApply"
          >
            {{ t('job.apply') }}
          </button>

          <div v-else class="application-status">
  <div class="status-pill" :data-status="myApp.status">
    <span class="dot"></span>
    <span class="text">
      {{ myApp.status === 'pending' ? t('job.applicationSent') : myApp.status }}
    </span>
  </div>

  <NuxtLink
    class="chat-link"
    :to="`/user/chat/${myApp.applicationId}`"
  >
    {{ t('job.openChat') }}
  </NuxtLink>
</div>

        </div>
      </div>

      <!-- DESCRIPTION -->
      <div class="description">
        {{ vacancy.description }}
      </div>
    </div>

    <!-- =====================
     | COMPANY
     ===================== -->
    <div v-if="company" class="company-card" :style="{'--brand': company.brandColor || '#6366f1'}">
      <div class="company-head">
        <img
          :src="companyAvatar"
          class="company-avatar"
        />

        <div>
          <div class="company-name">{{ company.displayName }}</div>
          <div class="company-meta">
            {{ company.location || t('profileHeader.notSpecified') }} · {{ company.industry || t('profileHeader.notSpecified') }}
          </div>

          <div>{{ company.about || t('job.companyDescriptionMissing') }}</div>

          <a
            v-if="company.website"
            class="company-link"
            :href="company.website"
            target="_blank"
          >
            {{ company.website }}
          </a>
        </div>
      </div>
    </div>

    <!-- =====================
     | STICKY MOBILE APPLY
     ===================== -->
    <div class="sticky-apply">
      <button
        v-if="!myApp"
        class="btn primary"
        @click="openApply"
      >
        {{ t('job.apply') }}
      </button>

      <NuxtLink
        v-else
        class="btn"
        :to="`/user/chat/${myApp.applicationId}`"
      >
        {{ t('job.openChat') }}
      </NuxtLink>
    </div>

    <!-- =====================
     | APPLY BOTTOM SHEET
     ===================== -->
    <div
      v-if="apply.open"
      class="sheet-backdrop"
      @click.self="apply.open = false"
    >
      <div class="sheet">
        <div class="sheet-title">
          {{ t('job.applyToVacancy') }}
          <button @click="apply.open = false">✕</button>
        </div>

        <div v-if="hint" class="alert hint">
          {{ hint }}
        </div>

        <div v-if="resumes.length === 0" class="alert hint">
          {{ t('job.noResumesAlert') }}
        </div>

        <label>
          <span>{{ t('jobs.resume') }}</span>
          <select v-model="apply.resumeId">
            <option value="" disabled>{{ t('jobs.selectResume') }}</option>
            <option
              v-for="r in resumes"
              :key="r.resumeId"
              :value="r.resumeId"
            >
              {{ r.title }}
            </option>
          </select>
        </label>

        <!-- MATCH PREVIEW -->
        <div
          v-if="apply.resumeId && matchPreview !== null"
          class="match-preview"
        >
          <div>{{ t('job.resumeMatch') }} <b>{{ matchPreview }}%</b></div>
        </div>

        <label>
          <span>{{ t('job.messageOptional') }}</span>
          <textarea
            v-model="apply.message"
            rows="4"
            :placeholder="t('job.whyFit')"
          />
        </label>

        <button
          class="btn primary full"
          @click="sendApply"
          :disabled="loading"
        >
          {{ loading ? t('job.sending') : t('job.sendApplication') }}
        </button>
      </div>
    </div>
  </div>

  <Footer />
</template>

<script setup lang="ts">
import axios from 'axios'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'
import { tagLabelMap } from '~/utils/tags'

const { t } = useI18n()

function tagLabel(value: string) {
  const key = tagLabelMap[value]
  return key ? t(key) : value
}

const route = useRoute()
const id = computed(() => String(route.params.id || ''))

const err = ref('')
const hint = ref('')
const loading = ref(false)

const vacancy = ref<any>(null)
const company = ref<any>(null)
const resumes = ref<any[]>([])
const myApp = ref<any>(null)
const companyAvatar = computed(() => company.value?.avatarUrl || '/images/com-base.webp')

const appliesCount = computed(() => {
  return Number(vacancy.value?.responsesCount || 0)
})

const publicMeta = reactive({
  responsesCount: 0,
  createdAt: '' as string,
})

async function loadPublicMetaFromList() {
  try {
    const r = await axios.get('/vacancies')
    const list = r.data.items || []
    const item = list.find((x: any) => x.vacancyId === id.value)

    if (!item) return

    publicMeta.responsesCount = Number(item.responsesCount || 0)
    publicMeta.createdAt = String(item.createdAt || '')
  } catch {
    // молча: метрика не должна ломать страницу
  }
}

/* analytics */
const analytics = reactive({
  daysAgo: 0,
})


/* SAVE — SSR SAFE */
const saved = ref(false)

/* APPLY */
const apply = reactive({
  open: false,
  vacancy: null as any,
  resumeId: '',
  message: '',
})

/* MATCH */
const match = reactive<{ score: number | null }>({ score: null })
const matchPreview = ref<number | null>(null)

/* =====================
 * HELPERS
 * ===================== */
function scoreClass(score: number) {
  if (score >= 80) return 'high'
  if (score >= 50) return 'mid'
  return 'low'
}

function calcDaysAgo(date: string) {
  const ms = new Date(date).getTime()
  if (!Number.isFinite(ms)) return 0
  return Math.max(0, Math.floor((Date.now() - ms) / 86400000))
}



/* =====================
 * SAVE (CLIENT ONLY)
 * ===================== */
onMounted(() => {
  try {
    const arr = JSON.parse(localStorage.getItem('savedVacancies') || '[]')
    saved.value = arr.includes(id.value)
  } catch {}
})

function toggleSave() {
  if (!process.client) return

  const arr: string[] = JSON.parse(
    localStorage.getItem('savedVacancies') || '[]'
  )

  if (saved.value) {
    localStorage.setItem(
      'savedVacancies',
      JSON.stringify(arr.filter(x => x !== id.value))
    )
  } else {
    localStorage.setItem(
      'savedVacancies',
      JSON.stringify([...arr, id.value])
    )
  }

  saved.value = !saved.value
}

/* =====================
 * LOADERS
 * ===================== */

async function loadResumes() {
  try {
    const r = await axios.get('/resumes/my')
    resumes.value = r.data.items || []
  } catch (e: any) {
    const code = e?.response?.data?.error
    if (code === 'mfa_required') {
      hint.value = getErrorMessage(code, t)
    }
  }
}

async function loadMyStatus() {
  try {
    const r = await axios.get('/applications/my')
    myApp.value =
      r.data.items.find((x: any) => x.vacancyId === id.value) || null
  } catch {}
}

/* =====================
 * APPLY FLOW
 * ===================== */
function openApply() {
  apply.open = true
  apply.resumeId = resumes.value[0]?.resumeId || ''
  apply.message = ''
  apply.vacancy = vacancy.value
  matchPreview.value = null
}

watch(
  () => apply.resumeId,
  async rid => {
    if (!rid) return
    try {
      const r = await axios.post('/match/score', {
        vacancyId: id.value,
        resumeId: rid,
      })
      matchPreview.value = r.data.score
      match.score = r.data.score
    } catch {
      matchPreview.value = null
    }
  }
)

async function sendApply() {
  if (!apply.resumeId) {
    hint.value = t('jobs.selectResume')
    return
  }

  loading.value = true
  err.value = ''
  hint.value = ''

  try {
    const r = await axios.post('/applications', {
      vacancyId: id.value,
      resumeId: apply.resumeId,
      message: apply.message,
    })

    myApp.value = {
      applicationId: r.data.applicationId,
      vacancyId: id.value,
      companyId: vacancy.value.companyId,
      status: r.data.status,
    }

    apply.open = false
    hint.value = t('job.applicationSent')
  } catch (e: any) {
    const code = e?.response?.data?.error || 'apply_failed'
    hint.value = getErrorMessage(code, t, 'jobs.errors.applyFailed')
  } finally {
    loading.value = false
  }
}

async function loadVacancy() {
  const r = await axios.get(`/vacancies/${id.value}`)
  vacancy.value = r.data.vacancy

  if (!vacancy.value) {
    throw new Error('vacancy_not_found')
  }

  // ВАЖНО: берём responsesCount и createdAt из списка /api/vacancies
  await loadPublicMetaFromList()

  // daysAgo: сначала из списка, если там есть createdAt; иначе из детальной вакансии
  const createdAt = publicMeta.createdAt || vacancy.value.createdAt || ''
  analytics.daysAgo = createdAt ? calcDaysAgo(createdAt) : 0

  // responsesCount: тоже проталкиваем в vacancy, чтобы шаблон был простым
  vacancy.value.responsesCount = publicMeta.responsesCount

  const p = await axios.get(`/profile/${vacancy.value.companyId}`)
  company.value = p.data.profile
}



/* =====================
 * INIT
 * ===================== */
onMounted(async () => {
  try {
    await loadVacancy()
    await loadResumes()
    await loadMyStatus()
  } catch (e: any) {
    const code = e?.response?.data?.error || 'load_failed'
    err.value = getErrorMessage(code, t, 'errors.loadFailed')
  }
})
</script>

<style scoped>
/* BASE */
.job-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 16px;
  box-sizing: border-box;
}

.back-link {
  display: inline-block;
  margin-bottom: 12px;
}

.alert {
  padding: 12px;
  border-radius: 14px;
}
.alert.error { background: #fee2e2; color: #991b1b }
.alert.hint { background: #fef3c7; color: #92400e }

/* VACANCY */
.vacancy-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 20px;
  padding: 16px;
}

.vacancy-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
}

.title-row {
  display: flex;
  gap: 10px;
  align-items: center;
}

.save-btn {
  border: none;
  background: none;
  cursor: pointer;
  opacity: .4;
}
.save-btn.active { opacity: 1 }

.tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-top: 6px;
}

.tag {
  background: #f3f4f6;
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 11px;
}

.analytics {
  margin-top: 6px;
  font-size: 12px;
  color: #6b7280;
  display: flex;
  gap: 12px;
}

.match.high { color: #16a34a }
.match.mid { color: #ca8a04 }
.match.low { color: #6b7280 }

/* COMPANY */
.company-card {
  margin-top: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  padding: 16px;
}

/* STICKY APPLY */
.sticky-apply {
  display: none;
}

@media (max-width: 768px) {
  .vacancy-head {
    flex-direction: column;
  }

  .cta.desktop {
    display: none;
  }

  .sticky-apply {
    display: block;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background: #fff;
    border-top: 1px solid #e5e7eb;
    padding: 12px;
    z-index: 40;
  }
}

/* BOTTOM SHEET */
.sheet-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.4);
  z-index: 50;
}

.sheet {
  position: absolute;
  width: 100%;
  background: #fff;
  border-radius: 20px 20px 0 0;
  padding: 16px;
}
.sheet-title {
  display: flex;
  justify-content: space-between;
  font-weight: 600;
  margin-bottom: 12px;
}

.btn {
  padding: 10px 14px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}
.btn.primary {
  background: #111827;
  color: #fff;
}
.btn.full {
  width: 100%;
}

/* =====================
 * VACANCY MEDIA FIX
 * ===================== */

/* общий контейнер карточки */
.vacancy-card {
  overflow: hidden;
}

/* если у тебя есть картинка/медиа внутри описания */
.vacancy-card img {
  max-width: 100%;
  height: auto;
  display: block;
}

/* если это hero-картинка (большая) */
.vacancy-card .description img,
.vacancy-card .media img {
  max-height: 420px;
  object-fit: cover;
  border-radius: 16px;
  margin: 12px auto 0;
}

/* защита от огромных svg / png */
.vacancy-card img[src$=".svg"],
.vacancy-card img[src$=".png"],
.vacancy-card img[src$=".jpg"],
.vacancy-card img[src$=".webp"] {
  max-height: 420px;
}

/* =====================
 * MOBILE
 * ===================== */
@media (max-width: 768px) {
  .vacancy-card .description img,
  .vacancy-card .media img {
    max-height: 240px;
    border-radius: 14px;
  }
}
.company-avatar {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  object-fit: cover;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .company-avatar {
    width: 52px;
    height: 52px;
    border-radius: 14px;
  }
}

.description {
  max-width: 720px;
  margin-top: 12px;
  line-height: 1.6;
  white-space: pre-line;
}
/* =====================
 * COMPANY CARD — UPGRADE
 * ===================== */

.company-card {
  margin-top: 20px;
  border-radius: 20px;
  padding: 18px;
  background:
    linear-gradient(180deg, #fafafa, #ffffff);
  border: 1px solid #e5e7eb;

  transition: box-shadow .2s ease, transform .2s ease;
}

.company-card:hover {
  box-shadow: 0 10px 30px rgba(0,0,0,.06);
  transform: translateY(-1px);
}

/* header */
.company-head {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* avatar */
.company-avatar {
  width: 64px;
  height: 64px;
  min-width: 64px;
  border-radius: 16px;
  object-fit: cover;
  background: #f3f4f6;
}

/* text block */
.company-name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.company-meta {
  font-size: 13px;
  color: #6b7280;
  margin-top: 2px;
}

/* link */
.company-link {
  display: inline-block;
  margin-top: 6px;
  font-size: 13px;
  color: #2563eb;
  text-decoration: none;
}

.company-link:hover {
  text-decoration: underline;
}

/* =====================
 * MOBILE
 * ===================== */
@media (max-width: 768px) {
  .company-card {
    padding: 14px;
  }

  .company-avatar {
    width: 52px;
    height: 52px;
    border-radius: 14px;
  }

  .company-name {
    font-size: 15px;
  }
}
/* =====================
 * COMPANY — BRAND ACCENT
 * ===================== */

.company-card {
  position: relative;
  margin-top: 20px;
  border-radius: 20px;
  padding: 18px;
  background: linear-gradient(180deg, #fafafa, #ffffff);
  border: 1px solid #e5e7eb;

  /* subtle accent */
  box-shadow:
    inset 4px 0 0 var(--brand),
    0 0 0 transparent;

  transition:
    box-shadow .25s ease,
    transform .2s ease;
}

.company-card:hover {
  transform: translateY(-1px);
  box-shadow:
    inset 6px 0 0 var(--brand),
    0 10px 30px rgba(0,0,0,.06);
}

/* header layout */
.company-head {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* avatar with brand ring */
.company-avatar {
  width: 64px;
  height: 64px;
  min-width: 64px;
  border-radius: 16px;
  object-fit: cover;
  background: #f3f4f6;

  box-shadow: 0 0 0 2px color-mix(in srgb, var(--brand) 30%, #fff);
}

/* name accent */
.company-name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

/* subtle underline on hover */
.company-card:hover .company-name {
  text-decoration: underline;
  text-decoration-color: color-mix(in srgb, var(--brand) 70%, transparent);
  text-underline-offset: 3px;
}

/* meta */
.company-meta {
  font-size: 13px;
  color: #6b7280;
  margin-top: 2px;
}

/* link */
.company-link {
  display: inline-block;
  margin-top: 6px;
  font-size: 13px;
  color: var(--brand);
  text-decoration: none;
}

.company-link:hover {
  text-decoration: underline;
}
@media (max-width: 768px) {
  .company-card {
    padding: 14px;
    box-shadow:
      inset 3px 0 0 var(--brand);
  }

  .company-avatar {
    width: 52px;
    height: 52px;
    border-radius: 14px;
  }

  .company-name {
    font-size: 15px;
  }
}
.company-card[data-verified="true"]::after {
  content: "Verified";
  position: absolute;
  top: 14px;
  right: 14px;
  font-size: 11px;
  padding: 4px 8px;
  border-radius: 999px;
  background: color-mix(in srgb, var(--brand) 15%, #fff);
  color: var(--brand);
}
.sheet-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(17, 24, 39, 0.55);
  backdrop-filter: blur(4px);
  z-index: 50;

  display: flex;
  align-items: center;
  justify-content: center;

  animation: fadeIn .2s ease;
  
  box-sizing: border-box;
  overflow: hidden;
}

@keyframes fadeIn {
  from { opacity: 0 }
  to { opacity: 1 }
}
.sheet {
  width: 100%;
  max-width: 520px;
  background: #ffffff;
  border-radius: 20px;
  padding: 20px;

  box-shadow:
    0 30px 80px rgba(0,0,0,.25);

  animation: scaleIn .2s ease;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(.96) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
.sheet-title {
  display: flex;
  justify-content: space-between;
  align-items: center;

  font-size: 16px;
  font-weight: 600;
  color: #111827;

  margin-bottom: 16px;
}

.sheet-title button {
  border: none;
  background: none;
  font-size: 18px;
  cursor: pointer;
  opacity: .6;
}

.sheet-title button:hover {
  opacity: 1;
}
.sheet label {
  display: block;
  margin-bottom: 14px;
}

.sheet label span {
  display: block;
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 6px;
}

.sheet select,
.sheet textarea {
  width: 100%;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  padding: 10px 12px;
  font-size: 14px;

  transition: border-color .15s ease, box-shadow .15s ease;
}

.sheet select:focus,
.sheet textarea:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99,102,241,.15);
}
.match-preview {
  margin: 10px 0 14px;
  padding: 10px 12px;
  border-radius: 12px;

  background: linear-gradient(
    90deg,
    rgba(99,102,241,.08),
    rgba(99,102,241,.02)
  );

  font-size: 13px;
  color: #1f2937;
}
.sheet .btn.primary.full {
  margin-top: 8px;
  padding: 12px;
  font-size: 15px;
  font-weight: 600;
  border-radius: 14px;

  background: linear-gradient(180deg, #111827, #020617);
  border: none;

  transition: transform .15s ease, box-shadow .15s ease;
}

.sheet .btn.primary.full:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 24px rgba(0,0,0,.25);
}

/* ===== APPLY MODAL — DESKTOP & MOBILE ===== */

.sheet-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(17, 24, 39, 0.55);
  backdrop-filter: blur(4px);
  z-index: 50;

  display: flex;
  align-items: center;
  justify-content: center;

  animation: fadeIn .2s ease;
  
  box-sizing: border-box;
  overflow: hidden;
}

@keyframes fadeIn {
  from { opacity: 0 }
  to { opacity: 1 }
}

.sheet {
  position: relative;
  bottom: auto;
  width: 100%;
  max-width: 520px;

  background: #fff;
  border-radius: 20px;
  padding: 20px;

  box-shadow: 0 30px 80px rgba(0,0,0,.25);
  animation: scaleIn .2s ease;
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(.96) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

@media (max-width: 768px) {
  /* Предотвращаем скролл body когда модалка открыта */
  body:has(.sheet-backdrop) {
    overflow: hidden;
  }

  .sheet-backdrop {
    align-items: flex-end;
    padding: 0 !important;
    margin: 0 !important;
    width: 100%;
    height: 100%;
  }

  .sheet {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;

    max-width: 100% !important;
    width: 100% !important;
    border-radius: 20px 20px 0 0;
    padding: 16px !important;
    margin: 0 !important;
    box-sizing: border-box;
    overflow-x: hidden;

    animation: slideUp .25s ease;
  }

  @keyframes slideUp {
    from { transform: translateY(100%); }
    to { transform: translateY(0); }
  }

  /* drag handle */
  .sheet::before {
    content: "";
    display: block;
    width: 36px;
    height: 4px;
    border-radius: 999px;
    background: #d1d5db;
    margin: 0 auto 10px;
  }
  
  /* Все вложенные элементы используют border-box и не выходят за границы */
  .sheet *,
  .sheet select,
  .sheet textarea,
  .sheet label,
  .sheet button {
    box-sizing: border-box;
    max-width: 100%;
  }
  
  /* Убеждаемся что label и все инпуты тоже не создают overflow */
  .sheet label {
    width: 100%;
  }
  
  .sheet select,
  .sheet textarea,
  .sheet input {
    width: 100% !important;
    min-width: 0 !important;
  }
}

.application-status {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.status-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;

  padding: 8px 14px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 500;

  background: #f3f4f6;
  color: #374151;
}

.status-pill .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

/* статус — pending */
.status-pill[data-status="pending"] {
  background: rgba(245, 158, 11, 0.12);
  color: #b45309;
}

/* accepted */
.status-pill[data-status="accepted"] {
  background: rgba(34, 197, 94, 0.12);
  color: #166534;
}

/* rejected */
.status-pill[data-status="rejected"] {
  background: rgba(239, 68, 68, 0.12);
  color: #991b1b;
}

.chat-link {
  font-size: 14px;
  font-weight: 500;
  color: #2563eb;
  text-decoration: none;
}

.chat-link:hover {
  text-decoration: underline;
}
.analytics-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 12px;
}

.metric {
  display: flex;
  align-items: center;
  gap: 10px;

  padding: 10px 12px;
  border-radius: 14px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
}

.metric img {
  width: 20px;
  height: 20px;
  opacity: 0.7;
}

.metric .value {
  font-size: 15px;
  font-weight: 600;
  color: #111827;
}

.metric .label {
  font-size: 11px;
  color: #6b7280;
}

/* mobile */
@media (max-width: 768px) {
  .analytics-grid {
    grid-template-columns: 1fr;
  }
}

</style>
