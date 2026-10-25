import { expect, test } from '@playwright/test'
import { mockApi, result } from './fixtures'

test('warns when a speedtest request is rate limited', async ({ page }) => {
  await mockApi(page, { run: { status: 429, body: { message: 'rate limited' } } })

  await page.goto('/')
  await page.getByRole('button', { name: 'Speedtest' }).click()

  await expect(page.locator('[data-slot="title"]', { hasText: 'Too many speedtests' })).toBeVisible()
  await expect(page.locator('[data-slot="description"]', {
    hasText: 'You can run up to three speedtests per minute. Please try again shortly.'
  })).toBeVisible()
})

test('shows a sign in button instead of the speedtest action when auth is enabled', async ({ page }) => {
  await mockApi(page, {
    config: {
      body: {
        auth: { enabled: true, admin_group: 'admins' }
      }
    }
  })

  await page.goto('/')

  await expect(page.getByRole('button', { name: 'Sign in' })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Speedtest' })).toHaveCount(0)
})

test('shows a spinner and opens the running result after starting a speedtest', async ({ page }) => {
  await mockApi(page)
  await page.route('**/v1/speedtest/latest', async (route) => {
    await route.fulfill({
      contentType: 'application/json',
      body: JSON.stringify({ ...result, complete: false })
    })
  })

  await page.goto('/')
  const speedtest = page.getByRole('button', { name: 'Speedtest' })
  await speedtest.click()

  await expect(speedtest).toBeEnabled()
  await expect(speedtest.locator('.animate-spin')).toBeVisible()

  await speedtest.click()
  await expect(page).toHaveURL('/results/result-1')
})

test('stops loading when the started result first appears as complete', async ({ page }) => {
  let latestRequests = 0
  let testStarted = false
  const completedResult = { ...result, id: 'running-result' }

  await mockApi(page)
  await page.route('**/v1/speedtest/latest', async (route) => {
    latestRequests += 1
    await route.fulfill({
      contentType: 'application/json',
      body: JSON.stringify(testStarted ? completedResult : result)
    })
  })
  await page.route('**/v1/speedtest/run', async (route) => {
    testStarted = true
    await route.fulfill({
      contentType: 'application/json',
      body: JSON.stringify({ id: completedResult.id })
    })
  })

  await page.goto('/')
  await expect.poll(() => latestRequests).toBeGreaterThan(0)

  const speedtest = page.getByRole('button', { name: 'Speedtest' })
  await speedtest.click()
  await expect(speedtest.locator('.animate-spin')).toBeVisible()

  await expect(speedtest.locator('.animate-spin')).toHaveCount(0, { timeout: 8_000 })
})
