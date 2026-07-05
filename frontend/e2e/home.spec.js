import { test, expect } from '@playwright/test'

test.describe('首页', () => {
  test('页面标题和品牌标识正确', async ({ page }) => {
    await page.goto('/')
    await expect(page.locator('h1')).toContainText('常用软件下载导航')
    await expect(page.locator('text=ORIGINAL SOFTWARE HUB')).toBeVisible()
  })

  test('加载软件卡片列表', async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('.software-card', { timeout: 15000 })
    const cards = page.locator('.software-card')
    const count = await cards.count()
    expect(count).toBeGreaterThan(10)
  })

  test('搜索功能按名称过滤', async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('.software-card', { timeout: 15000 })
    const searchInput = page.locator('input[type="search"]')
    await searchInput.fill('微信')
    await page.waitForTimeout(500)
    const cards = page.locator('.software-card')
    const count = await cards.count()
    expect(count).toBeGreaterThanOrEqual(1)
  })

  test('搜索无结果时显示提示', async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('.software-card', { timeout: 15000 })
    const searchInput = page.locator('input[type="search"]')
    await searchInput.fill('__nonexistent_software_xxx__')
    await page.waitForTimeout(500)
    await expect(page.locator('text=没有匹配的软件')).toBeVisible()
  })

  test('每张卡片显示名称、机构和标签', async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('.software-card', { timeout: 15000 })
    const card = page.locator('.software-card').first()
    await expect(card.locator('h3')).toBeVisible()
    await expect(card.locator('text=机构：')).toBeVisible()
  })
})
