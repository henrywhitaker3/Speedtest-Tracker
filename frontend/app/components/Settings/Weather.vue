<template>
  <UModal v-model:open="open">
    <UButton
      label="Get Weather"
      color="info"
      variant="outline"
      icon="i-lucide-cloud-sun"
      :disabled="!config.weatherEnabled()"
    />

    <template #content>
      <div class="flex flex-col gap-6 p-6">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="text-lg font-semibold">
              Current weather
            </h2>
            <p class="text-muted text-sm">
              Weather at the configured location.
            </p>
          </div>
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-lucide-refresh-cw"
            aria-label="Refresh weather"
            :loading="isLoading"
            @click="loadWeather"
          />
        </div>

        <p
          v-if="isLoading"
          class="text-muted text-sm"
        >
          Loading weather...
        </p>
        <div
          v-else-if="error"
          class="flex flex-col items-start gap-3"
        >
          <p class="text-error">
            {{ error }}
          </p>
          <UButton
            color="neutral"
            variant="outline"
            @click="loadWeather"
          >
            Try again
          </UButton>
        </div>
        <dl
          v-else-if="weather"
          class="grid grid-cols-2 gap-4 text-sm"
        >
          <div>
            <dt class="text-muted">
              Conditions
            </dt>
            <dd>{{ weather.code?.emoji }} {{ weather.code?.name }}</dd>
          </div>
          <div>
            <dt class="text-muted">
              Temperature
            </dt>
            <dd>{{ weather.temperature }} °C</dd>
          </div>
          <div>
            <dt class="text-muted">
              Precipitation (past hour)
            </dt>
            <dd>{{ weather.precipitation }} mm</dd>
          </div>
          <div>
            <dt class="text-muted">
              Wind speed
            </dt>
            <dd>{{ weather.wind_speed != null ? `${formatWindSpeed(weather.wind_speed)} m/s` : '-' }}</dd>
          </div>
        </dl>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { WeatherWeather } from '~/client'

const api = useApi()
const open = ref(false)
const isLoading = ref(false)
const error = ref<string>()
const weather = ref<WeatherWeather>()
const config = useConfigStore()

const loadWeather = async () => {
  isLoading.value = true
  error.value = undefined

  try {
    const response = await api.get<{ 200: WeatherWeather }>({ url: '/v1/config/weather' })
    if (response.error || !response.data) {
      error.value = 'Unable to load current weather.'
      return
    }
    weather.value = response.data
  } catch {
    error.value = 'Unable to load current weather.'
  } finally {
    isLoading.value = false
  }
}

watch(open, (isOpen) => {
  if (isOpen) {
    void loadWeather()
  }
})
</script>
