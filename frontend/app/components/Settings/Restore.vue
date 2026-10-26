<template>
  <UModal
    v-model:open="open"
    @after:leave="reset"
  >
    <UButton
      variant="outline"
      icon="i-lucide-upload"
      color="info"
    >
      Restore
    </UButton>

    <template #content>
      <div class="flex flex-col gap-6 p-6">
        <div>
          <h2 class="text-lg font-semibold">
            Restore backup
          </h2>
          <p class="text-muted text-sm">
            Upload a v1 Speedtest-Tracker JSON backup to restore its results.
          </p>
        </div>

        <UFileUpload
          v-model="file"
          accept="application/json,.json"
          label="Backup file"
          description="Choose a JSON backup file"
          :disabled="restoreInProgress"
        />

        <p
          v-if="error"
          class="text-error text-sm"
        >
          {{ error }}
        </p>

        <div class="flex gap-3">
          <UButton
            type="button"
            color="neutral"
            variant="outline"
            :disabled="restoreInProgress"
            @click="open = false"
          >
            Cancel
          </UButton>
          <UButton
            type="button"
            :disabled="!file || restoreInProgress"
            :loading="restoreInProgress"
            @click="restore"
          >
            Restore
          </UButton>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { restoreJson, type RestoreLegacyResult } from '~/client'

const toast = useToast()
const open = ref(false)
const file = ref<File | null>(null)
const error = ref<string>()
const restoreInProgress = ref(false)

const reset = () => {
  file.value = null
  error.value = undefined
}

const restore = async () => {
  if (!file.value) return

  error.value = undefined

  let results: RestoreLegacyResult[]
  try {
    const backup = JSON.parse(await file.value.text()) as unknown
    if (!Array.isArray(backup)) {
      error.value = 'The backup must contain a JSON array of results.'
      return
    }
    results = backup as RestoreLegacyResult[]
  } catch {
    error.value = 'The selected file is not valid JSON.'
    return
  }

  restoreInProgress.value = true
  try {
    const response = await restoreJson({ body: results })
    if (!response.data) {
      error.value = 'Unable to restore this backup. Please try again.'
      return
    }

    toast.add({
      title: 'Backup restored',
      description: `Restored ${response.data.restored ?? 0} results and skipped ${response.data.skipped ?? 0}.`,
      color: 'success'
    })
    open.value = false
  } catch {
    error.value = 'Unable to restore this backup. Please check your connection and try again.'
  } finally {
    restoreInProgress.value = false
  }
}
</script>
