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
import { postV1SpeedtestRun } from '~/client'

const client = useApi()
const run = async () => {
  const resp = await postV1SpeedtestRun({ client: client })
  console.log(resp.data?.id)
}
</script>
