<template>
  <section class="hero">
    <!-- LEFT -->
    <div class="hero-left">
      <h1 class="hero-title reveal">{{ t('hero.title') }}</h1>
      <p class="hero-text reveal delay-1">{{ t('hero.text') }}</p>
      <p class="hero-motto reveal delay-2">— {{ t('hero.motto') }}</p>

      <div class="hero-stats">
        <div class="stat reveal delay-3">
          <div class="stat-value">{{ stats.resumes }}</div>
          <div class="stat-label">{{ t('stats.resumes') }}</div>
        </div>
        <div class="stat reveal delay-4">
          <div class="stat-value">{{ stats.companies }}</div>
          <div class="stat-label">{{ t('stats.employers') }}</div>
        </div>
        <div class="stat reveal delay-5">
          <div class="stat-value">{{ stats.vacancies }}</div>
          <div class="stat-label">{{ t('stats.vacancies') }}</div>
        </div>
      </div>
    </div>

    <!-- RIGHT -->
    <div class="hero-right reveal delay-3">
      <div
        ref="visualRef"
        class="hero-visual"
        :style="visualStyle"
        @mousemove="onMove"
        @mouseleave="onLeave"
      >
        <!-- очень мягкий свет позади (не квадратом) -->
        <div class="hero-blob" aria-hidden="true"></div>

        <!-- эффекты по силуэту -->
        <div class="hero-imgFx" aria-hidden="true"></div>

        <!-- логотип -->
        <img class="hero-img" :src="logoSrc" alt="Unicornstar" />
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref, computed } from 'vue'
import { useI18n } from '~/composables/useI18n'
const { t } = useI18n()

const logoSrc = '/images/mainpage/main.png'
const visualRef = ref<HTMLElement | null>(null)

// Stats: will be filled from backend; defaults match previous hardcoded values
const stats = ref({
  resumes: '12 480+',
  companies: '340+',
  vacancies: '1 240',
})

const visualStyle = computed(() => ({
  '--mask': `url(${logoSrc})`,
}))

// Caching configuration
const CACHE_KEY = 'home_stats_v1'
const CACHE_TTL = 1000 * 60 * 5 // 5 minutes
let refreshTimer: ReturnType<typeof setTimeout> | null = null

function formatNumber(n: number) {
  try {
    return n.toLocaleString('ru-RU')
  } catch (e) {
    return String(n)
  }
}

async function loadStats() {
  try {
    const raw = localStorage.getItem(CACHE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (parsed?.ts && Date.now() - parsed.ts < CACHE_TTL && parsed.data) {
        const d = parsed.data
        stats.value.resumes = `${formatNumber(d.total_resumes)}+`
        stats.value.companies = `${formatNumber(d.total_companies)}+`
        stats.value.vacancies = `${formatNumber(d.total_vacancies)}+`
        scheduleRefresh()
        return
      }
    }

    const res = await fetch('https://unicornstar.online/api/home/stats')
    if (!res.ok) throw new Error('Network error')
    const json = await res.json()
    const d = json as { total_resumes: number; total_companies: number; total_vacancies: number }

    // update UI
    stats.value.resumes = `${formatNumber(d.total_resumes)}+`
    stats.value.companies = `${formatNumber(d.total_companies)}+`
    stats.value.vacancies = `${formatNumber(d.total_vacancies)}`

    // cache
    try {
      localStorage.setItem(CACHE_KEY, JSON.stringify({ ts: Date.now(), data: d }))
    } catch (e) {
      // ignore localStorage errors
    }

    scheduleRefresh()
  } catch (e) {
    // keep defaults on error
    scheduleRefresh()
  }
}

function scheduleRefresh() {
  if (refreshTimer) clearTimeout(refreshTimer)
  refreshTimer = setTimeout(() => {
    loadStats()
  }, CACHE_TTL)
}

/* REVEAL */
onMounted(() => {
  const io = new IntersectionObserver(
    entries => {
      entries.forEach(e => {
        if (e.isIntersecting) {
          e.target.classList.add('is-visible')
          io.unobserve(e.target)
        }
      })
    },
    { threshold: 0.15 }
  )
  document.querySelectorAll('.reveal').forEach(el => io.observe(el))
})

/* CURSOR INERTIA (чтобы “не пропадало резко”) */
let raf = 0
let hover = false

const target = { mx: 50, my: 50, tx: 0, ty: 0 }
const cur = { mx: 50, my: 50, tx: 0, ty: 0 }

function tick() {
  const el = visualRef.value
  if (el) {
    const k = hover ? 0.10 : 0.04 // когда ушёл курсор — очень медленно возвращаемся
    cur.mx += (target.mx - cur.mx) * k
    cur.my += (target.my - cur.my) * k
    cur.tx += (target.tx - cur.tx) * k
    cur.ty += (target.ty - cur.ty) * k

    el.style.setProperty('--mx', `${cur.mx}%`)
    el.style.setProperty('--my', `${cur.my}%`)
    el.style.setProperty('--tx', `${cur.tx}px`)
    el.style.setProperty('--ty', `${cur.ty}px`)
  }
  raf = requestAnimationFrame(tick)
}

onMounted(() => {
  raf = requestAnimationFrame(tick)
})
onMounted(() => {
  // load stats cached or from backend
  loadStats()
})
onBeforeUnmount(() => {
  cancelAnimationFrame(raf)
  if (refreshTimer) clearTimeout(refreshTimer)
})

function onMove(e: MouseEvent) {
  const el = visualRef.value
  if (!el) return
  hover = true

  const rect = el.getBoundingClientRect()
  const x = ((e.clientX - rect.left) / rect.width) * 100
  const y = ((e.clientY - rect.top) / rect.height) * 100

  target.mx = x
  target.my = y

  const dx = (e.clientX - (rect.left + rect.width / 2)) / rect.width
  const dy = (e.clientY - (rect.top + rect.height / 2)) / rect.height
  target.tx = dx * 14
  target.ty = dy * 14
}

function onLeave() {
  hover = false
  target.mx = 50
  target.my = 50
  target.tx = 0
  target.ty = 0
}
</script>

<style scoped>
/* LAYOUT */
.hero {
  display: grid;
  grid-template-columns: 1.1fr 1fr;
  gap: 48px;
  align-items: center;
  padding: 64px 0 48px;
}

.hero-left { max-width: 560px; padding-left: 15px; padding-right: 40px; }
.hero-title { font-size: 36px; line-height: 1.15; font-weight: 600; margin: 0 0 8px; }
.hero-text { max-width: 520px; font-size: 14px; line-height: 1.6; color: #4b5563; margin-bottom: 8px; }
.hero-motto { font-size: 13px; color: #6b7280; margin-bottom: 16px; }
.hero-stats { display: flex; gap: 24px; margin-top: 12px; }
.stat-value { font-size: 20px; font-weight: 600; }
.stat-label { font-size: 12px; color: #6b7280; }

.hero-right { display: flex; justify-content: flex-end; padding-right: 15px; }

/* ✅ КЛЮЧ: контейнер НЕ рисует “карточку” и НЕ даёт белый прямоугольник */
.hero-visual {
  display: inline-grid;
  place-items: center;
  line-height: 0;

  width: min(265px, 100%);
  height: auto;
  margin-right: 10px;
  background: transparent !important;
  box-shadow: none !important;
  border-radius: 0 !important;
  overflow: visible !important;

  position: relative;
  isolation: isolate;

  --mx: 50%;
  --my: 50%;
  --tx: 0px;
  --ty: 0px;
}

/* мягкий свет позади (не квадратом) */
.hero-blob {
  position: absolute;
  inset: -18%;
  border-radius: 999px;
  pointer-events: none;

  background:
    radial-gradient(closest-side at var(--mx) var(--my),
      rgba(99,102,241,.14),
      transparent 66%
    ),
    radial-gradient(closest-side at 35% 55%,
      rgba(233,213,255,.12),
      transparent 60%
    );

  filter: blur(34px);
  opacity: .9;
  transform: translate(calc(var(--tx) * 0.6), calc(var(--ty) * 0.6));
}

/* логотип */
.hero-img {
  position: relative;
  z-index: 3;
  width: 100%;
  height: auto;
  display: block;

  transform: translate(var(--tx), var(--ty)) scale(1.02);
  transition: transform .35s cubic-bezier(.22,.61,.36,1);

  filter:
    drop-shadow(0 18px 55px rgba(99,102,241,.18))
    drop-shadow(0 6px 18px rgba(147,197,253,.10));
}

/* эффекты по силуэту (mask) — + лёгкая “вечная жизнь”, чтобы не пропадало */
.hero-imgFx {
  position: absolute;
  inset: 0;
  z-index: 4;
  pointer-events: none;

  mix-blend-mode: screen;
  opacity: .85;

  -webkit-mask-image: var(--mask);
  -webkit-mask-repeat: no-repeat;
  -webkit-mask-position: center;
  -webkit-mask-size: contain;

  mask-image: var(--mask);
  mask-repeat: no-repeat;
  mask-position: center;
  mask-size: contain;
}

.hero-imgFx::before {
  content: "";
  position: absolute;
  inset: -22%;
  border-radius: 999px;

  background: conic-gradient(
    from 0deg,
    rgba(255,255,255,0),
    rgba(255,255,255,.14),
    rgba(99,102,241,.16),
    rgba(255,255,255,0)
  );

  filter: blur(22px);
  opacity: .55;
  animation: spin 12s linear infinite;
}

.hero-imgFx::after {
  content: "";
  position: absolute;
  inset: 0;

  background:
    radial-gradient(340px 240px at var(--mx) var(--my),
      rgba(255,255,255,.45),
      transparent 62%
    ),
    radial-gradient(520px 380px at calc(var(--mx) + 10%) calc(var(--my) + 10%),
      rgba(99,102,241,.20),
      transparent 72%
    );

  opacity: .70;
}

@keyframes spin { to { transform: rotate(360deg); } }

/* REVEAL */
.reveal {
  opacity: 0;
  transform: translateY(20px);
  filter: blur(6px);
  transition:
    opacity .8s cubic-bezier(.22,.61,.36,1),
    transform .8s cubic-bezier(.22,.61,.36,1),
    filter .8s cubic-bezier(.22,.61,.36,1);
}
.reveal.is-visible { opacity: 1; transform: translateY(0); filter: blur(0); }
.delay-1 { transition-delay: .1s }
.delay-2 { transition-delay: .2s }
.delay-3 { transition-delay: .3s }
.delay-4 { transition-delay: .4s }
.delay-5 { transition-delay: .5s }

@media (max-width: 1024px) {
  .hero { grid-template-columns: 1fr; gap: 28px; padding: 48px 0 48px; }
  .hero-right { display: none; }
  .hero-visual { width: min(460px, 100%); }
}

@media (max-width: 768px) {
  .hero { padding: 32px 0 40px; }
  .hero-stats { margin-bottom: 24px; }
}
</style>
