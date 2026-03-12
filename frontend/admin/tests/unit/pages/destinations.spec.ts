import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import Page from '../../../app/pages/destinations/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
const mockCreateObjectURL = vi.fn(() => 'blob:ports')
const mockClick = vi.fn()
const originalCreateElement = document.createElement.bind(document)

vi.stubGlobal('URL', { createObjectURL: mockCreateObjectURL, revokeObjectURL: vi.fn() })

describe('Port city dictionary page', () => {
	beforeEach(() => {
		mockRequest.mockReset()
		mockCreateObjectURL.mockClear()
		mockClick.mockClear()
		vi.spyOn(document, 'createElement').mockImplementation(((tagName: string) => {
			if (tagName === 'a') {
				return { click: mockClick, set href(_v: string) {}, set download(_v: string) {} } as any
			}
			return originalCreateElement(tagName)
		}) as any)
		mockRequest.mockImplementation((url: string, options?: any) => {
			if (url === '/custom-destinations' && !options) {
				return Promise.resolve({ data: [{ id: 1, name: '迈阿密', country: '美国', latitude: 25.7617, longitude: -80.1918, keywords: '迈阿密,miami', sort_order: 100, status: 1 }] })
			}
			if (url === '/custom-destinations/export') {
				return Promise.resolve('name,country,latitude,longitude,keywords,sort_order,status,description\n迈阿密,美国,25.7617,-80.1918,"迈阿密,miami",100,1,seed\n')
			}
			if (url === '/custom-destinations/import' && options?.method === 'POST') {
				return Promise.resolve({ data: { imported: 1 } })
			}
			if (url === '/custom-destinations' && options?.method === 'POST') {
				return Promise.resolve({ data: { id: 2 } })
			}
			return Promise.resolve({ data: [] })
		})
	})

	it('shows port city dictionary copy and loaded seeded ports', async () => {
		const wrapper = mount(Page)
		await flushPromises()

		expect(wrapper.text()).toContain('港口城市词典')
		expect(wrapper.text()).toContain('城市坐标会直接用于行程地图与站点坐标补全')
		expect((wrapper.find('tbody tr input').element as HTMLInputElement).value).toBe('迈阿密')
	})

	it('blocks save when coordinates are missing', async () => {
		const wrapper = mount(Page)
		await flushPromises()

		await wrapper.find('button.dest-new-btn').trigger('click')
		const inputs = wrapper.findAll('tbody tr input')
		await inputs[0]!.setValue('布宜诺斯艾利斯')
		await inputs[1]!.setValue('阿根廷')
		await inputs[4]!.setValue('布宜诺斯艾利斯,buenos aires')
		await wrapper.findAll('tbody tr .admin-btn')[0]!.trigger('click')
		await flushPromises()

		expect(mockRequest).not.toHaveBeenCalledWith('/custom-destinations', expect.objectContaining({ method: 'POST' }))
		expect(wrapper.text()).toContain('纬度和经度不能为空')
	})

	it('exports dictionary csv and triggers download', async () => {
		const wrapper = mount(Page)
		await flushPromises()

		await wrapper.find('[data-test="dest-export"]').trigger('click')
		await flushPromises()

		expect(mockRequest).toHaveBeenCalledWith('/custom-destinations/export', expect.objectContaining({ responseType: 'text' }))
		expect(mockCreateObjectURL).toHaveBeenCalled()
		expect(mockClick).toHaveBeenCalled()
	})

	it('imports dictionary csv via form-data upload', async () => {
		const wrapper = mount(Page)
		await flushPromises()

		const file = new File(['name,country,latitude,longitude,keywords,sort_order,status,description\n温哥华,加拿大,49.2827,-123.1207,"温哥华,vancouver",88,1,manual\n'], 'ports.csv', { type: 'text/csv' })
		const input = wrapper.find('[data-test="dest-import-input"]')
		Object.defineProperty(input.element, 'files', { value: [file] })
		await input.trigger('change')
		await flushPromises()

		expect(mockRequest).toHaveBeenCalledWith('/custom-destinations/import', expect.objectContaining({ method: 'POST', body: expect.any(FormData) }))
	})

	it('filters dictionary rows by name country or keywords', async () => {
		mockRequest.mockImplementationOnce(() => Promise.resolve({
			data: [
				{ id: 1, name: '迈阿密', country: '美国', latitude: 25.7617, longitude: -80.1918, keywords: '迈阿密,miami', sort_order: 100, status: 1 },
				{ id: 2, name: '布宜诺斯艾利斯', country: '阿根廷', latitude: -34.6037, longitude: -58.3816, keywords: '布宜诺斯艾利斯,buenos aires', sort_order: 90, status: 1 },
			],
		}))

		const wrapper = mount(Page)
		await flushPromises()

		await wrapper.find('[data-test="dest-search"]').setValue('buenos')
		await flushPromises()

		const rows = wrapper.findAll('tbody tr')
		expect(rows).toHaveLength(1)
		expect((wrapper.find('tbody tr input').element as HTMLInputElement).value).toBe('布宜诺斯艾利斯')
		expect(wrapper.text()).not.toContain('迈阿密')
	})
})