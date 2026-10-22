import { filesize } from 'filesize'

export default function useSpeedFormat(bytesPerSecond: number) {
  const { value, symbol } = filesize(bytesPerSecond, {
    bits: true,
    output: 'object'
  })

  return {
    value,
    unit: `${symbol}/s`
  }
}
