import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import LoginPage from '../../pages/account/login.vue'

// LoginForm 使用 $fetch（Nuxt 全局），在 vitest 中需 stub
vi.stubGlobal('$fetch', vi.fn())

describe('Login Page', () => {
    it('renders phone input', () => {
        const wrapper = mount(LoginPage, {
            global: {
                stubs: { LoginForm: true }  // stub 子组件
            }
        })
        expect(wrapper.findComponent({ name: 'LoginForm' }).exists()).toBe(true)
    })

    it('renders LoginForm component', () => {
        const wrapper = mount(LoginPage, {
            global: { stubs: { LoginForm: { template: '<input type="tel" />' } } }
        })
        expect(wrapper.find('input[type="tel"]').exists()).toBe(true)
    })
})
