import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/index.vue'

describe('Cruise list', () => {
    const UButtonStub = {
        name: 'UButton',
        props: ['to'],
        template: '<button :data-to="to"><slot /></button>'
    }

    test('renders a table', () => {
        const wrapper = mount(Page, {
            global: {
                stubs: { UButton: UButtonStub }
            }
        })
        expect(wrapper.find('table').exists()).toBe(true)
    })

    test('renders create cruise button', () => {
        const wrapper = mount(Page, {
            global: {
                stubs: { UButton: UButtonStub }
            }
        })
        const btn = wrapper.findComponent({ name: 'UButton' })
        expect(btn.exists()).toBe(true)
        expect(btn.attributes('data-to')).toBe('/cruises/create')
    })
})
