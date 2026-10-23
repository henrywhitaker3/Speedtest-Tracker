<template>
  <UButton
    icon="i-lucide-trash-2"
    :square="true"
    aria-label="Expand"
    :size="props.size ?? 'sm'"
    color="error"
    variant="subtle"
  />
</template>

<script setup lang="ts">
const props = defineProps<{
  size?: 'sm' | 'md' | 'lg' | 'xl'
}>()
</script>
