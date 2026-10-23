<template>
  <UButton
    color="primary"
    variant="outline"
    class="cursor-pointer"
    :disabled="disabled"
    :loading="isRunning"
    @click="run"
  >
    <UIcon
      v-if="!isRunning"
      class="size-5"
      name="i-ic-round-speed"
    />
    Speedtest
  </UButton>
</template>

<script setup lang="ts">
import { postV1SpeedtestRun } from '~/client/sdk.gen'
import { useSpeedtestCompleted, useSpeedtestStarted } from '~/events/speedtest'

const { emit } = useSpeedtestStarted()
const { on: onCompleted } = useSpeedtestCompleted()
const toast = useToast()
const isRunning = ref(false)

const config = useConfigStore()
const auth = useAuthStore()

const disabled = computed(() => {
  if (isRunning.value) {
    return true
  }
  if (config.authEnabled()) {
    return !auth.authenticated
  }
  return false
})

const run = async () => {
  if (isRunning.value) return

  isRunning.value = true

  try {
    const resp = await postV1SpeedtestRun()
    if (resp.data) {
      emit({
        id: resp.data.id ?? ''
      })
      return
    }

    if (resp.response?.status === 429) {
      toast.add({
        title: 'Too many speedtests',
        description: 'You can run up to three speedtests per minute. Please try again shortly.',
        color: 'warning'
      })
    }

    isRunning.value = false
  } catch {
    isRunning.value = false
  }
}

onCompleted(() => {
  isRunning.value = false
})
</script>
