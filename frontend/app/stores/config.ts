import { getV1Config, type ConfigConfigResponse } from '~/client'

interface ConfigState {
  config?: ConfigConfigResponse
}

export const useConfigStore = defineStore('config', {
  state: (): ConfigState => {
    return {}
  },
  actions: {
    async load() {
      const resp = await getV1Config()
      if (!resp.data) {
        throw new Error('load config')
      }
      this.config = resp.data
    },
    authEnabled(): boolean {
      return this.config?.auth?.enabled
    }
  }
})
