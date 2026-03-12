import { describe, it, expect, vi } from 'vitest'
import { render } from '@testing-library/vue'
import Page from '../../../../pages/voyage/detail.vue'

const mockRequest = vi.fn()

vi.mock('../../../../src/utils/request', async () => {
  const actual = await vi.importActual<typeof import('../../../../src/utils/request')>('../../../../src/utils/request')
  return {
    ...actual,
    request: (...args: any[]) => mockRequest(...args),
  }
})

describe('Miniapp voyage detail page', () => {
  it('renders backend-driven voyage detail content', async () => {
    mockRequest.mockReset()
    mockRequest
      .mockResolvedValueOnce({
        data: {
          id: 101,
          cruise_id: 11,
          code: 'RC101',
          brief_info: '海洋光谱号 上海-福冈-上海 5天4晚',
          depart_date: '2026-05-23T00:00:00Z',
          return_date: '2026-05-27T00:00:00Z',
          min_price_cents: 409900,
          sold_count: 23,
          fee_note: {
            included: [{ text: '邮轮船票' }],
            excluded: [{ text: '签证费用', emphasis: true }],
          },
          booking_notice: {
            sections: [
              { key: 'booking_limit', title: '预订限制', items: [{ text: '请至少提前 1 天到港' }] },
              { key: 'documents', title: '出行证件', items: [{ text: '请携带护照原件', emphasis: true }] },
            ],
          },
          cruise: { id: 11, name: '海洋光谱号', company: { id: 1, name: '皇家加勒比' } },
          itineraries: [
            { id: 1, day_no: 1, stop_index: 1, city: '上海（中国）', summary: '登船', eta_time: '08:30', etd_time: '17:30', has_breakfast: true, has_lunch: true, has_dinner: false, has_accommodation: true, accommodation_text: '船上住宿' },
            { id: 2, day_no: 2, stop_index: 1, city: '海上巡游', summary: '全天巡游', eta_time: '', etd_time: '', has_breakfast: true, has_lunch: false, has_dinner: true, has_accommodation: true, accommodation_text: '船上住宿' },
            { id: 3, day_no: 3, stop_index: 1, city: '福冈（日本）', summary: '岸上观光', eta_time: '09:00', etd_time: '18:00', has_breakfast: true, has_lunch: false, has_dinner: true, has_accommodation: true, accommodation_text: '船上住宿' },
          ],
        },
      })
      .mockResolvedValueOnce({
        data: { list: [{ id: 301, name: '阳台房', min_price_cents: 188000 }] },
      })
      .mockResolvedValueOnce({
        data: [{ id: 401, category_id: 1, name: '海上剧院', extra_charge: false }],
      })

    const { findByText, findAllByText, container } = render(Page, { props: { voyageId: 101 } })

    expect(await findByText('航次详情')).toBeTruthy()
    expect(await findByText('RC101')).toBeTruthy()
    expect(await findByText('出发 2026-05-23')).toBeTruthy()
    expect(await findByText('皇家加勒比 · 海洋光谱号')).toBeTruthy()
    expect((await findAllByText('早餐')).length).toBeGreaterThan(0)
    expect(await findByText('午餐')).toBeTruthy()
    expect((await findAllByText('住宿')).length).toBeGreaterThan(0)
    expect((await findAllByText('船上住宿')).length).toBeGreaterThan(0)
    expect(await findByText('费用包含')).toBeTruthy()
    expect(await findByText('费用不包含')).toBeTruthy()
    expect(await findByText('预订须知')).toBeTruthy()
    expect(await findByText('出行证件')).toBeTruthy()
    expect(await findByText('阳台房')).toBeTruthy()
    expect(await findByText('海上剧院')).toBeTruthy()
    const scheduleTable = container.querySelector('[data-test="schedule-table"]') as HTMLElement | null
    expect(scheduleTable).toBeTruthy()
    expect(scheduleTable?.textContent).toContain('海上巡游')
    expect(scheduleTable?.textContent).toContain('上海')
    expect(scheduleTable?.textContent).toContain('福冈')
    expect(scheduleTable?.textContent).not.toContain('上海（中国）')
    expect(scheduleTable?.textContent).not.toContain('福冈（日本）')
    expect(scheduleTable?.textContent).toContain('--')
    expect(container.textContent).not.toContain('08:30 / 17:30')
    expect(await findByText('08:30')).toBeTruthy()
    expect(await findByText('17:30')).toBeTruthy()
  })
})