<template>
  <div class="flex gap-2">
    <TableLink :to="'/results/' + props.id" />
    <ResultsDeleteResult
      :id="props.id"
      variant="table"
    />
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  id: string
}>()
</script>
