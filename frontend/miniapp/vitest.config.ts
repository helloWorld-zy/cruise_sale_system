import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

const uniElements = new Set([
    'view',
    'text',
    'image',
    'scroll-view',
    'swiper',
    'swiper-item',
    'navigator',
])

// 小程序单元测试配置：使用 Vue 插件并基于 happy-dom 运行。
export default defineConfig({
    plugins: [vue({
        template: {
            compilerOptions: {
                isCustomElement: (tag) => uniElements.has(tag),
            },
        },
    })],
    test: {
        environment: 'happy-dom',
    },
})
