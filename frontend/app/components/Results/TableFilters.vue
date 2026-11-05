<template>
  <UCard>
    <div class="flex items-center justify-between gap-4">
      <div>
        <h2 class="text-lg font-semibold">
          Results
        </h2>
      </div>

      <div class="flex w-full flex-1 flex-col">
        <div class="flex justify-end">
          <UDropdownMenu
            :items="columnItems"
            :content="{ align: 'end' }"
          >
            <UButton
              label="Columns"
              color="neutral"
              variant="outline"
              trailing-icon="i-lucide-chevron-down"
            />
          </UDropdownMenu>
        </div>
      </div>
    </div>
  </UCard>
</template>

<script setup lang="ts">
type HideableColumn = {
  id: string
  columnDef: {
    header?: unknown
  }
  getCanHide: () => boolean
  getIsVisible: () => boolean
  toggleVisibility: (visible?: boolean) => void
}

const props = defineProps<{
  table?: {
    tableApi?: {
      getAllColumns: () => HideableColumn[]
    }
  }
}>()

const columnItems = computed(() => {
  return props.table?.tableApi?.getAllColumns()
    .filter(column => column.getCanHide())
    .map(column => ({
      label: typeof column.columnDef.header === 'string' ? column.columnDef.header : column.id,
      type: 'checkbox' as const,
      checked: column.getIsVisible(),
      onUpdateChecked(checked: boolean) {
        column.toggleVisibility(checked)
      },
      onSelect(event: Event) {
        event.preventDefault()
      }
    }))
})
</script>
