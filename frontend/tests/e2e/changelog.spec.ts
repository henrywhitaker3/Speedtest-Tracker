import { expect, test, type Page } from '@playwright/test'
import { mockApi } from './fixtures'

const latestReleaseUrl = 'https://api.github.com/repos/henrywhitaker3/Speedtest-Tracker/releases/latest'

async function mockLatestRelease(page: Page, tagName: string) {
  await page.route(latestReleaseUrl, async (route) => {
    await route.fulfill({
      contentType: 'application/json',
      body: JSON.stringify({ tag_name: tagName })
    })
  })
}

test.describe('app changelog', () => {
  test.beforeEach(async ({ page }) => {
    await mockApi(page)
  })

  test('shows an update link for a newer release', async ({ page }) => {
    await mockLatestRelease(page, 'v1.3.0')

    await page.goto('/')

    await expect(page.getByRole('link', { name: 'View version v1.3.0 release' })).toHaveAttribute(
      'href',
      'https://github.com/henrywhitaker3/Speedtest-Tracker/releases/tag/v1.3.0'
    )
  })

  test('hides the update link when the latest release is not newer', async ({ page }) => {
    await mockLatestRelease(page, 'v1.2.0')

    await page.goto('/')

    await expect(page.getByRole('link', { name: /View version .* release/ })).toHaveCount(0)
  })
})
