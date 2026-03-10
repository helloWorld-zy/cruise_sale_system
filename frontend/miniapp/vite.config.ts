import { defineConfig } from 'vite'
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

export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => uniElements.has(tag),
        },
      },
    }),
  ],
})
