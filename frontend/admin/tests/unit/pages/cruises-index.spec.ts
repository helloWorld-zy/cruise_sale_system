import { describe, test, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../app/pages/cruises/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

describe('Cruise list', () => {
    beforeEach(() => {
        mockRequest.mockReset()
        mockRequest.mockResolvedValue({ data: [{ id: 1, name: '示例邮轮', company_id: 9, status: 'on' }] })
    })

    const UButtonStub = {
        name: 'UButton',
        props: ['to'],
        template: '<button :data-to="to"><slot /></button>'
    }

    test('renders a table', () => {
        const wrapper = mount(Page, {
            global: {
                stubs: { UButton: UButtonStub }
            }
        })
        expect(wrapper.find('table').exists()).toBe(true)
    })

    test('renders create cruise button', () => {
        const wrapper = mount(Page, {
            global: {
                stubs: { UButton: UButtonStub }
            }
        })
        const btn = wrapper.findComponent({ name: 'UButton' })
        expect(btn.exists()).toBe(true)
        expect(btn.attributes('data-to')).toBe('/cruises/create')
    })

    test('calls cruises api and renders data', async () => {
        const wrapper = mount(Page, {
            global: {
                stubs: { UButton: UButtonStub }
            }
        })
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/cruises')
        expect(wrapper.text()).toContain('示例邮轮')
    })

    test('shows error when api failed', async () => {
        mockRequest.mockRejectedValueOnce(new Error('load error'))
        const wrapper = mount(Page, {
            global: {
                stubs: { UButton: UButtonStub }
            }
        })
        await flushPromises()
        expect(wrapper.text()).toContain('load error')
    })
})
