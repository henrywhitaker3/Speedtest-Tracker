<template>
  <UButton
    v-if="props.variant == 'page' && (!config.authEnabled() || auth.isAdmin())"
    label="Delete"
    icon="i-lucide-trash-2"
    color="error"
    variant="outline"
    @click="deleteTestResult"
  />
  <TableDelete
    v-if="props.variant == 'table' && (!config.authEnabled() || auth.isAdmin())"
    @click="deleteTestResult"
  />
</template>

<script setup lang="ts">
import { deleteResult } from '~/client'
import { useRefreshData } from '~/events/speedtest'

const config = useConfigStore()
const auth = useAuthStore()

const props = defineProps<{
  id: string
  variant: 'table' | 'page'
}>()
const { emit: refresh } = useRefreshData()

const deleteTestResult = async () => {
  await deleteResult({
    path: {
      id: props.id
    }
  })
  refresh()
  if (props.variant == 'page') {
    return navigateTo('/')
  }
}
</script>
