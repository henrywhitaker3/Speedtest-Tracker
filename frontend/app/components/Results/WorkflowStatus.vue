<template>
  <UCard>
    <template #header>
      <h2 class="text-lg font-semibold">
        Timeline
      </h2>
    </template>

    <UTimeline :items="timeline" />
  </UCard>
</template>

<script setup lang="ts">
import type { TimelineItem } from '@nuxt/ui'
import type { SpeedtestWorkflowStatus, SpeedtestWorkflowStep } from '~/client'

const props = defineProps<{
  status: SpeedtestWorkflowStatus
}>()

const buildTimeline = (): TimelineItem[] => {
  const out: TimelineItem[] = []
  props.status.steps?.forEach((s) => {
    if (s.name && showStep(s.name)) {
      out.push(buildTimelineItem(s))
    }
  })
  return out
}

const { formatDatetime } = useDate()

const buildTimelineItem = (step: SpeedtestWorkflowStep): TimelineItem => {
  return {
    title: prettyName(step.name),
    date: formatDatetime(new Date(step.started_at as string)),
    icon: stepIcon(step.status),
    description: description(step.name)
  }
}

const showStep = (name: string): boolean => {
  switch (name) {
    case 'healthcheck':
      return false
    default:
      return true
  }
}

const prettyName = (name?: string) => {
  switch (name) {
    case 'start healthcheck':
      return 'Start Healthcheck'
    case 'run test':
      return 'Run Test'
    case 'store result':
      return 'Store Result'
    case 'finish healthcheck':
      return 'Finish Healthcheck'
    default:
      return name
  }
}

const description = (name?: string) => {
  switch (name) {
    case 'start healthcheck':
      return 'Sending the start signal to healthchecks.io'
    case 'run test':
      return 'Run the speedtest'
    case 'store result':
      return 'Store the speedtest results'
    case 'finish healthcheck':
      return 'Sending the finish signal the healthchecks.io'
    default:
      return undefined
  }
}

const stepIcon = (status?: string): string | undefined => {
  switch (status) {
    case 'pending':
      return 'i-lucide-hourglass'
    case 'running':
      return 'i-codicon-session-in-progress'
    case 'completed':
      return 'i-akar-icons-check'
    case 'failed':
      return 'i-fluent-mdl2-error'
    case 'skipped':
      return 'i-codicon-skip'
    default:
      return undefined
  }
}

const timeline = computed(buildTimeline)
</script>
