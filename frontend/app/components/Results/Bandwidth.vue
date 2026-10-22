<template>
  <UCard class="w-full">
    <div class="mb-6 flex items-start justify-between gap-4">
      <div>
        <h2 class="text-lg font-semibold">
          Bandwidth
        </h2>
      </div>
    </div>

    <UChart
      class="h-64 w-full"
      :definition="definition"
      :height="256"
      aria-label="Bandwidth over time"
      aria-description="A line graph showing download/upload bandwidth in mbps"
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
          <span class="ml-auto font-medium">{{ formatMbitPerSecond(Number(row.value)) }}</span>
        </div>
      </template>
    </UChart>
  </UCard>
</template>

<script setup lang="ts">
import { colorLegend, colorLegendItems, defineChart, lineY } from '@tanstack/charts'
import { crosshair } from '@tanstack/charts/crosshair'
import { d3Curve } from '@tanstack/charts/d3/shape'
import { scaleLinear } from '@tanstack/charts/scales/linear'
import { scaleOrdinal } from '@tanstack/charts/scales/ordinal'
import type { ResultsStat } from '~/client'
import { scaleTime } from 'd3-scale'
import { curveMonotoneX } from 'd3-shape'

const props = defineProps<{
  data: ResultsStat[]
}>()

const bitsPerSecondToMbps = (bitsPerSecond: number) => bitsPerSecond / 1_000_000

const chartData = computed(() => {
  return props.data.flatMap((sample) => {
    if (!sample.created_at) {
      return []
    }

    const createdAt = new Date(sample.created_at)

    if (Number.isNaN(createdAt.getTime())) {
      return []
    }

    return [
      {
        created_at: createdAt,
        direction: 'Download' as const,
        megabitsPerSecond: bitsPerSecondToMbps(sample.download ?? 0)
      },
      {
        created_at: createdAt,
        direction: 'Upload' as const,
        megabitsPerSecond: bitsPerSecondToMbps(sample.upload ?? 0)
      }
    ]
  })
})

const formatMbitPerSecond = (value: number) => {
  return `${new Intl.NumberFormat('en-GB', { maximumFractionDigits: 1 }).format(value)} Mbit/s`
}

const formatTimestamp = (value: Date | number | string) => {
  const date = value instanceof Date ? value : new Date(value)

  return new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const definition = computed(() => {
  return defineChart({
    marks: [
      lineY(chartData.value, {
        x: 'created_at',
        y: 'megabitsPerSecond',
        z: 'direction',
        color: 'direction',
        points: true,
        curve: d3Curve(curveMonotoneX),
        strokeWidth: 2
      }),
      crosshair({ x: true, y: false })
    ],
    scales: {
      x: {
        scale: scaleTime,
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
          label: 'Bandwidth',
          ticks: { format: formatMbitPerSecond }
        }
      }
    },
    color: {
      scale: scaleOrdinal(
        ['Download', 'Upload'],
        ['#fca503', '#3e95cd']
      ),
      legend: colorLegend({
        items: colorLegendItems({ justify: 'center', gap: 16 })
      })
    },
    focus: 'group-x',
    maxFocusDistance: Number.POSITIVE_INFINITY
  })
})
</script>
