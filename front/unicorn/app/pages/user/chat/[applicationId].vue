<template>
  <div class="page">
    <Header />

    <section class="wrap">
      <div class="hero">
        <div class="hero-left">
          <h1 class="hero-title">{{ t('chat.title') }}</h1>
          <p class="hero-sub">{{ t('chat.application') }} <b>#{{ applicationId }}</b></p>
        </div>

        <div class="hero-actions">
          <NuxtLink class="hero-btn" to="/user/chat">? {{ t('chatList.title') }}</NuxtLink>
          <NuxtLink class="hero-btn" to="/user/applications">{{ t('chatList.backToApps') }}</NuxtLink>
        </div>
      </div>

      <div v-if="err" class="banner danger">{{ err }}</div>
      <div v-if="hint" class="banner warn">{{ hint }}</div>

      <div class="chat-shell">
        <div class="chat-head">
          <div class="chat-head-left">
            <div class="chip">{{ t('chat.application') }}
              <span class="chip-id">{{ applicationId }}</span>
            </div>
            <div class="muted small">{{ t('chat.updateInterval') }}</div>
          </div>

          <button class="btn" :disabled="loading" @click="load(true)">
            {{ loading ? t('chat.refreshing') : t('chat.refresh') }}
          </button>
        </div>

        <div ref="box" class="chat-box">
          <div v-if="loading && !messages.length" class="loading">{{ t('chat.loadingMessages') }}</div>

          <div v-else-if="!messages.length" class="empty-in-chat">
            <div class="empty-title">{{ t('chat.noMessagesTitle') }}</div>
            <div class="empty-text">{{ t('chat.noMessagesText') }}</div>
          </div>

          <div v-else class="msg-list">
            <div
              v-for="m in messages"
              :key="m.messageId"
              class="msg-row"
              :class="m.senderType === 'user' ? 'me' : 'them'"
            >
              <div class="bubble">
                <div class="meta">
                  <span class="who">{{ m.senderType === 'user' ? t('chat.you') : t('chat.company') }}</span>
                  <span class="sep">·</span>
                  <time class="time">{{ formatDateTime(m.createdAt) }}</time>
                </div>

                <div class="text">{{ m.text }}</div>
              </div>
            </div>
          </div>
        </div>

        <div class="composer">
          <textarea
            v-model="text"
            class="inp"
            :placeholder="t('chat.writePlaceholder')"
            rows="2"
            :disabled="sending"
            @keydown="onKeydown"
          />
          <button class="send" @click="send" :disabled="sending || !text.trim()">
            {{ sending ? t('chat.sending') : t('chat.send') }}
          </button>
        </div>
      </div>
    </section>

    <Footer />
  </div>
</template>

<script setup lang="ts">
import axios from 'axios'
import Header from '~/components/Header.vue'
import Footer from '~/components/Footer.vue'
import { useI18n } from '~/composables/useI18n'
import { useAuthStore } from '~/stores/auth'
import { getErrorMessage } from '~/utils/errorMessages'

definePageMeta({ middleware: ['user-auth'] })

const route = useRoute()
const { t, lang } = useI18n()
const auth = useAuthStore()
const err = ref('')
const hint = ref('')
const loading = ref(false)
const sending = ref(false)
const text = ref('')
const box = ref<HTMLElement | null>(null)

type Msg = {
  messageId: string
  senderType: 'user' | 'company'
  text: string
  createdAt: string
}

const messages = ref<Msg[]>([])
const applicationId = computed(() => String(route.params.applicationId || ''))


function chatLastSeenKey(uid: string, appId: string) {
  return `unicorn:user:chat:last_seen:${uid || 'anon'}:${appId}`
}

function saveChatLastSeen(uid: string, appId: string, ts: number) {
  if (!process.client) return
  try {
    localStorage.setItem(chatLastSeenKey(uid, appId), String(ts || 0))
  } catch {}
}

function markChatSeen() {
  if (!process.client) return
  const uid = String((auth as any)?.user?.userId || (auth as any)?.profile?.id || (auth as any)?.me?.id || 'anon')
  const last = messages.value[messages.value.length - 1]
  if (!last?.createdAt) return
  const ts = Date.parse(String(last.createdAt))
  if (Number.isNaN(ts)) return
  saveChatLastSeen(uid, applicationId.value, ts)
}

let timer: any = null

function isNearBottom(el: HTMLElement, px = 140) {
  return el.scrollHeight - el.scrollTop - el.clientHeight < px
}

async function scrollToBottom(force = false) {
  await nextTick()
  if (!box.value) return
  const el = box.value
  if (force || isNearBottom(el)) {
    el.scrollTop = el.scrollHeight
  }
}

async function load(forceScroll = false) {
  err.value = ''
  hint.value = ''
  loading.value = true

  try {
    const res = await axios.get(`/chat/${applicationId.value}/messages`)
    messages.value = res.data.items || []
    await scrollToBottom(forceScroll)
    markChatSeen()
  } catch (e: any) {
    const code = e?.response?.data?.error
    if (code === 'mfa_required') {
      hint.value = t('chat.enable2fa')
    } else {
      err.value = getErrorMessage(code || 'load_failed', t, 'chat.errors.loadFailed')
    }
  } finally {
    loading.value = false
  }
}

async function send() {
  const t = text.value.trim()
  if (!t) return

  sending.value = true
  err.value = ''

  try {
    await axios.post(`/chat/${applicationId.value}/messages`, { text: t })
    text.value = ''
    await load(true)
  } catch (e: any) {
    const code = e?.response?.data?.error
    err.value = getErrorMessage(code || 'send_failed', t, 'chat.errors.sendFailed')
  } finally {
    sending.value = false
  }
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
}

function formatDateTime(v?: string) {
  if (!v) return ''
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return v
  return d.toLocaleString(lang.value === 'ru' ? 'ru-RU' : 'en-US', {
    day: '2-digit',
    month: '2-digit',
    year: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(async () => {
  await load(true)
  timer = setInterval(() => load(false), 3000)
})

onBeforeUnmount(() => timer && clearInterval(timer))

watch(
  () => applicationId.value,
  async () => {
    await load(true)
  }
)
</script>

<style scoped>
.page {
  min-height: 100vh;
  background:
    radial-gradient(1200px 600px at 20% -10%, rgba(79,70,229,0.16), transparent 60%),
    radial-gradient(900px 500px at 90% 10%, rgba(99,102,241,0.12), transparent 60%),
    #ffffff;
}

.wrap {
  max-width: 1100px;
  margin: 0 auto;
  padding: 28px 16px 80px;
}

/* HERO */
.hero {
  position: relative;
  overflow: hidden;

  border-radius: 20px;
  padding: 22px;
  margin-bottom: 14px;

  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;

  background:
    radial-gradient(600px 240px at 85% 30%, rgba(99,102,241,0.16), transparent 60%),
    linear-gradient(135deg, rgba(79,70,229,0.10), rgba(99,102,241,0.06));
  border: 1px solid rgba(99,102,241,0.14);
  box-shadow: 0 12px 40px rgba(17,24,39,.06);
}

.hero::after {
  content: '';
  position: absolute;
  pointer-events: none;

  left: -45px;
  top: 50%;
  transform: translateY(-52%);

  width: 1200px;
  height: 1200px;

  background-image: url('/images/chat_panel.webp');
  background-repeat: no-repeat;
  background-size: contain;
  background-position: center;

  opacity: 0.42;
  z-index: 0;
}

.hero > * {
  position: relative;
  z-index: 1;
}

@media (max-width: 900px) {
  .hero::after {
    width: 768px;
    height: 768px;
    right: -80px;
    top: 40%;
    transform: translateY(-40%);
    bottom: -80px;
    opacity: 0.5;
  }
}

.hero-title {
  margin: 0;
  font-size: 26px;
  font-weight: 900;
  letter-spacing: -0.02em;
}
.hero-sub {
  margin: 6px 0 0;
  font-size: 14px;
  color: #6b7280;
}

.hero-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.hero-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  padding: 10px 14px;
  font-weight: 900;
  background: rgba(17,24,39,0.03);
  color: #111827;
  border: 1px solid rgba(17,24,39,0.06);
  transition: transform .15s ease, background .15s ease;
}
.hero-btn:hover { transform: translateY(-1px); background: rgba(79,70,229,0.06); }

.banner {
  padding: 12px 14px;
  border-radius: 16px;
  border: 1px solid rgba(17,24,39,0.08);
  margin: 10px 0;
  font-weight: 800;
}
.banner.danger { background: rgba(239,68,68,0.08); border-color: rgba(239,68,68,0.18); color: #991b1b; }
.banner.warn { background: rgba(245,158,11,0.10); border-color: rgba(245,158,11,0.22); color: #92400e; }

.chat-shell {
  border-radius: 20px;
  background: rgba(255,255,255,0.92);
  border: 1px solid rgba(17,24,39,0.08);
  box-shadow: 0 18px 50px rgba(17,24,39,0.08);
  overflow: hidden;
  backdrop-filter: blur(10px);
}

.chat-head {
  padding: 14px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border-bottom: 1px solid rgba(17,24,39,0.06);
}

.chat-head-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}

.chip {
  display: inline-flex;
  gap: 8px;
  align-items: center;
  font-weight: 900;
  padding: 8px 10px;
  border-radius: 999px;
  background: rgba(79,70,229,0.07);
  border: 1px solid rgba(99,102,241,0.16);
  color: #3730a3;
}
.chip-id {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-weight: 900;
  color: #111827;
  background: rgba(255,255,255,0.7);
  border: 1px solid rgba(17,24,39,0.06);
  padding: 2px 8px;
  border-radius: 999px;
}

.small { font-size: 12px; }
.muted { color: #6b7280; }

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  padding: 10px 12px;
  font-weight: 900;
  background: rgba(17,24,39,0.03);
  color: #111827;
  border: 1px solid rgba(17,24,39,0.06);
  transition: transform .15s ease, background .15s ease;
}
.btn:hover { transform: translateY(-1px); background: rgba(79,70,229,0.06); }
.btn:disabled { opacity: .6; cursor: not-allowed; }

.chat-box {
  height: min(64vh, 680px);
  overflow: auto;
  padding: 16px;
  background:
    radial-gradient(700px 260px at 20% 0%, rgba(99,102,241,0.08), transparent 60%),
    #ffffff;
}

.loading {
  padding: 16px;
  color: #6b7280;
  font-weight: 800;
}

.empty-in-chat {
  padding: 20px;
  border-radius: 18px;
  border: 1px dashed rgba(17,24,39,0.14);
  background: rgba(17,24,39,0.02);
  text-align: center;
  max-width: 520px;
  margin: 0 auto;
}
.empty-title { font-weight: 900; margin-bottom: 6px; }
.empty-text { color: #6b7280; }

.msg-list { display: grid; gap: 10px; }

.msg-row {
  display: flex;
  width: 100%;
}
.msg-row.me { justify-content: flex-end; }
.msg-row.them { justify-content: flex-start; }

.bubble {
  max-width: min(720px, 92%);
  border-radius: 18px;
  padding: 10px 12px;
  border: 1px solid rgba(17,24,39,0.06);
  box-shadow: 0 10px 24px rgba(17,24,39,0.06);
}

.msg-row.me .bubble {
  background: linear-gradient(135deg, rgba(79,70,229,0.10), rgba(99,102,241,0.06));
  border-color: rgba(99,102,241,0.18);
}
.msg-row.them .bubble {
  background: rgba(17,24,39,0.03);
}

.meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #6b7280;
  font-weight: 800;
  margin-bottom: 6px;
}
.who { color: #111827; font-weight: 900; }
.sep { opacity: .5; }
.time { opacity: .8; }

.text {
  white-space: pre-wrap;
  line-height: 1.55;
  color: #0f172a;
  font-weight: 600;
}

/* COMPOSER */
.composer {
  display: flex;
  gap: 10px;
  padding: 12px 12px;
  border-top: 1px solid rgba(17,24,39,0.06);
  background: rgba(255,255,255,0.95);
}

.inp {
  flex: 1 1 auto;
  border-radius: 16px;
  padding: 12px 12px;
  border: 1px solid rgba(17,24,39,0.10);
  outline: none;
  resize: none;
  font-size: 14px;
  line-height: 1.4;
  background: #fff;
}
.inp:focus {
  border-color: rgba(99,102,241,0.40);
  box-shadow: 0 0 0 4px rgba(99,102,241,0.12);
}

.send {
  flex: 0 0 auto;
  border-radius: 16px;
  padding: 0 16px;
  font-weight: 900;
  color: #fff;
  border: none;
  background: linear-gradient(90deg, #4f46e5, #6366f1);
  box-shadow: 0 12px 26px rgba(79,70,229,0.22);
  transition: transform .15s ease, opacity .15s ease;
}
.send:hover { transform: translateY(-1px); }
.send:disabled { opacity: .6; cursor: not-allowed; }

@media (max-width: 900px) {
  .hero { flex-direction: column; align-items: flex-start; }
  .hero-actions { justify-content: flex-start; }
  .chat-box { height: 62vh; }
}
</style>
