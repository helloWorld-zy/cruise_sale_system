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
    const anchor = originalCreateElement('a') as HTMLAnchorElement
    anchor.click = clickMock as any
    return anchor
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
      { id: 1, status: 'pending_payment', total_cents: 19900, voyage_id: 2, user_id: 9, booking_no: 'BK20260307001', phone: '13800000001', voyage_code: 'VOY-2026-Q1', cruise_name: '海洋光谱号' },
      { id: 2, status: 'paid', total_cents: 38000, voyage_id: 3, user_id: 9, booking_no: 'BK20260307002', phone: '13800000002', voyage_code: 'VOY-2026-Q2', cruise_name: '海洋量子号' },
    ]),
  )
})

describe('Admin Bookings', () => {
  const mountOptions = {
    global: {
      stubs: {
        AdminActionLink: { template: '<a><slot /></a>' },
        AdminPageHeader: { props: ['title', 'subtitle'], template: '<div>{{ title }} {{ subtitle }}<slot /><slot name="actions" /></div>' },
        AdminFilterBar: { template: '<div><slot /></div>' },
        AdminDataCard: { template: '<div><slot /></div>' },
        AdminStatusTag: { props: ['text'], template: '<span>{{ text }}</span>' },
      },
    },
  }

  it('renders enhanced filters, search bar and balanced header actions', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()
    expect(wrapper.find('[data-test="booking-search-input"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-booking-no"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-phone"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-voyage-code"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="filter-cruise-name"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="export"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="header-actions"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="export"]').classes()).toContain('bookings-header__action')
    expect(wrapper.find('[data-test="create-booking"]').classes()).toContain('bookings-header__action')
  })

  it('loads list with query filters and keyword search', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()

    await wrapper.find('[data-test="booking-search-input"]').setValue('光谱')
    await wrapper.find('[data-test="filter-booking-no"]').setValue('1001')
    await wrapper.find('[data-test="filter-phone"]').setValue('138')
    await wrapper.find('[data-test="filter-voyage-code"]').setValue('VOY-2026')
    await wrapper.find('[data-test="filter-cruise-name"]').setValue('海洋')
    await wrapper.find('[data-test="filter-form"]').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/bookings', expect.objectContaining({ query: expect.objectContaining({ booking_no: '1001', phone: '138', voyage_code: 'VOY-2026', cruise_name: '海洋', keyword: '光谱' }) }))
  })

  it('supports standalone keyword search submission', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()

    await wrapper.find('[data-test="booking-search-input"]').setValue('BK20260307001')
    await wrapper.find('[data-test="booking-search-submit"]').trigger('click')
    await flushPromises()

    expect(mockRequest).toHaveBeenLastCalledWith('/bookings', expect.objectContaining({ query: expect.objectContaining({ keyword: 'BK20260307001' }) }))
  })

  it('status tab triggers status query', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()

    await wrapper.find('[data-test="tab-paid"]').trigger('click')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/bookings', expect.objectContaining({ query: expect.objectContaining({ status: 'paid' }) }))
  })

  it('shows pay and refund actions on row', async () => {
    const wrapper = mount(Page, mountOptions)
    await flushPromises()

    expect(wrapper.text()).toContain('处理支付')
    expect(wrapper.text()).toContain('处理退改')
    expect(wrapper.text()).toContain('查看详情')
  })

  it('exports current rows as csv', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/bookings/export') {
        return Promise.resolve(new Blob(['id,status\n1,paid'], { type: 'text/csv;charset=utf-8' }))
      }
      return Promise.resolve(successList([
        { id: 1, status: 'pending_payment', total_cents: 19900, voyage_id: 2, user_id: 9, booking_no: 'BK20260307001', phone: '13800000001', voyage_code: 'VOY-2026-Q1', cruise_name: '海洋光谱号' },
      ], 1))
    })
    const wrapper = mount(Page, mountOptions)
    await flushPromises()

    await wrapper.find('[data-test="booking-search-input"]').setValue('光谱')
    await wrapper.find('[data-test="filter-voyage-code"]').setValue('VOY-2026')
    await wrapper.find('[data-test="export"]').trigger('click')

    expect(mockRequest).toHaveBeenCalledWith('/bookings/export', expect.objectContaining({
      query: expect.objectContaining({ keyword: '光谱', voyage_code: 'VOY-2026' }),
      responseType: 'blob',
    }))
    expect(clickMock).toHaveBeenCalledTimes(1)
  })

  it('shows empty state when list is empty', async () => {
    mockRequest.mockResolvedValueOnce(successList([], 0))
    const wrapper = mount(Page, mountOptions)
    await flushPromises()

    expect(wrapper.find('[data-test="empty"]').exists()).toBe(true)
  })

  it('shows error when load failed', async () => {
    mockRequest.mockRejectedValueOnce(new Error('network error'))
    const wrapper = mount(Page, mountOptions)
    await flushPromises()

    expect(wrapper.find('[data-test="error"]').text()).toContain('network error')
  })

  it('deletes row and reloads', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/bookings/1' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve(successList([{ id: 2, status: 'paid', total_cents: 38000 }], 1))
    })

    const wrapper = mount(Page, { ...mountOptions, attachTo: document.body })
    await flushPromises()
    await (wrapper.vm as any).handleDelete(1)
    await (wrapper.vm as any).confirmDelete()
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/bookings/1', { method: 'DELETE' })
    expect((wrapper.vm as any).deleteDialogVisible).toBe(false)
    wrapper.unmount()
  })
})

afterAll(() => {
  createElementSpy.mockRestore()
})
