import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabin-types/pricing.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

const globalConfig = {
  stubs: {
    AdminActionLink: {
      props: ['to'],
      template: '<a :href="String(to)"><slot /></a>',
    },
    NuxtLink: {
      props: ['to'],
      template: '<a :href="String(to)"><slot /></a>',
    },
  },
}

describe('CabinTypesPricingPage', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/companies') {
        return Promise.resolve({ data: { list: [{ id: 1, name: 'Oceanic' }] } })
      }
      if (url === '/cruises') {
        return Promise.resolve({ data: { list: [{ id: 2, name: 'Ocean Nova' }] } })
      }
      if (url === '/cabin-pricing/voyages') {
        return Promise.resolve({ data: { list: [{ id: 10, cruise_id: 2, name: 'VN001', depart_date: '2026-05-01' }] } })
      }
      if (url === '/cabin-types' && options?.query?.cruise_id === 2) {
        return Promise.resolve({ data: { list: [{ id: 99, name: '豪华阳台房', code: 'BAL' }] } })
      }
      if (url === '/cabin-pricing/history') {
        return Promise.resolve({ data: { list: [] } })
      }
      if (url === '/cabin-pricing/batch-apply' && options?.method === 'POST') {
        return Promise.resolve({ data: { applied: 1, failed: 0, errors: [] } })
      }
      return Promise.resolve({ data: {} })
    })
  })

  it('挂载时加载公司、邮轮和航次数据', async () => {
    mount(Page, { global: globalConfig })
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/companies', expect.any(Object))
    expect(mockRequest).toHaveBeenCalledWith('/cruises', expect.any(Object))
    expect(mockRequest).toHaveBeenCalledWith('/cabin-pricing/voyages', expect.any(Object))
  })

  it('渲染批量应用价格入口', async () => {
    const wrapper = mount(Page, { global: globalConfig })
    await flushPromises()

    expect(wrapper.text()).toContain('舱型价格管理')
    expect(wrapper.text()).toContain('批量应用价格')
  })

  it('选择航次和舱型后提交 batch apply', async () => {
    const wrapper = mount(Page, { global: globalConfig })
    await flushPromises()

    const checkbox = wrapper.find('input[type="checkbox"]')
    await checkbox.setValue(true)

    const selects = wrapper.findAll('select')
    await selects[2]!.setValue('99')

    const numbers = wrapper.findAll('input[type="number"]')
    await numbers[0]!.setValue('20')
    await numbers[1]!.setValue('120000')
    await numbers[2]!.setValue('150000')

    const applyBtn = wrapper.findAll('button').find((btn) => btn.text().includes('批量应用价格'))
    await applyBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/cabin-pricing/batch-apply', {
      method: 'POST',
      body: {
        voyage_ids: [10],
        cabin_type_id: 99,
        inventory_total: 20,
        settlement_price_cents: 120000,
        sale_price_cents: 150000,
        effective_at: undefined,
      },
    })
  })
})
