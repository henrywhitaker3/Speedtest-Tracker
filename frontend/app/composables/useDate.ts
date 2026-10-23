import { useTimeAgo } from '@vueuse/core'

export default function useDate() {
  const formatDatetime = (date: Date): string => {
    const d = new Date(date)
    return (
      d.getFullYear()
      + '-'
      + zeroPad(d.getMonth() + 1)
      + '-'
      + zeroPad(d.getDate())
      + ' '
      + zeroPad(d.getHours())
      + ':'
      + zeroPad(d.getMinutes())
      + ':'
      + zeroPad(d.getSeconds())
    )
  }

  const formatTimestamp = (value: Date | number | string) => {
    const date = value instanceof Date ? value : new Date(value)

    return new Intl.DateTimeFormat('en-GB', {
      day: 'numeric',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit'
    }).format(date)
  }

  const zeroPad = (input: string | number) => {
    return ('0' + input).slice(-2)
  }

  const relativeDate = (d: Date): string => {
    return useTimeAgo(d).value
  }

  return {
    relativeDate,
    formatDatetime,
    formatTimestamp
  }
}
