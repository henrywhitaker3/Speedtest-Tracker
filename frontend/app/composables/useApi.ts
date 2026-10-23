import { createClient } from '~/client/client'

export default function useApi() {
  const url = useRuntimeConfig().public.apiUrl
  const client = createClient({
    baseUrl: url,
    credentials: 'include'
  })
  return client
}
