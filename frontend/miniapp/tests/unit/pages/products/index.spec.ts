import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { cleanup, fireEvent, render } from '@testing-library/vue'
import Page from '../../../../pages/products/index.vue'

const mockRequest = vi.fn()

vi.mock('../../../../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockImplementation((path: string) => {
    if (path === '/companies?page=1&page_size=50') {
      return Promise.resolve({
        data: {
          list: [
            { id: 1, name: '皇家加勒比', english_name: 'Royal Caribbean' },
            { id: 2, name: 'MSC地中海邮轮', english_name: 'MSC Cruises' },
          ],
        },
      })
    }
    if (path === '/cruises?page=1&page_size=100') {
      return Promise.resolve({
        data: {
          list: [
            { id: 11, name: '海洋光谱号', english_name: 'Spectrum of the Seas', company_id: 1 },
            { id: 12, name: '海洋量子号', english_name: 'Quantum of the Seas', company_id: 1 },
            { id: 21, name: '欧罗巴号', english_name: 'MSC Euribia', company_id: 2 },
          ],
        },
      })
    }
    if (path === '/voyages?page=1&page_size=100') {
      return Promise.resolve({
        data: {
          list: [
            { id: 101, cruise_id: 11, code: 'RC101', brief_info: '海洋光谱号 上海-福冈-上海 5天4晚', image_url: '', depart_date: '2026-05-23T00:00:00Z', min_price_cents: 409900 },
            { id: 102, cruise_id: 12, code: 'RC102', brief_info: '海洋量子号 上海-冲绳-上海 5天4晚', image_url: '', depart_date: '2026-05-27T00:00:00Z', min_price_cents: 379900 },
            { id: 201, cruise_id: 21, code: 'MSC201', brief_info: '欧罗巴号 上海-济州-上海 5天4晚', image_url: '', depart_date: '2026-05-31T00:00:00Z', min_price_cents: 389900 },
          ],
          total: 3,
        },
      })
    }
    if (path === '/voyages?cruise_id=11&page=1&page_size=100') {
      return Promise.resolve({
        data: {
          list: [
            { id: 101, cruise_id: 11, code: 'RC101', brief_info: '海洋光谱号 上海-福冈-上海 5天4晚', image_url: '', depart_date: '2026-05-23T00:00:00Z', min_price_cents: 409900 },
          ],
          total: 1,
        },
      })
    }
    if (path === '/voyages?keyword=%E5%85%89%E8%B0%B1&page=1&page_size=100') {
      return Promise.resolve({
        data: {
          list: [
            { id: 101, cruise_id: 11, code: 'RC101', brief_info: '海洋光谱号 上海-福冈-上海 5天4晚', image_url: '', depart_date: '2026-05-23T00:00:00Z', min_price_cents: 409900 },
          ],
          total: 1,
        },
      })
    }
    return Promise.reject(new Error(`unexpected path: ${path}`))
  })
})

afterEach(() => {
  cleanup()
})

describe('Miniapp all-products page', () => {
  it('renders company list and all voyages by default', async () => {
    const { findByText, queryByText } = render(Page)

    expect(await findByText('全部商品')).toBeTruthy()
    expect(await findByText('皇家加勒比')).toBeTruthy()
    expect(await findByText('MSC地中海邮轮')).toBeTruthy()
    expect(await findByText('海洋光谱号 上海-福冈-上海 5天4晚')).toBeTruthy()
    expect(await findByText('欧罗巴号 上海-济州-上海 5天4晚')).toBeTruthy()
    expect(queryByText('海洋光谱号')).toBeNull()
  })

  it('toggles cruises after clicking a company and reloads voyages after clicking a cruise', async () => {
    const { findByText, queryByText } = render(Page)

    await fireEvent.click(await findByText('皇家加勒比'))

    expect(await findByText('海洋光谱号')).toBeTruthy()
    expect(await findByText('海洋量子号')).toBeTruthy()
    expect(await findByText('欧罗巴号 上海-济州-上海 5天4晚')).toBeTruthy()

    await fireEvent.click(await findByText('海洋光谱号'))

    expect(mockRequest).toHaveBeenCalledWith('/voyages?cruise_id=11&page=1&page_size=100')
    expect(await findByText('海洋光谱号 上海-福冈-上海 5天4晚')).toBeTruthy()
    expect(queryByText('海洋量子号 上海-冲绳-上海 5天4晚')).toBeNull()
    expect(queryByText('欧罗巴号 上海-济州-上海 5天4晚')).toBeNull()

    await fireEvent.click(await findByText('皇家加勒比'))

    expect(queryByText('海洋光谱号')).toBeNull()
    expect(queryByText('海洋量子号')).toBeNull()
  })

  it('emits selected voyage after clicking a voyage card', async () => {
    const { emitted, findByText } = render(Page)

    await fireEvent.click(await findByText('海洋光谱号 上海-福冈-上海 5天4晚'))

    expect(emitted()['open-voyage']).toEqual([[101]])
  })

  it('requests voyages from backend when searching keyword', async () => {
    const { findByPlaceholderText, findByText, queryByText } = render(Page)

    const input = await findByPlaceholderText('请输入搜索的商品')
    await fireEvent.update(input, '光谱')

    expect(mockRequest).toHaveBeenCalledWith('/voyages?keyword=%E5%85%89%E8%B0%B1&page=1&page_size=100')
    expect(await findByText('海洋光谱号 上海-福冈-上海 5天4晚')).toBeTruthy()
    expect(queryByText('海洋量子号 上海-冲绳-上海 5天4晚')).toBeNull()
  })
})