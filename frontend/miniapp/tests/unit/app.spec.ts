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
    if (path === '/cruises?page=1&page_size=100') {
      return Promise.resolve({ data: { list: [{ id: 11, name: '海洋光谱号', english_name: 'Spectrum of the Seas', company_id: 1 }] } })
    }
    if (path === '/voyages?page=1&page_size=100') {
      return Promise.resolve({ data: { list: [{ id: 101, cruise_id: 11, code: 'RC101', brief_info: '海洋光谱号 上海-福冈-上海 5天4晚', min_price_cents: 409900 }] } })
    }
    if (path === '/voyages/101') {
      return Promise.resolve({ data: { id: 101, cruise_id: 11, code: 'RC101', brief_info: '海洋光谱号 上海-福冈-上海 5天4晚', depart_date: '2026-05-23T00:00:00Z', return_date: '2026-05-27T00:00:00Z', min_price_cents: 409900 } })
    }
    if (path === '/cruises/11') {
      return Promise.resolve({ data: { id: 11, name: '海洋绿洲号', english_name: 'Oasis of the Seas', tonnage: 236000, build_year: 2009, passenger_capacity: 6700, length: 362, deck_count: 18, description: '海洋绿洲号提供丰富的家庭度假体验。' } })
    }
    if (path === '/cabin-types?cruise_id=11&page=1&page_size=20') {
      return Promise.resolve({ data: { list: [] } })
    }
    if (path === '/facilities?cruise_id=11') {
      return Promise.resolve({ data: [{ id: 201, category_id: 3, category_name: '自助餐厅', name: '海上盛宴餐厅', extra_charge: false }] })
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

    expect(await findByText('基本参数')).toBeTruthy()
    expect(await findByText('邮轮设施')).toBeTruthy()
    expect(await findByText('海洋绿洲号')).toBeTruthy()
  })

  it('opens voyage detail after clicking an all-products voyage card', async () => {
    const { getByText, findByText } = render(App)

    await fireEvent.click(getByText('全部商品'))
    await fireEvent.click(await findByText('海洋光谱号 上海-福冈-上海 5天4晚'))

    expect(await findByText('航次详情')).toBeTruthy()
    expect(await findByText('RC101')).toBeTruthy()
  })
})