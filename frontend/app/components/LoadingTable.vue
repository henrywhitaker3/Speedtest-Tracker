<template>
  <div class="flex flex-col gap-2">
    <USkeleton
      v-for="i in props.count ?? 5"
      :key="i"
      :class="props.size ?? 'h-6'"
      class="w-full"
    />
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  count?: number
  size?: string
}>()
</script>
