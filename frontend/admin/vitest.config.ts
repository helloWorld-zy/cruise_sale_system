// admin/vitest.config.ts — 管理后台的 Vitest 测试配置
// 配置 Vue 插件、JSDOM 测试环境和代码覆盖率相关选项
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
    plugins: [vue()],
    test: {
        environment: 'jsdom',
        coverage: {
            include: [
                'app/**/*.vue',
                'app/**/*.ts',
                'pages/**/*.vue',
                'components/**/*.vue',
            ],
            exclude: [
                '.nuxt/**',
                'coverage/**',
                'tests/**',
                'vitest.config.ts',
            ],
        },
    }
})
