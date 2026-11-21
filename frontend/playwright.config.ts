import { defineConfig, devices } from '@playwright/test'

const port = process.env.PLAYWRIGHT_PORT ?? '3000'
const baseURL = process.env.PLAYWRIGHT_BASE_URL ?? `http://localhost:${port}`

export default defineConfig({
  testDir: './tests/e2e',
  workers: 1,
  globalSetup: process.env.PLAYWRIGHT_BASE_URL ? './tests/e2e/global-setup.ts' : undefined,
  webServer: process.env.PLAYWRIGHT_BASE_URL
    ? undefined
    : {
        command: `NUXT_PUBLIC_API_URL=http://localhost:${port} NUXT_PUBLIC_VERSION=v1.2.0 pnpm exec nuxt dev --port ${port}`,
        url: baseURL,
        timeout: 120_000,
        reuseExistingServer: !process.env.CI
      },
  timeout: 30_000,
  expect: {
    timeout: 5_000
  },
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  reporter: process.env.CI ? 'github' : 'list',
  use: {
    baseURL,
    trace: 'on-first-retry'
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] }
    }
  ]
})
