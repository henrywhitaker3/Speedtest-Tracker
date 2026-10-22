<template>
  <UButton
    color="primary"
    variant="outline"
    class="cursor-pointer"
    :disabled="isRunning"
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
const isRunning = ref(false)

const run = async () => {
  if (isRunning.value) return

  isRunning.value = true

  try {
    const resp = await postV1SpeedtestRun()
    if (resp.data) {
      emit({
        id: resp.data.id ?? ''
      })
    }
  } catch {
    isRunning.value = false
  }
}

onCompleted(() => {
  isRunning.value = false
})
</script>
