<template>
  <UModal
    v-model:open="open"
    :ui="{ content: 'w-[calc(100vw-2rem)] max-w-6xl' }"
  >
    <UButton
      label="List"
      color="info"
      variant="outline"
      icon="i-lucide-menu"
    />

    <template #content>
      <div class="flex min-h-0 flex-1 flex-col gap-6 p-6">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="text-lg font-semibold">
              Speedtest servers
            </h2>
            <p class="text-muted text-sm">
              Available Ookla Speedtest servers.
            </p>
          </div>
          <UButton
            color="neutral"
            variant="ghost"
            icon="i-lucide-refresh-cw"
            aria-label="Refresh servers"
            :loading="isLoading"
            @click="loadServers"
          />
        </div>

        <div
          v-if="isLoading"
          class="min-h-0 flex-1 overflow-auto py-2"
        >
          <LoadingTable :count="8" />
        </div>

        <div
          v-else-if="error"
          class="flex min-h-0 flex-1 flex-col items-start gap-3 overflow-auto py-8"
        >
          <p class="text-error">
            {{ error }}
          </p>
          <UButton
            color="neutral"
            variant="outline"
            @click="loadServers"
          >
            Try again
          </UButton>
        </div>

        <p
          v-else-if="!servers.length"
          class="text-muted min-h-0 flex-1 overflow-auto py-8 text-sm"
        >
          No Speedtest servers are available.
        </p>

        <div
          v-else
          class="min-h-0 flex-1"
        >
          <UTable
            :columns="columns"
            :data="servers"
            sticky
            :ui="{ root: 'h-full' }"
          />
        </div>
      </div>
    </template>
  </UModal>
</template>

<script lang="ts" setup>
import { listServers, type OoklaServer } from '~/client'
import type { TableColumn } from '@nuxt/ui'

const open = ref(false)
const isLoading = ref(false)
const error = ref<string>()
const servers = ref<OoklaServer[]>([])

const columns: TableColumn<OoklaServer>[] = [
  {
    accessorKey: 'id',
    header: 'ID'
  },
  {
    accessorKey: 'name',
    header: 'Name'
  },
  {
    accessorKey: 'location',
    header: 'Location'
  },
  {
    accessorKey: 'country',
    header: 'Country'
  },
  {
    accessorKey: 'host',
    header: 'Host'
  },
  {
    accessorKey: 'port',
    header: 'Port'
  }
]

const loadServers = async () => {
  isLoading.value = true
  error.value = undefined

  try {
    const response = await listServers()

    if (response.error) {
      error.value = 'Unable to load Speedtest servers.'
      return
    }

    servers.value = response.data ?? []
  } catch {
    error.value = 'Unable to load Speedtest servers.'
  } finally {
    isLoading.value = false
  }
}

watch(open, (isOpen) => {
  if (isOpen) {
    void loadServers()
  }
})
</script>
