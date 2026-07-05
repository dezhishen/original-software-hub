import { test, expect } from '@playwright/test'

test.describe('详情页', () => {
  test('从首页点击卡片进入详情页', async ({ page }) => {
    await page.goto('/')
    await page.waitForSelector('.software-card', { timeout: 15000 })
    await page.locator('.software-card').first().click()
    await page.waitForTimeout(3000)
    await expect(page.locator('text=版本信息').first()).toBeVisible({ timeout: 15000 })
  })

  test('直接访问已知软件详情', async ({ page }) => {
    await page.goto('/#/software/7zip')
    await page.waitForTimeout(3000)
    await expect(page.locator('text=7-Zip').first()).toBeVisible({ timeout: 15000 })
    await expect(page.locator('text=版本信息').first()).toBeVisible({ timeout: 15000 })
  })

  test('无效软件 ID 显示错误', async ({ page }) => {
    await page.goto('/#/software/__invalid__')
    await page.waitForTimeout(3000)
    await expect(
      page.locator('text=未找到软件').or(page.locator('text=加载失败'))
    ).toBeVisible({ timeout: 15000 })
  })

  test('详情页有返回按钮', async ({ page }) => {
    await page.goto('/#/software/7zip')
    await page.waitForTimeout(3000)
    const backBtn = page.locator('button:has-text("返回")')
    await expect(backBtn).toBeVisible({ timeout: 15000 })
  })
})
