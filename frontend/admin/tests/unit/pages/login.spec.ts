import { describe, test, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import LoginPage from '../../../app/pages/login.vue'

const mockRequest = vi.fn()
const mockSetToken = vi.fn()
const mockNavigate = vi.fn().mockResolvedValue(undefined)
vi.stubGlobal('useRuntimeConfig', () => ({ public: { apiBase: '/api/v1' } }))
vi.stubGlobal('$fetch', mockRequest)
vi.mock('../../../app/stores/auth', () => ({
    useAuthStore: () => ({ setToken: mockSetToken }),
}))
vi.stubGlobal('navigateTo', mockNavigate)

describe('Login page', () => {
    beforeEach(() => {
        mockRequest.mockReset()
        mockRequest.mockResolvedValue({ token: 'admin-token' })
        mockSetToken.mockClear()
        mockNavigate.mockClear()
    })

    test('renders login form', () => {
        const wrapper = mount(LoginPage)
        expect(wrapper.find('form').exists()).toBe(true)
    })

    test('renders username and password inputs', () => {
        const wrapper = mount(LoginPage)
        const inputs = wrapper.findAll('input')
        expect(inputs.length).toBeGreaterThanOrEqual(2)
    })

    test('validates form on submit', async () => {
        const wrapper = mount(LoginPage)
        await wrapper.find('form').trigger('submit.prevent')
        expect(wrapper.text()).toContain('请填写用户名和密码')
        expect(wrapper.find('p.text-red-500').exists()).toBe(true)
    })

    test('handles valid submit', async () => {
        const wrapper = mount(LoginPage)
        const inputs = wrapper.findAll('input')
        await inputs[0]!.setValue('admin')
        await inputs[1]!.setValue('password')

        await wrapper.find('form').trigger('submit.prevent')
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/api/v1/admin/auth/login', expect.objectContaining({ method: 'POST' }))
        expect(mockSetToken).toHaveBeenCalledWith('admin-token')
        expect(mockNavigate).toHaveBeenCalledWith('/cruises')
    })

    test('api 失败时展示错误', async () => {
        mockRequest.mockRejectedValueOnce(new Error('login fail'))
        const wrapper = mount(LoginPage)
        const inputs = wrapper.findAll('input')
        await inputs[0]!.setValue('admin')
        await inputs[1]!.setValue('password')
        await wrapper.find('form').trigger('submit.prevent')
        await flushPromises()
        expect(wrapper.text()).toContain('login fail')
    })
})
