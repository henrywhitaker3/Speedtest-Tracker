<template>
  <UDropdownMenu
    :items="items"
    :content="{ align: 'end' }"
  >
    <UButton
      aria-label="Graph settings"
      color="neutral"
      variant="outline"
      icon="i-lucide-settings-2"
      trailing-icon="i-lucide-chevron-down"
    />
  </UDropdownMenu>
</template>

<script setup lang="ts">
const showWeather = defineModel<boolean>({ required: true })

const items = computed(() => [{
  label: 'Weather emojis',
  type: 'checkbox' as const,
  checked: showWeather.value,
  onUpdateChecked(checked: boolean) {
    showWeather.value = checked
  },
  onSelect(event: Event) {
    event.preventDefault()
  }
}])
</script>
