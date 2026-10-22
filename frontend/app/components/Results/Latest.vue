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
    <div class="flex justify-center">
      <p class="text-muted">
        Last run {{ relativeDate(latest.created_at) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { ResultsResult } from '~/client'
import { getV1SpeedtestLatest } from '~/client/sdk.gen'

const latest = ref<ResultsResult>()
onMounted(async () => {
  const result = await getV1SpeedtestLatest()
  latest.value = result.data
})

const { relativeDate } = useDate()
</script>
