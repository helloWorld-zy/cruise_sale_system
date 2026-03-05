import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Layout from '../../../app/layouts/default.vue'

describe('Default layout', () => {
    test('renders a header', () => {
        const wrapper = mount(Layout, {
            slots: { default: '<div />' },
            global: {
                stubs: { NuxtLink: true }
            }
        })
        expect(wrapper.find('header').exists()).toBe(true)
    })

    test('renders nav links, slot content and footer text', () => {
        const wrapper = mount(Layout, {
            slots: { default: '<div data-test="slot-content">content</div>' },
            global: {
                stubs: { NuxtLink: { template: '<a><slot /></a>' } }
            }
        })

        expect(wrapper.text()).toContain('Azure Deck')
        expect(wrapper.text()).toContain('搜索')
        expect(wrapper.text()).toContain('预订')
        expect(wrapper.text()).toContain('登录')
        expect(wrapper.find('[data-test="slot-content"]').exists()).toBe(true)
        expect(wrapper.text()).toContain('© 2026 Azure Deck')
    })
})
