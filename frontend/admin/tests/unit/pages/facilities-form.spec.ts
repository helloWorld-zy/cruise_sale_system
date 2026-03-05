import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import EditPage from '../../../app/pages/facilities/[id].vue'

const mockRequest = vi.fn()
const mockNavigateTo = vi.fn().mockResolvedValue(undefined)
const routeMock = { params: { id: '9' } }

vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('navigateTo', mockNavigateTo)
vi.stubGlobal('useRoute', () => routeMock)

const mountOptions = {
  global: {
    stubs: {
      AdminActionLink: { template: '<a><slot /></a>' },
      NuxtLink: { template: '<a><slot /></a>' },
    },
  },
}

describe('Facilities edit page', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockNavigateTo.mockClear()

    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/cruises') return Promise.resolve({ data: [{ id: 1, name: 'Ocean Nova' }] })
      if (url === '/facility-categories') return Promise.resolve({ data: [{ id: 2, name: '娱乐' }] })
      if (url === '/facilities/9' && !options) {
        return Promise.resolve({
          data: {
            id: 9,
            cruise_id: 1,
            category_id: 2,
            name: '海上剧院',
            extra_charge: true,
            charge_price_tip: '¥200/人',
            target_audience: '家庭,情侣',
            status: 1,
          },
        })
      }
      if (url === '/facilities/9' && options?.method === 'PUT') return Promise.resolve({ data: { ok: true } })
      if (url === '/facilities/9' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
      return Promise.resolve({ data: {} })
    })
  })

  it('loads detail and renders extra charge section', async () => {
    const wrapper = mount(EditPage, mountOptions)
    await flushPromises()

    expect(wrapper.text()).toContain('编辑设施')
    expect(wrapper.text()).toContain('收费说明')
  })

  it('submits update and navigates', async () => {
    const wrapper = mount(EditPage, mountOptions)
    await flushPromises()

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/facilities/9', expect.objectContaining({ method: 'PUT' }))
    expect(mockNavigateTo).toHaveBeenCalledWith('/facilities')
  })

  it('deletes when confirmed', async () => {
    const wrapper = mount(EditPage, { ...mountOptions, attachTo: document.body })
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除'))
    await deleteBtn!.trigger('click')
    await flushPromises()

    const confirmBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
      (btn.textContent || '').includes('确认删除'),
    )
    expect(confirmBtn).toBeTruthy()
    ;(confirmBtn as HTMLButtonElement).click()
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/facilities/9', { method: 'DELETE' })
    expect(mockNavigateTo).toHaveBeenCalledWith('/facilities')
    wrapper.unmount()
  })
})
