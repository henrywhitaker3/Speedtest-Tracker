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
  await expect(page.getByText('Restored 1 results and skipped 2.')).toBeVisible()
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
