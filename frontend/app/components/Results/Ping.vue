<template>
  <UCard class="w-full">
    <div class="mb-6 flex items-start justify-between gap-4">
      <div>
        <h2 class="text-lg font-semibold">
          Ping
        </h2>
      </div>
    </div>

    <UChart
      class="h-64 w-full"
      :definition="definition"
      :height="256"
      aria-label="Ping over time"
      aria-description="A line graph showing ping results in milliseconds."
    >
      <template #tooltip="{ title, rows }">
        <p
          v-if="title"
          class="text-muted mb-1 text-xs"
        >
          {{ title }}
        </p>
        <div
          v-for="row in rows"
          :key="row.label"
          class="flex items-center gap-2 text-sm"
        >
          <span
            class="size-2 rounded-full"
            :style="{ backgroundColor: row.color }"
          />
          <span class="text-muted">{{ row.label }}</span>
          <span class="ml-auto font-medium">{{ formatPing(Number(row.value)) }}</span>
        </div>
      </template>
    </UChart>
  </UCard>
</template>

<script setup lang="ts">
import { defineChart, lineY } from '@tanstack/charts'
import { crosshair } from '@tanstack/charts/crosshair'
import { d3Curve } from '@tanstack/charts/d3/shape'
import { scaleLinear } from '@tanstack/charts/scales/linear'
import { scaleTime } from 'd3-scale'
import { curveMonotoneX } from 'd3-shape'

type PingSample = {
  created_at: string
  ping: number
}

const samples: PingSample[] = [
  { created_at: '2026-10-20T08:00:00Z', ping: 21 },
  { created_at: '2026-10-20T14:00:00Z', ping: 18 },
  { created_at: '2026-10-20T20:00:00Z', ping: 24 },
  { created_at: '2026-10-21T08:00:00Z', ping: 20 },
  { created_at: '2026-10-21T16:00:00Z', ping: 28 },
  { created_at: '2026-10-22T09:00:00Z', ping: 22 },
  { created_at: '2026-10-22T18:00:00Z', ping: 19 }
]

const chartData = samples.map(sample => ({
  ...sample,
  created_at: new Date(sample.created_at)
}))

const formatPing = (value: number) => `${value} ms`

const formatTimestamp = (value: Date | number | string) => {
  const date = value instanceof Date ? value : new Date(value)

  return new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const definition = defineChart({
  marks: [
    lineY(chartData, {
      x: 'created_at',
      y: 'ping',
      points: true,
      curve: d3Curve(curveMonotoneX),
      stroke: 'var(--ui-primary)',
      strokeWidth: 2
    }),
    crosshair({ x: true, y: false })
  ],
  scales: {
    x: {
      scale: scaleTime,
      nice: true,
      axis: {
        label: 'Time',
        ticks: { format: formatTimestamp }
      }
    },
    y: {
      scale: scaleLinear,
      nice: true,
      grid: true,
      axis: {
        label: 'Ping',
        ticks: { format: formatPing }
      }
    }
  },
  focus: 'group-x',
  maxFocusDistance: Number.POSITIVE_INFINITY
})
</script>
