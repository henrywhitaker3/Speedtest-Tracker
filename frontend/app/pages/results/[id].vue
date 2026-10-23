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
      <UButton
        to="/"
        color="neutral"
        variant="outline"
        icon="i-lucide-arrow-left"
      >
        All results
      </UButton>
    </div>

    <template v-if="isLoading">
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

      <UCard>
        <template #header>
          <h2 class="text-lg font-semibold">
            Test details
          </h2>
        </template>

        <div
          v-if="result.error"
          class="bg-error/10 text-error mb-6 flex gap-3 rounded-md p-4"
        >
          <UIcon
            name="i-lucide-circle-alert"
            class="size-5 shrink-0"
          />
          <div>
            <p class="font-medium">
              This speedtest did not complete successfully.
            </p>
            <p class="text-sm">
              {{ result.error }}
            </p>
          </div>
        </div>

        <dl class="grid gap-x-8 gap-y-5 sm:grid-cols-2">
          <div>
            <dt class="text-muted text-sm">
              Status
            </dt>
            <dd class="mt-1 font-medium">
              {{ result.complete ? 'Complete' : 'In progress' }}
            </dd>
          </div>
          <div>
            <dt class="text-muted text-sm">
              Server
            </dt>
            <dd class="mt-1 font-medium">
              {{ result.server_name || 'Unknown' }}
            </dd>
          </div>
          <div>
            <dt class="text-muted text-sm">
              Server host
            </dt>
            <dd class="mt-1 font-medium break-all">
              {{ result.server_host || 'Unknown' }}
            </dd>
          </div>
          <div>
            <dt class="text-muted text-sm">
              Server ID
            </dt>
            <dd class="mt-1 font-medium">
              {{ result.server_id ?? 'Unknown' }}
            </dd>
          </div>
          <div v-if="result.url">
            <dt class="text-muted text-sm">
              Ookla Result URL
            </dt>
            <dd class="mt-1 font-medium break-all">
              <a
                :href="result.url"
                target="_blank"
                rel="noopener noreferrer"
                class="text-primary hover:underline"
              >
                {{ result.url }}
              </a>
            </dd>
          </div>
          <div>
            <dt class="text-muted text-sm">
              Result ID
            </dt>
            <dd class="mt-1 font-mono text-sm break-all">
              {{ result.id }}
            </dd>
          </div>
        </dl>
      </UCard>
    </template>
  </div>
</template>

<script setup lang="ts">
import type { ResultsResult } from '~/client'
import { getV1SpeedtestById } from '~/client/sdk.gen'
import { useSpeedtestCompleted } from '~/events/speedtest'

const route = useRoute()
const { formatDatetime } = useDate()

const result = ref<ResultsResult>()
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

watch(() => route.params.id, () => {
  void loadResult()
}, { immediate: true })
</script>
