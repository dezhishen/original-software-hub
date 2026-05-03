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
})();
