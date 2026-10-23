<template>
  <UModal
    v-model:open="open"
    @after:leave="reset"
  >
    <UButton variant="outline">
      New
    </UButton>
    <template #content>
      <div
        v-if="secret"
        class="flex flex-col gap-4 p-6"
      >
        <div>
          <h2 class="text-lg font-semibold">
            Auth Token
          </h2>
          <p class="text-muted text-sm">
            Save this token somewhere safe, you won't be able to see it again
          </p>
        </div>
        <div
          class="focus-visible:outline-primary cursor-pointer overflow-auto rounded-md bg-zinc-800 p-4 transition-colors hover:bg-zinc-700 focus-visible:outline-2 focus-visible:outline-offset-2"
          role="button"
          tabindex="0"
          @click="copySecret"
          @keydown.enter="copySecret"
          @keydown.space.prevent="copySecret"
        >
          <pre>{{ secret }}</pre>
        </div>
        <div>
          <UButton
            type="button"
            color="neutral"
            variant="outline"
            @click="open = false"
          >
            Close
          </UButton>
        </div>
      </div>
      <div
        v-else
        class="flex flex-col gap-6 p-6"
      >
        <div>
          <h2 class="text-lg font-semibold">
            Create token
          </h2>
          <p class="text-muted text-sm">
            Give this token a name that helps you identify it later.
          </p>
        </div>

        <UForm
          :state="form"
          class="flex flex-col gap-6"
          @submit="submit"
        >
          <UFormField
            label="Token name"
            name="name"
            required
          >
            <UInput
              v-model="form.name"
              placeholder="e.g. Home Assistant"
              autocomplete="off"
              autofocus
            />
          </UFormField>

          <div class="flex gap-3">
            <UButton
              type="button"
              color="neutral"
              variant="outline"
              @click="open = false"
            >
              Cancel
            </UButton>
            <UButton
              type="submit"
              variant="outline"
            >
              Create token
            </UButton>
          </div>
        </UForm>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import { postV1AuthTokens } from '~/client'

const open = ref(false)
const form = reactive({
  name: ''
})
const toast = useToast()

const secret = ref<string>()

const reset = () => {
  form.name = ''
  secret.value = undefined
}

const emit = defineEmits<{
  created: []
}>()

const submit = async () => {
  try {
    const resp = await postV1AuthTokens({
      body: {
        name: form.name
      }
    })
    if (!resp.data) {
      toast.add({
        title: 'Unable to create token',
        description: 'Please try again.',
        color: 'error'
      })
      return
    }
    secret.value = resp.data.secret
    emit('created')
  } catch {
    toast.add({
      title: 'Unable to create token',
      description: 'Please check your connection and try again.',
      color: 'error'
    })
  }
}

const copySecret = async () => {
  if (!secret.value) return

  try {
    await navigator.clipboard.writeText(secret.value)
    toast.add({
      title: 'Token copied',
      description: 'The token secret has been copied to your clipboard.',
      color: 'success'
    })
  } catch {
    toast.add({
      title: 'Unable to copy token',
      description: 'Please copy the token secret manually.',
      color: 'error'
    })
  }
}
</script>
