export default defineNuxtRouteMiddleware(async () => {
  const config = useConfigStore()
  const auth = useAuthStore()

  if (!config.config) {
    await config.load()
  }

  if (!config.authEnabled()) {
    return
  }

  await auth.initialize()
})
