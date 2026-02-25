import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import DefaultLayout from '../../../app/layouts/default.vue'

describe('Default layout', () => {
    test('renders a main content slot and navigation links', () => {
        const wrapper = mount(DefaultLayout, {
            slots: { default: '<div data-test="content">Hello</div>' },
            global: {
                stubs: {
                    NuxtLink: {
                        name: 'NuxtLink',
                        props: ['to'],
                        template: '<a :href="String(to)"><slot /></a>'
                    }
                }
            }
        })
        expect(wrapper.find('[data-test="content"]').exists()).toBe(true)
        const links = wrapper.findAllComponents({ name: 'NuxtLink' })
        expect(links.length).toBeGreaterThan(0)
    })
})
