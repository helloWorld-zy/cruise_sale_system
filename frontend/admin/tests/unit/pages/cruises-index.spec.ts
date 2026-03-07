import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

describe('Cruise list enhanced', () => {
    const globalStubs = {
        NuxtLink: { template: '<a><slot /></a>' },
        AdminActionLink: { template: '<a><slot /></a>' },
        AdminPageHeader: { template: '<div><slot /><slot name="actions" /></div>' },
        AdminFilterBar: { template: '<div><slot /></div>' },
        AdminDataCard: { template: '<div><slot /></div>' },
        AdminStatusTag: { props: ['text'], template: '<span>{{ text }}</span>' },
    }

    beforeEach(() => {
        mockRequest.mockReset()
        mockRequest.mockImplementation((url: string, options?: any) => {
            if (url === '/companies' && !options) {
                return Promise.resolve({ data: { list: [{ id: 9, name: '示例公司' }] } })
            }
            if (url === '/cruises/batch-status' && options?.method === 'PUT') {
                return Promise.resolve({ data: { ok: true } })
            }
            return Promise.resolve({ data: { list: [{ id: 1, name: '示例邮轮', code: 'HARMONY', company_id: 9, status: 1 }], total: 1 } })
        })
    })

    it('renders filter controls', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: globalStubs,
            }
        })
        await flushPromises()
        expect(wrapper.find('[data-test="filter-keyword"]').exists()).toBe(true)
        expect(wrapper.find('[data-test="filter-status"]').exists()).toBe(true)
        expect(wrapper.text()).toContain('示例公司')
        expect(wrapper.find('[data-test="batch-action"]').exists()).toBe(false)
    })

    it('shows batch action after selecting rows', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: globalStubs,
            }
        })
        await flushPromises()
        const firstRowCheckbox = wrapper.find('tbody input[type="checkbox"]')
        await firstRowCheckbox.setValue(true)
        expect(wrapper.find('[data-test="batch-action"]').exists()).toBe(true)
    })

    it('maps filter params and submits trimmed keyword/status/sort', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: globalStubs,
            }
        })
        await flushPromises()

        const keyword = wrapper.find('[data-test="filter-keyword"]')
        await keyword.setValue('  nova  ')

        const status = wrapper.find('[data-test="filter-status"]')
        await status.setValue('-1')

        const sort = wrapper.findAll('select')[1]
        await sort!.setValue('tonnage_desc')

        const filterButton = wrapper.findAll('button').find((btn) => btn.text().includes('筛选'))
        await filterButton!.trigger('click')
        await flushPromises()

        expect(mockRequest).toHaveBeenLastCalledWith('/cruises', {
            query: expect.objectContaining({
                page: 1,
                page_size: 10,
                keyword: 'nova',
                status: 0,
                sort_by: 'tonnage_desc',
            }),
        })
    })

    it('calls batch status endpoint after selecting item', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: globalStubs,
            }
        })
        await flushPromises()

        const rowCheckbox = wrapper.find('tbody input[type="checkbox"]')
        await rowCheckbox.setValue(true)

        const downBtn = wrapper.findAll('[data-test="batch-action"] button').find((btn) => btn.text().includes('批量下架'))
        await downBtn!.trigger('click')
        await flushPromises()

        expect(mockRequest).toHaveBeenCalledWith('/cruises/batch-status', {
            method: 'PUT',
            body: { ids: [1], status: -1 },
        })
    })

    it('shows invalid id error when delete is called with bad id', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: globalStubs,
            }
        })
        await flushPromises()

        await (wrapper.vm as any).handleDelete('abc')
        await flushPromises()

        expect(wrapper.text()).toContain('无效记录 ID，无法删除')
    })

    it('opens delete modal and sends delete request after confirming', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: globalStubs,
            },
            attachTo: document.body,
        })
        await flushPromises()

        const deleteBtn = wrapper.findAll('button').find((btn) => btn.text().trim() === '删除')
        expect(deleteBtn).toBeTruthy()
        await deleteBtn!.trigger('click')
        await flushPromises()

        expect(document.body.textContent || '').toContain('确认删除邮轮')
        const confirmBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
            (btn.textContent || '').includes('确认删除'),
        )
        expect(confirmBtn).toBeTruthy()
        ;(confirmBtn as HTMLButtonElement).click()
        await flushPromises()

        expect(mockRequest).toHaveBeenCalledWith('/cruises/1', { method: 'DELETE' })
        expect(document.body.textContent || '').not.toContain('确认删除邮轮')
        wrapper.unmount()
    })

    it('shows explicit message when cruise delete is blocked by voyages', async () => {
        mockRequest.mockImplementation((url: string, options?: any) => {
            if (url === '/companies' && !options) {
                return Promise.resolve({ data: { list: [{ id: 9, name: '示例公司' }] } })
            }
            if (url === '/cruises/1' && options?.method === 'DELETE') {
                return Promise.reject({ code: 42204, message: 'cruise has voyages' })
            }
            return Promise.resolve({ data: { list: [{ id: 1, name: '示例邮轮', code: 'HARMONY', company_id: 9, status: 1 }], total: 1 } })
        })

        const wrapper = mount(Page, {
            global: {
                stubs: globalStubs,
            },
            attachTo: document.body,
        })
        await flushPromises()

        const deleteBtn = wrapper.findAll('button').find((btn) => btn.text().trim() === '删除')
        expect(deleteBtn).toBeTruthy()
        await deleteBtn!.trigger('click')
        await flushPromises()

        const confirmBtn = Array.from(document.body.querySelectorAll('button')).find((btn) =>
            (btn.textContent || '').includes('确认删除'),
        )
        expect(confirmBtn).toBeTruthy()
        ;(confirmBtn as HTMLButtonElement).click()
        await flushPromises()

        expect(wrapper.text()).toContain('删除失败：该邮轮下存在航次，请先处理关联航次后再删除。')
        wrapper.unmount()
    })
})
