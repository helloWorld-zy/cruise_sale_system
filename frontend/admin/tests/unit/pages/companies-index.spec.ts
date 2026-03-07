import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/companies/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

describe('Companies index page', () => {
  const globalStubs = {
    NuxtLink: { template: '<a><slot /></a>' },
    AdminActionLink: { template: '<a class="admin-btn"><slot /></a>' },
    AdminPageHeader: { props: ['title'], template: '<div>{{ title }}<slot /><slot name="actions" /></div>' },
    AdminFormCard: { template: '<div><slot /></div>' },
    AdminDataCard: { template: '<div><slot /></div>' },
  }

  beforeEach(() => {
    mockRequest.mockReset()
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/companies' && !options) {
        return Promise.resolve({
          data: {
            list: [
              {
                id: 1,
                name: '皇家加勒比',
                english_name: 'Royal Caribbean',
                logo_url: 'https://img/logo.png',
                description: 'test desc',
              },
            ],
            total: 1,
          },
        })
      }
      if (url === '/companies' && options?.method === 'POST') {
        return Promise.resolve({ data: { id: 2 } })
      }
      if (url === '/upload/image' && options?.method === 'POST') {
        return Promise.resolve({ data: { url: 'http://127.0.0.1:8080/uploads/logo.png' } })
      }
      if (url === '/companies/1' && options?.method === 'DELETE') {
        return Promise.resolve({ data: { ok: true } })
      }
      return Promise.resolve({ data: {} })
    })
  })

  it('renders company columns and row data', async () => {
    const wrapper = mount(Page, { global: { stubs: globalStubs } })
    await flushPromises()

    expect(wrapper.text()).toContain('邮轮公司管理')
    expect(wrapper.text()).toContain('中文名')
    expect(wrapper.text()).toContain('英文名')
    expect(wrapper.text()).toContain('文字介绍')
    expect(wrapper.text()).toContain('皇家加勒比')
    expect(wrapper.text()).toContain('Royal Caribbean')
    const logoCell = wrapper.find('.company-logo-cell')
    expect(logoCell.exists()).toBe(true)
    const logoImage = logoCell.find('.company-logo-image')
    expect(logoImage.exists()).toBe(true)
    const editAction = wrapper.findAll('a').find((el) => el.text().trim() === '编辑')
    expect(editAction?.classes()).toContain('admin-btn')
  })

  it('creates company using form submit', async () => {
    const wrapper = mount(Page, { global: { stubs: globalStubs } })
    await flushPromises()

    const inputs = wrapper.findAll('input')
    await inputs[0]!.setValue('地中海邮轮')
    await inputs[1]!.setValue('MSC Cruises')

    const fileInput = wrapper.find('input[type="file"]')
    const file = new File(['abc'], 'logo.png', { type: 'image/png' })
    Object.defineProperty(fileInput.element, 'files', { value: [file], configurable: true })
    await fileInput.trigger('change')

    await wrapper.find('textarea').setValue('desc')

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/companies', expect.objectContaining({ method: 'POST' }))
    const createCall = mockRequest.mock.calls.find((call) => call[0] === '/companies' && call[1]?.method === 'POST')
    expect(createCall?.[1]?.body?.logo_url).toBe('http://127.0.0.1:8080/uploads/logo.png')
  })

  it('opens delete dialog and sends delete request after confirm', async () => {
    const wrapper = mount(Page, { global: { stubs: globalStubs }, attachTo: document.body })
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((btn) => btn.text().trim() === '删除')
    expect(deleteBtn).toBeTruthy()
    await deleteBtn!.trigger('click')
    await flushPromises()

    expect(document.body.textContent || '').toContain('确认删除公司')
    const confirmDeleteBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
      (btn.textContent || '').includes('确认删除'),
    )
    expect(confirmDeleteBtn).toBeTruthy()
    ;(confirmDeleteBtn as HTMLButtonElement).click()
    await flushPromises()

    expect(mockRequest).toHaveBeenCalledWith('/companies/1', { method: 'DELETE' })
    expect(document.body.textContent || '').not.toContain('确认删除公司')
    wrapper.unmount()
  })

  it('shows explicit chinese message when backend returns company-has-cruises code', async () => {
    mockRequest.mockImplementation((url: string, options?: any) => {
      if (url === '/companies' && !options) {
        return Promise.resolve({
          data: {
            list: [{ id: 1, name: '皇家加勒比', english_name: 'Royal Caribbean', logo_url: '', description: '' }],
            total: 1,
          },
        })
      }
      if (url === '/companies/1' && options?.method === 'DELETE') {
        return Promise.reject({ code: 42202, message: 'company has cruises' })
      }
      return Promise.resolve({ data: {} })
    })

    const wrapper = mount(Page, { global: { stubs: globalStubs }, attachTo: document.body })
    await flushPromises()

    const deleteBtn = wrapper.findAll('button').find((btn) => btn.text().trim() === '删除')
    expect(deleteBtn).toBeTruthy()
    await deleteBtn!.trigger('click')
    await flushPromises()

    const confirmDeleteBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
      (btn.textContent || '').includes('确认删除'),
    )
    expect(confirmDeleteBtn).toBeTruthy()
    ;(confirmDeleteBtn as HTMLButtonElement).click()
    await flushPromises()

    expect(wrapper.text()).toContain('删除失败：该公司下存在邮轮，请先处理关联邮轮后再删除。')
    wrapper.unmount()
  })

  it('blocks create when company name is empty', async () => {
    const wrapper = mount(Page, { global: { stubs: globalStubs } })
    await flushPromises()

    await wrapper.find('form').trigger('submit.prevent')
    await flushPromises()

    expect(wrapper.text()).toContain('请先修正表单校验错误')
    expect(wrapper.text()).toContain('请填写公司中文名')
    expect(mockRequest).not.toHaveBeenCalledWith('/companies', expect.objectContaining({ method: 'POST' }))
  })
})
