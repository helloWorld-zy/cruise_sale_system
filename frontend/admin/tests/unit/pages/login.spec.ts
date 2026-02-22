import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import LoginPage from '../../../app/pages/login.vue'

// ME-09 FIX: Removed vi.stubGlobal('ref', ...) hack that broke Vue reactivity testing.
// Stubs are used for Nuxt UI components (UInput, UButton) which are not loaded in test env.
describe('Login page', () => {
    test('renders login form', () => {
        const wrapper = mount(LoginPage, {
            global: {
                stubs: { UInput: true, UButton: true },
            },
        })
        expect(wrapper.find('form').exists()).toBe(true)
    })

    test('renders username and password inputs', () => {
        const wrapper = mount(LoginPage, {
            global: {
                stubs: { UInput: true, UButton: true },
            },
        })
        // Stubbed UInputs should still render
        const inputs = wrapper.findAllComponents({ name: 'UInput' })
        expect(inputs.length).toBeGreaterThanOrEqual(2)
    })
})
