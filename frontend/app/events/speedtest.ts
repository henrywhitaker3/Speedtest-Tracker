import { useEventBus, type EventBusKey } from '@vueuse/core'

export type SpeedtestStarted = {
  id: string
}

export type RefreshData = boolean

const speedtestStartedKey: EventBusKey<SpeedtestStarted> = Symbol('speedtest:started')
const refreshDataKey: EventBusKey<RefreshData> = Symbol('speedtest:refresh')

export function useSpeedtestStarted() {
  return useEventBus(speedtestStartedKey)
}

export function useRefreshData() {
  return useEventBus(refreshDataKey)
}
