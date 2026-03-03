import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { cleanup, render } from '@testing-library/vue'
import Detail from '../pages/cabin/detail.vue'

const mockRequest = vi.fn()

vi.mock('../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

describe('Miniapp Cabin Detail', () => {
  beforeEach(() => {
    mockRequest.mockReset()
    mockRequest.mockImplementation((url: string) => {
      if (url === '/cabins/7') {
        return Promise.resolve({
          data: {
            id: 7,
            code: 'A701',
            bed_type: '双床',
            has_window: true,
            has_balcony: true,
            position: '中部',
            orientation: '左舷',
            amenities: '早餐,WiFi',
            price_cents: 19900,
          },
        })
      }
      if (url === '/cabins/7/prices') {
        return Promise.resolve({ data: [{ date: '2026-03-03', price_cents: 19900 }] })
      }
      if (url === '/cabins/7/inventory') {
        return Promise.resolve({ data: { total: 20, locked: 4, sold: 5, available: 11 } })
      }
      return Promise.resolve({ data: {} })
    })
  })

  afterEach(() => {
    cleanup()
  })

  it('renders calendar and inventory status', async () => {
    const { findByText } = render(Detail, { props: { cabinSkuId: 7 } })

    expect(await findByText('库存充足')).toBeTruthy()
    expect(await findByText('日期：2026-03-03')).toBeTruthy()
    expect(await findByText('价格：¥199')).toBeTruthy()
  })
})
