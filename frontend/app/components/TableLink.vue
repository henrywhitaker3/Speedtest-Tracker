<template>
  <UButton
    v-if="props.to === '#'"
    :size="props.size"
    :color="props.color"
    :icon="props.icon"
    :variant="props.variant"
    class="cursor-pointer"
    @click="emits('click')"
  />
  <NuxtLink
    v-else
    :to="props.to"
    @click="emits('click')"
  >
    <UButton
      :size="props.size"
      :color="props.color"
      :icon="props.icon"
      :variant="props.variant"
      class="cursor-pointer"
    />
  </NuxtLink>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  to: string
  variant?: 'subtle' | 'link' | 'solid' | 'outline' | 'soft' | 'ghost'
  color?: 'info' | 'warning' | 'error' | 'success' | 'secondary' | 'neutral' | 'primary'
  icon?: string
  size?: 'sm' | 'xs' | 'md' | 'lg' | 'xl'
}>(), {
  variant: 'subtle',
  color: 'neutral',
  icon: 'i-lucide-external-link',
  size: 'sm'
})

const emits = defineEmits<{
  click: []
}>()
</script>
