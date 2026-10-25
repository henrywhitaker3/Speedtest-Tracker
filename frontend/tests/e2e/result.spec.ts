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

test('does not show the timeline when the result has no workflow status', async ({ page }) => {
  await mockApi(page, {
    result: { body: { ...result, status: undefined } }
  })

  await page.goto(`/results/${result.id}`)

  await expect(page.getByRole('heading', { name: 'Test details' })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Timeline' })).toHaveCount(0)
})

test('shows the raw test result in a collapsible section', async ({ page }) => {
  await mockApi(page)

  await page.goto(`/results/${result.id}`)
  await page.getByRole('button', { name: 'Raw result' }).click()

  await expect(page.getByText('"result-1"').last()).toBeVisible()
  await expect(page.locator('pre code.hljs.json')).toBeVisible()
})

test('does not show the raw result section when the result has no raw output', async ({ page }) => {
  await mockApi(page, {
    result: { body: { ...result, raw: undefined } }
  })

  await page.goto(`/results/${result.id}`)

  await expect(page.getByRole('button', { name: 'Raw result' })).toHaveCount(0)
})

test('updates the result details as the page polls an in-progress test', async ({ page }) => {
  let requests = 0
  let releaseUpdate: () => void
  const updateAvailable = new Promise<void>((resolve) => {
    releaseUpdate = resolve
  })

  await mockApi(page)
  await page.route(`**/v1/speedtest/${result.id}`, async (route) => {
    if (route.request().method() !== 'GET') {
      await route.fallback()
      return
    }

    requests += 1
    if (requests > 1) {
      await updateAvailable
    }

    const body = requests === 1
      ? { ...result, complete: false, download: null, upload: null, ping: null }
      : result

    await route.fulfill({
      contentType: 'application/json',
      body: JSON.stringify(body)
    })
  })

  await page.goto(`/results/${result.id}`)
  await expect(page.getByText('In progress')).toBeVisible()

  await expect.poll(() => requests, { timeout: 10_000 }).toBeGreaterThanOrEqual(2)
  releaseUpdate!()
  await expect(page.getByText('Complete')).toBeVisible()
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
