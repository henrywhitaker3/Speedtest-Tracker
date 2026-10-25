const baseURL = process.env.PLAYWRIGHT_BASE_URL ?? 'http://localhost:3000'
const timeout = 60_000
const pollInterval = 1_000

export default async function globalSetup() {
  const deadline = Date.now() + timeout

  while (Date.now() < deadline) {
    try {
      const response = await fetch(baseURL)

      if (response.ok) {
        return
      }
    } catch {
      // The application container is still starting.
    }

    await new Promise(resolve => setTimeout(resolve, pollInterval))
  }

  throw new Error(`Timed out waiting for the application at ${baseURL}`)
}
