<template>
  <UButton
    v-if="props.variant == 'page' && (!config.authEnabled() || auth.isAdmin())"
    label="Delete"
    icon="i-lucide-trash-2"
    color="error"
    variant="outline"
    @click="deleteResult"
  />
  <TableDelete
    v-if="props.variant == 'table' && (!config.authEnabled() || auth.isAdmin())"
    @click="deleteResult"
  />
</template>

<script setup lang="ts">
import { deleteV1SpeedtestById } from '~/client'
import { useRefreshData } from '~/events/speedtest'

const config = useConfigStore()
const auth = useAuthStore()

const props = defineProps<{
  id: string
  variant: 'table' | 'page'
}>()
const { emit: refresh } = useRefreshData()

const deleteResult = async () => {
  await deleteV1SpeedtestById({
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
