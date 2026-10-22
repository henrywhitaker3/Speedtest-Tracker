<template>
  <div
    v-if="!latest"
    class="flex items-center justify-center gap-6"
  >
    <LoadingCard />
    <LoadingCard />
    <LoadingCard />
  </div>
  <div
    v-else
    class="flex flex-col gap-2"
  >
    <div
      class="flex items-center justify-center gap-6"
    >
      <ResultsCard
        variant="ping"
        :result="latest"
      />
      <ResultsCard
        variant="download"
        :result="latest"
      />
      <ResultsCard
        variant="upload"
        :result="latest"
      />
    </div>
    <div
      v-if="latest.created_at"
      class="flex justify-center"
    >
      <p class="text-muted">
        Last run {{ relativeDate(new Date(latest.created_at)) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ResultsResult } from '~/client'
import { getV1SpeedtestLatest } from '~/client/sdk.gen'
import { useRefreshData } from '~/events/speedtest'

const latest = ref<ResultsResult>()

const refresh = async () => {
  const result = await getV1SpeedtestLatest()
  latest.value = result.data
  return result
}

onMounted(() => {
  void refresh()
})

const { on } = useRefreshData()
on(() => {
  void refresh()
})

const { relativeDate } = useDate()
</script>
