<template>
  <div>
    <UButton
      label="Backfill Weather"
      color="info"
      variant="outline"
      icon="i-lucide-cloud-download"
      :loading="loading"
      :disabled="disabled"
      @click="startBackfill"
    />
    <p
      v-if="message"
      class="text-muted mt-2 font-sans text-sm"
      role="status"
    >
      {{ message }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { historicalWeatherBackfill, historicalWeatherStatus } from '~/client'

const loading = ref(true)
const disabled = ref(true)
const message = ref('')

onMounted(async () => {
  try {
    const response = await historicalWeatherStatus()
    if (response.response?.status === 404) {
      disabled.value = false
      return
    }
    if (response.error || !response.data || typeof response.data !== 'object' || !('Status' in response.data)) {
      message.value = 'Unable to check weather backfill status.'
      return
    }
    const status = response.data.Status
    if (status === 0 || status === 1) {
      message.value = 'Weather backfill is already in progress.'
    } else if (status === 2 || status === 3) {
      disabled.value = false
    } else {
      message.value = 'Unable to check weather backfill status.'
    }
  } catch {
    message.value = 'Unable to check weather backfill status.'
  } finally {
    loading.value = false
  }
})

const startBackfill = async () => {
  loading.value = true
  disabled.value = true
  await historicalWeatherBackfill()
}
</script>
