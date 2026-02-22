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
})
