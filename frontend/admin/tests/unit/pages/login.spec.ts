import { describe, test, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import LoginPage from '../../../app/pages/login.vue'

const mockRequest = vi.fn()
const mockSetToken = vi.fn()
const mockNavigate = vi.fn().mockResolvedValue(undefined)
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.mock('../../../app/stores/auth', () => ({
    useAuthStore: () => ({ setToken: mockSetToken }),
}))
vi.stubGlobal('navigateTo', mockNavigate)

const UInputStub = {
    name: 'UInput',
    props: ['modelValue'],
    emits: ['update:modelValue'],
    template: '<input :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />'
}

const UButtonStub = {
    name: 'UButton',
    props: ['loading', 'type'],
    template: '<button :data-loading="String(loading)" :type="type"><slot /></button>'
}

// ME-09 FIX: Removed vi.stubGlobal('ref', ...) hack that broke Vue reactivity testing.
// Stubs are used for Nuxt UI components (UInput, UButton) which are not loaded in test env.
describe('Login page', () => {
    beforeEach(() => {
        mockRequest.mockReset()
        mockRequest.mockResolvedValue({ token: 'admin-token' })
        mockSetToken.mockClear()
        mockNavigate.mockClear()
    })

    test('renders login form', () => {
        const wrapper = mount(LoginPage, {
            global: {
                stubs: { UInput: UInputStub, UButton: UButtonStub },
            },
        })
        expect(wrapper.find('form').exists()).toBe(true)
    })

    test('renders username and password inputs', () => {
        const wrapper = mount(LoginPage, {
            global: {
                stubs: { UInput: UInputStub, UButton: UButtonStub },
            },
        })
        // Stubbed UInputs should still render
        const inputs = wrapper.findAllComponents({ name: 'UInput' })
        expect(inputs.length).toBeGreaterThanOrEqual(2)
    })

    test('validates form on submit', async () => {
        const wrapper = mount(LoginPage, {
            global: { stubs: { UInput: UInputStub, UButton: UButtonStub } },
        })
        await wrapper.find('form').trigger('submit.prevent')
        expect(wrapper.text()).toContain('请填写用户名和密码')
        expect(wrapper.find('p.text-red-500').exists()).toBe(true)
    })

    test('handles valid submit', async () => {
        const wrapper = mount(LoginPage, {
            global: { stubs: { UInput: UInputStub, UButton: UButtonStub } },
        })
        // Set values via instances or directly on data if possible. Since we stubbed UInput, finding them to setValue is tricky without knowing implementation.
        // We can interact with Vue's reactivity properly.
        const inputs = wrapper.findAllComponents({ name: 'UInput' })
        await inputs[0]!.vm.$emit('update:modelValue', 'admin')
        await inputs[1]!.vm.$emit('update:modelValue', 'password')

        await wrapper.find('form').trigger('submit.prevent')
        await flushPromises()
        expect(mockRequest).toHaveBeenCalledWith('/admin/auth/login', expect.objectContaining({ method: 'POST' }))
        expect(mockSetToken).toHaveBeenCalledWith('admin-token')
        expect(mockNavigate).toHaveBeenCalledWith('/dashboard')
    })

    test('api 失败时展示错误', async () => {
        mockRequest.mockRejectedValueOnce(new Error('login fail'))
        const wrapper = mount(LoginPage, {
            global: { stubs: { UInput: UInputStub, UButton: UButtonStub } },
        })
        const inputs = wrapper.findAllComponents({ name: 'UInput' })
        await inputs[0]!.vm.$emit('update:modelValue', 'admin')
        await inputs[1]!.vm.$emit('update:modelValue', 'password')
        await wrapper.find('form').trigger('submit.prevent')
        await flushPromises()
        expect(wrapper.text()).toContain('login fail')
    })
})
