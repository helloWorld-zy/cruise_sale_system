import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import LoginForm from '../../../components/LoginForm.vue'

const mockFetch = vi.fn()
vi.stubGlobal('$fetch', mockFetch)

describe('LoginForm', () => {
  beforeEach(() => {
    mockFetch.mockReset()
    sessionStorage.clear()
    vi.useFakeTimers()
  })

  afterEach(() => {
    vi.runOnlyPendingTimers()
    vi.useRealTimers()
  })

  it('validates phone and enables send code', async () => {
    const wrapper = mount(LoginForm)
    const sendButton = wrapper.find('button.send-btn')

    expect(sendButton.attributes('disabled')).toBeDefined()

    await wrapper.find('input#phone').setValue('13800000000')
    expect(sendButton.attributes('disabled')).toBeUndefined()
  })

  it('sends sms code and starts countdown', async () => {
    mockFetch.mockResolvedValueOnce({ ok: true })
    const wrapper = mount(LoginForm)

    await wrapper.find('input#phone').setValue('13800000000')
    await wrapper.find('button.send-btn').trigger('click')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/users/sms-code', {
      method: 'POST',
      body: { phone: '13800000000' },
    })
    expect(wrapper.text()).toContain('60s 后重新发送')

    vi.advanceTimersByTime(1000)
    await flushPromises()
    expect(wrapper.text()).toContain('59s 后重新发送')
  })

  it('shows send error message when sms request fails', async () => {
    mockFetch.mockRejectedValueOnce({ message: 'sms error' })
    const wrapper = mount(LoginForm)

    await wrapper.find('input#phone').setValue('13800000000')
    await wrapper.find('button.send-btn').trigger('click')
    await flushPromises()

    expect(wrapper.find('.error').text()).toContain('sms error')
  })

  it('submits login, stores token and emits event', async () => {
    mockFetch.mockResolvedValueOnce({ data: { token: 'token-abc' } })
    const wrapper = mount(LoginForm)

    await wrapper.find('input#phone').setValue('13800000000')
    await wrapper.find('input#code').setValue('123456')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(mockFetch).toHaveBeenCalledWith('/api/v1/users/login', {
      method: 'POST',
      body: { phone: '13800000000', code: '123456' },
    })
    expect(sessionStorage.getItem('auth_token')).toBe('token-abc')
    expect(wrapper.emitted('logged-in')).toEqual([['token-abc']])
    expect(wrapper.text()).toContain('登录成功')
  })

  it('shows login error message when submit fails', async () => {
    mockFetch.mockRejectedValueOnce({ data: { message: 'invalid code' } })
    const wrapper = mount(LoginForm)

    await wrapper.find('input#phone').setValue('13800000000')
    await wrapper.find('input#code').setValue('123456')
    await wrapper.find('form').trigger('submit')
    await flushPromises()

    expect(wrapper.find('.error').text()).toContain('invalid code')
  })
})
