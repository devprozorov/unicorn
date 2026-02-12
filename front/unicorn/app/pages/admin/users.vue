<template>
  <div class="min-h-screen bg-gray-50">
    <div class="mx-auto max-w-7xl py-8 px-4">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-2xl font-semibold">{{ t('admin.users.title') }}</h1>
        <div class="flex items-center gap-2">
          <button @click="loadUsers" class="rounded-md bg-indigo-600 text-white px-3 py-2 hover:bg-indigo-700">
            {{ t('admin.users.refresh') }}
          </button>
        </div>
      </div>

      <!-- FILTERS -->
      <div class="bg-white rounded-lg shadow p-4 mb-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ t('admin.users.filters.search') }}</label>
            <input
              v-model="filters.search"
              @input="loadUsers"
              :placeholder="t('admin.users.filters.searchPlaceholder')"
              class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ t('admin.users.filters.type') }}</label>
            <select v-model="filters.type" @change="loadUsers" class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500">
              <option value="">{{ t('admin.common.all') }}</option>
              <option value="user">{{ t('admin.users.filters.typeUser') }}</option>
              <option value="company">{{ t('admin.users.filters.typeCompany') }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ t('admin.users.filters.blocked') }}</label>
            <select v-model="filters.blocked" @change="loadUsers" class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500">
              <option value="">{{ t('admin.common.all') }}</option>
              <option value="true">{{ t('admin.users.filters.blockedYes') }}</option>
              <option value="false">{{ t('admin.users.filters.blockedNo') }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ t('admin.users.filters.premium') }}</label>
            <select v-model="filters.premium" @change="loadUsers" class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500">
              <option value="">{{ t('admin.common.all') }}</option>
              <option value="true">{{ t('admin.users.filters.premiumYes') }}</option>
              <option value="false">{{ t('admin.users.filters.premiumNo') }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-700 mb-1">{{ t('admin.users.filters.deleted') }}</label>
            <select v-model="filters.deleted" @change="loadUsers" class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500">
              <option value="false">{{ t('admin.users.filters.deletedNo') }}</option>
              <option value="">{{ t('admin.common.all') }}</option>
              <option value="true">{{ t('admin.users.filters.deletedYes') }}</option>
            </select>
          </div>
        </div>
      </div>

      <div class="bg-white shadow rounded-lg overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-100">
            <tr>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500">{{ t('admin.users.table.id') }}</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500">{{ t('admin.users.table.login') }}</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500">{{ t('admin.users.table.displayName') }}</th>
              <th class="px-4 py-2 text-left text-xs font-medium text-gray-500">{{ t('admin.users.table.type') }}</th>
              <th class="px-4 py-2 text-right text-xs font-medium text-gray-500">{{ t('admin.users.table.actions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="u in users" :key="u.userId">
              <td class="px-4 py-2 font-mono text-xs">{{ u.userId }}</td>
              <td class="px-4 py-2">{{ u.login }}</td>
              <td class="px-4 py-2">{{ u.displayName }}</td>
              <td class="px-4 py-2">
                <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium"
                  :class="u.type === 'company' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800'">
                  {{ u.type }}
                </span>
              </td>
              <td class="px-4 py-2 text-right">
                <div class="inline-flex gap-2">
                  <button @click="openDetails(u.userId)" class="rounded-md bg-gray-500 text-white px-3 py-1 hover:bg-gray-600 text-xs">
                    {{ t('admin.users.actions.details') }}
                  </button>
                  <button @click="onBlock(u.userId)" class="rounded-md bg-yellow-500 text-white px-3 py-1 hover:bg-yellow-600 text-xs">
                    {{ t('admin.users.actions.block') }}
                  </button>
                  <button @click="onUnblock(u.userId)" class="rounded-md bg-green-600 text-white px-3 py-1 hover:bg-green-700 text-xs">
                    {{ t('admin.users.actions.unblock') }}
                  </button>
                  <button @click="onDelete(u.userId)" class="rounded-md bg-red-600 text-white px-3 py-1 hover:bg-red-700 text-xs">
                    {{ t('admin.users.actions.delete') }}
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!loading && users.length === 0">
              <td colspan="5" class="px-4 py-6 text-center text-gray-500">{{ t('admin.users.empty') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="loading" class="mt-4 text-sm text-gray-500">{{ t('admin.common.loading') }}</div>
      <div v-if="error" class="mt-2 text-sm text-red-600">{{ error }}</div>

      <!-- USER DETAILS MODAL -->
      <div v-if="detailsModal.open" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4" @click.self="closeDetails">
        <div class="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
          <div class="flex items-center justify-between p-4 border-b">
            <h2 class="text-lg font-semibold">{{ t('admin.users.details.title') }}</h2>
            <button @click="closeDetails" class="text-gray-500 hover:text-gray-700">x</button>
          </div>

          <div v-if="detailsModal.loading" class="p-8 text-center text-gray-500">{{ t('admin.common.loading') }}</div>
          <div v-else-if="detailsModal.error" class="p-4 text-red-600">{{ detailsModal.error }}</div>
          <div v-else-if="detailsModal.data" class="p-4 space-y-4">
            <!-- Basic Info -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.table.id') }}</label>
                <div class="font-mono text-sm">{{ detailsModal.data.userId }}</div>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.table.type') }}</label>
                <div>{{ detailsModal.data.type }}</div>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.table.login') }}</label>
                <div>{{ detailsModal.data.login }}</div>
              </div>
              <div>
                <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.table.displayName') }}</label>
                <div>{{ detailsModal.data.displayName }}</div>
              </div>
            </div>

            <!-- Status -->
            <div class="border-t pt-4">
              <h3 class="font-semibold mb-2">{{ t('admin.users.details.status.title') }}</h3>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.details.status.blocked') }}</label>
                  <div :class="detailsModal.data.status.blocked ? 'text-red-600 font-semibold' : 'text-green-600'">
                    {{ detailsModal.data.status.blocked ? t('admin.common.yes') : t('admin.common.no') }}
                  </div>
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.details.status.deleted') }}</label>
                  <div :class="detailsModal.data.status.deleted ? 'text-red-600 font-semibold' : 'text-green-600'">
                    {{ detailsModal.data.status.deleted ? t('admin.common.yes') : t('admin.common.no') }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Subscription -->
            <div class="border-t pt-4">
              <h3 class="font-semibold mb-2">{{ t('admin.users.details.subscription.title') }}</h3>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.details.subscription.active') }}</label>
                  <div :class="detailsModal.data.subscription.active ? 'text-green-600 font-semibold' : 'text-gray-500'">
                    {{ detailsModal.data.subscription.active ? t('admin.common.yes') : t('admin.common.no') }}
                  </div>
                </div>
                <div v-if="detailsModal.data.subscription.active">
                  <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.details.subscription.until') }}</label>
                  <div>{{ formatDate(detailsModal.data.subscription.until) }}</div>
                </div>
              </div>
              <div class="mt-3 flex gap-2">
                <input
                  v-model.number="subscriptionDays"
                  type="number"
                  min="1"
                  max="3650"
                  :placeholder="t('admin.users.details.subscription.daysPlaceholder')"
                  class="w-24 rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                />
                <button @click="activateSubscription" class="rounded-md bg-green-600 text-white px-3 py-1 hover:bg-green-700 text-sm">
                  {{ t('admin.users.details.subscription.activate') }}
                </button>
                <button @click="deactivateSubscription" class="rounded-md bg-red-600 text-white px-3 py-1 hover:bg-red-700 text-sm">
                  {{ t('admin.users.details.subscription.deactivate') }}
                </button>
              </div>
            </div>

            <!-- MFA -->
            <div class="border-t pt-4">
              <h3 class="font-semibold mb-2">{{ t('admin.users.details.mfa.title') }}</h3>
              <div>
                <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.details.mfa.totpEnabled') }}</label>
                <div>{{ detailsModal.data.mfa.totpEnabled ? t('admin.common.yes') : t('admin.common.no') }}</div>
              </div>
            </div>

            <!-- Timestamps -->
            <div class="border-t pt-4">
              <h3 class="font-semibold mb-2">{{ t('admin.users.details.timestamps.title') }}</h3>
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.details.timestamps.createdAt') }}</label>
                  <div class="text-sm">{{ formatDate(detailsModal.data.createdAt) }}</div>
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-500">{{ t('admin.users.details.timestamps.updatedAt') }}</label>
                  <div class="text-sm">{{ formatDate(detailsModal.data.updatedAt) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useAdminApi, type AdminUserItem, type AdminUserDetails, type UserFilters } from '~/composables/useAdminApi'
import { useToast } from '~/composables/useToast'
import { useI18n } from '~/composables/useI18n'
import { getErrorMessage } from '~/utils/errorMessages'

definePageMeta({
  middleware: ['admin-auth']
})

const toast = useToast()
const { t, lang } = useI18n()
const { getUsers, getUserDetails, blockUser, unblockUser, deleteUser, activateSubscription: activateSub, deactivateSubscription: deactivateSub } = useAdminApi()

const filters = reactive<UserFilters>({
  type: undefined,
  search: '',
  blocked: undefined,
  deleted: 'false',
  premium: undefined
})

const users = ref<AdminUserItem[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const detailsModal = reactive({
  open: false,
  loading: false,
  error: '',
  data: null as AdminUserDetails | null
})

const subscriptionDays = ref(30)

const loadUsers = async () => {
  loading.value = true
  error.value = null
  try {
    const cleanFilters: UserFilters = {}
    if (filters.type) cleanFilters.type = filters.type
    if (filters.search) cleanFilters.search = filters.search
    if (filters.blocked) cleanFilters.blocked = filters.blocked
    if (filters.deleted) cleanFilters.deleted = filters.deleted
    if (filters.premium) cleanFilters.premium = filters.premium

    users.value = await getUsers(cleanFilters)
  } catch (e: any) {
    error.value = getErrorMessage(e?.message || 'load_failed', t, 'admin.errors.loadUsers')
    toast.error(error.value)
  } finally {
    loading.value = false
  }
}

const openDetails = async (userId: string) => {
  detailsModal.open = true
  detailsModal.loading = true
  detailsModal.error = ''
  detailsModal.data = null

  try {
    detailsModal.data = await getUserDetails(userId)
  } catch (e: any) {
    detailsModal.error = getErrorMessage(e?.message || 'load_failed', t, 'admin.errors.loadDetails')
    toast.error(detailsModal.error)
  } finally {
    detailsModal.loading = false
  }
}

const closeDetails = () => {
  detailsModal.open = false
  detailsModal.data = null
}

const onBlock = async (userId: string) => {
  if (!confirm(t('admin.users.confirm.block'))) return

  try {
    const res = await blockUser(userId)
    if (res?.ok) {
      toast.success(t('admin.users.success.block'))
      await loadUsers()
      if (detailsModal.data?.userId === userId) {
        await openDetails(userId)
      }
    } else {
      toast.error(getErrorMessage(res?.error || 'block_failed', t, 'admin.errors.blockFailed'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'block_failed', t, 'admin.errors.blockFailed'))
  }
}

const onUnblock = async (userId: string) => {
  if (!confirm(t('admin.users.confirm.unblock'))) return

  try {
    const res = await unblockUser(userId)
    if (res?.ok) {
      toast.success(t('admin.users.success.unblock'))
      await loadUsers()
      if (detailsModal.data?.userId === userId) {
        await openDetails(userId)
      }
    } else {
      toast.error(getErrorMessage(res?.error || 'unblock_failed', t, 'admin.errors.unblockFailed'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'unblock_failed', t, 'admin.errors.unblockFailed'))
  }
}

const onDelete = async (userId: string) => {
  if (!confirm(t('admin.users.confirm.delete'))) return

  try {
    const res = await deleteUser(userId)
    if (res?.ok) {
      toast.success(t('admin.users.success.delete'))
      await loadUsers()
      if (detailsModal.data?.userId === userId) {
        closeDetails()
      }
    } else {
      toast.error(getErrorMessage(res?.error || 'delete_failed', t, 'admin.errors.deleteFailed'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'delete_failed', t, 'admin.errors.deleteFailed'))
  }
}

const activateSubscription = async () => {
  if (!detailsModal.data) return
  if (!subscriptionDays.value || subscriptionDays.value < 1 || subscriptionDays.value > 3650) {
    toast.error(t('admin.errors.invalidDays'))
    return
  }

  try {
    const res = await activateSub(detailsModal.data.userId, subscriptionDays.value)
    if (res?.ok) {
      toast.success(t('admin.users.success.activateSubscription').replace('{{days}}', String(subscriptionDays.value)))
      await openDetails(detailsModal.data.userId)
      await loadUsers()
    } else {
      toast.error(getErrorMessage(res?.error || 'activation_failed', t, 'admin.errors.activateSubscription'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'activation_failed', t, 'admin.errors.activateSubscription'))
  }
}

const deactivateSubscription = async () => {
  if (!detailsModal.data) return
  if (!confirm(t('admin.users.confirm.deactivateSubscription'))) return

  try {
    const res = await deactivateSub(detailsModal.data.userId)
    if (res?.ok) {
      toast.success(t('admin.users.success.deactivateSubscription'))
      await openDetails(detailsModal.data.userId)
      await loadUsers()
    } else {
      toast.error(getErrorMessage(res?.error || 'deactivation_failed', t, 'admin.errors.deactivateSubscription'))
    }
  } catch (e: any) {
    toast.error(getErrorMessage(e?.message || 'deactivation_failed', t, 'admin.errors.deactivateSubscription'))
  }
}

const formatDate = (dateStr: string) => {
  if (!dateStr || dateStr === '0001-01-01T00:00:00Z') return '-'
  try {
    return new Date(dateStr).toLocaleString(lang.value === 'ru' ? 'ru-RU' : 'en-US', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

onMounted(loadUsers)
</script>
