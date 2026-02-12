<template>
  <Header />

  <div class="jobs-page">
    <h1 class="page-title">{{ t('jobs.title') }}</h1>

    <!-- Desktop filters live in the sidebar to reduce clutter -->

    <!-- =====================
     | MOBILE TOP BAR
     ===================== -->
    <div class="mobile-bar">
      <button class="mobile-btn" @click="showFilters = true">
        {{ t('jobs.filters') }}
      </button>

      <button
        v-if="savedIds.length"
        class="mobile-btn"
        @click="showSaved = !showSaved"
      >
        {{ savedLabel }}
      </button>
    </div>

   <div class="layout">
  <!-- =====================
   | FILTERS SIDEBAR (DESKTOP)
   ===================== -->
  <aside class="filters-panel">
    <div class="panel-card">
      <div class="panel-header">
        <div class="panel-title">{{ t('jobs.filtersTitle') }}</div>
        <div class="panel-count">{{ jobCountLabel }}</div>
      </div>

      <label class="panel-field">
        <span>{{ t('jobs.keywords') }}</span>
        <div class="input-wrap">
          <input
            v-model="q"
            type="text"
            :placeholder="t('jobs.titleDescriptionPlaceholder')"
          />
          <button
            v-if="q"
            class="clear-btn"
            type="button"
            @click="q = ''"
            aria-label="Clear"
          >
            x
          </button>
        </div>
      </label>

      <label class="panel-field">
        <span>{{ t('jobs.location') }}</span>
        <div class="location-row">
          <select v-model="locRegion">
            <option value="">{{ t('jobs.anyRegion') }}</option>
            <option
              v-for="r in locRegionOptions"
              :key="r.key"
              :value="r.key"
            >
              {{ r.label }}
            </option>
          </select>

          <div class="country-select">
            <select
              v-model="locCountry"
              :disabled="!locRegion"
            >
              <option value="">{{ t('jobs.anyCountry') }}</option>
              <option
                v-for="c in locCountryOptions"
                :key="c"
                :value="c"
              >
                {{ c }}
              </option>
            </select>
            <span v-if="!locRegion" class="field-hint">{{ t('jobs.pickRegionHint') }}</span>
          </div>
        </div>
      </label>

      <div class="panel-section">
        <div class="panel-label-row">
          <span class="panel-label">{{ t('jobs.featuredCategories') }}</span>
          <button
            v-if="tagSelection.length"
            class="panel-reset"
            type="button"
            @click="tagSelection = []"
          >
            {{ t('jobs.clearAll') }}
          </button>
        </div>
        <div class="featured-row">
          <button
            class="featured-chip"
            :class="{ active: tagSelection.includes('IT') }"
            type="button"
            @click="toggleFeaturedCategory('IT')"
          >
            {{ t('jobs.featuredIt') }}
          </button>
          <button
            class="featured-chip"
            :class="{ active: tagSelection.includes('Startup') }"
            type="button"
            @click="toggleFeaturedCategory('Startup')"
          >
            {{ t('jobs.featuredStartup') }}
          </button>
          <button
            class="featured-chip"
            :class="{ active: tagSelection.includes('GameDev') }"
            type="button"
            @click="toggleFeaturedCategory('GameDev')"
          >
            {{ t('jobs.featuredGamedev') }}
          </button>
          <button
            class="featured-chip"
            :class="{ active: tagSelection.includes('Other') }"
            type="button"
            @click="toggleFeaturedCategory('Other')"
          >
            {{ t('jobs.featuredOther') }}
          </button>
        </div>
      </div>

      <div class="panel-section tags-panel">
        <span class="panel-label">{{ t('jobs.tags') }}</span>
        <HierarchicalTagSelector
          :key="lang"
          :selected-tags="tagSelection"
          :selected-category="selectedTagCategory"
          enable-search
          @toggle-tag="toggleFilterTag"
          @select-category="selectedTagCategory = $event"
          @clear-all="tagSelection = []"
        />
      </div>

      <label class="panel-field">
        <span>{{ t('jobs.employment') }}</span>
        <select v-model="employment">
          <option value="any">{{ t('jobs.any') }}</option>
          <option value="remote">{{ t('jobs.remote') }}</option>
          <option value="office">{{ t('jobs.office') }}</option>
        </select>
      </label>

      <label class="panel-field">
        <span>{{ t('jobs.posted') }}</span>
        <select v-model="posted">
          <option value="any">{{ t('jobs.anyTime') }}</option>
          <option value="7">{{ t('jobs.last7days') }}</option>
          <option value="30">{{ t('jobs.last30days') }}</option>
        </select>
      </label>

      <label class="panel-field">
        <span>{{ t('jobs.sortBy') }}</span>
        <select v-model="sort">
          <option value="relevance">{{ t('jobs.relevance') }}</option>
          <option value="match">{{ t('jobs.match') }}</option>
          <option value="title">{{ t('jobs.titleSort') }}</option>
        </select>
      </label>
    </div>
  </aside>
  <!-- =====================
   | JOB LIST
   ===================== -->
  <div class="list">
    <div v-if="err" class="alert error">
      {{ err }}
    </div>

    <transition name="fade-list" mode="out-in">
      <div :key="filtersKey">
        <NuxtLink
          v-for="v in filteredVisible"
          :key="v.vacancyId"
          :to="`/jobs/${v.vacancyId}`"
          class="job-card link-card"
          :class="{ 'premium': v.isPremium }"
        >
          <div class="job-card-inner">
            <!-- LEFT -->
            <div class="job-main">
              <h2 class="job-title">
                {{ v.title }}
              </h2>

              <div class="job-meta">
                <span class="company">
                  {{ companyName(v.companyId) }}
                </span>
                <span class="dot">&middot;</span>
                <span class="location">
                  {{ v.location || t('jobs.remote') }}
                </span>
              </div>

              <p class="job-desc">
                {{ v.description }}
              </p>

              <div class="job-tags">
                <span
                  v-for="t in (v.tags || [])"
                  :key="t"
                  class="tag"
                >
                  {{ tagLabel(t) }}
                </span>
              </div>
            </div>

            <!-- RIGHT -->
            <div class="job-side">
              <div class="company-avatar">
                <img
                  :src="companyAvatar(v.companyId)"
                  alt=""
                />
              </div>

              <button
                class="apply-cta"
                @click.stop.prevent="openApply(v)"
              >
                {{ t('jobs.apply') }}
              </button>
            </div>
          </div>
        </NuxtLink>
      </div>
    </transition>
  </div>
</div>
</div>
  <!-- =====================
 | MOBILE FILTER SHEET
 ===================== -->
<div
  v-if="showFilters"
  class="sheet-backdrop"
  @click.self="showFilters = false"
>
  <div class="filters-sheet">
    <!-- HEADER -->
    <div class="filters-header">
      <span class="filters-title">{{ t('jobs.filtersTitle') }}</span>
      <button
        class="close-x"
        :aria-label="t('jobs.filtersTitle')"
        @click="showFilters = false"
      >
        x
      </button>
    </div>

    <!-- BODY -->
    <div class="filters-body">
      <!-- Keywords -->
      <label class="filter-field">
        <span>{{ t('jobs.keywords') }}</span>
        <div class="input-wrap">
          <input
            v-model="q"
            type="text"
            :placeholder="t('jobs.titleDescriptionPlaceholder')"
          />
          <button
            v-if="q"
            class="clear-btn"
            type="button"
            @click="q = ''"
            aria-label="Clear"
          >
            x
          </button>
        </div>
      </label>

      <!-- Location -->
      <label class="filter-field">
        <span>{{ t('jobs.location') }}</span>
        <div class="location-row">
          <select v-model="locRegion">
            <option value="">{{ t('jobs.anyRegion') }}</option>
            <option
              v-for="r in locRegionOptions"
              :key="r.key"
              :value="r.key"
            >
              {{ r.label }}
            </option>
          </select>

          <div class="country-select">
            <select
              v-model="locCountry"
              :disabled="!locRegion"
            >
              <option value="">{{ t('jobs.anyCountry') }}</option>
              <option
                v-for="c in locCountryOptions"
                :key="c"
                :value="c"
              >
                {{ c }}
              </option>
            </select>
            <span v-if="!locRegion" class="field-hint">{{ t('jobs.pickRegionHint') }}</span>
          </div>
        </div>
      </label>

      <!-- Quick tags -->
      <div class="filter-group quick-group">
        <span class="filter-label">{{ t('jobs.quickFilters') }}</span>
        <div class="chip-row">
          <button
            class="chip"
            :class="{ active: tagSelection.includes('IT') }"
            @click="toggleFeaturedCategory('IT')"
            type="button"
          >
            {{ t('jobs.featuredIt') }}
          </button>
          <button
            class="chip"
            :class="{ active: tagSelection.includes('Startup') }"
            @click="toggleFeaturedCategory('Startup')"
            type="button"
          >
            {{ t('jobs.featuredStartup') }}
          </button>
          <button
            class="chip"
            :class="{ active: tagSelection.includes('GameDev') }"
            @click="toggleFeaturedCategory('GameDev')"
            type="button"
          >
            {{ t('jobs.featuredGamedev') }}
          </button>
          <button
            class="chip"
            :class="{ active: tagSelection.includes('Other') }"
            @click="toggleFeaturedCategory('Other')"
            type="button"
          >
            {{ t('jobs.featuredOther') }}
          </button>
          <button
            v-if="tagSelection.length"
            class="chip ghost"
            type="button"
            @click="tagSelection = []"
          >
            {{ t('jobs.reset') }}
          </button>
        </div>
      </div>

      <!-- Tag -->
      <div class="filter-field tags-filter">
        <span class="filter-title">{{ t('jobs.tags') }}</span>
        <HierarchicalTagSelector
          :key="lang"
          :selected-tags="tagSelection"
          :selected-category="selectedTagCategory"
          enable-search
          @toggle-tag="toggleFilterTag"
          @select-category="selectedTagCategory = $event"
          @clear-all="tagSelection = []"
        />
      </div>

      <!-- Employment -->
      <div class="filter-group">
        <span class="filter-label">{{ t('jobs.employment') }}</span>
        <div class="chip-row">
          <button class="chip active" type="button">{{ t('jobs.any') }}</button>
          <button class="chip" type="button">{{ t('jobs.remote') }}</button>
          <button class="chip" type="button">{{ t('jobs.office') }}</button>
        </div>
      </div>

      <!-- Posted -->
      <label class="filter-field">
        <span>{{ t('jobs.posted') }}</span>
        <select>
          <option value="any">{{ t('jobs.anyTime') }}</option>
          <option value="7">{{ t('jobs.last7days') }}</option>
          <option value="30">{{ t('jobs.last30days') }}</option>
        </select>
      </label>

      <!-- Competition -->
      <label class="filter-field">
        <span>{{ t('jobs.competition') }}</span>
        <select>
          <option value="any">{{ t('jobs.any') }}</option>
          <option value="low">{{ t('jobs.low') }}</option>
          <option value="high">{{ t('jobs.high') }}</option>
        </select>
      </label>
    </div>

    <!-- FOOTER -->
    <div class="filters-footer">
      <button
        class="btn primary"
        @click="showFilters = false"
      >
        {{ t('jobs.applyFilters') }}
      </button>
    </div>
  </div>
</div>


  <!-- =====================
   | APPLY MODAL
   ===================== -->
  <div
    v-if="apply.open"
    class="sheet-backdrop"
    @click.self="closeApply"
  >
    <div class="apply-modal">
      <div class="apply-header">
        <div class="apply-title">
          {{ t('jobs.applyToVacancy') }}: {{ apply.vacancy?.title }}
        </div>
        <button class="close-x" @click="closeApply">x</button>
      </div>

      <div class="apply-body">
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

        <label>
          <span>{{ t('jobs.message') }}</span>
          <textarea v-model="apply.message" />
        </label>
      </div>

      <div class="apply-footer">
        <button class="btn ghost" @click="closeApply">
          {{ t('jobs.cancel') }}
        </button>
        <button
          class="btn primary"
          :disabled="sending"
          @click="sendApply"
        >
          {{ sending ? t('jobs.sending') : t('jobs.sendApplication') }}
        </button>
      </div>
    </div>
  </div>

  <section class="seo-links-section">
    <div class="seo-links-inner">
      <h2>{{ t('seoLanding.linksTitle') }}</h2>
      <p>{{ t('seoLanding.linksText') }}</p>
      <div class="seo-links-grid">
        <NuxtLink to="/seo/startup-jobs">{{ t('seoLanding.startupJobs.linkTitle') }}</NuxtLink>
        <NuxtLink to="/seo/gamedev-jobs">{{ t('seoLanding.gamedevJobs.linkTitle') }}</NuxtLink>
      </div>
    </div>
  </section>

  <Footer />
</template>


<script setup lang="ts">
import axios from 'axios'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import HierarchicalTagSelector from '~/components/tags/HierarchicalTagSelector.vue'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'
import { useToast } from '~/composables/useToast'
import { tagLabelMap, matchesTagFilter } from '~/utils/tags'
import { locationRegions } from '~/utils/locations'

const { t, lang } = useI18n()
const toast = useToast()

useHead(() => ({
  title: t('seo.jobs.title'),
  htmlAttrs: {
    lang: lang.value
  },
  meta: [
    { name: 'description', content: t('seo.jobs.description') },
    { name: 'keywords', content: t('seo.keywords.jobs') },
    { name: 'robots', content: 'index, follow' },
    { property: 'og:title', content: t('seo.jobs.title') },
    { property: 'og:description', content: t('seo.jobs.description') },
    { property: 'og:type', content: 'website' },
    { property: 'og:url', content: 'https://unicornstar.online/jobs' },
    { property: 'og:site_name', content: 'Unicornstar' },
    { property: 'og:locale', content: lang.value === 'ru' ? 'ru_RU' : 'en_US' },
    { property: 'og:image', content: 'https://unicornstar.online/images/mainpage/main.png' },
    { name: 'twitter:card', content: 'summary_large_image' },
    { name: 'twitter:title', content: t('seo.jobs.title') },
    { name: 'twitter:description', content: t('seo.jobs.description') },
    { name: 'twitter:image', content: 'https://unicornstar.online/images/mainpage/main.png' }
  ],
  link: [
    { rel: 'canonical', href: 'https://unicornstar.online/jobs' }
  ],
  script: [
    {
      type: 'application/ld+json',
      children: JSON.stringify({
        '@context': 'https://schema.org',
        '@type': 'CollectionPage',
        name: t('seo.jobs.title'),
        description: t('seo.jobs.description'),
        url: 'https://unicornstar.online/jobs'
      })
    }
  ]
}))

function calcDaysAgo(date: string) {
  return Math.max(
    1,
    Math.floor((Date.now() - new Date(date).getTime()) / 86400000)
  )
}

/* =====================
 * TYPES
 * ===================== */
type Vacancy = {
  vacancyId: string
  companyId: string
  title: string
  description: string
  location?: string
  tags?: string[]
  isPremium?: boolean
  analytics?: {
    views: number
    applies: number
    daysAgo?: number
  }
}

type Profile = {
  userId: string
  displayName: string
  avatarUrl?: string
}

type Resume = {
  resumeId: string
  title: string
}

type MatchScore = {
  score: number
  matched: string[]
  missing: string[]
}

/* =====================
 * STATE
 * ===================== */
const err = ref('')
const items = ref<Vacancy[]>([])

/* filters */
const q = ref('')
const locRegion = ref('')
const locCountry = ref('')
const sort = ref<'relevance' | 'match' | 'title'>('relevance')

/* mobile */
const showFilters = ref(false)
const showSaved = ref(false)

const locRegionOptions = locationRegions
const locCountryOptions = computed(() => {
  const region = locationRegions.find(r => r.key === locRegion.value)
  return region?.countries || []
})

const selectedRegion = computed(() =>
  locationRegions.find(r => r.key === locRegion.value) || null
)

/* saved vacancies (SSR-safe) */
const savedIds = ref<string[]>([])

const savedLabel = computed(() =>
  t('jobs.saved').replace('{{count}}', String(savedIds.value.length))
)

/* apply modal */
const apply = reactive({
  open: false,
  vacancy: null as Vacancy | null,
  resumeId: '',
  message: '',
})

const resumes = ref<Resume[]>([])
const applyPreview = ref<MatchScore | null>(null)
const sending = ref(false)

/* match score cache */
const scores = reactive<Record<string, MatchScore>>({})

/* company cache */
const companyCache = reactive<Record<string, Profile | null>>({})

/* =====================
 * SSR-SAFE INIT
 * ===================== */
onMounted(() => {
  try {
    const raw = localStorage.getItem('savedVacancies')
    savedIds.value = raw ? JSON.parse(raw) : []
  } catch {
    savedIds.value = []
  }
})

/* =====================
 * HELPERS
 * ===================== */
function norm(s: string) {
  return (s || '').toLowerCase().trim()
}

function companyName(id: string) {
  return companyCache[id]?.displayName || t('companyPanel.company')
}

function companySearchText(id: string) {
  return norm(companyCache[id]?.displayName || '')
}


function companyAvatar(id: string) {
  return companyCache[id]?.avatarUrl || '/images/com-base.webp'
}

function scoreClass(score: number) {
  if (score >= 80) return 'high'
  if (score >= 50) return 'mid'
  return 'low'
}

function tagLabel(value: string) {
  if (value.startsWith('tags.')) {
    const translated = t(value)
    return translated === value ? formatTagKey(value) : translated
  }
  const key = tagLabelMap[value]
  if (!key) return value
  const translated = t(key)
  return translated === key ? value : translated
}

function formatTagKey(value: string) {
  const tail = value.split('.').pop() || value
  return tail
    .replace(/([a-z])([A-Z0-9])/g, '$1 $2')
    .replace(/\bue(\d+)/gi, 'UE$1')
    .replace(/Csharp/gi, 'C#')
    .replace(/Cpp/gi, 'C++')
    .replace(/MlAi/gi, 'ML/AI')
    .replace(/Ui/gi, 'UI')
    .replace(/Ux/gi, 'UX')
    .replace(/Vfx/gi, 'VFX')
    .replace(/Rpg/gi, 'RPG')
    .replace(/Smm/gi, 'SMM')
    .replace(/Seo/gi, 'SEO')
    .replace(/Ppc/gi, 'PPC')
}

/* =====================
 * SAVE VACANCY
 * ===================== */
function toggleSave(id: string) {
  const set = new Set(savedIds.value)
  set.has(id) ? set.delete(id) : set.add(id)
  savedIds.value = [...set]

  if (import.meta.client) {
    localStorage.setItem(
      'savedVacancies',
      JSON.stringify(savedIds.value)
    )
  }
}

function toggleFilterTag(tag: string) {
  if (tagSelection.value.includes(tag)) {
    tagSelection.value = tagSelection.value.filter(t => t !== tag)
  } else {
    tagSelection.value = [...tagSelection.value, tag]
  }
}

function toggleFeaturedCategory(category: 'IT' | 'Startup' | 'GameDev' | 'Other') {
  if (tagSelection.value.includes(category)) {
    tagSelection.value = tagSelection.value.filter(tag => tag !== category)
    if (selectedTagCategory.value === category) {
      selectedTagCategory.value = null
    }
    return
  }

  tagSelection.value = [...tagSelection.value, category]
  selectedTagCategory.value = category
}

/* =====================
 * LOADERS
 * ===================== */
async function ensureCompany(id: string) {
  if (companyCache[id] !== undefined) return
  companyCache[id] = null
  try {
    const r = await axios.get(`/profile/${id}`)
    companyCache[id] = r.data.profile
  } catch {
    companyCache[id] = null
  }
}

async function load() {
  err.value = ''
  try {
    const r = await axios.get('/vacancies')

    items.value = (r.data.items || []).map((v: any) => ({
      ...v,
      analytics: {
        applies: v.responsesCount || 0,
        daysAgo: calcDaysAgo(v.createdAt),
      },
    }))
    const ids = Array.from(new Set(items.value.map(v => String(v.companyId || '')).filter(Boolean)))
    await Promise.all(ids.map(id => ensureCompany(id)))
  } catch (e: any) {
    const code = e?.response?.data?.error || 'load_failed'
    err.value = getErrorMessage(code, t, 'jobs.errors.loadFailed')
  }
}


/* =====================
 * FILTER + SORT
 * ===================== */
const employment = ref<'any' | 'remote' | 'office'>('any')
const posted = ref<'any' | '7' | '30'>('any')
const competition = ref<'any' | 'low' | 'mid' | 'high'>('any')
const popularity = ref<'any' | 'hot'>('any')
const tagSelection = ref<string[]>([])
const selectedTagCategory = ref<string | null>(null)

const hasActiveFilters = computed(() => {
  return (
    q.value ||
    locRegion.value ||
    locCountry.value ||
    tagSelection.value.length > 0 ||
    employment.value !== 'any' ||
    posted.value !== 'any' ||
    competition.value !== 'any'
  )
})


const filtered = computed(() => {
  let list = items.value.filter(v => {
    const text = norm(v.title + ' ' + v.description)
    const companyText = companySearchText(v.companyId)

    // поиск по:
    // - названию вакансии
    // - описанию
    // - названию компании
    if (
        q.value &&
        !text.includes(norm(q.value)) &&
        !companyText.includes(norm(q.value))
    ) {
        return false
    }

    const locText = norm(v.location || '')
    if (locCountry.value) {
      if (!locText.includes(norm(locCountry.value))) return false
    } else if (locRegion.value) {
      const region = selectedRegion.value
      if (!region || !region.countries.some(c => locText.includes(norm(c)))) return false
    }
    
    // НОВАЯ ИЕРАРХИЧЕСКАЯ ФИЛЬТРАЦИЯ ТЕГОВ
    // Если выбран родительский тег, показываем все вакансии с дочерними тегами
    if (tagSelection.value.length) {
      const vacancyTags = v.tags || []
      if (!matchesTagFilter(vacancyTags, tagSelection.value)) {
        return false
      }
    }

    // employment
    if (employment.value === 'remote') {
      if (
        v.location !== 'Remote' &&
        !text.includes('remote')
      ) return false
    }

    if (employment.value === 'office') {
      if (v.location === 'Remote') return false
    }

    // posted
    if (posted.value !== 'any') {
      const limit =
        posted.value === '7' ? 7 :
        posted.value === '30' ? 30 :
        Infinity

      if ((v.analytics?.daysAgo ?? 999) > limit) return false
    }

    // competition
    if (competition.value !== 'any') {
      const a = v.analytics?.applies ?? 0
      if (competition.value === 'low' && a > 3) return false
      if (competition.value === 'mid' && (a <= 3 || a > 10)) return false
      if (competition.value === 'high' && a <= 10) return false
    }

    // popularity
    if (popularity.value === 'hot') {
      if ((v.analytics?.views ?? 0) < 100) return false
    }

    return true
  })

  if (showSaved.value) {
    list = list.filter(v => savedIds.value.includes(v.vacancyId))
  }

  if (sort.value === 'match') {
    list = [...list].sort(
      (a, b) =>
        (scores[b.vacancyId]?.score || 0) -
        (scores[a.vacancyId]?.score || 0)
    )
  }

  if (sort.value === 'title') {
    list = [...list].sort((a, b) =>
      a.title.localeCompare(b.title)
    )
  }

  return list
})

const jobCountLabel = computed(() =>
  t('jobs.jobCount').replace('{{count}}', String(filtered.value.length))
)

const isFiltering = ref(false)

watch(
  () => [
    q.value,
    locRegion.value,
    locCountry.value,
    employment.value,
    posted.value,
    competition.value,
    sort.value
  ],
  () => {
    isFiltering.value = true
    setTimeout(() => {
      isFiltering.value = false
    }, 280)
  }
)


const filtersKey = computed(() =>
  [
    q.value,
    locRegion.value,
    locCountry.value,
    tagSelection.value.join(','),
    employment.value,
    posted.value,
    competition.value,
    sort.value,
    showSaved.value
  ].join('|')
)


const filteredVisible = computed(() => filtered.value)

/* =====================
 * APPLY FLOW
 * ===================== */
async function openApply(v: Vacancy) {
  apply.open = true
  apply.vacancy = v
  apply.resumeId = ''
  apply.message = ''
  applyPreview.value = null

  try {
    const r = await axios.get('/resumes/my')
    resumes.value = r.data.items || []
    if (!resumes.value.length) {
      const label = t('jobs.noResumes')
      toast.info(label === 'jobs.noResumes' ? 'No resumes available' : label)
    }
  } catch (e: any) {
    const code = e?.response?.data?.error || 'resume_load_failed'
    toast.error(getErrorMessage(code, t, 'jobs.errors.resumeLoadFailed'))
  }
}

function closeApply() {
  apply.open = false
}

watch(
  () => apply.resumeId,
  async (id) => {
    if (!id || !apply.vacancy) return

    const r = await axios.post('/match/score', {
      vacancyId: apply.vacancy.vacancyId,
      resumeId: id,
    })

    applyPreview.value = r.data
    scores[apply.vacancy.vacancyId] = r.data
  }
)

async function sendApply() {
  if (!apply.vacancy || !apply.resumeId) return
  sending.value = true

  try {
    await axios.post('/applications', {
      vacancyId: apply.vacancy.vacancyId,
      resumeId: apply.resumeId,
      message: apply.message,
    })
    const successText = t('jobs.applicationSent')
    toast.success(
      successText === 'jobs.applicationSent' ? 'Application sent' : successText
    )
    closeApply()
  } catch (e: any) {
    const code = e?.response?.data?.error || 'apply_failed'
    toast.error(getErrorMessage(code, t, 'jobs.errors.applyFailed'))
  } finally {
    sending.value = false
  }
}

/* =====================
 * SUBSCRIBE
 * ===================== */
function subscribe() {
  if (!import.meta.client) return

  localStorage.setItem(
    'jobSubscription',
    JSON.stringify({
      q: q.value,
      locRegion: locRegion.value,
      locCountry: locCountry.value,
      tags: tagSelection.value,
    })
  )
  alert(t('jobs.subscribed'))
}

onMounted(load)
</script>


<style scoped>
/* =====================
 * BASE
 * ===================== */
.jobs-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 16px 80px;
}

.page-title {
  font-size: 26px;
  font-weight: 600;
  margin-bottom: 16px;
}

/* =====================
 * MOBILE BAR
 * ===================== */
.mobile-bar {
  display: none;
  gap: 10px;
  margin-bottom: 14px;
}

.mobile-btn {
  flex: 1;
  padding: 12px;
  border-radius: 14px;
  border: 1px solid #e5e7eb;
  background: #fff;
  font-weight: 500;
}

/* =====================
 * LAYOUT
 * ===================== */
.layout {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 24px;
  width: 100%;
  align-items: start;
}

/* =====================
 * FILTERS SIDEBAR
 * ===================== */
.filters-panel {
  position: sticky;
  top: 18px;
  align-self: start;
}

.panel-card {
  background: #fff;
  border-radius: 18px;
  border: 1px solid #e5e7eb;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-shadow: 0 10px 24px rgba(15, 23, 42, 0.06);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 10px;
}

.panel-title {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
}

.panel-count {
  font-size: 12px;
  color: #4f46e5;
  font-weight: 600;
}

.panel-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.panel-field span {
  font-size: 11px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.panel-field input,
.panel-field select {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  font-size: 14px;
  background: #fff;
}

.panel-field input:focus,
.panel-field select:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.12);
}

.location-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 10px;
  align-items: start;
}

.location-row select {
  width: 100%;
  min-width: 0;
}

.country-select {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.field-hint {
  font-size: 11px;
  color: #9ca3af;
}

.location-row select:disabled {
  background: #f3f4f6;
  color: #9ca3af;
  cursor: not-allowed;
}

.panel-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.panel-label {
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.panel-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.panel-reset {
  border: none;
  background: transparent;
  color: #4f46e5;
  font-size: 12px;
  cursor: pointer;
}

.panel-reset:hover {
  text-decoration: underline;
}

.featured-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.featured-chip {
  padding: 8px 14px;
  border-radius: 999px;
  border: 1px solid rgba(79, 70, 229, 0.25);
  background: #eef2ff;
  color: #4338ca;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.featured-chip.active {
  background: #4338ca;
  color: #fff;
  border-color: #4338ca;
  box-shadow: 0 10px 20px rgba(79, 70, 229, 0.3);
}

/* =====================
 * JOB CARD
 * ===================== */
.link-card {
  text-decoration: none;
  color: inherit;
  display: block;
}

.job-card {
  background: #fff;
  border: 1px solid #e5e7eb;
  border-radius: 22px;
  padding: 20px;
  margin-bottom: 16px;
  transition: transform .15s ease, box-shadow .15s ease, border-color .15s;
}

.job-card.premium {
  background: linear-gradient(135deg, rgba(79,70,229,0.08) 0%, #ffffff 60%);
  border: 1px solid rgba(99,102,241,0.35);
  box-shadow: 0 12px 28px rgba(79,70,229,0.12);
  position: relative;
  padding-top: 32px;
}

.job-card.premium::before {
  content: 'Premium';
  position: absolute;
  top: 12px;
  right: 12px;
  background: linear-gradient(90deg, #4f46e5, #6366f1);
  color: #fff;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.6px;
  z-index: 1;
  pointer-events: none;
}

.job-card:hover {
  transform: translateY(-2px);
  border-color: #c7d2fe;
  box-shadow: 0 18px 42px rgba(0,0,0,.08);
}

.job-card.premium:hover {
  border-color: rgba(79,70,229,0.55);
  box-shadow: 0 18px 36px rgba(79,70,229,0.18);
  transform: translateY(-3px);
}

.job-card-inner {
  display: flex;
  gap: 24px;
}

.job-main {
  flex: 1;
}

.job-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  line-height: 1.25;
  color: #111827;
}

.job-card:hover .job-title {
  color: #4f46e5;
}

.job-meta {
  display: flex;
  gap: 6px;
  font-size: 13px;
  color: #475569;
  margin-top: 4px;
}

.job-meta .dot {
  color: #cbd5e1;
}

.job-desc {
  margin-top: 10px;
  font-size: 14px;
  line-height: 1.5;
  color: #334155;
  white-space: pre-line;

  display: -webkit-box;
  -webkit-line-clamp: 4;
  line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.job-tags {
  display: flex;
  gap: 6px;
  margin-top: 10px;
  flex-wrap: wrap;
}

.job-tags .tag {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: 999px;
  background: #eef2ff;
  color: #4338ca;
  font-weight: 500;
}

.tag-groups {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.tag-group-title {
  font-weight: 600;
  font-size: 13px;
}

.tag-select-row {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 8px;
}

.input-wrap {
  position: relative;
  display: grid;
}

.input-wrap input {
  padding-right: 34px;
}

.clear-btn {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  border: 1px solid rgba(17,24,39,0.08);
  background: #fff;
  color: #6b7280;
  width: 22px;
  height: 22px;
  border-radius: 999px;
  line-height: 20px;
  font-weight: 700;
  cursor: pointer;
}

.clear-btn:hover {
  color: #111827;
  border-color: rgba(99,102,241,0.35);
}

.quick-group .chip.ghost {
  background: #fff;
  border: 1px dashed rgba(99,102,241,0.4);
  color: #4f46e5;
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

/* =====================
 * JOB SIDE
 * ===================== */
.job-side {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  justify-content: space-between;
  gap: 14px;
}

.company-avatar {
  width: 52px;
  height: 52px;
  border-radius: 14px;
  border: 1px solid #e5e7eb;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.company-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-fallback {
  font-weight: 600;
  font-size: 20px;
  color: #475569;
}

.apply-cta {
  padding: 10px 18px;
  border-radius: 14px;
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  color: #fff;
  font-weight: 600;
  border: none;
  cursor: pointer;
  box-shadow: 0 8px 20px rgba(79,70,229,.35);
  transition: transform .15s ease, box-shadow .15s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 36px;
  line-height: 1.1;
  white-space: nowrap;
}

.apply-cta:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 26px rgba(79,70,229,.45);
}

.job-card.premium .apply-cta {
  background: linear-gradient(135deg, #111827, #4f46e5);
  box-shadow: 0 10px 22px rgba(79,70,229,.35);
}

.job-card.premium .apply-cta:hover {
  box-shadow: 0 14px 28px rgba(79,70,229,.45);
}

/* =====================
 * BACKDROP
 * ===================== */
.sheet-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,.45);
  z-index: 50;
}

/* =====================
 * MOBILE FILTER SHEET
 * ===================== */
.sheet {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  max-height: 90vh;
  background: #fff;
  border-radius: 22px 22px 0 0;
  display: flex;
  flex-direction: column;
}

.sheet-header {
  display: flex;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.sheet-title {
  font-weight: 600;
}

.sheet-body {
  padding: 16px;
  overflow-y: auto;
}

.sheet-footer {
  padding: 14px 16px 20px;
  border-top: 1px solid #e5e7eb;
}

.apply-btn {
  width: 100%;
  padding: 14px;
  border-radius: 16px;
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  color: #fff;
  font-weight: 600;
  border: none;
}

/* =====================
 * APPLY MODAL
 * ===================== */
.apply-modal {
  position: fixed;
  left: 50%;
  right: auto;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 520px;
  max-width: calc(100vw - 32px);
  max-height: 90vh;
  background: #fff;
  border-radius: 20px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.apply-header {
  display: flex;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.apply-title {
  font-size: 16px;
  font-weight: 600;
}

.close-x {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  opacity: .6;
}

.apply-body {
  padding: 16px;
  overflow-y: auto;
}

.apply-body label {
  display: block;
  margin-bottom: 14px;
}

.apply-body span {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 4px;
  display: block;
}

.apply-body select,
.apply-body textarea {
  width: 100%;
  padding: 12px 14px;
  border-radius: 14px;
  border: 1px solid #e5e7eb;
}

.apply-body textarea {
  min-height: 160px;
  resize: none;
}

.apply-footer {
  display: flex;
  gap: 10px;
  padding: 14px 16px;
  border-top: 1px solid #e5e7eb;
}

.btn {
  flex: 1;
  padding: 12px;
  border-radius: 14px;
}

.btn.primary {
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  color: #fff;
  border: none;
}

.btn.ghost {
  background: #fff;
  border: 1px solid #e5e7eb;
}

/* =====================
 * MOBILE
 * ===================== */
@media (max-width: 1024px) {
  .layout {
    grid-template-columns: 1fr;
  }

  .search {
    display: none;
  }

  .filters-panel {
    display: none;
  }

  .mobile-bar {
    display: flex;
  }
}

@media (max-width: 768px) {
  body:has(.apply-modal) {
    overflow: hidden;
  }

  .jobs-page {
    padding: 16px 14px 120px;
  }

  /* Увеличиваем кнопки для мобильных */
  .mobile-btn {
    padding: 14px 16px;
    font-size: 15px;
    border-radius: 16px;
    /* Увеличенная область тапа */
    min-height: 48px;
    -webkit-tap-highlight-color: transparent;
    touch-action: manipulation;
  }

  .mobile-btn:active {
    transform: scale(0.98);
    background: #f3f4f6;
  }

  /* Карточки вакансий */
  .job-card {
    padding: 16px;
    border-radius: 16px;
  }

  .job-card-inner {
    flex-direction: column;
    gap: 14px;
  }

  .job-title {
    font-size: 17px;
    line-height: 1.4;
  }

  .job-meta {
    font-size: 14px;
    margin-top: 6px;
  }

  .job-desc {
    font-size: 15px;
    line-height: 1.6;
    margin-top: 12px;
  }

  /* Теги на мобильных */
  .job-tags {
    margin-top: 12px;
    gap: 8px;
  }

  .job-tags .tag {
    font-size: 13px;
    padding: 6px 12px;
    /* Увеличенная область тапа */
    min-height: 28px;
  }

  /* Кнопка Apply */
  .apply-cta {
    padding: 12px 20px;
    font-size: 15px;
    /* Увеличенная область тапа */
    min-height: 44px;
    border-radius: 12px;
  }

  .job-side {
    flex-direction: row;
    width: 100%;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
  }

  .company-avatar {
    width: 44px;
    height: 44px;
  }

  .company-avatar img {
    width: 44px;
    height: 44px;
  }

  /* Фильтры в sheet */
  .sheet {
    border-radius: 24px 24px 0 0;
    box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.15);
  }

  .sheet-header {
    padding: 18px 20px;
  }

  .sheet-title {
    font-size: 18px;
  }

  .sheet-body {
    padding: 20px;
    /* Улучшенный скроллинг на iOS */
    -webkit-overflow-scrolling: touch;
  }

  /* Чипы и кнопки фильтров */
  .chip {
    padding: 10px 16px;
    font-size: 14px;
    /* Увеличенная область тапа */
    min-height: 40px;
    border-radius: 12px;
    -webkit-tap-highlight-color: transparent;
  }

  .chip:active {
    transform: scale(0.97);
  }

  .add-chip {
    padding: 8px 14px;
    font-size: 13px;
    /* Увеличенная область тапа */
    min-height: 38px;
  }

  /* Селекты и инпуты */
  .filter-field input,
  .filter-field select {
    padding: 14px 16px;
    font-size: 16px; /* Предотвращает зум на iOS */
    border-radius: 12px;
    /* Увеличенная область тапа */
    min-height: 48px;
  }

  /* Tag chip в выбранных */
  .tag-chip {
    padding: 10px 14px;
    font-size: 14px;
    /* Увеличенная область тапа */
    min-height: 40px;
    border-radius: 10px;
  }

  /* Кнопка Apply в sheet */
  .apply-btn {
    padding: 16px;
    font-size: 16px;
    /* Увеличенная область тапа */
    min-height: 52px;
    border-radius: 16px;
  }

  .apply-btn:active {
    transform: scale(0.98);
  }

  /* Модальное окно Apply */
  .location-row {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .sheet-backdrop {
    padding: 0 !important;
    margin: 0 !important;
  }

  .apply-modal {
    top: auto;
    bottom: 0;
    left: 0;
    right: 0;
    transform: none;
    width: 100% !important;
    max-width: 100% !important;
    max-height: 92vh;
    border-radius: 24px 24px 0 0;
    margin: 0;
    box-sizing: border-box;
    overflow-x: hidden;
  }

  .apply-header {
    padding: 18px 20px;
  }

  .apply-title {
    font-size: 18px;
  }

  .close-x {
    font-size: 24px;
    padding: 8px;
    /* Увеличенная область тапа */
    min-width: 40px;
    min-height: 40px;
  }

  .apply-body {
    padding: 20px;
    max-height: 60vh;
    box-sizing: border-box;
    /* Улучшенный скроллинг на iOS */
    -webkit-overflow-scrolling: touch;
  }

  .apply-body *,
  .apply-body select,
  .apply-body textarea {
    box-sizing: border-box;
    max-width: 100%;
  }

  .apply-body select,
  .apply-body textarea {
    width: 100% !important;
    font-size: 16px; /* Предотвращает зум на iOS */
    padding: 14px 16px;
    border-radius: 12px;
  }

  .apply-body textarea {
    min-height: 120px;
  }

  .apply-footer {
    padding: 16px 20px 24px;
    flex-direction: column;
    gap: 10px;
    box-sizing: border-box;
  }

  .apply-footer .btn {
    width: 100%;
    padding: 14px;
    font-size: 16px;
    /* Увеличенная область тапа */
    min-height: 48px;
    border-radius: 12px;
  }
}

/* Очень маленькие экраны (iPhone SE и подобные) */
@media (max-width: 375px) {
  .jobs-page {
    padding: 12px 10px 100px;
  }

  .job-card {
    padding: 14px;
  }

  .job-title {
    font-size: 16px;
  }

  .job-meta,
  .job-desc {
    font-size: 14px;
  }

  .job-tags .tag {
    font-size: 12px;
    padding: 5px 10px;
  }

  .mobile-btn {
    font-size: 14px;
    padding: 12px 14px;
  }

  .chip {
    font-size: 13px;
    padding: 8px 14px;
  }
}

/* Landscape на мобильных */
@media (max-width: 768px) and (orientation: landscape) {
  .apply-modal,
  .sheet {
    max-height: 85vh;
  }

  .apply-body,
  .sheet-body {
    max-height: 55vh;
  }
}

/* Улучшения для touch-устройств */
@media (hover: none) and (pointer: coarse) {
  /* Убираем hover эффекты на touch-устройствах */
  .job-card:hover {
    transform: none;
  }

  /* Добавляем активное состояние при тапе */
  .job-card:active {
    transform: scale(0.99);
    opacity: 0.95;
  }

  .chip:hover {
    transform: none;
    background: inherit;
  }

  .btn:hover {
    transform: none;
  }

  .btn:active,
  .chip:active {
    opacity: 0.9;
  }

  /* Улучшаем читаемость */
  body {
    -webkit-text-size-adjust: 100%;
    -webkit-font-smoothing: antialiased;
  }
}

.job-card:hover {
  transform: translateY(-2px);
  border-color: #c7d2fe;
  box-shadow: 0 18px 42px rgba(0,0,0,.08);
}

.job-card:hover .job-title {
  color: #4f46e5;
}
.apply-cta:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 26px rgba(79,70,229,.45);
}
.chip:hover {
  border-color: #c7d2fe;
  background: #f8fafc;
}

.chip.active:hover {
  background: #e0e7ff;
}
.subscribe-btn:hover {
  background: #030712;
}
.btn.primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 24px rgba(79,70,229,.45);
}
.btn.ghost:hover {
  background: #f9fafb;
}
.close-x:hover {
  opacity: .9;
}
.search-box input:hover,
.search-box select:hover,
.apply-body textarea:hover,
.apply-body select:hover {
  border-color: #c7d2fe;
}
/* =====================
 * BACKDROP
 * ===================== */
.sheet-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  z-index: 100;
  display: flex;
  align-items: flex-end;
}

/* =====================
 * FILTERS SHEET
 * ===================== */
.filters-sheet {
  width: 100%;
  max-height: 92vh;              /* ключевое: не улетает вниз */
  background: #ffffff;
  border-radius: 22px 22px 0 0;
  box-shadow: 0 -20px 50px rgba(0, 0, 0, 0.25);

  display: flex;
  flex-direction: column;
  overflow: hidden;

  padding-bottom: env(safe-area-inset-bottom);
}

/* =====================
 * HEADER
 * ===================== */
.filters-header {
  display: flex;
  align-items: center;
  justify-content: space-between;

  padding: 14px 18px;
  border-bottom: 1px solid #e5e7eb;
  background: #ffffff;
}

.filters-title {
  font-size: 16px;
  font-weight: 600;
}

.close-x {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  opacity: 0.6;
}
.close-x:hover {
  opacity: 0.9;
}

/* =====================
 * BODY (SCROLLABLE)
 * ===================== */
.filters-body {
  padding: 16px 18px;
  overflow-y: auto;
  flex: 1;

  display: flex;
  flex-direction: column;
  gap: 14px;
  padding-bottom: calc(24px + env(safe-area-inset-bottom));
}

/* =====================
 * FIELD
 * ===================== */
.filter-field {
  display: block;
}

.filter-field span {
  display: block;
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 4px;
}

.filter-field input,
.filter-field select {
  width: 100%;
  padding: 10px 12px;
  border-radius: 14px;
  border: 1px solid #e5e7eb;
  font-size: 14px;
  background: #fff;
}

.filter-field input:focus,
.filter-field select:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

/* =====================
 * GROUP
 * ===================== */
.filter-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.filter-label {
  font-size: 12px;
  color: #6b7280;
}

/* =====================
 * CHIPS
 * ===================== */
.chip-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.chip {
  padding: 6px 14px;
  font-size: 13px;
  border-radius: 999px;
  border: 1px solid #e5e7eb;
  background: #ffffff;
  cursor: pointer;
  color: #374151;
  transition: all 0.15s ease;
}

.chip:hover {
  border-color: #c7d2fe;
}

.chip.active {
  background: #eef2ff;
  color: #4338ca;
  border-color: #c7d2fe;
  font-weight: 500;
}

/* =====================
 * FOOTER (STICKY)
 * ===================== */
.filters-footer {
  padding: 14px 18px;
  border-top: 1px solid #e5e7eb;
  background: #ffffff;
}

.filters-footer .btn {
  width: 100%;
  padding: 12px;
  border-radius: 16px;
  font-weight: 600;
}

/* =====================
 * PRIMARY BUTTON
 * ===================== */
.btn.primary {
  background: linear-gradient(135deg, #4f46e5, #6366f1);
  color: #ffffff;
  border: none;
  cursor: pointer;
  box-shadow: 0 10px 24px rgba(79, 70, 229, 0.35);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.btn.primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 14px 30px rgba(79, 70, 229, 0.45);
}
/* =====================
 * ACTIVE FIELD
 * ===================== */
.filter-field.active input,
.filter-field.active select {
  border-color: #6366f1;
  background: #f5f7ff;
}

.filter-field.active span {
  color: #4338ca;
  font-weight: 500;
}

/* =====================
 * TAGS FILTER
 * ===================== */
.filter-field.tags-filter {
  margin: 1rem 0;
}

.filter-field.tags-filter .filter-title {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 0.75rem;
}

/* =====================
 * ACTIVE FILTERS BAR
 * ===================== */
.active-filters {
  display: flex;
  align-items: center;
  justify-content: space-between;

  padding: 10px 14px;
  margin-bottom: 12px;

  background: #f5f7ff;
  border: 1px solid #c7d2fe;
  border-radius: 14px;

  font-size: 13px;
  color: #4338ca;
}

.reset-btn {
  background: none;
  border: none;
  font-size: 13px;
  font-weight: 500;
  color: #4338ca;
  cursor: pointer;
}

.reset-btn:hover {
  text-decoration: underline;
}

/* =====================
 * CHIP ACTIVE (ещё раз)
 * ===================== */
.chip.active {
  background: #eef2ff;
  color: #4338ca;
  border-color: #c7d2fe;
  font-weight: 500;
}
/* =====================
 * FADE LIST
 * ===================== */
.fade-list-enter-active,
.fade-list-leave-active {
  transition: opacity .18s ease, transform .18s ease;
}

.fade-list-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.fade-list-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
/* =====================
 * SKELETON
 * ===================== */
.skeleton-card {
  padding: 20px;
}

.skeleton-line {
  height: 14px;
  border-radius: 8px;
  background: linear-gradient(
    90deg,
    #f3f4f6 25%,
    #e5e7eb 37%,
    #f3f4f6 63%
  );
  background-size: 400% 100%;
  animation: shimmer 1.2s ease infinite;
}

.skeleton-line.title {
  height: 18px;
  width: 60%;
  margin-bottom: 10px;
}

.skeleton-line.meta {
  width: 40%;
  margin-bottom: 12px;
}

.skeleton-line.desc {
  width: 90%;
  height: 14px;
  margin-bottom: 14px;
}

.skeleton-tags {
  display: flex;
  gap: 8px;
}

.skeleton-pill {
  width: 60px;
  height: 22px;
  border-radius: 999px;
  background: #e5e7eb;
}

.seo-links-section {
  padding: 40px 16px 24px;
}

.seo-links-inner {
  max-width: 1080px;
  margin: 0 auto;
  padding: 24px;
  border: 1px solid #e5e7eb;
  border-radius: 18px;
  background: #fff;
}

.seo-links-inner h2 {
  margin-bottom: 8px;
}

.seo-links-grid {
  margin-top: 12px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.seo-links-grid a {
  color: #3b35c7;
  font-weight: 600;
  text-decoration: none;
}

/* shimmer */
@keyframes shimmer {
  0% {
    background-position: 100% 0;
  }
  100% {
    background-position: -100% 0;
  }
}

</style>






