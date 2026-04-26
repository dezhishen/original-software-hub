// 全局数据源配置，需在 bundle 之前加载
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
