import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/voyages/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockResolvedValue({
    data: [
      { id: 1, code: 'V-1', brief_info: '日韩 5晚6天', depart_date: '2026-03-01T00:00:00Z', return_date: '2026-03-06T00:00:00Z', first_stop_city: '天津' },
      { id: 2, code: 'V-2', brief_info: '地中海 7晚8天', depart_date: '2026-04-01T00:00:00Z', return_date: '2026-04-08T00:00:00Z', first_stop_city: '巴塞罗那' },
    ],
  })
})

const mountOptions = {
  global: {
    stubs: {
      AdminActionLink: { template: '<a><slot /></a>' },
      NuxtLink: { template: '<a><slot /></a>' },
    },
  },
}

describe('VoyagesIndexPage', () => {
  it('调用 voyages API', async () => {
    mount(Page, mountOptions)
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/voyages')
  })

  it('渲染航次行', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()
    expect(wrapper.text()).toContain('V-1')
    expect(wrapper.text()).toContain('2026-03-01')
  })

  it('支持按代码、出发/结束日期、港口关键字筛选', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()

    await wrapper.find('[data-test="filter-code"]').setValue('V-2')
    await wrapper.find('[data-test="filter-depart"]').setValue('2026-04-01')
    await wrapper.find('[data-test="filter-return"]').setValue('2026-04-08')
    await wrapper.find('[data-test="filter-port"]').setValue('巴塞罗那')
    await flushPromises()

    const rows = wrapper.findAll('tbody tr')
    expect(rows.length).toBe(1)
    expect(wrapper.text()).toContain('V-2')
    expect(wrapper.text()).not.toContain('V-1')
  })

  it('空数据时显示 No data', async () => {
    mockRequest.mockResolvedValueOnce({ data: [] })
    const wrapper = mount(Page, mountOptions)
    await flushPromises()
    expect(wrapper.text()).toContain('No data')
  })

  it('删除无效 id 时显示错误', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()
    await (wrapper.vm as any).handleDelete('bad')
    await flushPromises()
    expect(wrapper.text()).toContain('无效记录 ID，无法删除')
  })

  it('删除确认取消时不调用 DELETE', async () => {
    const wrapper = mount(Page, { ...mountOptions, attachTo: document.body })
    await flushPromises()
    await (wrapper.vm as any).handleDelete(1)
    await flushPromises()

    const cancelBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
      (btn.textContent || '').includes('取消'),
    )
    expect(cancelBtn).toBeTruthy()
    ;(cancelBtn as HTMLButtonElement).click()
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/voyages/1', { method: 'DELETE' })
    wrapper.unmount()
  })

  it('删除成功后调用 DELETE', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/voyages/1' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: [{ id: 1, code: 'V-1', brief_info: '日韩 5晚6天', depart_date: '2026-03-01T00:00:00Z', return_date: '2026-03-06T00:00:00Z', first_stop_city: '天津' }] })
    })
    const wrapper = mount(Page, { ...mountOptions, attachTo: document.body })
    await flushPromises()
    await (wrapper.vm as any).handleDelete(1)
    await flushPromises()

    const confirmBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
      (btn.textContent || '').includes('确认删除'),
    )
    expect(confirmBtn).toBeTruthy()
    ;(confirmBtn as HTMLButtonElement).click()
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/voyages/1', { method: 'DELETE' })
    wrapper.unmount()
  })
})
