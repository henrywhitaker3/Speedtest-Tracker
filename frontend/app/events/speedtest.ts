import { useEventBus, type EventBusKey } from '@vueuse/core'

export type SpeedtestStarted = {
  id: string
}

export type SpeedtestCompleted = boolean

export type RefreshData = boolean

const speedtestStartedKey: EventBusKey<SpeedtestStarted> = Symbol('speedtest:started')
const speedtestCompletedKey: EventBusKey<SpeedtestCompleted> = Symbol('speedtest:completed')
const refreshDataKey: EventBusKey<RefreshData> = Symbol('speedtest:refresh')

export function useSpeedtestStarted() {
  return useEventBus(speedtestStartedKey)
}

export function useSpeedtestCompleted() {
  return useEventBus(speedtestCompletedKey)
}

export function useRefreshData() {
  return useEventBus(refreshDataKey)
}
