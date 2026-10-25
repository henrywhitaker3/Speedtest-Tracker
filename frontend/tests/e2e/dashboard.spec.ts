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
})
