<template>
  <UApp>
    <main class="bg-default text-default flex min-h-screen items-center justify-center p-6">
      <UCard class="w-full max-w-md">
        <div class="flex flex-col items-center gap-4 py-6 text-center">
          <UIcon
            name="i-lucide-circle-alert"
            class="size-12 text-error"
          />
          <div class="space-y-2">
            <p class="text-muted text-sm font-medium">
              Error {{ error.statusCode || 500 }}
            </p>
            <h1 class="text-2xl font-semibold">
              {{ title }}
            </h1>
            <p class="text-muted">
              {{ description }}
            </p>
          </div>
          <UButton
            color="primary"
            icon="i-lucide-house"
            @click="handleError"
          >
            Back to dashboard
          </UButton>
        </div>
      </UCard>
    </main>
  </UApp>
</template>

<script setup lang="ts">
const props = defineProps<{
  error: {
    statusCode?: number
    statusMessage?: string
    message?: string
  }
}>()

const title = computed(() => {
  return props.error.statusCode === 404 ? 'Page not found' : 'Something went wrong'
})

const description = computed(() => {
  if (props.error.statusCode === 404) {
    return 'The page you requested does not exist or may have moved.'
  }

  return props.error.statusMessage || props.error.message || 'Please try again in a moment.'
})

const handleError = () => clearError({ redirect: '/' })
</script>
