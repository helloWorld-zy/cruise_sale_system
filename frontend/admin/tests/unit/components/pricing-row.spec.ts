import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PricingRow from '../../../app/components/PricingRow.vue'

describe('PricingRow', () => {
  it('renders date occupancy and price', () => {
    const wrapper = mount(PricingRow, {
      props: {
        row: { date: '2026-06-01', occupancy: 2, price: 29900 },
      },
    })

    expect(wrapper.text()).toContain('2026-06-01')
    expect(wrapper.text()).toContain('2')
    expect(wrapper.text()).toContain('29900')
  })
})
