// web/nuxt.config.ts — 前台 Web 应用的 Nuxt 框架配置文件
// 参考文档: https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2024-04-03', // 兼容性日期标记
    devtools: { enabled: true },     // 开发工具开关
    future: {
        compatibilityVersion: 4,     // 启用 Nuxt 4 兼容模式
    },
    runtimeConfig: {
        public: {
            // Default to same-origin API path to avoid browser CORS issues in dev.
            apiBase: process.env.NUXT_PUBLIC_API_BASE || '/api/v1',
        },
    },
    nitro: {
        // Proxy API calls during local development to the Go backend.
        devProxy: {
            '/api/v1': {
                target: 'http://localhost:8080/api/v1',
                changeOrigin: true,
            },
        },
    },
    modules: ['@pinia/nuxt'], // 仅保留状态管理模块，避免 Nuxt UI 运行时导入冲突
})
