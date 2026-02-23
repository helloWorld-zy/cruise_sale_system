import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import LoginPage from '../../pages/account/login.vue'

describe('Login Page', () => {
    it('renders phone input', () => {
        const wrapper = mount(LoginPage)
        expect(wrapper.find('input').exists()).toBe(true)
    })
})
