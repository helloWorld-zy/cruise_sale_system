import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { cleanup, fireEvent, render } from '@testing-library/vue'
import Page from '../../../../pages/wiki/index.vue'

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
    if (path === '/cruises?page=1&page_size=30') {
      return Promise.resolve({
        data: {
          list: [
            { id: 11, name: '海洋绿洲号', english_name: 'Oasis of the Seas', company_id: 1 },
            { id: 12, name: '欧罗巴号', english_name: 'MSC Euribia', company_id: 2 },
          ],
        },
      })
    }
    if (path === '/cruises?company_id=1&page=1&page_size=30') {
      return Promise.resolve({
        data: {
          list: [
            { id: 11, name: '海洋绿洲号', english_name: 'Oasis of the Seas', company_id: 1 },
          ],
        },
      })
    }
    return Promise.reject(new Error(`unexpected path: ${path}`))
  })
})

afterEach(() => {
  cleanup()
})

describe('Miniapp wiki home page', () => {
  it('renders all cruises by default', async () => {
    const { findByText, queryByText } = render(Page)

    expect(await findByText('全部邮轮')).toBeTruthy()
    expect(await findByText('海洋绿洲号')).toBeTruthy()
    expect(await findByText('欧罗巴号')).toBeTruthy()
    expect(queryByText('加载失败')).toBeNull()
  })

  it('filters cruises after clicking a company', async () => {
    const { findByText, queryByText } = render(Page)

    await fireEvent.click(await findByText('皇家加勒比'))

    expect(mockRequest).toHaveBeenCalledWith('/cruises?company_id=1&page=1&page_size=30')
    expect(await findByText('海洋绿洲号')).toBeTruthy()
    expect(queryByText('欧罗巴号')).toBeNull()
  })

  it('emits selected cruise id after clicking card', async () => {
    const { emitted, findByText } = render(Page)

    await fireEvent.click(await findByText('海洋绿洲号'))

    expect(emitted()['open-cruise']).toEqual([[11]])
  })
})