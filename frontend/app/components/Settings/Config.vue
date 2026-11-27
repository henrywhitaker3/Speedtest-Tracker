<template>
  <UCard>
    <template #header>
      <h2 class="text-lg font-semibold">
        Speedtest
      </h2>
    </template>

    <dl class="grid gap-x-8 gap-y-5 sm:grid-cols-2">
      <div>
        <dt class="text-muted text-sm">
          Schedule
        </dt>
        <dd class="mt-1 font-mono text-sm">
          {{ configStore.config?.schedule ?? 'Not configured' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Server ID
        </dt>
        <dd class="mt-1 font-medium">
          {{ serverIds }}
        </dd>
      </div>
    </dl>
  </UCard>

  <UCard v-if="configStore.weatherEnabled()">
    <template #header>
      <h2 class="text-lg font-semibold">
        Weather
      </h2>
    </template>

    <dl class="grid gap-x-8 gap-y-5 sm:grid-cols-2">
      <div>
        <dt class="text-muted text-sm">
          Provider
        </dt>
        <dd class="mt-1 font-medium">
          {{ configStore.config?.weather?.provider || '-' }}
        </dd>
      </div>
      <div v-if="configStore.config?.weather?.historical?.enabled">
        <dt class="text-muted text-sm">
          Historical provider
        </dt>
        <dd class="mt-1 font-medium">
          {{ configStore.config?.weather?.historical?.provider || '-' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Latitude
        </dt>
        <dd class="mt-1 font-medium">
          {{ configStore.config?.weather?.latitude ?? '-' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Longitude
        </dt>
        <dd class="mt-1 font-medium">
          {{ configStore.config?.weather?.longitude ?? '-' }}
        </dd>
      </div>
    </dl>
  </UCard>

  <UCard v-if="configStore.authEnabled()">
    <template #header>
      <h2 class="text-lg font-semibold">
        Authentication
      </h2>
    </template>

    <dl class="grid gap-x-8 gap-y-5 sm:grid-cols-2">
      <div>
        <dt class="text-muted text-sm">
          Status
        </dt>
        <dd class="mt-1 font-medium">
          {{ configStore.config?.auth?.enabled ? 'Enabled' : 'Disabled' }}
        </dd>
      </div>
      <div>
        <dt class="text-muted text-sm">
          Admin group
        </dt>
        <dd class="mt-1 font-medium">
          {{ configStore.config?.auth?.admin_group ?? 'Not configured' }}
        </dd>
      </div>
    </dl>
  </UCard>
</template>

<script setup lang="ts">
const configStore = useConfigStore()

const serverIds = computed(() => {
  const serverIds = configStore.config?.server_id
  return serverIds?.length ? serverIds.join(', ') : '-'
})
</script>
