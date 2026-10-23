interface AuthState {
  authenticated: boolean
  groups: string[]
}

interface SessionResponse {
  groups?: string[]
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => {
    return {
      authenticated: false,
      groups: []
    }
  },
  actions: {
    async session() {
      const runtimeConfig = useRuntimeConfig()
      const apiBaseURL = runtimeConfig.public.apiUrl || window.location.origin
      const response = await fetch(new URL('/v1/auth/session', apiBaseURL), {
        credentials: 'include'
      })
      if (!response.ok) {
        this.authenticated = false
        this.groups = []
        return
      }

      const session = await response.json() as SessionResponse
      this.authenticated = true
      this.groups = session.groups ?? []
    },
    async login() {
      const runtimeConfig = useRuntimeConfig()
      const apiBaseURL = runtimeConfig.public.apiUrl || window.location.origin
      window.location.assign(new URL('/v1/auth/login', apiBaseURL).toString())
    },
    isAdmin() {
      const config = useConfigStore()
      return this.groups.includes(config.config.auth.admin_group)
    }
  }
})
