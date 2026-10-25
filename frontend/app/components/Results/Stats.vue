<template>
  <ResultsRangePicker v-model="period" />

  <UCard v-if="!stats">
    <LoadingTable :count="7" />
  </UCard>
  <UCard v-if="!stats">
    <LoadingTable :count="7" />
  </UCard>

  <ResultsPing
    v-if="stats"
    :data="stats"
    :time-domain="timeDomain"
  />
  <ResultsBandwidth
    v-if="stats"
    :data="stats"
    :time-domain="timeDomain"
  />
</template>

<script setup lang="ts">
import { timeseries, type ResultsStat } from '~/client'
import { useRefreshData } from '~/events/speedtest'
import { useLocalStorage } from '@vueuse/core'

const stats = ref<ResultsStat[]>()

const period = useLocalStorage<PickerPeriod>('speedtest-tracker:chart-range', '24h')
const fetchedAt = ref(new Date())

const rangeDurations = {
  '24h': 24 * 60 * 60 * 1000,
  '7d': 7 * 24 * 60 * 60 * 1000,
  '30d': 30 * 24 * 60 * 60 * 1000
} as const

const timeDomain = computed(() => {
  const rangeEnd = fetchedAt.value.getTime()
  return [new Date(rangeEnd - rangeDurations[period.value]), new Date(rangeEnd)] as const
})

const refresh = async () => {
  fetchedAt.value = new Date()

  const resp = await timeseries({
    query: {
      period: period.value
    }
  })
  stats.value = resp.data
  fetchedAt.value = new Date()
}

onMounted(() => {
  void refresh()
})

watch(period, () => {
  void refresh()
})

const { on: onRefresh } = useRefreshData()
onRefresh(() => {
  void refresh()
})
</script>
