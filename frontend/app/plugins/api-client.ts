import { client } from '~/client/client.gen'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  client.setConfig({
    baseUrl: config.public.apiUrl,
    // The UI and API may use different origins during development. Fetch only
    // sends cookies cross-origin when this is explicitly enabled.
    credentials: 'include'
  })
})
