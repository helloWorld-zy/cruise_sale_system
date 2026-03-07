import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import PricingPage from '../../app/pages/cabins/pricing.vue'

vi.stubGlobal('useApi', () => ({ request: vi.fn().mockResolvedValue({ data: [] }) }))
vi.stubGlobal('useRoute', () => ({ query: { skuId: '1' } }))
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
vi.stubGlobal('navigateTo', mockNavigateTo)

describe('Cabin Pricing Page', () => {
    it('redirects to cabin-types pricing page', async () => {
        const Wrapper = defineComponent({
            components: { PricingPage },
            template: '<Suspense><PricingPage /></Suspense>',
        })
        const wrapper = mount(Wrapper)
        await Promise.resolve()

        expect(mockNavigateTo).toHaveBeenCalledWith('/cabin-types/pricing', { replace: true })
        expect(wrapper.exists()).toBe(true)
    })
})
