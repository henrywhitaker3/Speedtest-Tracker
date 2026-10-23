export default defineNuxtRouteMiddleware(async () => {
  const config = useConfigStore()
  const auth = useAuthStore()

  if (!config.authEnabled()) {
    return
  }

  if (!auth.authenticated) {
    return navigateTo('/')
  }
})
