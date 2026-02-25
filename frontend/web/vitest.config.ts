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
                'nuxt.config.ts',
                '.nuxt/**',
                'coverage/**',
                'tests/**',
                'vitest.config.ts',
            ],
        }
    }
})
