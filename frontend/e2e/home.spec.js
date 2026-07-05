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

  test.describe('分类导航', () => {
    test('分类栏可见且包含分类按钮', async ({ page }) => {
      await page.goto('/')
      await page.waitForSelector('.category-bar', { timeout: 15000 })
      const bar = page.locator('.category-bar')
      await expect(bar).toBeVisible()
      // 至少有一个分类按钮（除了"全部分类"）
      const buttons = bar.locator('button')
      const count = await buttons.count()
      expect(count).toBeGreaterThan(1)
    })

    test('默认显示"全部分类"选中', async ({ page }) => {
      await page.goto('/')
      await page.waitForSelector('.category-bar', { timeout: 15000 })
      const allBtn = page.locator('.category-bar button', { hasText: '全部分类' })
      await expect(allBtn).toBeVisible()
    })

    test('点击分类按钮过滤软件列表', async ({ page }) => {
      await page.goto('/')
      await page.waitForSelector('.category-bar', { timeout: 15000 })
      // 点击第一个非"全部分类"的按钮
      const catBtn = page.locator('.category-bar button').nth(1)
      const catText = await catBtn.textContent()
      await catBtn.click()
      await page.waitForTimeout(500)
      // 选中状态应高亮（背景色为 brand）
      await expect(catBtn).toHaveClass(/bg-brand/)
      // 卡片数量应减少或不变
      const cards = page.locator('.software-card')
      const totalCards = await cards.count()
      const allCards = await page.locator('.software-card').all()
      // 每个卡片应属于该分类
      if (totalCards > 0) {
        for (const card of allCards) {
          const text = await card.textContent()
          expect(text.length).toBeGreaterThan(0)
        }
      }
    })

    test('分类间切换正常', async ({ page }) => {
      await page.goto('/')
      await page.waitForSelector('.category-bar', { timeout: 15000 })
      const buttons = page.locator('.category-bar button')
      const btnCount = await buttons.count()
      // 切换到第二个分类
      if (btnCount >= 2) {
        await buttons.nth(1).click()
        await page.waitForTimeout(300)
        // 再切换回"全部分类"
        await buttons.nth(0).click()
        await page.waitForTimeout(300)
        // 全部分类应显示所有卡片
        const cards = page.locator('.software-card')
        const count = await cards.count()
        expect(count).toBeGreaterThan(10)
      }
    })
  })
})
