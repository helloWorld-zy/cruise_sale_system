import { afterAll, beforeEach, describe, expect, it, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../app/pages/bookings/index.vue'

const mockRequest = vi.fn()
const confirmMock = vi.fn(() => true)
const clickMock = vi.fn()
const originalCreateElement = document.createElement.bind(document)

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('confirm', confirmMock)

const createElementSpy = vi.spyOn(document, 'createElement').mockImplementation((tagName: string) => {
  if (tagName.toLowerCase() === 'a') {
    return {
      href: '',
      download: '',
      click: clickMock,
    } as any
  }
  return originalCreateElement(tagName)
})

vi.spyOn(URL, 'createObjectURL').mockReturnValue('blob:mock')
vi.spyOn(URL, 'revokeObjectURL').mockImplementation(() => {})

function successList(list: any[], total = list.length) {
  return { data: { list, total } }
}

beforeEach(() => {
  mockRequest.mockReset()
  confirmMock.mockClear()
  clickMock.mockClear()
  mockRequest.mockResolvedValue(
    successList([
      { id: 1, status: 'pending_payment', total_cents: 19900, voyage_id: 2, user_id: 9 },
      { id: 2, status: 'paid', total_cents: 38000, voyage_id: 3, user_id: 9 },
    ]),
  )
})

describe('Admin Bookings', () => {
  it('renders enhanced filters and export button', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.find('[data-test="filter-booking-no"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-phone"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="export"]').exists()).toBe(true)
  })

  it('loads list with query filters', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('[data-test="filter-booking-no"]').setValue('1001')
    await wrapper.find('[data-test="filter-phone"]').setValue('138')
    await wrapper.find('[data-test="filter-submit"]').trigger('submit')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/bookings', expect.objectContaining({ query: expect.objectContaining({ booking_no: '1001', phone: '138' }) }))
  })

  it('status tab triggers status query', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('[data-test="tab-paid"]').trigger('click')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/bookings', expect.objectContaining({ query: expect.objectContaining({ status: 'paid' }) }))
  })

  it('shows pay and refund actions on row', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('去支付')
    expect(wrapper.text()).toContain('申请退改')
    expect(wrapper.text()).toContain('查看详情')
  })

  it('exports current rows as csv', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('[data-test="export"]').trigger('click')
    expect(clickMock).toHaveBeenCalledTimes(1)
  })

  it('shows empty state when list is empty', async () => {
    mockRequest.mockResolvedValueOnce(successList([], 0))
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.find('[data-test="empty"]').exists()).toBe(true)
  })

  it('shows error when load failed', async () => {
    mockRequest.mockRejectedValueOnce(new Error('network error'))
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.find('[data-test="error"]').text()).toContain('network error')
  })

  it('deletes row and reloads', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/bookings/1' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve(successList([{ id: 2, status: 'paid', total_cents: 38000 }], 1))
    })

    const wrapper = mount(Page)
    await flushPromises()
    await (wrapper.vm as any).handleDelete(1)
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/bookings/1', { method: 'DELETE' })
  })
})

afterAll(() => {
  createElementSpy.mockRestore()
})
