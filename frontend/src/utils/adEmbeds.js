function createPromoHtml({ eyebrow, title, body, cta, accent, glow, darkGlow }) {
  return `
    <div style="height:100%;min-height:320px;display:flex;flex-direction:column;justify-content:space-between;padding:24px;border-radius:18px;overflow:hidden;background:
      radial-gradient(circle at top right, ${glow} 0%, transparent 44%),
      linear-gradient(160deg, rgba(255,255,255,0.98), rgba(241,245,249,0.94));
      color:#0f172a;box-shadow:inset 0 1px 0 rgba(255,255,255,0.7);font-family:'Space Grotesk', 'Noto Sans SC', sans-serif;">
      <div style="display:flex;flex-direction:column;gap:14px;">
        <span style="display:inline-flex;width:max-content;border-radius:999px;padding:6px 10px;background:rgba(255,255,255,0.88);border:1px solid rgba(15,23,42,0.08);font-size:11px;font-weight:700;letter-spacing:0.18em;text-transform:uppercase;color:${accent};">${eyebrow}</span>
        <div>
          <h3 style="margin:0;font-size:28px;line-height:1.1;font-weight:700;">${title}</h3>
          <p style="margin:12px 0 0;font-size:14px;line-height:1.7;color:rgba(15,23,42,0.72);">${body}</p>
        </div>
      </div>
      <div style="display:flex;align-items:flex-end;justify-content:space-between;gap:16px;">
        <div style="width:124px;height:124px;border-radius:999px;background:radial-gradient(circle, ${darkGlow} 0%, rgba(255,255,255,0) 72%);"></div>
        <span style="display:inline-flex;align-items:center;justify-content:center;border-radius:999px;padding:12px 18px;background:${accent};color:#fff;font-size:13px;font-weight:700;letter-spacing:0.04em;white-space:nowrap;">${cta}</span>
      </div>
    </div>
  `
}

export const homeAdEmbeds = [
  {
    id: 'home-pro-workflow',
    title: '高效工作流套件',
    html: createPromoHtml({
      eyebrow: 'Workflow',
      title: '把下载、安装、更新放进一条工作流',
      body: '示例广告一，用来验证首页广告位轮转。你可以替换成自有推广、合作软件位或真实嵌入代码。',
      cta: '查看方案',
      accent: '#0f766e',
      glow: 'rgba(45, 212, 191, 0.34)',
      darkGlow: 'rgba(15, 118, 110, 0.26)'
    })
  },
  {
    id: 'home-security-stack',
    title: '终端安全工具集合',
    html: createPromoHtml({
      eyebrow: 'Security',
      title: '统一收口常用安全软件与审计工具',
      body: '示例广告二，强调不同内容卡片在同一广告容器内自动切换时的展示状态。',
      cta: '立即了解',
      accent: '#1d4ed8',
      glow: 'rgba(96, 165, 250, 0.34)',
      darkGlow: 'rgba(29, 78, 216, 0.22)'
    })
  },
  {
    id: 'home-dev-pack',
    title: '开发环境精选包',
    html: createPromoHtml({
      eyebrow: 'Developers',
      title: '为新设备快速补齐 IDE、浏览器与同步工具',
      body: '示例广告三，用于观察多条 embed 在首页连续轮转时的节奏和占位稳定性。',
      cta: '开始配置',
      accent: '#7c3aed',
      glow: 'rgba(196, 181, 253, 0.34)',
      darkGlow: 'rgba(124, 58, 237, 0.22)'
    })
  }
]

export const detailAdEmbeds = [
  {
    id: 'detail-sync-suite',
    title: '跨端同步方案',
    html: createPromoHtml({
      eyebrow: 'Sync',
      title: '让版本下载、账号配置与文件同步一步到位',
      body: '详情页示例广告一，适合验证广告位在版本信息区域旁边的可读性与轮转体验。',
      cta: '查看套餐',
      accent: '#be123c',
      glow: 'rgba(251, 113, 133, 0.3)',
      darkGlow: 'rgba(190, 24, 93, 0.2)'
    })
  },
  {
    id: 'detail-cloud-space',
    title: '云存储加速位',
    html: createPromoHtml({
      eyebrow: 'Cloud',
      title: '安装包存储、分发与版本归档集中管理',
      body: '详情页示例广告二，模拟偏 B 端的嵌入内容，帮助你确认不同素材风格下的切换表现。',
      cta: '查看详情',
      accent: '#c2410c',
      glow: 'rgba(251, 191, 36, 0.3)',
      darkGlow: 'rgba(194, 65, 12, 0.2)'
    })
  },
  {
    id: 'detail-member-offer',
    title: '会员促销横幅',
    html: createPromoHtml({
      eyebrow: 'Membership',
      title: '把合作推广位直接挂到详情页右侧长期展示',
      body: '详情页示例广告三，用于确认轮转中的按钮、文案和背景层次是否达到你要的效果。',
      cta: '立即体验',
      accent: '#2563eb',
      glow: 'rgba(59, 130, 246, 0.28)',
      darkGlow: 'rgba(37, 99, 235, 0.2)'
    })
  }
]