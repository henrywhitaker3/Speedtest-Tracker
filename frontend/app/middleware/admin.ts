export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()

  if (!auth.authenticated || !auth.isAdmin()) {
    return navigateTo('/')
  }
})
