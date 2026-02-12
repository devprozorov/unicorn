// app/plugins/admin-init.client.ts
export default defineNuxtPlugin(() => {
  const token = useState<string | null>('adminToken')
  if (token.value) return
  try {
    const t = localStorage.getItem('adminToken')
    if (t) token.value = t
  } catch {}
})
