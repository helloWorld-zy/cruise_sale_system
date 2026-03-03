import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cabins/pricing.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  vi.stubGlobal('useRoute', () => ({ query: { skuId: '9' } }))
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({ data: [{ date: '2026-05-01', occupancy: 2, price_cents: 19900, price_type: 'base' }] })
})

describe('CabinsPricingPage', () => {
  it('挂载时加载价格列表', async () => {
    mount(Page)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/cabins/9/prices')
  })

  it('渲染价格类型与日历价格', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('基础')
    expect(wrapper.text()).toContain('批量设价')
    expect(wrapper.text()).toContain('保存单日价格')
  })

  it('skuId 缺失时显示错误', async () => {
    vi.stubGlobal('useRoute', () => ({ query: {} }))
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('缺少 skuId 参数')
    expect(mockRequest).not.toHaveBeenCalled()
  })

  it('保存单日价格会调用 POST 并刷新', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cabins/9/prices' && options?.method === 'POST') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: [{ date: '2026-05-02', occupancy: 2, price_cents: 28800, price_type: 'base' }] })
    })

    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('2026-05-02')
    await inputs[1]!.setValue('28800')

    const saveBtn = wrapper.findAll('button').find((btn) => btn.text().includes('保存单日价格'))
    await saveBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/cabins/9/prices', {
      method: 'POST',
      body: {
        date: '2026-05-02',
        price_cents: 28800,
        price_type: 'base',
      },
    })
  })

  it('批量设价日期非法时给出错误', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    const openBatchBtn = wrapper.findAll('button').find((btn) => btn.text().includes('批量设价'))
    await openBatchBtn!.trigger('click')
    await flushPromises()

    const dateInputs = wrapper.findAll('input[type="date"]')
    const numberInputs = wrapper.findAll('input[type="number"]')
    await dateInputs[1]!.setValue('2026-05-10')
    await dateInputs[2]!.setValue('2026-05-01')
    await numberInputs[1]!.setValue('20000')

    const submitBtn = wrapper.findAll('button').find((btn) => btn.text() === '提交')
    await submitBtn!.trigger('click')
    await flushPromises()

    expect(wrapper.text()).toContain('无效日期范围')
  })

  it('批量设价会按天循环提交', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cabins/9/prices' && options?.method === 'POST') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: [] })
    })

    const openBatchBtn = wrapper.findAll('button').find((btn) => btn.text().includes('批量设价'))
    await openBatchBtn!.trigger('click')
    await flushPromises()

    const dateInputs = wrapper.findAll('input[type="date"]')
    const numberInputs = wrapper.findAll('input[type="number"]')
    await dateInputs[1]!.setValue('2026-05-01')
    await dateInputs[2]!.setValue('2026-05-03')
    await numberInputs[1]!.setValue('30000')

    const submitBtn = wrapper.findAll('button').find((btn) => btn.text() === '提交')
    await submitBtn!.trigger('click')
    await flushPromises()

    const postCalls = mockRequest.mock.calls.filter((call) => call[0] === '/cabins/9/prices' && call[1]?.method === 'POST')
    expect(postCalls.length).toBe(3)
  })
})
