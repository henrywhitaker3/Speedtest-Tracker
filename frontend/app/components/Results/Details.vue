<template>
  <UCard>
    <template #header>
      <h2 class="text-lg font-semibold">
        Test details
      </h2>
    </template>

    <div
      v-if="props.result.error"
      class="bg-error/10 text-error mb-6 flex gap-3 rounded-md p-4"
    >
      <UIcon
        name="i-lucide-circle-alert"
        class="size-5 shrink-0"
      />
      <div>
        <p class="font-medium">
          This speedtest did not complete successfully.
        </p>
        <p class="text-sm">
          {{ props.result.error }}
        </p>
      </div>
    </div>

    <dl class="grid gap-x-8 gap-y-5 sm:grid-cols-2">
      <div>
        <dt class="text-muted text-sm">
          Status
        </dt>
        <dd class="mt-1 font-medium">
          {{ props.result.complete ? 'Complete' : 'In progress' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Server
        </dt>
        <dd class="mt-1 font-medium">
          {{ props.result.server_name || 'Unknown' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Server host
        </dt>
        <dd class="mt-1 font-medium break-all">
          {{ props.result.server_host || 'Unknown' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Server ID
        </dt>
        <dd class="mt-1 font-medium">
          {{ props.result.server_id ?? 'Unknown' }}
        </dd>
      </div>
      <div v-if="props.result.url">
        <dt class="text-muted text-sm">
          Ookla Result URL
        </dt>
        <dd class="mt-1 font-medium break-all">
          <a
            :href="props.result.url"
            target="_blank"
            rel="noopener noreferrer"
            class="text-primary hover:underline"
          >
            {{ props.result.url }}
          </a>
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Result ID
        </dt>
        <dd class="mt-1 font-mono text-sm break-all">
          {{ props.result.id }}
        </dd>
      </div>
    </dl>

    <UCollapsible
      v-if="props.result.raw != null"
      class="border-default mt-6 border-t pt-4"
    >
      <template #default="{ open }">
        <UButton
          color="neutral"
          variant="ghost"
          size="sm"
          class="-ml-2"
        >
          Raw result
          <UIcon
            name="i-lucide-chevron-down"
            class="size-4 transition-transform"
            :class="{ 'rotate-180': open }"
          />
        </UButton>
      </template>

      <template #content>
        <HighlightJs
          language="json"
          :code="rawResult"
          class="mt-3 max-h-96 overflow-auto rounded-md text-xs wrap-break-word [&_code]:whitespace-pre-wrap"
        />
      </template>
    </UCollapsible>
  </UCard>
</template>

<script setup lang="ts">
import hljsVuePlugin from '@highlightjs/vue-plugin'
import hljs from 'highlight.js/lib/core'
import json from 'highlight.js/lib/languages/json'
import 'highlight.js/styles/tokyo-night-dark.min.css'
import type { SpeedtestGetResponse } from '~/client'

hljs.registerLanguage('json', json)

const props = defineProps<{
  result: SpeedtestGetResponse
}>()

const HighlightJs = hljsVuePlugin.component
const rawResult = computed(() => JSON.stringify(props.result.raw ?? {}, null, 2))
</script>
