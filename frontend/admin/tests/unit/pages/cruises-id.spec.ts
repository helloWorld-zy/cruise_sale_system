import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const routeMock = { params: { id: '5' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('useRoute', () => routeMock)

describe('Cruise edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/companies' && !options) {
        return Promise.resolve({ data: { list: [{ id: 1, name: '皇家加勒比' }] } })
      }
      if (url === '/cruises/5' && !options) {
        return Promise.resolve({
          data: {
            id: 5,
            name: 'Ocean Nova',
            code: 'ONOVA',
            company_id: 1,
            tonnage: 90000,
            passenger_capacity: 2400,
            crew_count: 900,
            build_year: 2018,
            refurbish_year: 2024,
            length: 310,
            width: 38,
            deck_count: 14,
            status: 1,
          },
        })
      }
      if (url === '/cruises/5' && options?.method === 'PUT') {
        return Promise.resolve({ data: { ok: true } })
      }
      if (url === '/cruises/5' && options?.method === 'DELETE') {
        return Promise.resolve({ data: { ok: true } })
      }
      return Promise.resolve({ data: {} })
    })
  })

  it('loads detail on mounted', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('编辑邮轮 #5')
    expect(wrapper.text()).toContain('所属公司')
    expect(mockRequest).toHaveBeenCalledWith('/cruises/5')
    expect((wrapper.find('input').element as HTMLInputElement).value).toContain('Ocean Nova')
  })

  it('saves form then navigates', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    const form = wrapper.find('form')
    await form.trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/cruises/5', expect.objectContaining({ method: 'PUT' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/cruises')
  })

  it('deletes after confirm then navigates', async () => {
    const wrapper = mount(Page, { attachTo: document.body })
    await flushPromises()

    const deleteButton = wrapper.findAll('button').find((btn) => btn.text().includes('删除'))
    await deleteButton!.trigger('click')
    await flushPromises()

    const confirmBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
      (btn.textContent || '').includes('确认删除'),
    )
    expect(confirmBtn).toBeTruthy()
    ;(confirmBtn as HTMLButtonElement).click()
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/cruises/5', { method: 'DELETE' })
    expect(document.body.textContent || '').not.toContain('确认删除邮轮')
    expect(mockNavigateTo).toHaveBeenCalledWith('/cruises')
    wrapper.unmount()
  })

  it('does not delete when confirm is cancelled', async () => {
    const wrapper = mount(Page, { attachTo: document.body })
    await flushPromises()

    const deleteButton = wrapper.findAll('button').find((btn) => btn.text().includes('删除'))
    await deleteButton!.trigger('click')
    await flushPromises()

    const cancelBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
      (btn.textContent || '').includes('取消'),
    )
    expect(cancelBtn).toBeTruthy()
    ;(cancelBtn as HTMLButtonElement).click()
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/cruises/5', { method: 'DELETE' })
    wrapper.unmount()
  })

  it('shows save error when update request fails', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/companies' && !options) {
        return Promise.resolve({ data: { list: [{ id: 1, name: '皇家加勒比' }] } })
      }
      if (url === '/cruises/5' && !options) {
        return Promise.resolve({ data: { id: 5, name: 'Ocean Nova', code: 'ONOVA', company_id: 1 } })
      }
      if (url === '/cruises/5' && options?.method === 'PUT') {
        return Promise.reject(new Error('update failed'))
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.text()).toContain('update failed')
  })

  it('handles load error', async () => {
    mockRequest.mockRejectedValueOnce(new Error('load failed'))
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('load failed')
  })

  it('blocks save when required fields are invalid', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/companies' && !options) {
        return Promise.resolve({ data: { list: [] } })
      }
      if (url === '/cruises/5' && !options) {
        return Promise.resolve({ data: { id: 5, name: '', code: 'ONOVA', company_id: 0 } })
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(Page)
    await flushPromises()

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.text()).toContain('请先修正表单校验错误')
    expect(wrapper.text()).toContain('请填写邮轮名称')
    expect(wrapper.text()).toContain('请选择所属公司')
    expect(mockRequest).not.toHaveBeenCalledWith('/cruises/5', expect.objectContaining({ method: 'PUT' }))
  })
})
