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
      <UTable
        ref="table"
        v-model:pagination="pagination"
        :columns="columns"
        :data="data.data ?? []"
      />
      <div
        v-if="(data?.total_pages ?? 0) > 1"
        class="border-default flex justify-center border-t pt-4"
      >
        <UPagination
          variant="outline"
          active-variant="outline"
          class="[&_button]:cursor-pointer"
          :default-page="(table?.tableApi?.getState().pagination.pageIndex || 0) + 1"
          :items-per-page="table?.tableApi?.getState().pagination.pageSize || 15"
          :total="data.total_rows"
          @update:page="(p: number) => {
            table?.tableApi?.setPageIndex(p)
            page = p
            void refresh()
          }"
        />
      </div>
    </div>
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

const pagination = ref({
  pageIndex: data.value?.page ?? 1,
  pageSize: data.value?.limit ?? 15
})

const refresh = async () => {
  const results = await getV1Speedtest({
    query: {
      page: page.value,
      limit: pagination.value.pageSize
    }
  })
  data.value = results.data
  if (results.data?.page) {
    page.value = results.data.page
  }
  pagination.value.pageIndex = page.value
  pagination.value.pageSize = data.value?.limit || 15
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
