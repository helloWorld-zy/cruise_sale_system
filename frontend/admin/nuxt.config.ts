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
      // 默认直连本地后端，避免 runtimeConfig 缺失导致请求地址异常。
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:18080/api/v1',
    },
  },
  modules: ['@pinia/nuxt'], // 仅保留必要模块，避免 UI 模块与 Tailwind 版本冲突
})
