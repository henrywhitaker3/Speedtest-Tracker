<template>
  <UCard>
    <template #header>
      <h2 class="text-lg font-semibold">
        Actions
      </h2>
    </template>

    <dl class="grid gap-x-8 gap-y-5 sm:grid-cols-2">
      <div>
        <dt class="text-muted text-sm">
          Purge Errors
        </dt>
        <dd class="mt-2 font-mono text-sm">
          <UButton
            variant="outline"
            color="error"
            label="Purge"
            icon="i-lucide-trash"
            :disabled="purgeInProgress"
            :loading="purgeInProgress"
            @click="purgeResults"
          />
        </dd>
      </div>
      <div>
        <dl class="text-muted-text-sm">
          Restore JSON Backup
        </dl>
        <dd class="mt-2 font-mono text-sm">
          <SettingsRestore />
        </dd>
      </div>
    </dl>
  </UCard>
</template>

<script setup lang="ts">
import { purge } from '~/client'

const toast = useToast()

const purgeInProgress = ref(false)

const purgeResults = async () => {
  purgeInProgress.value = true
  const resp = await purge()
  purgeInProgress.value = false

  if (!resp.data) {
    toast.add({
      title: 'Something went wronf',
      description: 'Errored results could not be pruges, please try again',
      color: 'warning'
    })
    return
  }

  toast.add({
    title: 'Errored results purged',
    color: 'success'
  })
}
</script>
