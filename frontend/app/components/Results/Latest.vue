<template>
  <div
    v-if="!latest"
    class="flex flex-col items-center justify-center gap-6 sm:flex-row"
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
      class="flex flex-col items-center justify-center gap-6 sm:flex-row"
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
    <!-- <div -->
    <!--   v-if="latest.created_at" -->
    <!--   class="flex justify-center" -->
    <!-- > -->
    <!--   <p class="text-muted text-sm"> -->
    <!--     Last run {{ relativeDate(new Date(latest.created_at)) }} -->
    <!--   </p> -->
    <!-- </div> -->
  </div>
</template>

<script setup lang="ts">
import { latestResult, type SpeedtestResult } from '~/client'
import { useRefreshData } from '~/events/speedtest'

const latest = ref<SpeedtestResult>()

const refresh = async () => {
  const result = await latestResult()
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
</script>
