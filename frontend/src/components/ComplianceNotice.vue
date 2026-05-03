<template>
  <section class="w-full rounded-xl border border-slate-200/80 bg-white/70 px-3 py-2.5 text-[11px] text-slate-500 shadow-[0_4px_14px_rgba(15,70,56,0.08)] backdrop-blur-sm dark:border-slate-700/70 dark:bg-slate-900/55 dark:text-slate-400 md:ml-auto md:max-w-[340px]">
    <p class="leading-5">下载入口来自官方渠道或官方镜像，仅作导航参考。</p>
    <details class="group mt-1.5">
      <summary class="cursor-pointer select-none text-[11px] font-medium text-brand-700 marker:hidden dark:text-brand-400">
        <span class="inline-flex items-center gap-1">
          合规与隐私说明
          <span class="transition group-open:rotate-180">▾</span>
        </span>
      </summary>
      <div class="mt-2 space-y-1.5 leading-5 text-slate-500 dark:text-slate-400">
        <p>本站不托管安装包，不替代官方发布页；链接可用性、文件安全性和版本准确性请自行核验。</p>
        <p>遵循数据最小化实践（参考 GDPR 原则），仅请求软件目录与版本数据，不要求登录，不采集身份信息。</p>
        <p>仅在浏览器本地保存主题偏好（localStorage）；可随时清除站点数据。</p>
      </div>
    </details>
  </section>
</template>
