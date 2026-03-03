import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/vue'
import Page from '../../../../pages/cruise/list.vue'

const mockRequest = vi.fn()
vi.mock('../../../../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

describe('Miniapp cruise list page', () => {
  it('renders cards from cruises api', async () => {
    mockRequest.mockResolvedValueOnce({ data: { list: [{ id: 1, name: 'Ocean Nova', tonnage: 120000, passenger_capacity: 2800, deck_count: 16 }] } })
    const { findByText } = render(Page)
    expect(await findByText('Ocean Nova')).toBeTruthy()
    expect(await findByText('查看详情 >')).toBeTruthy()
  })
})
