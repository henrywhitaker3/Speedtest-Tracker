<template>
  <UCard v-if="!stats">
    <LoadingTable :count="7" />
  </UCard>
  <UCard v-if="!stats">
    <LoadingTable :count="7" />
  </UCard>

  <ResultsPing
    v-if="stats"
    :data="stats"
  />
  <ResultsBandwidth
    v-if="stats"
    :data="stats"
  />
</template>

<script setup lang="ts">
import { getV1SpeedtestTimeseries, type ResultsStat } from '~/client'
import { useRefreshData } from '~/events/speedtest'

const stats = ref<ResultsStat[]>()

const period = ref('1d')

const refresh = async () => {
  const resp = await getV1SpeedtestTimeseries({
    query: {
      period: period.value
    }
  })
  stats.value = resp.data
}

onMounted(() => {
  void refresh()
})

const { on: onRefresh } = useRefreshData()
onRefresh(() => {
  void refresh()
})
</script>
