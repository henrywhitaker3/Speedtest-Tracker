import { expect, test } from '@playwright/test'
import { mockApi } from './fixtures'

test('shows configuration, notifications, and the settings actions', async ({ page }) => {
  await mockApi(page, {
    notifications: {
      body: [
        {
          kind: 'webhook',
          url: 'https://hooks.example.test/speedtest',
          phases: ['pre-run', 'post-run']
        }
      ]
    }
  })

  await page.goto('/settings')

  await expect(page.getByRole('heading', { name: 'Settings' })).toBeVisible()
  await expect(page.getByText('0 * * * *')).toBeVisible()
  await expect(page.getByText('123', { exact: true })).toBeVisible()
  await expect(page.getByText('https://hooks.example.test/speedtest')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Purge' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Restore', exact: true })).toBeVisible()

  await page.getByRole('link', { name: 'Dashboard' }).click()
  await expect(page).toHaveURL('/')
})

test('restores results from a JSON backup', async ({ page }) => {
  await mockApi(page, {
    restore: { body: { restored: 1, skipped: 2 } }
  })

  await page.goto('/settings')
  await page.getByRole('button', { name: 'Restore', exact: true }).click()

  const modal = page.getByRole('dialog')
  await expect(modal).toBeVisible()
  await modal.locator('input[type="file"]').setInputFiles({
    name: 'backup.json',
    mimeType: 'application/json',
    buffer: Buffer.from('[{"id": 1}]')
  })

  const request = page.waitForRequest('**/v1/restore/json')
  await modal.getByRole('button', { name: 'Restore', exact: true }).click()

  expect((await request).postDataJSON()).toEqual([{ id: 1 }])
  await expect(page.getByText('Restored 1 results and skipped 2.', { exact: true })).toBeVisible()
})

test('does not show the settings link when an authenticated user is not an admin', async ({ page }) => {
  await mockApi(page, {
    config: {
      body: {
        auth: { enabled: true, admin_group: 'admins' }
      }
    },
    session: { body: { groups: ['readers'] } }
  })

  await page.goto('/')

  await expect(page.getByRole('link', { name: 'Settings' })).toHaveCount(0)
})

test('shows the settings link when an authenticated user is an admin', async ({ page }) => {
  await mockApi(page, {
    config: {
      body: {
        auth: { enabled: true, admin_group: 'admins' }
      }
    },
    session: { body: { groups: ['admins'] } }
  })

  await page.goto('/')

  await expect(page.getByRole('link', { name: 'Settings' })).toBeVisible()
})

test('loads and refreshes current weather and handles errors', async ({ page }) => {
  await mockApi(page, { config: { body: { auth: { enabled: false }, weather: { enabled: true } } } })
  let requests = 0
  await page.route('**/v1/config/weather', async (route) => {
    requests++
    await route.fulfill({
      status: requests === 2 ? 503 : 200,
      json: requests === 2
        ? { message: 'weather is not configured' }
        : { code: { name: 'Slight rain', emoji: '🌧️' }, temperature: 12.5, precipitation: 1.3, wind_speed: 3.2 }
    })
  })

  await page.goto('/settings')
  await page.getByRole('button', { name: 'Get Weather' }).click()
  const dialog = page.getByRole('dialog')
  await expect(dialog.getByText('Slight rain')).toBeVisible()
  await expect(dialog.getByText('12.5 °C')).toBeVisible()
  await expect(dialog.getByText('1.3 mm')).toBeVisible()
  await expect(dialog.getByText('3.2 m/s')).toBeVisible()
  await dialog.getByRole('button', { name: 'Refresh weather' }).click()
  await expect(dialog.getByText('Unable to load current weather.')).toBeVisible()
  await dialog.getByRole('button', { name: 'Try again' }).click()
  await expect(dialog.getByText('Slight rain')).toBeVisible()
  expect(requests).toBe(3)
})

for (const scenario of [
  { name: 'no existing job', status: 404, body: {}, enabled: true },
  { name: 'pending job', status: 200, body: { Status: 0 }, enabled: false },
  { name: 'active job', status: 200, body: { Status: 1 }, enabled: false },
  { name: 'completed job', status: 200, body: { Status: 2 }, enabled: true },
  { name: 'failed job', status: 200, body: { Status: 3 }, enabled: true },
  { name: 'status error', status: 500, body: {}, enabled: false }
]) {
  test(`weather backfill availability with ${scenario.name}`, async ({ page }) => {
    await mockApi(page, { config: { body: { auth: { enabled: false }, weather: { historical: { enabled: true } } } } })
    let requests = 0
    await page.route('**/v1/config/weather/backfill', async (route) => {
      requests++
      await route.fulfill({ status: scenario.status, json: scenario.body })
    })
    await page.goto('/settings')
    const button = page.getByRole('button', { name: 'Backfill Weather' })
    if (scenario.enabled) {
      await expect(button).toBeEnabled()
      expect(requests).toBe(1)
    } else {
      await expect(page.getByRole('status')).toBeVisible()
      await expect(button).toBeDisabled()
    }
  })
}

test('disables weather backfill when historical weather is not configured', async ({ page }) => {
  await mockApi(page)
  await page.goto('/settings')
  await expect(page.getByRole('button', { name: 'Backfill Weather' })).toBeDisabled()
  await expect(page.getByText('Historical weather is not configured.')).toBeVisible()
})

test('updates backfill status after starting and allows retry after failure', async ({ page }) => {
  await mockApi(page, { config: { body: { auth: { enabled: false }, weather: { historical: { enabled: true } } } } })
  let starts = 0
  let checks = 0
  await page.route('**/v1/config/weather/backfill', async (route) => {
    if (route.request().method() === 'POST') {
      starts++
      await route.fulfill({ status: starts === 1 ? 500 : 200, json: starts === 1 ? {} : { job_id: 'job-1', created: true } })
    } else {
      checks++
      await route.fulfill({ status: checks === 1 ? 404 : 200, json: checks === 1 ? {} : { Status: checks === 2 ? 1 : 2 } })
    }
  })
  await page.goto('/settings')
  const button = page.getByRole('button', { name: 'Backfill Weather' })
  await expect(button).toBeEnabled()
  await button.click()
  await expect(page.getByText('Unable to start weather backfill. Please try again.')).toBeVisible()
  await expect(button).toBeEnabled()
  await button.click()
  await expect(page.getByText('Weather backfill is already in progress.')).toBeVisible()
  await expect(button).toBeDisabled()
  await expect(page.getByText('Weather backfill completed.')).toBeVisible({ timeout: 10_000 })
  await expect(button).toBeEnabled()
})

test('shows configured weather providers and coordinates', async ({ page }) => {
  await mockApi(page, {
    config: { body: { auth: { enabled: false }, weather: { enabled: true, provider: 'openweathermap', latitude: 0, longitude: -0.12, historical: { enabled: true, provider: 'apifreaks' } } } }
  })
  await page.goto('/settings')
  const card = page.locator('div').filter({ has: page.getByRole('heading', { name: 'Weather', exact: true }) }).filter({ has: page.getByText('Latitude', { exact: true }) }).last()
  await expect(card.getByText('openweathermap')).toBeVisible()
  await expect(card.getByText('apifreaks')).toBeVisible()
  await expect(card.getByText('0', { exact: true })).toBeVisible()
  await expect(card.getByText('-0.12', { exact: true })).toBeVisible()
})

test('shows placeholders for missing weather coordinates and hides disabled weather', async ({ page }) => {
  await mockApi(page, { config: { body: { auth: { enabled: false }, weather: { enabled: true, provider: 'openweathermap', latitude: null, longitude: null } } } })
  await page.goto('/settings')
  const card = page.locator('div').filter({ has: page.getByRole('heading', { name: 'Weather', exact: true }) }).filter({ has: page.getByText('Latitude', { exact: true }) }).last()
  await expect(card.getByText('-', { exact: true })).toHaveCount(2)
  await expect(card.getByText('Historical provider')).toHaveCount(0)
  await mockApi(page, { config: { body: { auth: { enabled: false }, weather: { enabled: false } } } })
  await page.reload()
  await expect(page.getByRole('heading', { name: 'Weather', exact: true })).toHaveCount(0)
})
