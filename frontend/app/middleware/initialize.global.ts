export default defineNuxtRouteMiddleware(async () => {
  const auth = useAuthStore()
  await auth.initialize()
})
