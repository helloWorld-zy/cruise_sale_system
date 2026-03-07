import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AuthLayout from '../../../app/layouts/auth.vue'

describe('Auth layout', () => {
  test('renders only auth container and default slot content', () => {
    const wrapper = mount(AuthLayout, {
      slots: { default: '<div data-test="auth-content">Login</div>' },
    })

    expect(wrapper.find('[data-test="auth-content"]').exists()).toBe(true)
    expect(wrapper.find('.admin-sidebar').exists()).toBe(false)
    expect(wrapper.find('.admin-header').exists()).toBe(false)
  })
})
