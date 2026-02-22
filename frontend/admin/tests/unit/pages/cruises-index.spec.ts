import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/index.vue'

describe('Cruise list', () => {
    test('renders a table', () => {
        const wrapper = mount(Page, {
            global: {
                stubs: { UButton: true }
            }
        })
        expect(wrapper.find('table').exists()).toBe(true)
    })
})
