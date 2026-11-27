<template>
  <div
    ref="overlay"
    data-weather-overlay
    class="pointer-events-none absolute inset-0"
    aria-hidden="true"
  >
    <span
      v-for="marker in markers"
      :key="marker.key"
      class="absolute bottom-full mb-1 -translate-x-1/2 text-base leading-none"
      :style="{ left: `${marker.position}px` }"
    >{{ marker.emoji }}</span>
  </div>
</template>

<script setup lang="ts">
import { useElementSize } from '@vueuse/core'
import type { ResultsStat } from '~/client'

const props = defineProps<{
  data: ResultsStat[]
  timeDomain?: readonly [Date, Date]
}>()

const overlay = ref<HTMLElement | null>(null)
const { width } = useElementSize(overlay)

const markers = computed(() => {
  const samples = props.data.flatMap((sample, index) => {
    const time = sample.created_at ? new Date(sample.created_at).getTime() : NaN
    const emoji = sample.weather?.code?.emoji
    return Number.isFinite(time) && emoji ? [{ time, emoji, key: `${time}-${index}` }] : []
  }).sort((a, b) => a.time - b.time)

  if (!samples.length || width.value <= 0) return []

  const allTimes = props.data.flatMap((sample) => {
    const time = sample.created_at ? new Date(sample.created_at).getTime() : NaN
    return Number.isFinite(time) ? [time] : []
  })
  const start = props.timeDomain?.[0].getTime() ?? Math.min(...allTimes)
  const end = props.timeDomain?.[1].getTime() ?? Math.max(...allTimes)
  const spacing = 28
  let lastPosition = -spacing

  return samples.flatMap((sample) => {
    if (sample.time < start || sample.time > end) return []
    const position = end === start ? width.value / 2 : (sample.time - start) / (end - start) * width.value
    const clamped = Math.max(8, Math.min(width.value - 8, position))
    if (clamped - lastPosition < spacing) return []
    lastPosition = clamped
    return [{ ...sample, position: clamped }]
  })
})
</script>
