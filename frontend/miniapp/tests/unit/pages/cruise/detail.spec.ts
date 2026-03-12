import { describe, it, expect, vi } from 'vitest'
import { fireEvent, render } from '@testing-library/vue'
import Page from '../../../../pages/cruise/detail.vue'

const mockRequest = vi.fn()
vi.mock('../../../../src/utils/request', () => ({
  request: (...args: any[]) => mockRequest(...args),
}))

describe('Miniapp cruise detail page', () => {
  it('renders redesigned sections with dynamic basic specs and description', async () => {
    mockRequest.mockReset()
    mockRequest
      .mockResolvedValueOnce({
        data: {
          id: 9,
          name: 'Atlantic Dream',
          english_name: 'Atlantic Dream',
          tonnage: 140000,
          build_year: 2019,
          passenger_capacity: 3200,
          deck_count: 16,
          length: 315,
          width: 41,
          description: '这是一艘主打家庭娱乐体验的邮轮。',
        },
      })
      .mockResolvedValueOnce({ data: { list: [{ id: 1, name: '阳台房', area_min: 24, max_capacity: 3, min_price_cents: 188000 }] } })
      .mockResolvedValueOnce({
        data: [
          { id: 3, category_id: 1, category_name: '自助餐厅', name: '星海自助餐厅', extra_charge: false },
          { id: 4, category_id: 2, category_name: '剧院', name: '星幕剧场', extra_charge: false },
        ],
      })

    const { container, findAllByText, findByText, queryByText } = render(Page, { props: { cruiseId: 9 } })
    expect((await findAllByText('Atlantic Dream')).length).toBeGreaterThan(0)
    expect(await findByText('基本参数')).toBeTruthy()
    expect(await findByText('邮轮介绍')).toBeTruthy()
    expect(await findByText('邮轮设施')).toBeTruthy()
    expect(await findByText('总吨位')).toBeTruthy()
    expect(await findByText('建造年份')).toBeTruthy()
    expect(await findByText('这是一艘主打家庭娱乐体验的邮轮。')).toBeTruthy()
    expect(queryByText('船员数')).toBeNull()

    expect(mockRequest).not.toHaveBeenCalledWith('/routes')
  })

  it('toggles grouped facility detail card on category click', async () => {
    mockRequest.mockReset()
    mockRequest
      .mockResolvedValueOnce({
        data: {
          id: 9,
          name: 'Atlantic Dream',
          english_name: 'Atlantic Dream',
          tonnage: 140000,
          passenger_capacity: 3200,
          deck_count: 16,
          description: '这是一艘主打家庭娱乐体验的邮轮。',
        },
      })
      .mockResolvedValueOnce({ data: { list: [] } })
      .mockResolvedValueOnce({
        data: [
          { id: 3, category_id: 1, category_name: '自助餐厅', name: '星海自助餐厅', extra_charge: false, open_hours: '07:00-22:00' },
          { id: 4, category_id: 2, category_name: '酒吧', name: '暮色酒吧', extra_charge: true, charge_price_tip: '特调酒水单点收费' },
        ],
      })

    const { container, findAllByText, findByText, queryByText } = render(Page, { props: { cruiseId: 9 } })

    const complimentaryDining = (await findByText('免费餐厅')).closest('button') as HTMLElement
    expect(queryByText('星海自助餐厅')).toBeNull()

    await fireEvent.click(complimentaryDining)
    expect(await findByText('星海自助餐厅')).toBeTruthy()
    expect(await findByText('服务亮点')).toBeTruthy()

    await fireEvent.click(complimentaryDining)
    expect(queryByText('星海自助餐厅')).toBeNull()

    expect((await findAllByText('酒吧')).length).toBeGreaterThan(0)
  })
})
