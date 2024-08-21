const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://158.160.122.246:8080',
      changeOrigin: true,
    })
  );
};