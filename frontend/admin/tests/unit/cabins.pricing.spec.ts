import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import PricingPage from '../../app/pages/cabins/pricing.vue'

vi.stubGlobal('useApi', () => ({ request: vi.fn().mockResolvedValue({ data: [] }) }))
vi.stubGlobal('useRoute', () => ({ query: { skuId: '1' } }))

describe('Cabin Pricing Page', () => {
    it('shows pricing title', () => {
        const wrapper = mount(PricingPage)
        expect(wrapper.text()).toContain('Pricing Matrix')
    })
})
