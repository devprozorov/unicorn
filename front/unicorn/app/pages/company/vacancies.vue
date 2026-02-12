<template>
  <Header />

  <div class="page p-6 space-y-6">
    <!-- HEADER -->
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold">{{ t('companyVacancies.title') }}</h1>
      <button class="btn primary" @click="openCreate">
        {{ t('companyVacancies.newButton') }}
      </button>
    </div>

    <!-- STATS -->
    <div class="stats">
      <div class="stat">
        <div class="label">{{ t('companyVacancies.totalResponses') }}</div>
        <div class="value">{{ totalResponses }}</div>
      </div>
      <div class="stat">
        <div class="label">{{ t('companyVacancies.newResponses') }}</div>
        <div class="value accent">{{ totalNew }}</div>
      </div>
      <div class="stat">
        <div class="label">{{ t('companyVacancies.responseRate') }}</div>
        <div class="value">{{ responseRate }}%</div>
      </div>
      <div class="stat">
        <div class="label">{{ t('companyVacancies.avgInterest') }}</div>
        <div class="value">{{ avgInterest }}</div>
      </div>
    </div>

    <!-- ALERTS -->
    <div v-if="err" class="alert error">{{ err }}</div>
    <div v-if="hint" class="alert hint">{{ hint }}</div>

    <!-- LIST -->
    <div class="grid gap-4" :class="{ compact: isCompact }">
      <div
        v-for="v in myVacancies"
        :key="v.vacancyId"
        class="swipe-wrap"
        @touchstart.passive="onTouchStart($event, v)"
        @touchmove.passive="onTouchMove($event)"
        @touchend="onTouchEnd()"
        @touchcancel="onTouchCancel()"
      >
        <!-- swipe background actions (mobile only visually) -->
        <div class="swipe-actions swipe-left">
          <div class="swipe-pill">{{ t('companyVacancies.applications') }}</div>
        </div>
        <div class="swipe-actions swipe-right">
          <div class="swipe-pill danger">{{ t('companyVacancies.delete') }}</div>
        </div>

        <div
          class="vacancy-card"
          :class="{ 'swiping': isSwipingThis(v.vacancyId) }"
          :style="swipeStyle(v.vacancyId)"
          @click="onCardTap(v)"
        >
          <div class="vacancy-body">
            <!-- INFO -->
            <div class="info">
              <div class="title-row">
                <div class="title">{{ v.title }}</div>

                <!-- HEALTH -->
                <img
                  class="health"
                  :src="healthIcon(v.applicationsCount)"
                  alt="health"
                  :title="healthTitle(v.applicationsCount)"
                  @click.stop="openHealthExplain(v.applicationsCount)"
                />
              </div>

              <div class="meta">
                {{ v.location || t('companyVacancies.remote') }}
              </div>

              <div class="description">
                {{ v.description }}
              </div>

              <div class="tags">
                {{ formatTags(v.tags) }}
              </div>
            </div>

            <!-- SIDE -->
            <div class="side">
              <!-- BADGES -->
              <div class="badges">
                <span class="badge neutral">
                  {{ v.applicationsCount }} {{ t('companyVacancies.total') }}
                </span>

                <span
                  v-if="v.newApplicationsCount"
                  class="badge accent"
                >
                  {{ v.newApplicationsCount }} {{ t('companyVacancies.new') }}
                </span>
              </div>

              <!-- ACTIONS -->
              <div class="actions">
                <button class="btn" @click.stop="openEdit(v)">
                  {{ t('companyVacancies.edit') }}
                </button>

                <button
                  v-if="v.newApplicationsCount"
                  class="btn primary"
                  @click.stop="openFirstChat(v)"
                >
                  {{ t('companyVacancies.chat') }}
                </button>

                <NuxtLink
                  v-else
                  class="btn"
                  :to="`/company/applications?vacancyId=${v.vacancyId}`"
                >
                  {{ t('companyVacancies.applications') }}
                </NuxtLink>

                <button
                  class="btn danger"
                  @click.stop="remove(v)"
                >
                  {{ t('companyVacancies.deleteBtn') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL (REDESIGNED, LOGIC PRESERVED) -->
<div v-if="modal.open" class="overlay" :class="{ mobile: isMobile }">
  <div class="modal modern" :class="{ sheet: isMobile }">

    <!-- HEADER -->
    <div class="modal-header">
      <div class="modal-title">
        {{ modal.mode === 'create' ? t('companyVacancies.createTitle') : t('companyVacancies.editTitle') }}
      </div>
      <button class="close" @click="closeModal">x</button>
    </div>

    <!-- BODY -->
    <div class="modal-body">
      <div class="field">
        <label class="label">{{ t('companyVacancies.modal.titleLabel') }}</label>
        <input
          v-model="modal.form.title"
          :placeholder="t('companyVacancies.modal.titlePlaceholder')"
        />
      </div>

      <div class="field">
        <label class="label">{{ t('companyVacancies.modal.locationLabel') }}</label>
        <div class="location-row">
          <select v-model="modal.form.region">
            <option value="">{{ t('companyVacancies.modal.selectRegion') }}</option>
            <option
              v-for="r in regionOptions"
              :key="r.key"
              :value="r.key"
            >
              {{ r.label }}
            </option>
            <option value="custom">{{ t('companyVacancies.modal.other') }}</option>
          </select>

          <select
            v-model="modal.form.country"
            :disabled="!modal.form.region || modal.form.region === 'custom'"
          >
            <option value="">{{ t('companyVacancies.modal.selectCountry') }}</option>
            <option
              v-for="c in countryOptions"
              :key="c"
              :value="c"
            >
              {{ c }}
            </option>
          </select>
        </div>

        <input
          v-if="modal.form.region === 'custom'"
          v-model="modal.form.customLocation"
          class="mt-2"
          :placeholder="t('companyVacancies.modal.cityPlaceholder')"
        />
      </div>

      <div class="field">
        <label class="label">{{ t('companyVacancies.modal.tagLabel') }}</label>
        
        <!-- Hierarchical Tag Selector -->
        <HierarchicalTagSelector
          :selected-tags="modal.tagsSelected"
          @toggle-tag="toggleTag"
          @clear-all="clearAllTags"
          :enable-search="true"
        />
      </div>

      <div class="field">
        <label class="label">{{ t('companyVacancies.modal.workFormatLabel') }}</label>
        <div class="tag-options">
          <button
            class="tag-chip"
            type="button"
            :class="{ active: modal.tagsSelected.includes('Remote') }"
            @click="toggleTag('Remote')"
          >
            {{ t('tags.options.remote') }}
          </button>
          <button
            class="tag-chip"
            type="button"
            :class="{ active: modal.tagsSelected.includes('On-site') }"
            @click="toggleTag('On-site')"
          >
            {{ t('tags.options.onSite') }}
          </button>
        </div>
      </div>

      <div class="field">
        <label class="label">{{ t('companyVacancies.modal.descriptionLabel') }}</label>
        <textarea
            v-model="modal.form.description"
            rows="6"
            :placeholder="t('companyVacancies.modal.descriptionPlaceholder')"
        />
      </div>
    </div>

    <!-- FOOTER -->
    <div class="modal-actions">
      <button class="btn ghost" @click="closeModal">
        {{ t('companyVacancies.modal.cancel') }}
      </button>
      <button
        class="btn primary"
        @click="saveVacancy"
        :disabled="loading"
      >
        {{ loading ? t('companyVacancies.saving') : t('companyVacancies.modal.saveChanges') }}
      </button>
    </div>

  </div>
</div>



    <!-- HEALTH EXPLAIN (mobile bottom sheet) -->
    <div v-if="healthExplain.open" class="overlay mobile" @click.self="closeHealthExplain">
      <div class="modal sheet">
        <div class="modal-header">
          <div class="font-semibold">{{ t('companyVacancies.healthExplain.title') }}</div>
          <button class="close" @click="closeHealthExplain">x</button>
        </div>

        <div class="health-sheet">
          <div class="health-row">
            <img class="health-big" :src="nothing" alt="none" />
            <div>
              <div class="health-title">{{ t('companyVacancies.health.noResponses') }}</div>
              <div class="health-text">{{ t('companyVacancies.healthExplain.noResponsesText') }}</div>
            </div>
          </div>

          <div class="health-row">
            <img class="health-big" :src="lowInterest" alt="low" />
            <div>
              <div class="health-title">{{ t('companyVacancies.health.lowInterest') }}</div>
              <div class="health-text">{{ t('companyVacancies.healthExplain.lowInterestText') }}</div>
            </div>
          </div>

          <div class="health-row">
            <img class="health-big" :src="highInterest" alt="high" />
            <div>
              <div class="health-title">{{ t('companyVacancies.health.highInterest') }}</div>
              <div class="health-text">{{ t('companyVacancies.healthExplain.highInterestText') }}</div>
            </div>
          </div>

          <div class="health-current">
            {{ t('companyVacancies.healthExplain.current') }} <b>{{ healthTitle(healthExplain.count) }}</b>
          </div>
        </div>

        <div class="modal-actions">
          <button class="btn primary" @click="closeHealthExplain">
            {{ t('companyVacancies.healthExplain.ok') }}
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- MOBILE STICKY ACTIONS -->
  <div v-if="isMobile && selectedVacancy" class="sticky-actions">
    <button class="btn" @click="openEdit(selectedVacancy)">
      {{ t('companyVacancies.edit') }}
    </button>

    <NuxtLink
      class="btn"
      :to="`/company/applications?vacancyId=${selectedVacancy.vacancyId}`"
    >
      {{ t('companyVacancies.applications') }}
    </NuxtLink>

    <button
      v-if="selectedVacancy.newApplicationsCount"
      class="btn primary"
      @click="openFirstChat(selectedVacancy)"
    >
      {{ t('companyVacancies.chat') }}
    </button>

    <button class="btn danger" @click="remove(selectedVacancy)">
      {{ t('companyVacancies.deleteBtn') }}
    </button>
  </div>

  <Footer />
</template>

<script setup lang="ts">
import axios from 'axios'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import HierarchicalTagSelector from '~/components/tags/HierarchicalTagSelector.vue'
import nothing from '~/images/nothing.png'
import lowInterest from '~/images/low-interest.png'
import highInterest from '~/images/high-interest.png'
import { locationRegions, findRegionByCountry } from '~/utils/locations'
import { tagsGroups, allTags, tagLabelMap } from '~/utils/tags'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'

definePageMeta({ middleware: ['company-auth'] })

const auth = useAuthStore()
const { t } = useI18n()
const err = ref('')
const hint = ref('')
const loading = ref(false)

type Vacancy = {
  vacancyId: string
  companyId: string
  title: string
  description: string
  location?: string
  tags?: string[]
  applicationsCount: number
  newApplicationsCount: number
}

const all = ref<Vacancy[]>([])
const myVacancies = computed(() =>
  all.value.filter(v => v.companyId === auth.userId)
)

/* STATS */
const totalResponses = computed(() =>
  myVacancies.value.reduce((s, v) => s + v.applicationsCount, 0)
)

const totalNew = computed(() =>
  myVacancies.value.reduce((s, v) => s + v.newApplicationsCount, 0)
)

const responseRate = computed(() => {
  if (!totalResponses.value) return 0
  return Math.round((totalNew.value / totalResponses.value) * 100)
})

const avgInterest = computed(() => {
  if (!myVacancies.value.length) return t('profileHeader.notSpecified')
  const avg = totalResponses.value / myVacancies.value.length
  return avg.toFixed(1)
})

/* COMPACT MODE (CSS only) */
const isCompact = computed(() => myVacancies.value.length > 10)

/* MOBILE DETECT (reactive) */
const isMobile = ref(false)
function updateMobile() {
  if (process.server) return
  isMobile.value = window.innerWidth <= 768
}

onMounted(() => {
  updateMobile()
  window.addEventListener('resize', updateMobile)
})
onBeforeUnmount(() => {
  if (process.server) return
  window.removeEventListener('resize', updateMobile)
})

/* HEALTH ICON */
function healthIcon(count: number) {
  if (count === 0) return nothing
  if (count <= 10) return lowInterest
  return highInterest
}

function healthTitle(count: number) {
  if (count === 0) return t('companyVacancies.health.noResponses')
  if (count <= 10) return t('companyVacancies.health.lowInterest')
  return t('companyVacancies.health.highInterest')
}

function tagLabel(value: string) {
  const key = tagLabelMap[value]
  return key ? t(key) : value
}

function formatTags(tags?: string[]) {
  if (!tags || !tags.length) return t('profileHeader.notSpecified')
  return tags.map(tagLabel).join(', ')
}

/* HEALTH EXPLAIN (mobile-only sheet) */
const healthExplain = reactive({
  open: false,
  count: 0,
})

function openHealthExplain(count: number) {
  if (!isMobile.value) return
  healthExplain.open = true
  healthExplain.count = count
}

function closeHealthExplain() {
  healthExplain.open = false
}

/* SELECTED VACANCY for sticky actions */
const selectedVacancy = ref<Vacancy | null>(null)
function onCardTap(v: Vacancy) {
  if (!isMobile.value) return
  // if user is swiping, don’t select
  if (swipe.active && Math.abs(swipe.dx) > 10) return
  selectedVacancy.value = v
}

/* SWIPE ACTIONS (mobile-only) */
const swipe = reactive({
  active: false,
  id: '' as string,
  startX: 0,
  startY: 0,
  dx: 0,
  dy: 0,
  locked: '' as '' | 'h' | 'v',
})

const swipeDxMap = reactive<Record<string, number>>({})

function isSwipingThis(id: string) {
  return swipe.active && swipe.id === id && Math.abs(swipe.dx) > 2
}

function swipeStyle(id: string) {
  const dx = swipeDxMap[id] || 0
  return {
    transform: `translateX(${dx}px)`,
  }
}

function onTouchStart(e: TouchEvent, v: Vacancy) {
  if (!isMobile.value) return
  const t = e.touches[0]
  swipe.active = true
  swipe.id = v.vacancyId
  swipe.startX = t.clientX
  swipe.startY = t.clientY
  swipe.dx = 0
  swipe.dy = 0
  swipe.locked = ''
  swipeDxMap[v.vacancyId] = 0
}

function onTouchMove(e: TouchEvent) {
  if (!isMobile.value) return
  if (!swipe.active) return
  const t = e.touches[0]
  swipe.dx = t.clientX - swipe.startX
  swipe.dy = t.clientY - swipe.startY

  // lock direction
  if (!swipe.locked) {
    const ax = Math.abs(swipe.dx)
    const ay = Math.abs(swipe.dy)
    if (ax < 6 && ay < 6) return
    swipe.locked = ax > ay * 1.2 ? 'h' : 'v'
  }

  // if vertical scroll -> ignore swipe
  if (swipe.locked === 'v') {
    swipeDxMap[swipe.id] = 0
    return
  }

  // horizontal swipe clamp
  const clamp = (n: number, min: number, max: number) => Math.max(min, Math.min(max, n))
  swipeDxMap[swipe.id] = clamp(swipe.dx, -140, 140)
}

function onTouchEnd() {
  if (!isMobile.value) return
  if (!swipe.active) return

  const id = swipe.id
  const dx = swipeDxMap[id] || 0

  // threshold actions
  if (dx > 90) {
    // right swipe -> Applications
    const v = myVacancies.value.find(x => x.vacancyId === id)
    if (v) navigateTo(`/company/applications?vacancyId=${v.vacancyId}`)
  } else if (dx < -90) {
    // left swipe -> Delete
    const v = myVacancies.value.find(x => x.vacancyId === id)
    if (v) remove(v)
  }

  // reset
  swipeDxMap[id] = 0
  swipe.active = false
  swipe.id = ''
  swipe.locked = ''
}

function onTouchCancel() {
  if (!isMobile.value) return
  if (!swipe.active) return
  if (swipe.id) swipeDxMap[swipe.id] = 0
  swipe.active = false
  swipe.id = ''
  swipe.locked = ''
}

/* MODAL + LOAD (НЕ МЕНЯЛ ЛОГИКУ) */
const modal = reactive({
  open: false,
  mode: 'create' as 'create' | 'edit',
  editingId: '',
  form: {
    title: '',
    description: '',
    region: '',
    country: '',
    customLocation: '',
  },
  tagsSelected: [] as string[],
})

const regionOptions = locationRegions
const countryOptions = computed(() => {
  const r = locationRegions.find(x => x.key === modal.form.region)
  return r?.countries || []
})

const tagToAdd = ref('')

const allTagOptions = computed(() =>
  tagsGroups.flatMap(g =>
    g.options.map(o => ({ group: g.label, value: o }))
  )
)

const tagPickers = reactive<Record<string, string>>(
  Object.fromEntries(tagsGroups.map(g => [g.key, '']))
)

watch(
  () => modal.form.region,
  val => {
    modal.form.country = ''
    if (val !== 'custom') {
      modal.form.customLocation = ''
    }
  }
)

function hydrateLocationFields(loc: string) {
  const reg = findRegionByCountry(loc || '')
  if (reg) {
    modal.form.region = reg.key
    modal.form.country =
      reg.countries.find(c => c.toLowerCase() === (loc || '').toLowerCase()) || ''
    modal.form.customLocation = ''
  } else {
    modal.form.region = ''
    modal.form.country = ''
    modal.form.customLocation = loc || ''
  }
}

function hydrateTags(tags: string[] | undefined) {
  modal.tagsSelected = (tags || []).filter(t => allTags.includes(t))
}

function openCreate() {
  modal.open = true
  modal.mode = 'create'
  modal.editingId = ''
  modal.form = {
    title: '',
    description: '',
    region: '',
    country: '',
    customLocation: '',
  }
  modal.tagsSelected = []
  Object.keys(tagPickers).forEach(k => { tagPickers[k] = '' })
}

function openEdit(v: Vacancy) {
  modal.open = true
  modal.mode = 'edit'
  modal.editingId = v.vacancyId
  modal.form.title = v.title
  modal.form.description = v.description
  hydrateLocationFields(v.location || '')
  hydrateTags(v.tags || [])
  Object.keys(tagPickers).forEach(k => { tagPickers[k] = '' })
}

function closeModal() {
  modal.open = false
}

async function load() {
  err.value = ''
  try {
    if (!auth.userId) {
      const me = await axios.get('/auth/me')
      auth.userId = me.data.userId
    }

    const [vacRes, appRes] = await Promise.all([
      axios.get('/vacancies'),
      axios.get('/applications/inbox'),
    ])

    const stats: Record<string, { total: number; new: number }> = {}
    for (const a of appRes.data.items || []) {
      if (!stats[a.vacancyId]) stats[a.vacancyId] = { total: 0, new: 0 }
      stats[a.vacancyId].total++
      if (a.status === 'pending' && !a.viewed) stats[a.vacancyId].new++
    }

    all.value = (vacRes.data.items || []).map((v: any) => ({
      ...v,
      applicationsCount: stats[v.vacancyId]?.total || 0,
      newApplicationsCount: stats[v.vacancyId]?.new || 0,
    }))
  } catch (e: any) {
    err.value = getErrorMessage(e?.response?.data?.error || 'load_failed', t, 'errors.loadFailed')
  }
}

function openFirstChat(v: Vacancy) {
  navigateTo(`/company/applications?vacancyId=${v.vacancyId}&autochat=1`)
}

function toggleTag(tag: string) {
  if (modal.tagsSelected.includes(tag)) {
    modal.tagsSelected = modal.tagsSelected.filter(t => t !== tag)
  } else {
    modal.tagsSelected = [...modal.tagsSelected, tag]
  }
}

function clearAllTags() {
  modal.tagsSelected = []
}

function addTag() {
  const val = tagToAdd.value
  if (!val) return
  if (!modal.tagsSelected.includes(val) && allTags.includes(val)) {
    modal.tagsSelected = [...modal.tagsSelected, val]
  }
  tagToAdd.value = ''
}

function addTagFromGroup(key: string) {
  const val = tagPickers[key]
  if (!val) return
  if (!modal.tagsSelected.includes(val)) {
    modal.tagsSelected = [...modal.tagsSelected, val]
  }
  tagPickers[key] = ''
}

async function saveVacancy() {
  loading.value = true
  const locationValue = modal.form.country || modal.form.customLocation || ''
  if (modal.form.region && !modal.form.country && !modal.form.customLocation) {
    hint.value = t('companyVacancies.locationHint')
    loading.value = false
    return
  }

  const payload = {
    title: modal.form.title,
    description: modal.form.description,
    location: locationValue,
    tags: modal.tagsSelected,
  }

  try {
    if (modal.mode === 'create')
      await axios.post('/vacancies', payload)
    else
      await axios.patch(`/vacancies/${modal.editingId}`, payload)

    closeModal()
    await load()
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
    await axios.delete(`/vacancies/${v.vacancyId}`)
    await load()
  } catch (e: any) {
    hint.value = getErrorMessage(e?.response?.data?.error || 'delete_failed', t, 'errors.deleteFailed')
  }
}

onMounted(load)
</script>

<style scoped>
.page {
  max-width: 1100px;
  margin: 0 auto;
}

/* =====================
 * STATS
 * ===================== */
.stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.stat {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 16px;
  padding: 12px;
}

.stat .label {
  font-size: 12px;
  opacity: .6;
}

.stat .value {
  font-size: 18px;
  font-weight: 600;
}

.stat .accent {
  color: #16a34a;
}

/* =====================
 * VACANCY CARD
 * ===================== */
.vacancy-card {
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  background: #fff;
}

.vacancy-body {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 16px;
}

/* =====================
 * INFO
 * ===================== */
.info {
  flex: 1;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title {
  font-weight: 600;
  font-size: 16px;
}

.health {
  width: 22px;
  height: 22px;
  flex-shrink: 0;
}

.meta {
  font-size: 13px;
  color: #6b7280;
  margin-top: 2px;
}

.description {
  margin-top: 8px;
  font-size: 14px;
  white-space: pre-line;
  line-height: 1.6;
}

.tags {
  margin-top: 6px;
  font-size: 12px;
  color: #9ca3af;
}

/* =====================
 * SIDE
 * ===================== */
.side {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 140px;
}

.badges {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.badge {
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 999px;
  border: 1px solid #e5e7eb;
}

.badge.accent {
  background: rgba(34,197,94,.15);
  color: #166534;
}

/* =====================
 * ACTIONS
 * ===================== */
.actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.btn {
  padding: 8px 12px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  background: #fff;
  text-align: center;
}

.btn.primary {
  background: #111827;
  color: #fff;
}

.btn.danger {
  color: #b91c1c;
  border-color: #fecaca;
}

/* =====================
 * MODAL
 * ===================== */
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.4);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

.modal {
  background: #fff;
  border-radius: 18px;
  padding: 16px;
  max-width: 520px;
  width: 100%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

/* =====================
 * MOBILE
 * ===================== */
@media (max-width: 768px) {
  /* PAGE */
  .page {
    padding-left: 12px;
    padding-right: 12px;
  }

  /* HEADER */
  .flex.items-center.justify-between {
    flex-direction: column;
    align-items: stretch;
    gap: 10px;
  }

  .btn.primary {
    width: 100%;
  }

  /* STATS */
  .stats {
    grid-template-columns: 1fr 1fr;
  }

  .stat {
    padding: 10px;
  }

  .stat .value {
    font-size: 16px;
  }

  /* CARD */
  .vacancy-body {
    flex-direction: column;
  }

  .health {
    width: 26px;
    height: 26px;
  }

  /* BADGES */
  .badges {
    margin-top: 6px;
  }

  /* ACTIONS */
  .actions {
    flex-direction: row;
    flex-wrap: wrap;
    gap: 8px;
  }

  .actions .btn {
    flex: 1 1 calc(50% - 4px);
  }

  /* DESCRIPTION */
  .description {
    font-size: 14px;
    line-height: 1.45;
    white-space: pre-line;
  }
}

/* =====================
 * SMALL PHONES
 * ===================== */
@media (max-width: 420px) {
  .stats {
    grid-template-columns: 1fr;
  }

  .actions .btn {
    flex: 1 1 100%;
  }

  .title {
    font-size: 15px;
  }
}

/* =========================================================
 * ДОБАВЛЕНО: MOBILE-ONLY UX (НИЧЕГО ВЫШЕ НЕ УДАЛЕНО)
 * ========================================================= */

/* swipe wrapper */
.swipe-wrap {
  position: relative;
  border-radius: 18px;
}

/* swipe background actions */
.swipe-actions {
  position: absolute;
  inset: 0;
  border-radius: 18px;
  display: none; /* show only on mobile */
  align-items: center;
  padding: 0 14px;
  pointer-events: none;
}

.swipe-left {
  justify-content: flex-start;
  background: rgba(17,24,39,.06);
  border: 1px solid #e5e7eb;
}

.swipe-right {
  justify-content: flex-end;
  background: rgba(185,28,28,.06);
  border: 1px solid #fecaca;
}

.swipe-pill {
  font-size: 12px;
  padding: 8px 10px;
  border-radius: 999px;
  border: 1px solid #e5e7eb;
  background: #fff;
}

.swipe-pill.danger {
  border-color: #fecaca;
  color: #b91c1c;
}

/* animated card while swiping */
.vacancy-card {
  will-change: transform;
  transition: transform .18s ease;
}
.vacancy-card.swiping {
  transition: none;
}

@media (max-width: 768px) {
  .swipe-actions {
    display: flex;
  }
}

/* sticky actions bar */
.sticky-actions {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 60;
  display: flex;
  gap: 8px;
  padding: 12px;
  background: rgba(255,255,255,.96);
  backdrop-filter: blur(10px);
  border-top: 1px solid #e5e7eb;
}

.sticky-actions .btn {
  flex: 1;
}

/* add bottom padding so sticky doesn't cover content */
@media (max-width: 768px) {
  .page {
    padding-bottom: 92px;
  }
}

/* bottom sheet modal on mobile */
.overlay.mobile {
  align-items: flex-end;
  justify-content: center;
}

.modal.sheet {
  max-width: 720px;
  width: 100%;
  border-radius: 18px 18px 0 0;
  padding-bottom: 18px;
  animation: sheetUp .18s ease;
}

@keyframes sheetUp {
  from { transform: translateY(18px); opacity: .9; }
  to { transform: translateY(0); opacity: 1; }
}

/* health explain sheet content */
.health-sheet {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 10px 2px 2px;
}

.health-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  border: 1px solid #e5e7eb;
  border-radius: 14px;
  padding: 12px;
  background: #fff;
}

.health-big {
  width: 34px;
  height: 34px;
  flex-shrink: 0;
}

.health-title {
  font-weight: 600;
  font-size: 14px;
}

.health-text {
  margin-top: 3px;
  font-size: 13px;
  color: #6b7280;
  line-height: 1.35;
}

.health-current {
  font-size: 13px;
  color: #111827;
  opacity: .9;
}

/* compact cards when >10 */
.grid.compact .vacancy-body {
  padding: 12px;
}
.grid.compact .description {
  font-size: 13px;
  white-space: pre-line;
}
.grid.compact .title {
  font-size: 15px;
}
.grid.compact .actions .btn {
  padding: 7px 10px;
  font-size: 13px;
}
/* ============================
 * MODAL OVERLAY
 * ============================ */
.overlay {
  position: fixed;
  inset: 0;
  z-index: 100;
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}

/* mobile bottom-sheet alignment */
.overlay.mobile {
  align-items: flex-end;
}

/* ============================
 * MODAL CONTAINER
 * ============================ */
.modal.modern {
  width: 100%;
  max-width: 560px;
  background: #ffffff;
  border-radius: 20px;
  box-shadow:
    0 30px 80px rgba(0,0,0,.25),
    0 4px 14px rgba(0,0,0,.12);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: modalIn .18s ease-out;
}

/* bottom-sheet variant */
.modal.modern.sheet {
  max-width: 100%;
  border-radius: 22px 22px 0 0;
  animation: sheetUp .22s ease-out;
}

/* ============================
 * HEADER
 * ============================ */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 18px;
  border-bottom: 1px solid #e5e7eb;
  background: linear-gradient(
    to bottom,
    #ffffff,
    #f9fafb
  );
}

.modal-title {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.close {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 18px;
  opacity: .6;
  transition: opacity .15s ease;
}
.close:hover {
  opacity: 1;
}

/* ============================
 * BODY
 * ============================ */
.modal-body {
  padding: 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  overflow-y: auto;
  overflow-x: hidden;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.label {
  font-size: 12px;
  font-weight: 500;
  color: #374151;
}

.hint-inline {
  font-size: 11px;
  color: #9ca3af;
  margin-left: 6px;
}

/* inputs */
.modal input,
.modal textarea {
  width: 100%;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  padding: 10px 12px;
  font-size: 14px;
  background: #ffffff;
  transition: border-color .15s ease, box-shadow .15s ease;
}

.modal textarea {
  resize: vertical;
  min-height: 120px;
}

.location-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.location-row select {
  width: 100%;
}

.tag-groups {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tag-group-title {
  font-weight: 600;
  font-size: 13px;
  margin-bottom: 6px;
}

.tag-select-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 8px;
}
.tag-select-row select {
  width: 100%;
  min-width: 0;
}

.tag-options {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-chip {
  padding: 8px 10px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  background: #fff;
  font-size: 13px;
  cursor: pointer;
  transition: all .15s ease;
}

.tag-chip.active {
  background: #111827;
  color: #fff;
  border-color: #111827;
}

.add-tag-btn {
  margin-top: 6px;
}

.selected-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 8px;
}

.remove-tag {
  background: none;
  border: none;
  color: inherit;
  margin-left: 6px;
  cursor: pointer;
  font-weight: 700;
}

.modal input::placeholder,
.modal textarea::placeholder {
  color: #9ca3af;
}

.modal input:focus,
.modal textarea:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99,102,241,.15);
}

/* ============================
 * FOOTER / ACTIONS
 * ============================ */
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 18px;
  border-top: 1px solid #e5e7eb;
  background: #fafafa;
}

/* buttons */
.modal .btn {
  padding: 10px 16px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  background: #ffffff;
  transition:
    background .15s ease,
    border-color .15s ease,
    transform .05s ease;
}

.modal .btn:hover {
  background: #f9fafb;
}

.modal .btn:active {
  transform: translateY(1px);
}

.modal .btn.primary {
  background: #111827;
  color: #ffffff;
  border-color: #111827;
}

.modal .btn.primary:hover {
  background: #030712;
}

.modal .btn.primary:disabled {
  opacity: .6;
  cursor: default;
}

.modal .btn.ghost {
  background: transparent;
  color: #374151;
}

/* ============================
 * ANIMATIONS
 * ============================ */
@keyframes modalIn {
  from {
    opacity: 0;
    transform: scale(.97) translateY(6px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

@keyframes sheetUp {
  from {
    transform: translateY(18px);
    opacity: .9;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

/* ============================
 * MOBILE TUNING
 * ============================ */
@media (max-width: 768px) {
  .modal-body {
    padding: 16px;
  }

  .location-row {
    grid-template-columns: 1fr;
  }

  .tag-options {
    gap: 6px;
  }

  .tag-select-row {
    grid-template-columns: 1fr;
  }

  .add-tag-btn {
    width: 100%;
  }

  .modal-actions {
    flex-direction: column;
    gap: 8px;
  }

  .modal-actions .btn {
    width: 100%;
  }
}

</style>




