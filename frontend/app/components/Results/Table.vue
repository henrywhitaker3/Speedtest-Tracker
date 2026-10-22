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
import { PulseLoader } from 'vue-spinner'
import { useRefreshData } from '~/events/speedtest'

const data = ref<Paginated>()
const page = ref(1)

const refresh = async () => {
  const results = await getV1Speedtest({
    query: {
      page: page.value
    }
  })
  data.value = results.data
}

onMounted(async () => {
  await refresh()
})

const { on } = useRefreshData()
on(() => {
  void refresh()
})

const { formatDatetime } = useDate()

const table = ref()
const columns: TableColumn<ResultsResult>[] = [
  {
    accessorKey: 'ping',
    header: 'Ping',
    cell: ({ row }) => {
      if (!row.original.complete) {
        return h(PulseLoader, {
          color: '#E6ECF4',
          size: '0.4rem'
        })
      }
      return row.getValue('ping') + ' ms'
    }
  },
  {
    accessorKey: 'download',
    header: 'Download',
    cell: ({ row }) => {
      if (!row.original.complete) {
        return h(PulseLoader, {
          color: '#E6ECF4',
          size: '0.4rem'
        })
      }
      const fmt = useSpeedFormat(row.getValue('download'))
      return `${fmt.value} ${fmt.unit}`
    }
  },
  {
    accessorKey: 'upload',
    header: 'Upload',
    cell: ({ row }) => {
      if (!row.original.complete) {
        return h(PulseLoader, {
          color: '#E6ECF4',
          size: '0.4rem'
        })
      }
      const fmt = useSpeedFormat(row.getValue('upload'))
      return `${fmt.value} ${fmt.unit}`
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
