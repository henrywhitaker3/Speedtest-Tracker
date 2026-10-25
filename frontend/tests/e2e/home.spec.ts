import { expect, test } from '@playwright/test'

test('loads the dashboard', async ({ page }) => {
  await page.goto('/')

  await expect(page).toHaveTitle('Speedtest-Tracker')
  await expect(page.getByRole('link', { name: 'Speedtest Tracker' })).toBeVisible()
})
