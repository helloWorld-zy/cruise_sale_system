import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/vue'
import Page from '../../../../pages/cruise/detail.vue'

const mockRequest = vi.fn()
vi.mock('../../../../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

describe('Miniapp cruise detail page', () => {
  it('renders key sections', async () => {
    mockRequest
      .mockResolvedValueOnce({ data: { id: 9, name: 'Atlantic Dream', tonnage: 140000, passenger_capacity: 3200, length: 315, deck_count: 16 } })
      .mockResolvedValueOnce({ data: { list: [{ id: 1, name: '阳台房', area_min: 24, max_capacity: 3, min_price_cents: 188000 }] } })
      .mockResolvedValueOnce({ data: [{ id: 3, category_id: 1, name: '海上剧院', extra_charge: false }] })
      .mockResolvedValueOnce({ data: [{ id: 5, name: '地中海之旅', departure_date: '2026-06-01', min_price_cents: 256000 }] })

    const { findByText } = render(Page, { props: { cruiseId: 9 } })
    expect(await findByText('Atlantic Dream')).toBeTruthy()
    expect(await findByText('舱房类型')).toBeTruthy()
    expect(await findByText('设施导览')).toBeTruthy()
    expect(await findByText('关联航线')).toBeTruthy()
  })
})
