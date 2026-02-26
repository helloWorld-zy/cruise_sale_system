import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Page from '../../pages/pay/[id].vue'

describe('Pay Page', () => {
    it('renders pay button', () => {
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('Pay Now')
    })
})
