import type { Page } from '@playwright/test'

type MockResponse = {
  body?: unknown
  status?: number
}

export type ApiMocks = {
  config?: MockResponse
  latest?: MockResponse
  notifications?: MockResponse
  result?: MockResponse
  results?: MockResponse
  run?: MockResponse
  session?: MockResponse
  stats?: MockResponse
}

export const result = {
  id: 'result-1',
  complete: true,
  created_at: '2026-10-25T10:00:00Z',
  download: 125_000_000,
  upload: 25_000_000,
  ping: 12.5,
  server_name: 'Example ISP',
  server_host: 'speedtest.example.test',
  server_id: 123,
  url: 'https://www.speedtest.net/result/result-1',
  raw: { result: { id: 'result-1' } },
  status: {
    steps: [
      {
        name: 'run test',
        status: 'completed',
        started_at: '2026-10-25T10:00:00Z'
      },
      {
        name: 'store result',
        status: 'completed',
        started_at: '2026-10-25T10:01:00Z'
      }
    ]
  }
}

const responses = {
  config: {
    body: {
      schedule: '0 * * * *',
      server_id: 123,
      auth: { enabled: false }
    }
  },
  latest: { body: result },
  notifications: { body: [] },
  result: { body: result },
  results: {
    body: {
      data: [result],
      page: 1,
      limit: 15,
      total_pages: 1,
      total_rows: 1
    }
  },
  run: { body: { id: result.id } },
  session: { status: 401, body: { message: 'unauthenticated' } },
  stats: {
    body: [
      {
        id: result.id,
        created_at: result.created_at,
        download: result.download,
        upload: result.upload,
        ping: result.ping
      }
    ]
  }
} satisfies Record<keyof ApiMocks, MockResponse>

export async function mockApi(page: Page, mocks: ApiMocks = {}) {
  await page.route('**/v1/**', async (route) => {
    const path = new URL(route.request().url()).pathname
    const response = mocks[responseKey(path, route.request().method())] ?? responseFor(path)

    await route.fulfill({
      status: response.status ?? 200,
      contentType: 'application/json',
      body: JSON.stringify(response.body ?? {})
    })
  })
}

function responseKey(path: string, method: string): keyof ApiMocks | undefined {
  if (path === '/v1/config') return 'config'
  if (path === '/v1/config/notifications') return 'notifications'
  if (path === '/v1/auth/session') return 'session'
  if (path === '/v1/speedtest/latest') return 'latest'
  if (path === '/v1/speedtest/timeseries') return 'stats'
  if (path === '/v1/speedtest/run') return 'run'
  if (path === '/v1/speedtest') return 'results'
  if (path.startsWith('/v1/speedtest/') && method === 'GET') return 'result'
}

function responseFor(path: string): MockResponse {
  const key = responseKey(path, 'GET')

  return key ? responses[key] : { body: {} }
}
