import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { cleanup, fireEvent, render } from '@testing-library/vue'
import App from '../../src/App.vue'

const mockRequest = vi.fn()

vi.mock('../../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockImplementation((path: string) => {
    if (path === '/companies?page=1&page_size=50') {
      return Promise.resolve({ data: { list: [{ id: 1, name: '皇家加勒比', english_name: 'Royal Caribbean' }] } })
    }
    if (path === '/cruises?page=1&page_size=30') {
      return Promise.resolve({ data: { list: [{ id: 11, name: '海洋绿洲号', english_name: 'Oasis of the Seas', company_id: 1 }] } })
    }
    if (path === '/cruises/11') {
      return Promise.resolve({ data: { id: 11, name: '海洋绿洲号', tonnage: 236000, passenger_capacity: 6700, length: 362, deck_count: 18 } })
    }
    if (path === '/cabin-types?cruise_id=11&page=1&page_size=20') {
      return Promise.resolve({ data: { list: [] } })
    }
    if (path === '/facilities?cruise_id=11') {
      return Promise.resolve({ data: [] })
    }
    return Promise.reject(new Error(`unexpected path: ${path}`))
  })
})

afterEach(() => {
  cleanup()
})

describe('Miniapp app shell', () => {
  it('opens cruise detail after clicking a wiki cruise card', async () => {
    const { getByText, findByText } = render(App)

    await fireEvent.click(getByText('邮轮百科'))
    await fireEvent.click(await findByText('海洋绿洲号'))

    expect(await findByText('舱房类型')).toBeTruthy()
    expect(await findByText('海洋绿洲号')).toBeTruthy()
  })
})