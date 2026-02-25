import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import App from '../../app/app.vue'

describe('Admin App root', () => {
    test('renders NuxtLayout and NuxtPage container', () => {
        const wrapper = mount(App, {
            global: {
                stubs: {
                    NuxtLayout: { name: 'NuxtLayout', template: '<div class="layout"><slot /></div>' },
                    NuxtPage: { name: 'NuxtPage', template: '<div class="page" />' },
                },
            },
        })

        expect(wrapper.find('.layout').exists()).toBe(true)
        expect(wrapper.find('.page').exists()).toBe(true)
    })
})
