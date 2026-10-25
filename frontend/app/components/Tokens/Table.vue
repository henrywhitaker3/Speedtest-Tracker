<template>
  <UCard>
    <LoadingTable
      v-if="!data"
      :count="5"
    />
    <div
      v-else
      class="flex flex-col gap-2"
    >
      <div class="flex justify-between">
        <h2 class="text-lg font-semibold">
          Tokens
        </h2>
        <TokensCreate @created="refresh" />
      </div>
      <UTable
        ref="table"
        :columns="columns"
        :data="data ?? []"
      />
    </div>
  </UCard>
</template>

<script setup lang="ts">
import { deleteToken, listTokens, type AuthToken, type SpeedtestResult } from '~/client'
import type { TableColumn } from '@nuxt/ui'
import TableDelete from '../TableDelete.vue'
import RelativeDate from '../RelativeDate.vue'

const data = ref<AuthToken[]>()

const refresh = async () => {
  const results = await listTokens()
  data.value = results.data
}

onMounted(async () => {
  await refresh()
})

const table = ref()
const columns: TableColumn<SpeedtestResult>[] = [
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) => {
      return row.getValue('name')
    }
  },
  {
    accessorKey: 'created_at',
    header: 'Created At',
    cell: ({ row }) => {
      return h(RelativeDate, {
        date: new Date(row.getValue('created_at'))
      }
      )
    }
  },
  {
    accessorKey: 'id',
    header: '',
    cell: ({ row }) => {
      return h(TableDelete, {
        onClick: async () => {
          await deleteToken({ path: {
            id: row.getValue('id')
          } })
          void refresh()
        }
      })
    }
  }
]
</script>
