import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import Page from '../../../app/pages/content-templates/index.vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))

describe('Content templates page', () => {
	beforeEach(() => {
		mockRequest.mockReset()
		mockRequest.mockImplementation((url: string, options?: any) => {
			if (url === '/content-templates' && !options) {
				return Promise.resolve({ data: [{ id: 1, name: '默认费用说明', kind: 'fee_note', status: 1, content: { included: [{ text: '船票' }, { text: '港务费' }], excluded: [{ text: '签证费' }] } }] })
			}
			if (url === '/content-templates?kind=booking_notice' && !options) {
				return Promise.resolve({ data: [{ id: 2, name: '默认预订须知', kind: 'booking_notice', status: 1, content: { sections: [{ key: 'documents', title: '出行证件', items: [{ text: '请携带护照' }, { text: '儿童需带出生证明' }] }, { key: 'cancel', title: '退改说明', items: [{ text: '开航前 7 天内不可退' }] }] } }] })
			}
			if (url === '/content-templates' && options?.method === 'POST') {
				return Promise.resolve({ data: { id: 3, name: '春节通用模板' } })
			}
			if (url === '/content-templates/1' && options?.method === 'PUT') {
				return Promise.resolve({ data: { id: 1, name: '默认费用说明-更新' } })
			}
			return Promise.resolve({ data: [] })
		})
	})

	it('loads fee-note templates by default and switches to booking notice templates', async () => {
		const wrapper = mount(Page)
		await flushPromises()

		expect(mockRequest).toHaveBeenCalledWith('/content-templates')
		expect(wrapper.text()).toContain('默认费用说明')

		await wrapper.find('[data-test="kind-booking-notice"]').trigger('click')
		await flushPromises()

		expect(mockRequest).toHaveBeenCalledWith('/content-templates?kind=booking_notice')
		expect(wrapper.text()).toContain('默认预订须知')
	})

	it('creates a fee-note template with structured content', async () => {
		const wrapper = mount(Page)
		await flushPromises()

		await wrapper.find('[data-test="template-name"]').setValue('春节通用模板')
		await wrapper.find('[data-test="fee-included-0"]').setValue('邮轮船票')
		await wrapper.find('[data-test="fee-excluded-0"]').setValue('签证费用')
		await wrapper.find('form').trigger('submit.prevent')
		await flushPromises()

		expect(mockRequest).toHaveBeenCalledWith(
			'/content-templates',
			expect.objectContaining({
				method: 'POST',
				body: expect.objectContaining({
					name: '春节通用模板',
					kind: 'fee_note',
					content: expect.objectContaining({
						included: [{ text: '邮轮船票' }],
						excluded: [{ text: '签证费用' }],
					}),
				}),
			}),
		)
	})

	it('supports multi-item fee notes and multi-section booking notices', async () => {
		const wrapper = mount(Page)
		await flushPromises()
		await wrapper.find('[data-test="edit-template-1"]').trigger('click')
		await flushPromises()

		expect(wrapper.find('[data-test="fee-included-1"]').exists()).toBe(true)
		expect(wrapper.find('[data-test="fee-excluded-0"]').exists()).toBe(true)

		await wrapper.find('[data-test="kind-booking-notice"]').trigger('click')
		await flushPromises()
		await wrapper.find('[data-test="edit-template-2"]').trigger('click')
		await flushPromises()

		expect(wrapper.find('[data-test="booking-section-title-0"]').element).toBeTruthy()
		expect(wrapper.find('[data-test="booking-item-0-1"]').exists()).toBe(true)
		expect(wrapper.find('[data-test="booking-section-title-1"]').exists()).toBe(true)
	})

	it('creates a booking notice template with multiple sections and items', async () => {
		const wrapper = mount(Page)
		await flushPromises()

		await wrapper.find('[data-test="kind-booking-notice"]').trigger('click')
		await flushPromises()
		await wrapper.find('[data-test="template-name"]').setValue('日韩线须知')
		await wrapper.find('[data-test="booking-section-title-0"]').setValue('出行证件')
		await wrapper.find('[data-test="booking-section-key-0"]').setValue('documents')
		await wrapper.find('[data-test="booking-item-0-0"]').setValue('请携带护照原件')
		await wrapper.find('[data-test="add-booking-item-0"]').trigger('click')
		await wrapper.find('[data-test="booking-item-0-1"]').setValue('儿童请携带出生证明')
		await wrapper.find('[data-test="add-booking-section"]').trigger('click')
		await wrapper.find('[data-test="booking-section-title-1"]').setValue('退改说明')
		await wrapper.find('[data-test="booking-section-key-1"]').setValue('cancel')
		await wrapper.find('[data-test="booking-item-1-0"]').setValue('开航前 7 天内不可退')
		await wrapper.find('form').trigger('submit.prevent')
		await flushPromises()

		expect(mockRequest).toHaveBeenCalledWith(
			'/content-templates',
			expect.objectContaining({
				method: 'POST',
				body: expect.objectContaining({
					kind: 'booking_notice',
					content: {
						sections: [
							{
								key: 'documents',
								title: '出行证件',
								items: [{ text: '请携带护照原件' }, { text: '儿童请携带出生证明' }],
							},
							{
								key: 'cancel',
								title: '退改说明',
								items: [{ text: '开航前 7 天内不可退' }],
							},
						],
					},
				}),
			}),
		)
	})

	it('loads an existing template into the form and updates it', async () => {
		const wrapper = mount(Page)
		await flushPromises()

		await wrapper.find('[data-test="edit-template-1"]').trigger('click')
		await flushPromises()
		await wrapper.find('[data-test="template-name"]').setValue('默认费用说明-更新')
		await wrapper.find('form').trigger('submit.prevent')
		await flushPromises()

		expect(mockRequest).toHaveBeenCalledWith(
			'/content-templates/1',
			expect.objectContaining({ method: 'PUT' }),
		)
	})
})