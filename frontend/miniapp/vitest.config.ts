import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

// 小程序单元测试配置：使用 Vue 插件并基于 happy-dom 运行。
export default defineConfig({
    plugins: [vue()],
    test: {
        environment: 'happy-dom',
    },
})
