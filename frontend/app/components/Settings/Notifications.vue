<template>
  <UCard>
    <template #header>
      <h2 class="text-lg font-semibold">
        Notifications
      </h2>
    </template>

    <div
      v-if="isLoading"
      class="text-muted flex items-center gap-3 py-2"
    >
      <UIcon
        name="i-lucide-loader-circle"
        class="size-5 animate-spin"
      />
      Loading notification configuration…
    </div>

    <div
      v-else-if="error"
      class="flex flex-col items-start gap-3 py-2"
    >
      <p class="text-error">
        {{ error }}
      </p>
      <UButton
        color="neutral"
        variant="outline"
        @click="loadNotifications"
      >
        Try again
      </UButton>
    </div>

    <p
      v-else-if="!notifications?.length"
      class="text-muted py-2 text-sm"
    >
      No notification channels are configured.
    </p>

    <div
      v-else
      class="divide-default divide-y"
    >
      <div
        v-for="(notification, index) in notifications"
        :key="`${notification.kind}-${notification.url}-${index}`"
        class="flex flex-col gap-2 py-4 first:pt-0 last:pb-0 sm:flex-row sm:items-center sm:justify-between"
      >
        <div>
          <p class="font-medium capitalize">
            {{ notification.kind ?? 'Unknown' }}
          </p>
          <p class="text-muted font-mono text-sm break-all">
            {{ notification.url }}
          </p>
        </div>
        <div class="flex flex-wrap gap-2">
          <UBadge
            v-for="phase in notification.phases"
            :key="phase"
            color="neutral"
            variant="subtle"
          >
            {{ phase }}
          </UBadge>
        </div>
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import { getNotifications, type ConfigNotification } from '~/client'

const notifications = ref<ConfigNotification[]>()
const error = ref<string>()
const isLoading = ref(true)

const loadNotifications = async () => {
  isLoading.value = true
  error.value = undefined

  try {
    const response = await getNotifications()
    if (response.data) {
      notifications.value = response.data
    } else {
      notifications.value = []
    }
  } catch {
    error.value = 'Unable to load notification configuration.'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  void loadNotifications()
})
</script>
