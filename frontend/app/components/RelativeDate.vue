<template>
  <p
    v-if="isYesterdayOrOlder"
    class="inline-block"
  >
    {{ formatDatetime(date) }}
  </p>
  <UTooltip
    v-else
    :text="formatDatetime(date)"
    :delay-duration="200"
  >
    <p class="inline-block">
      {{ relativeDate(date) }}
    </p>
  </UTooltip>
</template>

<script setup lang="ts">
const props = defineProps<{
  date: Date
}>()

const { formatDatetime, relativeDate } = useDate()

const dayFormatter = new Intl.DateTimeFormat('en-CA', {
  timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
  year: 'numeric',
  month: '2-digit',
  day: '2-digit'
})

const dayKey = (date: Date) => {
  const parts = dayFormatter.formatToParts(date)
  const getPart = (type: Intl.DateTimeFormatPartTypes) => parts.find(part => part.type === type)?.value

  return `${getPart('year')}-${getPart('month')}-${getPart('day')}`
}

const isYesterdayOrOlder = computed(() => {
  return dayKey(props.date) < dayKey(new Date())
})
</script>
