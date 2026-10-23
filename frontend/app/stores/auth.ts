interface AuthState {
  authenticated: boolean
  groups: string[]
  initialized: boolean
}

interface SessionResponse {
  groups?: string[]
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => {
    return {
      authenticated: false,
      groups: [],
      initialized: false
    }
  },
  actions: {
    async initialize() {
      if (this.initialized) {
        return
      }

      const config = useConfigStore()
      if (!config.config) {
        await config.load()
      }

      if (config.authEnabled()) {
        await this.session()
      }

      this.initialized = true
    },
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
    isAdmin(): boolean {
      const config = useConfigStore()
      const adminGroup = config.config?.auth?.admin_group
      return adminGroup !== undefined && this.groups.includes(adminGroup)
    }
  }
})
