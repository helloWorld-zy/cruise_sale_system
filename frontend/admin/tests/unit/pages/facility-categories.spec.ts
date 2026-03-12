import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/facility-categories/index.vue'

const mockRequest = vi.fn()
const confirmMock = vi.fn(() => true)
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('confirm', confirmMock)

beforeEach(() => {
  mockRequest.mockReset()
  confirmMock.mockClear()
  confirmMock.mockReturnValue(true)
  mockRequest.mockImplementation((url: string, options?: any) => {
    if (url === '/facility-categories' && !options) {
      return Promise.resolve({
        data: [{ id: 1, name: '餐饮', icon: 'utensils', sort_order: 1, status: 1 }],
      })
    }
    if (url === '/facility-categories' && options?.method === 'POST') return Promise.resolve({ data: { id: 2 } })
    if (url === '/facility-categories/1' && options?.method === 'PUT') return Promise.resolve({ data: { ok: true } })
    if (url === '/facility-categories/1' && options?.method === 'DELETE') return Promise.resolve({ data: { ok: true } })
    return Promise.resolve({ data: [] })
  })
})

describe('FacilityCategory list', () => {
  const globalStubs = {
    NuxtLink: { template: '<a><slot /></a>' },
    AdminActionLink: { template: '<a class="admin-btn"><slot /></a>' },
    AdminPageHeader: { props: ['title'], template: '<div>{{ title }}<slot /><slot name="actions" /></div>' },
    AdminDataCard: { template: '<div><slot /></div>' },
  }

  it('renders compact table with inline inputs', async () => {
    const wrapper = mount(Page, {
      global: { stubs: globalStubs },
    })
    await flushPromises()
    expect(wrapper.find('table').exists()).toBe(true)
    expect(wrapper.find('tbody input').exists()).toBe(true)
  })

  it('loads facility categories from api', async () => {
    const wrapper = mount(Page, {
      global: { stubs: globalStubs },
    })
    await flushPromises()
    expect(mockRequest).toHaveBeenCalledWith('/facility-categories')
    const nameInput = wrapper.find('tbody tr td input')
    expect((nameInput.element as HTMLInputElement).value).toBe('餐饮')
  })

  it('appends a new editable row', async () => {
    const wrapper = mount(Page, {
      global: { stubs: globalStubs },
    })
    await flushPromises()

    const addBtn = wrapper.findAll('button').find((b) => b.text().includes('新增分类'))
    await addBtn!.trigger('click')
    await flushPromises()

    const rows = wrapper.findAll('tbody tr')
    expect(rows.length).toBeGreaterThan(1)
  })

  it('saves existing row with PUT', async () => {
    const wrapper = mount(Page, {
      global: { stubs: globalStubs },
    })
    await flushPromises()

    const saveBtn = wrapper.findAll('button').find((b) => b.text().includes('保存'))
    await saveBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/facility-categories/1', expect.objectContaining({ method: 'PUT' }))
  })

  it('opens icon picker and selects icon', async () => {
    const wrapper = mount(Page, {
      global: { stubs: globalStubs },
    })
    await flushPromises()

    const selectBtn = wrapper.find('[data-test="facility-category-icon-trigger-1"]')
    await selectBtn!.trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-test="facility-category-icon-option-1-included-dining"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-icon-option-1-specialty-dining"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-icon-option-1-bar-lounge"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-icon-option-1-leisure-entertainment"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-icon-option-1-kids-family"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-icon-option-1-suite-privilege"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-icon-option-1-sports-fitness"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="facility-category-icon-option-1-other"]').exists()).toBe(true)

    const iconBtn = wrapper.find('[data-test="facility-category-icon-option-1-bar-lounge"]')
    expect(iconBtn.find('svg').exists()).toBe(true)
    await iconBtn!.trigger('click')
    await flushPromises()

    expect(wrapper.find('[data-test="facility-category-icon-current-1"]').attributes('data-icon')).toBe('bar-lounge')
    expect(wrapper.find('[data-test="facility-category-icon-input-1"]').exists()).toBe(false)
  })

  it('does not delete when confirm is cancelled', async () => {
    confirmMock.mockReturnValueOnce(false)
    const wrapper = mount(Page, {
      global: { stubs: globalStubs },
    })
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((b) => b.text().includes('删除'))
    await deleteBtn!.trigger('click')
    await flushPromises()

    expect(mockRequest).not.toHaveBeenCalledWith('/facility-categories/1', { method: 'DELETE' })
  })
})
