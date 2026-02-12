import { test, expect } from '@playwright/test'

const routes = ['/', '/jobs', '/contacts', '/career', '/doc']

for (const route of routes) {
  test(`page loads: ${route}`, async ({ page }) => {
    await page.goto(route)
    await expect(page).toHaveURL(new RegExp(`${route.replace('/', '\\/')}$`))
  })
}

test('jobs filter toggle RU/EN updates labels', async ({ page }) => {
  await page.goto('/jobs')

  const ruToggle = page.locator('button:has-text("RU"), a:has-text("RU"), [role="button"]:has-text("RU")').first()
  const enToggle = page.locator('button:has-text("EN"), a:has-text("EN"), [role="button"]:has-text("EN")').first()

  await expect(ruToggle).toBeVisible()
  await ruToggle.click()
  // Проверяем, что появился видимый русский текст (не скрытый мобильный элемент)
  await expect(page.locator('.panel-title:visible, .panel-label:visible').filter({ hasText: /Фильтры|Теги|Ключевые/i }).first()).toBeVisible()

  await expect(enToggle).toBeVisible()
  await enToggle.click()
  // Проверяем, что появился видимый английский текст
  await expect(page.locator('.panel-title:visible, .panel-label:visible').filter({ hasText: /Filters|Tags|Keywords/i }).first()).toBeVisible()
})

test('contacts page shows primary links', async ({ page }) => {
  await page.goto('/contacts')
  await expect(page.locator('a[href*="miraclescapedev"]')).toBeVisible()
  await expect(page.locator('a[href*="t.me/MiracleScape"]')).toBeVisible()
  await expect(page.locator('a[href^="mailto:"]')).toBeVisible()
})

test('career page shows empty roles block', async ({ page }) => {
  await page.goto('/career')
  await expect(page.locator('.career-empty-count')).toBeVisible()
  await expect(page.locator('.career-empty-count')).toHaveText('0')
})
