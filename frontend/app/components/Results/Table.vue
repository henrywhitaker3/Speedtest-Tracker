<template>
  <UCard>
    <LoadingTable
      v-if="!data"
      :count="5"
    />
    <UTable
      v-else
      ref="table"
      :columns="columns"
      :data="data.data ?? []"
    />
  </UCard>
</template>

<script setup lang="ts">
import type { TableColumn } from '@nuxt/ui'
import { getV1Speedtest, type DatabasePaginationGithubComHenrywhitaker3SpeedtestTrackerInternalResultsResult as Paginated, type ResultsResult } from '~/client'
import TableLink from '../TableLink.vue'

const data = ref<Paginated>()
const page = ref(1)
onMounted(async () => {
  const results = await getV1Speedtest({
    query: {
      page: page.value
    }
  })
  data.value = results.data
})

const { formatDatetime } = useDate()

const table = ref()
const columns: TableColumn<ResultsResult>[] = [
  {
    accessorKey: 'ping',
    header: 'Ping',
    cell: ({ row }) => {
      return row.getValue('ping') + ' ms'
    }
  },
  {
    accessorKey: 'download',
    header: 'Download',
    cell: ({ row }) => {
      return row.getValue('download') + ' Mbps'
    }
  },
  {
    accessorKey: 'upload',
    header: 'Upload',
    cell: ({ row }) => {
      return row.getValue('upload') + ' Mbps'
    }
  },
  {
    accessorKey: 'created_at',
    header: 'Date',
    cell: ({ row }) => {
      return formatDatetime(row.getValue('created_at'))
    }
  },
  {
    accessorKey: 'id',
    header: '',
    cell: ({ row }) => {
      return h(
        TableLink,
        { to: `/results/${row.getValue('id')}` }
      )
    }
  }
]
</script>
