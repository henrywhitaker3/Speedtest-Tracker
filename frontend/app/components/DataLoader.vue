<script setup lang="ts">
import type { ResultsResult } from '~/client'
import { getV1SpeedtestLatest } from '~/client/sdk.gen'
import { useRefreshData, useSpeedtestCompleted, useSpeedtestStarted } from '~/events/speedtest'

const { on: onStarted } = useSpeedtestStarted()
const { emit: refreshData } = useRefreshData()
const { emit: onCompleted } = useSpeedtestCompleted()

const latest = ref<ResultsResult>()
const activeResultID = ref<string>()
let isPolling = true

onStarted(({ id }) => {
  activeResultID.value = id
  refreshData()
})

const pollLatest = async () => {
  while (isPolling) {
    const result = await getV1SpeedtestLatest()

    if (!isPolling) {
      return
    }

    if (result.error) {
      await new Promise(resolve => setTimeout(resolve, 10_000))
      continue
    }

    const isNewResult = latest.value?.id !== result.data?.id
    const hasJustCompleted = result.data?.id === activeResultID.value
      && result.data.complete
    latest.value = result.data

    if (hasJustCompleted) {
      activeResultID.value = undefined
      onCompleted()
      refreshData()
    } else if (isNewResult || !result.data?.complete) {
      refreshData()
    }

    await new Promise(resolve => setTimeout(resolve, 5_000))
  }
}

onMounted(() => {
  void pollLatest()
})

onBeforeUnmount(() => {
  isPolling = false
})
</script>

<template>
  <div />
</template>
