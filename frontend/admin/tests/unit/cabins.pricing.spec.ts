import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PricingPage from '../../app/pages/cabins/pricing.vue'

describe('Cabin Pricing Page', () => {
    it('shows pricing title', () => {
        const wrapper = mount(PricingPage)
        expect(wrapper.text()).toContain('Pricing Matrix')
    })
})
