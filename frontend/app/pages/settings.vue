<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold">
          Settings
        </h1>
        <p class="text-muted text-sm">
          Current Speedtest Tracker configuration
        </p>
      </div>
      <UButton
        to="/"
        color="neutral"
        variant="outline"
        icon="i-lucide-arrow-left"
      >
        Dashboard
      </UButton>
    </div>

    <UCard v-if="isLoading">
      <div class="text-muted flex items-center gap-3 py-6">
        <UIcon
          name="i-lucide-loader-circle"
          class="size-5 animate-spin"
        />
        Loading configuration…
      </div>
    </UCard>

    <UCard v-else-if="error">
      <div class="flex flex-col items-start gap-3 py-2">
        <p class="text-error">
          {{ error }}
        </p>
        <UButton
          color="neutral"
          variant="outline"
          @click="loadConfig"
        >
          Try again
        </UButton>
      </div>
    </UCard>

    <SettingsConfig />

    <SettingsActions />

    <TokensTable />
  </div>
</template>

<script setup lang="ts">
const configStore = useConfigStore()
const isLoading = ref(true)
const error = ref<string>()

useSeoMeta({
  title: 'Settings'
})
definePageMeta({
  middleware: ['auth', 'admin']
})

const loadConfig = async () => {
  isLoading.value = true
  error.value = undefined

  try {
    await configStore.load()
  } catch {
    error.value = 'Unable to load the current configuration.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  void loadConfig()
})
</script>
