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
          {{ configStore.config?.server_id ?? 'Not configured' }}
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
</script>
