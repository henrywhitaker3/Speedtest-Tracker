<template>
  <div
    v-if="weather"
    class="mt-2 border-t border-default pt-2 text-sm"
  >
    <p v-if="weather.code">
      {{ weather.code.emoji }} {{ weather.code.name }}
    </p>
    <p v-if="weather.temperature != null">
      Temperature: {{ weather.temperature }} °C
    </p>
    <p v-if="weather.precipitation != null">
      Precipitation: {{ weather.precipitation }} mm
    </p>
    <p v-if="weather.wind_speed != null">
      Wind speed: {{ weather.wind_speed }} m/s
    </p>
  </div>
</template>

<script setup lang="ts">
import type { ResultsStat } from '~/client'

const props = defineProps<{ data: ResultsStat[], time: unknown }>()

const weather = computed(() => {
  const time = props.time instanceof Date ? props.time.getTime() : new Date(String(props.time)).getTime()
  if (!Number.isFinite(time)) return undefined
  return props.data.find(sample => sample.created_at && new Date(sample.created_at).getTime() === time)?.weather
})
</script>
