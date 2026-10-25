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

  await page.getByRole('link', { name: 'Dashboard' }).click()
  await expect(page).toHaveURL('/')
})
