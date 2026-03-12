import { describe, expect, it } from 'vitest'
import { buildDisplayRouteStops, buildScheduleTableStops, buildRouteMapModel, buildRouteMapPoints, resolvePortCoordinate } from '../../src/utils/route-map'

describe('route map utils', () => {
  it('resolves known port coordinates from real geography', () => {
    const shanghai = resolvePortCoordinate('上海')
    const fukuoka = resolvePortCoordinate('福冈')

    expect(shanghai).toEqual(expect.objectContaining({ latitude: expect.any(Number), longitude: expect.any(Number) }))
    expect(fukuoka).toEqual(expect.objectContaining({ latitude: expect.any(Number), longitude: expect.any(Number) }))
    expect(Math.abs((shanghai?.latitude || 0) - (fukuoka?.latitude || 0))).toBeGreaterThan(0.1)
    expect(Math.abs((shanghai?.longitude || 0) - (fukuoka?.longitude || 0))).toBeGreaterThan(0.1)
  })

  it('projects route points using geographic spread instead of alternating placeholders', () => {
    const points = buildRouteMapPoints([
      { city: '上海', dayNo: 1 },
      { city: '福冈', dayNo: 2 },
      { city: '济州', dayNo: 3 },
    ])

    expect(points).toHaveLength(3)
    expect(points[0]?.x).not.toBe(points[1]?.x)
    expect(points[0]?.y).not.toBe(points[1]?.y)
    expect(points[1]?.y).not.toBe(66)
    expect(points[2]?.y).not.toBe(42)
  })

  it('excludes sea-cruise placeholders from displayed route stops', () => {
    const stops = buildDisplayRouteStops([
      { city: '上海', dayNo: 1, summary: '登船' },
      { city: '海上巡游', dayNo: 2, summary: '全天巡游' },
      { city: '福冈', dayNo: 3, summary: '岸上观光' },
      { city: '海上巡游', dayNo: 4, summary: '返航' },
      { city: '上海', dayNo: 5, summary: '离船' },
    ])

    expect(stops.map((item) => item.city)).toEqual(['上海', '福冈', '上海'])
  })

  it('includes sea-cruise days in schedule table stops with empty times', () => {
    const stops = buildScheduleTableStops([
      { city: '上海', dayNo: 1, summary: '登船', etdTime: '17:00' },
      { city: '海上巡游', dayNo: 2, summary: '全天巡游', etaTime: '08:00', etdTime: '18:00' },
      { city: '福冈', dayNo: 3, summary: '岸上观光', etaTime: '09:00', etdTime: '19:00' },
      { city: '海上巡游', dayNo: 4, summary: '返航' },
      { city: '上海', dayNo: 5, summary: '离船', etaTime: '08:00' },
    ])

    expect(stops.map((item) => item.city)).toEqual(['上海', '海上巡游', '福冈', '海上巡游', '上海'])
    // Sea cruise days should have undefined (cleared) times
    expect(stops[1]?.etaTime).toBeUndefined()
    expect(stops[1]?.etdTime).toBeUndefined()
    // Port days should keep their times
    expect(stops[2]?.etaTime).toBe('09:00')
    expect(stops[2]?.etdTime).toBe('19:00')
  })

  it('strips country suffixes from schedule table port names', () => {
    const stops = buildScheduleTableStops([
      { city: '天津（中国）', dayNo: 1, summary: '登船', etdTime: '17:00' },
      { city: '济州（韩国）', dayNo: 2, summary: '岸上观光', etaTime: '09:00', etdTime: '19:00' },
      { city: '海上巡游', dayNo: 3, summary: '返航' },
    ])

    expect(stops.map((item) => item.city)).toEqual(['天津', '济州', '海上巡游'])
  })

  it('builds ordered maritime route segments with intermediate sea points', () => {
    const model = buildRouteMapModel([
      { city: '上海', dayNo: 1 },
      { city: '海上巡游', dayNo: 2 },
      { city: '福冈', dayNo: 3 },
      { city: '釜山', dayNo: 4 },
    ])

    expect(model.stops.map((item) => item.city)).toEqual(['上海', '福冈', '釜山'])
    expect(model.segments).toHaveLength(2)
    expect(model.segments[0]?.from.city).toBe('上海')
    expect(model.segments[0]?.to.city).toBe('福冈')
    expect(model.segments[0]?.projectedPoints.length).toBeGreaterThan(2)
    expect(model.segments[1]?.from.city).toBe('福冈')
    expect(model.segments[1]?.to.city).toBe('釜山')
  })

  it('prefers backend-provided route geometry when available', () => {
    const model = buildRouteMapModel(
      [
        { city: '上海', dayNo: 1 },
        { city: '福冈', dayNo: 3 },
      ],
      {
        geometryType: 'MultiLineString',
        coordinates: [
          [
            [121.4737, 31.2304],
            [124.8, 32.0],
            [130.4017, 33.5902],
          ],
        ],
      },
    )

    expect(model.segments).toHaveLength(1)
    expect(model.landPolygons).toBeTypeOf('object')
    expect(model.segments[0]?.projectedPoints.length).toBeGreaterThanOrEqual(3)
  })
})