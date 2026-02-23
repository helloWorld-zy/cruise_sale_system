// admin/nuxt.config.ts — Nuxt 框架配置文件
// 参考文档: https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-04-03', // 兼容性日期标记
  devtools: { enabled: true },     // 开发工具开关
  future: {
    compatibilityVersion: 4,       // 启用 Nuxt 4 兼容模式
  },
  modules: ['@pinia/nuxt']         // 注册 Pinia 状态管理模块
})
