<template>
  <div>
    <UButton
      label="Backfill Weather"
      color="info"
      variant="outline"
      icon="i-lucide-cloud-download"
      :loading="loading"
      :disabled="disabled || !config.historicalWeatherEnabled()"
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
import { useIntervalFn } from '@vueuse/core'

const config = useConfigStore()
const loading = ref(true)
const disabled = ref(true)
const message = ref('')

const checkStatus = async () => {
  try {
    const response = await historicalWeatherStatus()
    if (response.response?.status === 404) {
      disabled.value = false
      message.value = ''
      pausePolling()
      return
    }
    if (response.error || !response.data || typeof response.data !== 'object' || !('Status' in response.data)) {
      message.value = 'Unable to check weather backfill status.'
      return
    }
    const status = response.data.Status
    if (status === 0 || status === 1) {
      disabled.value = true
      message.value = 'Weather backfill is already in progress.'
      resumePolling()
    } else if (status === 2 || status === 3) {
      disabled.value = false
      message.value = status === 2 ? 'Weather backfill completed.' : 'Weather backfill failed. Please try again.'
      pausePolling()
    } else {
      disabled.value = true
      message.value = 'Unable to check weather backfill status.'
    }
  } catch {
    disabled.value = true
    message.value = 'Unable to check weather backfill status.'
  } finally {
    loading.value = false
  }
}

const { pause: pausePolling, resume: resumePolling } = useIntervalFn(() => {
  void checkStatus()
}, 3_000, { immediate: false })

watch(() => config.historicalWeatherEnabled(), (enabled) => {
  if (enabled) {
    loading.value = true
    void checkStatus()
  } else {
    pausePolling()
    loading.value = false
    disabled.value = true
    message.value = 'Historical weather is not configured.'
  }
}, { immediate: true })

const startBackfill = async () => {
  loading.value = true
  disabled.value = true
  message.value = ''
  try {
    const response = await historicalWeatherBackfill()
    if (response.error || !response.response?.ok) {
      disabled.value = false
      message.value = 'Unable to start weather backfill. Please try again.'
      return
    }
    message.value = 'Weather backfill is already in progress.'
    await checkStatus()
  } catch {
    disabled.value = false
    message.value = 'Unable to start weather backfill. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
