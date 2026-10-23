export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()

  if (!auth.authenticated) {
    return navigateTo('/')
  }
})
