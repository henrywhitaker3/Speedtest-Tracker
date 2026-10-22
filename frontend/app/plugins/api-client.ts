import { client } from '~/client/client.gen'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  client.setConfig({
    baseUrl: config.public.apiUrl
  })
})
