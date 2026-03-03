// admin/nuxt.config.ts — Nuxt 框架配置文件
// 参考文档: https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-04-03', // 兼容性日期标记
  devtools: { enabled: true },     // 开发工具开关
  future: {
    compatibilityVersion: 4,       // 启用 Nuxt 4 兼容模式
  },
  runtimeConfig: {
    public: {
      // 默认使用同源 API 路径，避免本地开发时浏览器跨域失败。
      apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api/v1',
    },
  },
  nitro: {
    // 本地开发代理到 Go 后端，保持前端始终同源请求。
    devProxy: {
      '/api/v1': {
        target: 'http://localhost:8080/api/v1',
        changeOrigin: true,
      },
    },
  },
  modules: ['@pinia/nuxt'], // 仅保留必要模块，避免 UI 模块与 Tailwind 版本冲突
})
