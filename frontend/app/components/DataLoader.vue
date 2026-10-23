<script setup lang="ts">
import type { ResultsResult } from '~/client'
import { getV1SpeedtestLatest } from '~/client/sdk.gen'
import { useRefreshData, useSpeedtestCompleted, useSpeedtestStarted } from '~/events/speedtest'

const { on: onStarted } = useSpeedtestStarted()
const { emit: refreshData } = useRefreshData()
const { emit: onCompleted } = useSpeedtestCompleted()

const latest = ref<ResultsResult>()
let isPolling = true

onStarted(async () => {
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

    if (!latest.value || latest.value.id !== result.data?.id) {
      latest.value = result.data
      refreshData()
      continue
    }

    latest.value = result.data
    if (!result.data?.complete) {
      refreshData()
    } else {
      onCompleted()
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
