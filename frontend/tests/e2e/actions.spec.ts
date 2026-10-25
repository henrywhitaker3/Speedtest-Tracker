import { expect, test } from '@playwright/test'
import { mockApi } from './fixtures'

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
