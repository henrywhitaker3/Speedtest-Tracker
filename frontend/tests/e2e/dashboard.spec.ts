import { expect, test } from '@playwright/test'
import { mockApi } from './fixtures'

test.describe('dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await mockApi(page)
  })

  test('shows the latest result, history, and charts', async ({ page }) => {
    await page.goto('/')

    await expect(page.getByRole('heading', { name: 'Ping' }).first()).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Download' }).first()).toBeVisible()
    await expect(page.getByRole('heading', { name: 'Upload' }).first()).toBeVisible()
    await expect(page.getByLabel('Ping over time')).toBeVisible()
    await expect(page.getByLabel('Bandwidth over time')).toBeVisible()
    await expect(page.getByRole('cell', { name: '12.5 ms' })).toBeVisible()
  })

  test('updates charts when the selected time range changes', async ({ page }) => {
    let requestedPeriod = ''
    await page.route('**/v1/speedtest/timeseries**', async (route) => {
      requestedPeriod = new URL(route.request().url()).searchParams.get('period') ?? ''
      await route.fulfill({ contentType: 'application/json', body: '[]' })
    })

    await page.goto('/')
    await page.getByLabel('Chart time range').click()
    await page.getByRole('option', { name: '7 days' }).click()

    await expect.poll(() => requestedPeriod).toBe('7d')
  })

  test('deletes a result using its concrete id', async ({ page }) => {
    let deletePathname: string | undefined
    await page.route('**/v1/speedtest/**', async (route) => {
      if (route.request().method() === 'DELETE') {
        deletePathname = new URL(route.request().url()).pathname
        await route.fulfill({ status: 200, contentType: 'application/json', body: '{}' })
        return
      }

      await route.fallback()
    })

    await page.goto('/')
    await page.getByLabel('Expand').click()

    await expect.poll(() => deletePathname).toBe('/v1/speedtest/result-1')
  })
})

test('shows spaced weather markers on both charts when weather is available', async ({ page }) => {
  const now = Date.now()
  await mockApi(page, {
    stats: {
      body: Array.from({ length: 100 }, (_, index) => ({
        id: `weather-${index}`,
        created_at: new Date(now - (100 - index) * 60_000).toISOString(),
        ping: 10,
        download: 100_000,
        upload: 50_000,
        weather: { code: { name: 'Clear sky', emoji: '☀️' }, temperature: 20 }
      }))
    }
  })

  await page.goto('/')
  const overlays = page.locator('[data-weather-overlay]')
  await expect(overlays).toHaveCount(2)
  for (const overlay of await overlays.all()) {
    await expect(overlay.locator('span').first()).toBeVisible()
    expect(await overlay.locator('span').count()).toBeLessThan(100)
  }

  await page.getByRole('button', { name: 'Graph settings' }).click()
  const option = page.getByRole('menuitemcheckbox', { name: 'Weather emojis' })
  await expect(option).toHaveAttribute('aria-checked', 'true')
  await option.click()
  await expect(overlays).toHaveCount(0)
  await page.reload()
  await expect(overlays).toHaveCount(0)
  await page.getByRole('button', { name: 'Graph settings' }).click()
  await expect(option).toHaveAttribute('aria-checked', 'false')
  await option.click()
  await expect(overlays).toHaveCount(2)
})
