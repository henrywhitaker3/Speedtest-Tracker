<template>
  <UButton
    color="success"
    variant="outline"
    class="cursor-pointer"
    @click="run"
  >
    <UIcon
      class="size-5"
      name="i-ic-round-speed"
    />
    Speedtest
  </UButton>
</template>

<script setup lang="ts">
import { postV1SpeedtestRun } from '~/client/sdk.gen'
import { useSpeedtestStarted } from '~/events/speedtest'

const { emit } = useSpeedtestStarted()

const run = async () => {
  const resp = await postV1SpeedtestRun()
  if (resp.data) {
    emit({
      id: resp.data.id ?? ''
    })
  }
}
</script>
