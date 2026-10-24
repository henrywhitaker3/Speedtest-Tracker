<template>
  <div class="flex flex-col gap-6">
    <div class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold">
          Speedtest result
        </h1>
        <p
          v-if="result?.created_at"
          class="text-muted text-sm"
        >
          {{ formatDatetime(new Date(result.created_at)) }}
        </p>
      </div>

      <div class="flex gap-2">
        <UButton
          to="/"
          color="neutral"
          variant="outline"
          icon="i-lucide-arrow-left"
        >
          All results
        </UButton>
        <ResultsDeleteResult
          :id="String(route.params.id)"
          variant="page"
        />
      </div>
    </div>

    <template v-if="isLoading && !result">
      <div class="flex justify-center gap-6">
        <LoadingCard />
        <LoadingCard />
        <LoadingCard />
      </div>
    </template>

    <UCard
      v-else-if="error"
      class="mx-auto w-full max-w-xl"
    >
      <div class="flex flex-col items-center gap-3 py-6 text-center">
        <UIcon
          name="i-lucide-circle-alert"
          class="text-error size-10"
        />
        <div>
          <h2 class="text-lg font-semibold">
            Unable to load this result
          </h2>
          <p class="text-muted text-sm">
            {{ error }}
          </p>
        </div>
      </div>
    </UCard>

    <template v-else-if="result">
      <div class="flex flex-col justify-center gap-6 md:flex-row">
        <ResultsCard
          variant="ping"
          :result="result"
        />
        <ResultsCard
          variant="download"
          :result="result"
        />
        <ResultsCard
          variant="upload"
          :result="result"
        />
      </div>

      <ResultsDetails :result="result" />

      <ResultsWorkflowStatus
        v-if="result.status"
        :status="result.status"
      />
    </template>
  </div>
</template>

<script setup lang="ts">
import type { SpeedtestGetResponse } from '~/client'
import { getV1SpeedtestById } from '~/client/sdk.gen'
import { useSpeedtestCompleted } from '~/events/speedtest'
import { useIntervalFn } from '@vueuse/core'

const route = useRoute()
const { formatDatetime } = useDate()

const result = ref<SpeedtestGetResponse>()
const error = ref<string>()
const isLoading = ref(true)

const loadResult = async () => {
  isLoading.value = true
  error.value = undefined

  const response = await getV1SpeedtestById({
    path: {
      id: String(route.params.id)
    }
  })

  if (response.data) {
    result.value = response.data
  } else if (response.response?.status === 404) {
    error.value = 'This speedtest result could not be found.'
  } else {
    error.value = 'Please check your connection and try again.'
  }

  isLoading.value = false
}

const { on: onCompleted } = useSpeedtestCompleted()
onCompleted(() => {
  void loadResult()
})

const { pause: pausePolling, resume: resumePolling } = useIntervalFn(() => {
  if (!isLoading.value) {
    void loadResult()
  }
}, 3_000, { immediate: false })

watch(() => result.value?.complete, (complete) => {
  if (complete === true) {
    pausePolling()
  } else {
    resumePolling()
  }
})

watch(() => route.params.id, () => {
  void loadResult()
}, { immediate: true })
</script>
