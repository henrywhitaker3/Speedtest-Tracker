<template>
  <UCard>
    <div class="flex items-center justify-between gap-4">
      <div>
        <h2 class="text-lg font-semibold">
          Stats
        </h2>
      </div>

      <USelect
        v-model="range"
        :items="ranges"
        aria-label="Chart time range"
        class="w-28"
      />
    </div>
  </UCard>
</template>

<script setup lang="ts">
const range = defineModel<PickerPeriod>({ required: true })

const ranges = [
  { label: '24 hours', value: '24h' },
  { label: '7 days', value: '7d' },
  { label: '30 days', value: '30d' }
]
</script>
