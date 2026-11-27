<template>
  <UCard>
    <template #header>
      <h2 class="text-lg font-semibold">
        Weather
      </h2>
    </template>

    <dl class="grid gap-x-8 gap-y-5 sm:grid-cols-2 lg:grid-cols-4">
      <div>
        <dt class="text-muted text-sm">
          Conditions
        </dt>
        <dd class="mt-1 font-medium">
          {{ weather.code?.emoji }} {{ weather.code?.name ?? '-' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Temperature
        </dt>
        <dd class="mt-1 font-medium">
          {{ weather.temperature != null ? `${weather.temperature} °C` : '-' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Precipitation (past hour)
        </dt>
        <dd class="mt-1 font-medium">
          {{ weather.precipitation != null ? `${weather.precipitation} mm` : '-' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Wind speed
        </dt>
        <dd class="mt-1 font-medium">
          {{ weather.wind_speed != null ? `${formatWindSpeed(weather.wind_speed)} m/s` : '-' }}
        </dd>
      </div>
    </dl>
  </UCard>
</template>

<script setup lang="ts">
import type { WeatherWeather } from '~/client'

defineProps<{ weather: WeatherWeather }>()
</script>
