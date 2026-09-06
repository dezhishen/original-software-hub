// 全局数据源配置 — 可在部署时修改此文件覆盖默认值
(function () {
  window.APP_DATA_SOURCE_CONFIG = {
    endpoint: {
      type: "json",
      url: "./data/json",
      indexPath: "index.json",
      timeoutMs: 8000
    }
  };

  // 广告位配置（可在部署时替换为真实平台参数）
  // 当前示例替换为公益主题素材，优先引导到真实公益/寻亲入口。
  // 支持平台：iframe / custom-html / script / adsense
  window.APP_AD_CONFIG = {
    defaultRotateMs: 5000,
    slots: {
      home: {
        rotateMs: 4800,
        providers: [
          {
            id: "home-missing-script",
            platform: "script",
            title: "寻亲公益入口",
            slotHtml: "<section id='home-ad-script-slot' style='height:100%;min-height:320px;border-radius:18px;background:linear-gradient(145deg,#fff7ed,#fffbeb);padding:22px;box-sizing:border-box;font-family:Space Grotesk,Noto Sans SC,sans-serif;color:#9a3412;display:flex;flex-direction:column;justify-content:space-between'><div><div style='font-size:11px;letter-spacing:.18em;text-transform:uppercase;opacity:.72'>Public Service</div><h3 style='margin:10px 0 0;font-size:24px;line-height:1.2'>紧急寻亲与线索转发</h3><p style='margin:10px 0 0;font-size:13px;line-height:1.7'>如发现疑似走失儿童、老人或救助线索，请优先提交到已有公益或官方平台。</p></div><div id='home-ad-script-links' style='display:grid;gap:10px'></div></section>",
            scripts: [
              {
                inline: `(function(){var host=document.getElementById('home-ad-script-links');if(!host)return;var links=[{label:'宝贝回家寻子网',desc:'民间公益寻亲平台，适合转发线索与查看协查信息',href:'https://www.baobeihuijia.com/'},{label:'全国救助寻亲网',desc:'民政系统救助寻亲入口，适合核验官方救助信息',href:'https://xunqin.mca.gov.cn/'},{label:'腾讯公益项目广场',desc:'支持向合规公益项目捐助或参与志愿行动',href:'https://gongyi.qq.com/succor/project_list.htm'}];host.innerHTML=links.map(function(item){return '<a href="'+item.href+'" target="_blank" rel="noopener noreferrer nofollow" style="display:block;text-decoration:none;border-radius:14px;padding:12px 14px;background:rgba(255,255,255,.68);border:1px solid rgba(251,146,60,.22);color:#7c2d12"><div style="font-size:14px;font-weight:700;line-height:1.4">'+item.label+'</div><div style="margin-top:4px;font-size:12px;line-height:1.6;opacity:.82">'+item.desc+'</div></a>';}).join('');})();`
              }
            ]
          },
          {
            id: "home-iframe-xunqin",
            platform: "iframe",
            title: "官方救助寻亲入口",
            srcdoc: "<html><body style='margin:0;padding:18px;font-family:Space Grotesk,Noto Sans SC,sans-serif;background:linear-gradient(145deg,#ecfeff,#f0fdfa);height:100%;box-sizing:border-box;color:#115e59'><div style='height:100%;border-radius:18px;background:rgba(255,255,255,.74);padding:22px;box-sizing:border-box;display:flex;flex-direction:column;justify-content:space-between;border:1px solid rgba(13,148,136,.12)'><div><div style='font-size:11px;letter-spacing:.18em;text-transform:uppercase;opacity:.74'>Official Route</div><h3 style='margin:12px 0 0;font-size:24px;line-height:1.2'>全国救助寻亲网</h3><p style='margin:10px 0 0;font-size:13px;line-height:1.7'>对接民政系统公开入口，适合查看救助寻亲、身份核验等官方信息。</p></div><a href='https://xunqin.mca.gov.cn/' target='_blank' rel='noopener noreferrer nofollow' style='display:inline-flex;align-items:center;justify-content:center;height:42px;border-radius:999px;background:#0f766e;color:#fff;text-decoration:none;font-size:13px;font-weight:700'>打开官方入口</a></div></body></html>"
          },
          {
            id: "home-html-charity",
            platform: "custom-html",
            title: "公益参与卡片",
            html: "<div style='height:100%;min-height:320px;border-radius:18px;background:linear-gradient(145deg,#eff6ff,#ecfeff);padding:24px;box-sizing:border-box;font-family:Space Grotesk,Noto Sans SC,sans-serif;color:#1d4ed8;display:flex;flex-direction:column;justify-content:space-between'><div><div style='font-size:11px;letter-spacing:.18em;text-transform:uppercase;opacity:.76'>Volunteer</div><h3 style='margin:12px 0 0;font-size:24px;line-height:1.2'>把一次下载，变成一次公益触达</h3><p style='margin:10px 0 0;font-size:13px;line-height:1.7;color:#334155'>你可以转发寻亲信息、核验公益机构、参与儿童保护或乡村教育项目。</p><div style='display:grid;gap:8px;margin-top:14px'><a href='https://www.baobeihuijia.com/' target='_blank' rel='noopener noreferrer nofollow' style='text-decoration:none;border-radius:12px;padding:10px 12px;background:rgba(255,255,255,.72);color:#1e3a8a;border:1px solid rgba(59,130,246,.16)'>宝贝回家：查看寻亲与协查</a><a href='https://gongyi.qq.com/' target='_blank' rel='noopener noreferrer nofollow' style='text-decoration:none;border-radius:12px;padding:10px 12px;background:rgba(255,255,255,.72);color:#1e3a8a;border:1px solid rgba(59,130,246,.16)'>腾讯公益：寻找合规项目</a></div></div><div style='font-size:12px;opacity:.72'>公益信息仅作导航，请以目标平台最新公告为准。</div></div>"
          }
        ]
      },
      detail: {
        rotateMs: 5200,
        providers: [
          {
            id: "detail-adsense-public-service",
            platform: "adsense",
            title: "公益合作广告位",
            fallbackHtml: "<div style='height:100%;min-height:320px;border-radius:18px;background:linear-gradient(145deg,#f0fdf4,#dcfce7);padding:24px;box-sizing:border-box;font-family:Space Grotesk,Noto Sans SC,sans-serif;color:#166534;display:flex;flex-direction:column;justify-content:space-between'><div><div style='font-size:11px;letter-spacing:.18em;text-transform:uppercase;opacity:.74'>Public Welfare Reserve</div><h3 style='margin:12px 0 0;font-size:24px;line-height:1.2'>公益合作投放位</h3><p style='margin:10px 0 0;font-size:13px;line-height:1.7'>预留给寻亲、救助、儿童保护、乡村教育等公益组织。当前展示的是站内公益导航卡片。</p><div style='margin-top:14px;display:flex;flex-wrap:wrap;gap:8px'><a href='https://gongyi.qq.com/' target='_blank' rel='noopener noreferrer nofollow' style='text-decoration:none;border-radius:999px;padding:8px 12px;background:#166534;color:#fff;font-size:12px'>腾讯公益</a><a href='https://www.mca.gov.cn/n157/n193/n455/index.html' target='_blank' rel='noopener noreferrer nofollow' style='text-decoration:none;border-radius:999px;padding:8px 12px;background:rgba(255,255,255,.82);color:#166534;font-size:12px;border:1px solid rgba(22,101,52,.14)'>民政部公益广告</a></div></div><div style='font-size:12px;opacity:.72'>如需接入真实公益流量平台，可在部署时补充 client/slot。</div></div>"
          },
          {
            id: "detail-script-missing",
            platform: "script-snippet",
            title: "失踪人员线索协作",
            slotHtml: "<section id='detail-ad-script-slot' style='height:100%;min-height:320px;border-radius:18px;background:linear-gradient(145deg,#f5f3ff,#ede9fe);padding:22px;box-sizing:border-box;font-family:Space Grotesk,Noto Sans SC,sans-serif;color:#5b21b6;display:flex;flex-direction:column;justify-content:space-between'><div><div style='font-size:11px;letter-spacing:.18em;text-transform:uppercase;opacity:.74'>Case Relay</div><h3 style='margin:10px 0 0;font-size:24px;line-height:1.2'>如果你看到相似线索</h3><p style='margin:10px 0 0;font-size:13px;line-height:1.7'>请不要在社交平台直接扩散个人隐私，优先走公益平台或官方渠道提交线索。</p></div><div id='detail-ad-script-actions' style='display:grid;gap:10px'></div></section>",
            scripts: [
              {
                inline: `(function(){var host=document.getElementById('detail-ad-script-actions');if(!host)return;var items=[{label:'提交到宝贝回家',href:'https://www.baobeihuijia.com/'},{label:'查看全国救助寻亲',href:'https://xunqin.mca.gov.cn/'},{label:'支持合规公益项目',href:'https://gongyi.qq.com/succor/project_list.htm'}];host.innerHTML=items.map(function(item,index){return '<a href="'+item.href+'" target="_blank" rel="noopener noreferrer nofollow" style="display:flex;align-items:center;justify-content:space-between;text-decoration:none;border-radius:14px;padding:12px 14px;background:rgba(255,255,255,.68);border:1px solid rgba(124,58,237,.14);color:#4c1d95"><span style="font-size:13px;font-weight:700">'+item.label+'</span><span style="font-size:12px;opacity:.72">0'+(index+1)+'</span></a>';}).join('');})();`
              }
            ]
          },
          {
            id: "detail-iframe-child-protection",
            platform: "iframe",
            title: "儿童保护与救助",
            srcdoc: "<html><body style='margin:0;padding:18px;font-family:Space Grotesk,Noto Sans SC,sans-serif;background:linear-gradient(145deg,#fef2f2,#fff7ed);height:100%;box-sizing:border-box;color:#9a3412'><div style='height:100%;border-radius:18px;background:rgba(255,255,255,.76);padding:22px;box-sizing:border-box;display:flex;flex-direction:column;justify-content:space-between;border:1px solid rgba(249,115,22,.14)'><div><div style='font-size:11px;letter-spacing:.18em;text-transform:uppercase;opacity:.74'>Child Protection</div><h3 style='margin:12px 0 0;font-size:24px;line-height:1.2'>关注儿童保护与走失预防</h3><p style='margin:10px 0 0;font-size:13px;line-height:1.7'>识别高风险场景、学习走失预防知识，也可以将公益入口转发给更多人。</p></div><div style='display:grid;gap:10px'><a href='https://www.baobeihuijia.com/' target='_blank' rel='noopener noreferrer nofollow' style='display:inline-flex;align-items:center;justify-content:center;height:42px;border-radius:999px;background:#ea580c;color:#fff;text-decoration:none;font-size:13px;font-weight:700'>查看寻亲入口</a><a href='https://www.mca.gov.cn/' target='_blank' rel='noopener noreferrer nofollow' style='display:inline-flex;align-items:center;justify-content:center;height:42px;border-radius:999px;background:rgba(255,255,255,.9);color:#9a3412;text-decoration:none;font-size:13px;font-weight:700;border:1px solid rgba(249,115,22,.18)'>查看民政部公开信息</a></div></div></body></html>"
          }
        ]
      }
    }
  };
})();
