import { afterEach, beforeEach, describe, it, expect, vi } from 'vitest'
import { cleanup, fireEvent, render } from '@testing-library/vue'
import { createPinia, setActivePinia } from 'pinia'
import Page from '../../../pages/login/login.vue'

const mockRequest = vi.fn()
vi.mock('../../../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

beforeEach(() => {
  setActivePinia(createPinia())
  mockRequest.mockReset()
})

afterEach(() => cleanup())

describe('Miniapp Login', () => {
  it('发送验证码调用接口', async () => {
    mockRequest.mockResolvedValue({})
    const { getByPlaceholderText, getByText } = render(Page)
    await fireEvent.update(getByPlaceholderText('手机号'), '13800138000')
    await fireEvent.click(getByText('发送验证码'))
    expect(mockRequest).toHaveBeenCalledWith('/users/sms-code', expect.objectContaining({ method: 'POST' }))
  })

  it('登录调用接口并成功', async () => {
    mockRequest.mockResolvedValueOnce({ token: 't-1' })
    const { getByPlaceholderText, getByText, findByText } = render(Page)
    await fireEvent.update(getByPlaceholderText('手机号'), '13800138000')
    await fireEvent.update(getByPlaceholderText('验证码'), '1234')
    await fireEvent.click(getByText('登录'))
    expect(mockRequest).toHaveBeenCalledWith('/users/login', expect.objectContaining({ method: 'POST' }))
    expect(await findByText('登录成功')).toBeTruthy()
  })
})
