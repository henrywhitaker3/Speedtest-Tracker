import { expect, test } from '@playwright/test'
import { mockApi, result } from './fixtures'

test('shows a completed result with its details and workflow', async ({ page }) => {
  await mockApi(page)

  await page.goto(`/results/${result.id}`)

  await expect(page.getByRole('heading', { name: 'Speedtest result' })).toBeVisible()
  await expect(page.getByText('Example ISP')).toBeVisible()
  await expect(page.getByText('speedtest.example.test')).toBeVisible()
  await expect(page.getByRole('link', { name: result.url })).toHaveAttribute('href', result.url)
  await expect(page.getByRole('heading', { name: 'Timeline' })).toBeVisible()
  await expect(page.getByText('Run Test')).toBeVisible()

  await page.getByRole('button', { name: 'Raw result' }).click()
  await expect(page.getByText('"result-1"').last()).toBeVisible()
})

test('shows a clear error for a missing result', async ({ page }) => {
  await mockApi(page, { result: { status: 404, body: { message: 'not found' } } })

  await page.goto('/results/missing-result')

  await expect(page.getByRole('heading', { name: 'Unable to load this result' })).toBeVisible()
  await expect(page.getByText('This speedtest result could not be found.')).toBeVisible()
})

test('deletes a result and returns to the dashboard', async ({ page }) => {
  let deleteRequested = false
  await mockApi(page)
  await page.route(`**/v1/speedtest/${result.id}`, async (route) => {
    if (route.request().method() === 'DELETE') {
      deleteRequested = true
      await route.fulfill({ status: 200, contentType: 'application/json', body: '{}' })
      return
    }

    await route.fallback()
  })

  await page.goto(`/results/${result.id}`)
  await page.getByRole('button', { name: 'Delete' }).click()

  await expect.poll(() => deleteRequested).toBe(true)
  await expect(page).toHaveURL('/')
})
