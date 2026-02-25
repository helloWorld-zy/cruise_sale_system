import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import App from '../../app/app.vue'

describe('App', () => {
    it('renders Nuxt elements', () => {
        const wrapper = mount(App, {
            global: {
                stubs: { NuxtLayout: true, NuxtPage: true }
            }
        })
        expect(wrapper.findComponent({ name: 'NuxtLayout' }).exists()).toBe(true)
    })
})
