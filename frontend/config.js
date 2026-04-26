// 全局数据源配置，需在 bundle 之前加载
(function () {
  const isLocalhost = ["localhost", "127.0.0.1"].includes(window.location.hostname);

  window.APP_DATA_SOURCE_CONFIG = {
    endpoint: {
      type: "jsonp",
      url: isLocalhost
        ? "./data/jsonp"
        : "https://cdn.jsdmirror.com/gh/dezhishen/original-software-hub@data/jsonp",
      indexPath: "index.js",
      callbackParam: "callback",
      timeoutMs: 8000
    }
  };
})();
