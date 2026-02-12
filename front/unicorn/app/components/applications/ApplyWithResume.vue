<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useResumesApi } from '~/services/resumesApi'
import { useApplicationsApi } from '~/services/applicationsApi'
import { useI18n } from '~/composables/useI18n'

const props = defineProps<{
  vacancyId: string
}>()

const resumesApi = useResumesApi()
const applicationsApi = useApplicationsApi()
const { t } = useI18n()

const resumes = ref<any[]>([])
const resumeId = ref<string | null>(null)
const message = ref('')

onMounted(async () => {
  const res = await resumesApi.getMyResumes()
  resumes.value = res.items
  if (resumes.value.length === 1) {
    resumeId.value = resumes.value[0].id
  }
})

const apply = async () => {
  if (!resumeId.value) return
  await applicationsApi.apply({
    vacancyId: props.vacancyId,
    resumeId: resumeId.value,
    message: message.value
  })
}
</script>

<template>
  <div>
    <select v-model="resumeId">
      <option
        v-for="r in resumes"
        :key="r.id"
        :value="r.id"
      >
        {{ r.title }}
      </option>
    </select>

    <textarea v-model="message" :placeholder="t('jobs.message')" />

    <button @click="apply">{{ t('jobs.apply') }}</button>
  </div>
</template>
