import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import DefaultLayout from '../../../app/layouts/default.vue'

describe('Default layout', () => {
    test('renders a main content slot', () => {
        const wrapper = mount(DefaultLayout, {
            slots: { default: '<div data-test="content">Hello</div>' },
            global: {
                stubs: { NuxtLink: true }
            }
        })
        expect(wrapper.find('[data-test="content"]').exists()).toBe(true)
    })
})
