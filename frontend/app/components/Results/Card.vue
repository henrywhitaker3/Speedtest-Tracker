<template>
  <UCard class="w-full">
    <div class="flex flex-col gap-4">
      <div class="flex items-center gap-2">
        <UIcon
          :name="icon"
          class="size-6"
        />
        <h2 class="text-xl">
          {{ title }}
        </h2>
      </div>
      <div
        v-if="props.result.error"
        class="flex gap-2 text-red-500"
      >
        <UIcon
          name="i-lucide-circle-alert"
          class="size-8 shrink-0"
        />
        <span class="min-w-0 flex-1 truncate">Test didn't run successfully</span>
      </div>
      <PulseLoader
        v-else-if="!props.result.complete"
        size="0.7rem"
        color="#E6ECF4"
      />
      <div
        v-else
        class="flex items-baseline gap-1"
      >
        <h3 class="text-lg">
          {{ value }}
        </h3>
        <p class="text-muted">
          {{ units }}
        </p>
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import type { ResultsResult } from '~/client'
import { PulseLoader } from 'vue-spinner'
import 'vue-spinner/style.css'

const props = defineProps<{
  variant: 'ping' | 'download' | 'upload'
  result: ResultsResult
}>()

const title = computed(() => {
  return props.variant.charAt(0).toUpperCase() + props.variant.slice(1)
})

const icon = computed(() => {
  switch (props.variant) {
    case 'ping':
      return 'i-lucide-audio-waveform'
    case 'download':
      return 'i-lucide-download'
    case 'upload':
      return 'i-lucide-upload'
    default:
      return 'i-lucide-circle-question-mark'
  }
})

const value = computed(() => {
  switch (props.variant) {
    case 'ping':
      return props.result.ping?.toFixed(2)
    case 'download':
      return useSpeedFormat(props.result.download as number).value
    case 'upload':
      return useSpeedFormat(props.result.upload as number).value
    default:
      return ''
  }
})

const units = computed(() => {
  switch (props.variant) {
    case 'ping':
      return 'ms'
    case 'download':
      return useSpeedFormat(props.result.download as number).unit
    case 'upload':
      return useSpeedFormat(props.result.upload as number).unit
    default:
      return ''
  }
})
</script>
