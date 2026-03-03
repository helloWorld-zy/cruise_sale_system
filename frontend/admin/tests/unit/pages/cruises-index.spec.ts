import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
const confirmMock = vi.fn(() => true)
vi.stubGlobal('confirm', confirmMock)

describe('Cruise list enhanced', () => {
    beforeEach(() => {
        mockRequest.mockReset()
        confirmMock.mockClear()
        confirmMock.mockReturnValue(true)
        mockRequest.mockResolvedValue({ data: { list: [{ id: 1, name: '示例邮轮', code: 'HARMONY', company_id: 9, status: 1 }], total: 1 } })
    })

    it('renders filter controls', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: { NuxtLink: { template: '<a><slot /></a>' } }
            }
        })
        await flushPromises()
        expect(wrapper.find('[data-test="filter-keyword"]').exists()).toBe(true)
        expect(wrapper.find('[data-test="filter-status"]').exists()).toBe(true)
        expect(wrapper.find('[data-test="batch-action"]').exists()).toBe(false)
    })

    it('shows batch action after selecting rows', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: { NuxtLink: { template: '<a><slot /></a>' } }
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
                stubs: { NuxtLink: { template: '<a><slot /></a>' } }
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
                stubs: { NuxtLink: { template: '<a><slot /></a>' } }
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
                stubs: { NuxtLink: { template: '<a><slot /></a>' } }
            }
        })
        await flushPromises()

        await (wrapper.vm as any).handleDelete('abc')
        await flushPromises()

        expect(wrapper.text()).toContain('无效记录 ID，无法删除')
    })
})
