export type PortCoordinate = {
  latitude: number
  longitude: number
}

export type RouteStopInput = {
  city: string
  dayNo: number
  summary?: string
  etaTime?: string
  etdTime?: string
}

export type RouteDisplayStop = {
  city: string
  dayNo: number
  summary?: string
  etaTime?: string
  etdTime?: string
}

export type RouteMapPoint = {
  city: string
  dayNo: number
  x: number
  y: number
  latitude: number
  longitude: number
  hasCoordinate: boolean
}

export type RouteMapSegment = {
  id: string
  from: RouteMapPoint
  to: RouteMapPoint
  projectedPoints: RouteMapPoint[]
  pathD: string
  arrow?: { x: number; y: number; angle: number }
}

export type RouteMapLandPolygon = {
  id: string
  points: string
}

export type RouteMapModel = {
  stops: RouteMapPoint[]
  segments: RouteMapSegment[]
  landPolygons: RouteMapLandPolygon[]
  hasKnownCoordinates: boolean
}

export type BackendRouteMapInput = {
  geometryType?: string
  coordinates?: number[][][]
}

type RawPoint = {
  id: string
  city: string
  dayNo: number
  latitude: number
  longitude: number
  hasCoordinate: boolean
}

type GeoPoint = {
  latitude: number
  longitude: number
}

import worldLandRaw from './world-land.json?raw'

let _WORLD_LAND_POLYGONS: number[][][] | null = null;
function getLandPolygons(): number[][][] {
  if (!_WORLD_LAND_POLYGONS) {
    _WORLD_LAND_POLYGONS = JSON.parse(worldLandRaw);
  }
  return _WORLD_LAND_POLYGONS;
}

function polygonBBoxIntersectsViewport(polygonCoords: number[][], proj: MapProjection): boolean {
  let polyMinLon = Infinity, polyMaxLon = -Infinity
  let polyMinLat = Infinity, polyMaxLat = -Infinity
  for (const [lon, lat] of polygonCoords) {
    if (lon < polyMinLon) polyMinLon = lon
    if (lon > polyMaxLon) polyMaxLon = lon
    if (lat < polyMinLat) polyMinLat = lat
    if (lat > polyMaxLat) polyMaxLat = lat
  }
  return polyMaxLon >= proj.minLon && polyMinLon <= proj.maxLon
      && polyMaxLat >= proj.minLat && polyMinLat <= proj.maxLat
}

function collectVisibleLandPolygons(proj: MapProjection): RouteMapLandPolygon[] {
  const result: RouteMapLandPolygon[] = []
  getLandPolygons().forEach((polygonCoords, idx) => {
    if (polygonBBoxIntersectsViewport(polygonCoords, proj)) {
      const pathPoints = polygonCoords.map((coord: number[]) => projectPoint(coord[0], coord[1], proj))
      const pStr = pathPoints.map((p) => `${p.x.toFixed(1)},${p.y.toFixed(1)}`).join(' ')
      result.push({ id: `land-${idx}`, points: pStr })
    }
  })
  return result
}

/**
 * Ray-casting 射线法点-in-polygon 测试（地理坐标：[lon, lat]）
 * 用于判断一个经纬度点是否位于陆地多边形内部
 */
function pointInGeoPolygon(lon: number, lat: number, polygon: number[][]): boolean {
  let inside = false
  for (let i = 0, j = polygon.length - 1; i < polygon.length; j = i++) {
    const xi = polygon[i]![0], yi = polygon[i]![1]
    const xj = polygon[j]![0], yj = polygon[j]![1]
    if ((yi > lat) !== (yj > lat) && lon < ((xj - xi) * (lat - yi)) / (yj - yi) + xi) {
      inside = !inside
    }
  }
  return inside
}

type LandPolygonWithBBox = { coords: number[][]; minLon: number; maxLon: number; minLat: number; maxLat: number }

/** 预过滤并缓存可见陆地多边形及其边界框 */
function getVisibleLandWithBBox(proj: MapProjection): LandPolygonWithBBox[] {
  const result: LandPolygonWithBBox[] = []
  for (const coords of getLandPolygons()) {
    let pMinLon = Infinity, pMaxLon = -Infinity, pMinLat = Infinity, pMaxLat = -Infinity
    for (const [lon, lat] of coords) {
      if (lon < pMinLon) pMinLon = lon
      if (lon > pMaxLon) pMaxLon = lon
      if (lat < pMinLat) pMinLat = lat
      if (lat > pMaxLat) pMaxLat = lat
    }
    if (pMaxLon >= proj.minLon && pMinLon <= proj.maxLon && pMaxLat >= proj.minLat && pMinLat <= proj.maxLat) {
      result.push({ coords, minLon: pMinLon, maxLon: pMaxLon, minLat: pMinLat, maxLat: pMaxLat })
    }
  }
  return result
}

/** 使用预过滤多边形和边界框检查进行快速的点-in-陆地测试 */
function isInsideLandFast(lon: number, lat: number, polys: LandPolygonWithBBox[]): boolean {
  for (const p of polys) {
    if (lon < p.minLon || lon > p.maxLon || lat < p.minLat || lat > p.maxLat) continue
    if (pointInGeoPolygon(lon, lat, p.coords)) return true
  }
  return false
}

/** 检查地理点是否位于任何可见陆地多边形内部 */
function isInsideLand(lon: number, lat: number, proj: MapProjection): boolean {
  for (const polygon of getLandPolygons()) {
    if (!polygonBBoxIntersectsViewport(polygon, proj)) continue
    if (pointInGeoPolygon(lon, lat, polygon)) return true
  }
  return false
}

/** 检查屏幕空间点是否位于陆地内部（使用预过滤多边形的快速路径） */
function isInsideLandScreenFast(x: number, y: number, proj: MapProjection, polys: LandPolygonWithBBox[]): boolean {
  const { lon, lat } = unprojectPoint(x, y, proj)
  return isInsideLandFast(lon, lat, polys)
}

const VIEWPORT = { width: 320, height: 208, paddingX: 30, paddingY: 30 }

const PORT_COORDINATES: Record<string, PortCoordinate> = {
  上海: { latitude: 31.2304, longitude: 121.4737 },
  天津: { latitude: 39.0842, longitude: 117.2009 },
  大连: { latitude: 38.914, longitude: 121.6147 },
  青岛: { latitude: 36.0671, longitude: 120.3826 },
  厦门: { latitude: 24.4798, longitude: 118.0894 },
  香港: { latitude: 22.3193, longitude: 114.1694 },
  福冈: { latitude: 33.5902, longitude: 130.4017 },
  长崎: { latitude: 32.7503, longitude: 129.8777 },
  鹿儿岛: { latitude: 31.5966, longitude: 130.5571 },
  济州: { latitude: 33.4996, longitude: 126.5312 },
  西归浦: { latitude: 33.2541, longitude: 126.5601 },
  西归浦市: { latitude: 33.2541, longitude: 126.5601 },
  釜山: { latitude: 35.1796, longitude: 129.0756 },
  仁川: { latitude: 37.4563, longitude: 126.7052 },
  横滨: { latitude: 35.4437, longitude: 139.638 },
  神户: { latitude: 34.6901, longitude: 135.1955 },
  大阪: { latitude: 34.6937, longitude: 135.5023 },
}

const SEA_ROUTE_PLACEHOLDERS = ['海上巡游', '海上巡航', '海上观光', '巡游日', '海上']

const MARITIME_NODES: Record<string, GeoPoint> = {
  shanghai_gate: { latitude: 30.95, longitude: 122.15 },
  east_china_1: { latitude: 31.2, longitude: 123.9 },
  east_china_2: { latitude: 31.85, longitude: 125.4 },
  yellow_sea_1: { latitude: 34.6, longitude: 124.4 },
  yellow_sea_2: { latitude: 36.6, longitude: 123.5 },
  bohai_gate: { latitude: 38.2, longitude: 120.8 },
  jeju_west: { latitude: 33.6, longitude: 125.6 },
  jeju_east: { latitude: 33.6, longitude: 127.35 },
  seogwipo_south: { latitude: 32.95, longitude: 126.55 },
  korea_strait_mid: { latitude: 34.4, longitude: 128.35 },
  busan_gate: { latitude: 34.95, longitude: 129.25 },
  fukuoka_gate: { latitude: 33.78, longitude: 130.05 },
  nagasaki_gate: { latitude: 32.78, longitude: 129.62 },
  kyushu_south: { latitude: 31.7, longitude: 130.1 },
  japan_sea_1: { latitude: 35.2, longitude: 131.4 },
  japan_sea_2: { latitude: 36.2, longitude: 133.6 },
  japan_sea_3: { latitude: 36.5, longitude: 136.0 },
}

const MARITIME_EDGES: Array<[string, string]> = [
  ['shanghai_gate', 'east_china_1'],
  ['east_china_1', 'east_china_2'],
  ['east_china_2', 'jeju_west'],
  ['jeju_west', 'jeju_east'],
  ['jeju_west', 'yellow_sea_1'],
  ['yellow_sea_1', 'yellow_sea_2'],
  ['yellow_sea_2', 'bohai_gate'],
  ['jeju_east', 'korea_strait_mid'],
  ['korea_strait_mid', 'busan_gate'],
  ['seogwipo_south', 'jeju_west'],
  ['seogwipo_south', 'jeju_east'],
  ['korea_strait_mid', 'fukuoka_gate'],
  ['fukuoka_gate', 'nagasaki_gate'],
  ['nagasaki_gate', 'kyushu_south'],
  ['fukuoka_gate', 'japan_sea_1'],
  ['japan_sea_1', 'japan_sea_2'],
  ['japan_sea_2', 'japan_sea_3'],
]

const PORT_CONNECTORS: Record<string, string> = {
  上海: 'shanghai_gate',
  天津: 'bohai_gate',
  大连: 'yellow_sea_2',
  青岛: 'yellow_sea_2',
  济州: 'jeju_west',
  西归浦: 'seogwipo_south',
  西归浦市: 'seogwipo_south',
  福冈: 'fukuoka_gate',
  长崎: 'nagasaki_gate',
  釜山: 'busan_gate',
  仁川: 'yellow_sea_1',
  鹿儿岛: 'kyushu_south',
}

const LAND_POLYGONS: Array<{ id: string; points: GeoPoint[] }> = [
  {
    id: 'china-coast',
    points: [
      { latitude: 39.9, longitude: 116.2 },
      { latitude: 39.6, longitude: 118.1 },
      { latitude: 38.7, longitude: 120.4 },
      { latitude: 37.5, longitude: 121.7 },
      { latitude: 35.4, longitude: 121.5 },
      { latitude: 33.3, longitude: 121.4 },
      { latitude: 31.1, longitude: 121.7 },
      { latitude: 29.3, longitude: 121.0 },
      { latitude: 26.8, longitude: 119.1 },
      { latitude: 24.0, longitude: 117.6 },
      { latitude: 22.2, longitude: 114.0 },
      { latitude: 21.6, longitude: 111.0 },
      { latitude: 20.4, longitude: 108.5 },
      { latitude: 40.6, longitude: 108.5 },
      { latitude: 39.9, longitude: 116.2 },
    ],
  },
  {
    id: 'korean-peninsula',
    points: [
      { latitude: 38.7, longitude: 124.5 },
      { latitude: 39.3, longitude: 126.2 },
      { latitude: 38.5, longitude: 128.2 },
      { latitude: 36.8, longitude: 129.0 },
      { latitude: 35.1, longitude: 129.4 },
      { latitude: 34.2, longitude: 128.8 },
      { latitude: 34.3, longitude: 126.2 },
      { latitude: 35.6, longitude: 125.1 },
      { latitude: 38.7, longitude: 124.5 },
    ],
  },
  {
    id: 'kyushu-west',
    points: [
      { latitude: 33.2, longitude: 129.4 },
      { latitude: 33.9, longitude: 130.1 },
      { latitude: 33.7, longitude: 131.1 },
      { latitude: 32.6, longitude: 131.4 },
      { latitude: 31.4, longitude: 130.9 },
      { latitude: 31.1, longitude: 129.9 },
      { latitude: 31.9, longitude: 129.3 },
      { latitude: 33.2, longitude: 129.4 },
    ],
  },
]

function normalizeCityName(city: string) {
  return city.trim().replace(/[（(].*?[）)]/g, '').replace(/港口|港|市$/g, '')
}

function isSeaCruiseStop(city: string, summary?: string) {
  const normalized = normalizeCityName(city)
  if (SEA_ROUTE_PLACEHOLDERS.some((item) => normalized.includes(item))) return true
  if (normalized === '') {
    const summaryText = (summary || '').trim()
    return SEA_ROUTE_PLACEHOLDERS.some((item) => summaryText.includes(item))
  }
  return false
}

function fallbackCoordinate(city: string): PortCoordinate {
  let hash = 0
  for (const char of city) {
    hash = (hash * 31 + char.charCodeAt(0)) >>> 0
  }
  return {
    latitude: 20 + (hash % 1800) / 100,
    longitude: 105 + ((hash >> 8) % 3500) / 100,
  }
}

function geoDistance(left: GeoPoint, right: GeoPoint) {
  const latScale = 111
  const lonScale = Math.cos(((left.latitude + right.latitude) / 2) * Math.PI / 180) * 111
  const dx = (left.longitude - right.longitude) * lonScale
  const dy = (left.latitude - right.latitude) * latScale
  return Math.sqrt(dx * dx + dy * dy)
}

export type MapProjection = {
  minLon: number
  maxLon: number
  minLat: number
  maxLat: number
  lonScale: number
  latScale: number
  width: number
  height: number
  offsetX: number
  offsetY: number
}

function computeProjection(points: GeoPoint[]): MapProjection {
  if (points.length === 0) {
    return { minLon: 0, maxLon: 1, minLat: 0, maxLat: 1, lonScale: 1, latScale: 1, width: 1, height: 1, offsetX: 0, offsetY: 0 }
  }
  const longitudes = points.map((point) => point.longitude)
  const latitudes = points.map((point) => point.latitude)
  let minLon = Math.min(...longitudes)
  let maxLon = Math.max(...longitudes)
  let minLat = Math.min(...latitudes)
  let maxLat = Math.max(...latitudes)

  // 添加内边距到边界 — 充足的内边距用于显示地理背景
  const lonPad = Math.max((maxLon - minLon) * 0.25, 3)
  const latPad = Math.max((maxLat - minLat) * 0.25, 3)
  minLon -= lonPad
  maxLon += lonPad
  minLat -= latPad
  maxLat += latPad

  const lonSpan = maxLon - minLon
  const latSpan = maxLat - minLat

  const midLat = (minLat + maxLat) / 2
  const cosLat = Math.cos((midLat * Math.PI) / 180)

  // 可用尺寸
  const availW = VIEWPORT.width - VIEWPORT.paddingX * 2
  const availH = VIEWPORT.height - VIEWPORT.paddingY * 2

  // 边界框在地图上的实际长宽比
  const boxW = lonSpan * cosLat
  const boxH = latSpan

  // 将 boxW x boxH 适配到 availW x availH
  let scale = 1
  if (boxW / boxH > availW / availH) {
    // 受宽度限制
    scale = availW / boxW
  } else {
    // 受高度限制
    scale = availH / boxH
  }

  const lonScale = scale * cosLat
  const latScale = scale

  const width = lonSpan * lonScale
  const height = latSpan * latScale
  
  const offsetX = VIEWPORT.paddingX + (availW - width) / 2
  const offsetY = VIEWPORT.paddingY + (availH - height) / 2

  return { minLon, maxLon, minLat, maxLat, lonScale, latScale, width, height, offsetX, offsetY }
}

function projectPoint(lon: number, lat: number, proj: MapProjection) {
  return {
    x: proj.offsetX + (lon - proj.minLon) * proj.lonScale,
    y: proj.offsetY + proj.height - (lat - proj.minLat) * proj.latScale, // Y 轴翻转
  }
}

/** 逆投影：屏幕坐标 (x,y) → 地理坐标 (lon,lat) */
function unprojectPoint(x: number, y: number, proj: MapProjection) {
  return {
    lon: proj.minLon + (x - proj.offsetX) / proj.lonScale,
    lat: proj.minLat + (proj.offsetY + proj.height - y) / proj.latScale,
  }
}

function projectGeoPoints(points: GeoPoint[], proj?: MapProjection) {
  const p = proj || computeProjection(points)
  return points.map((point) => projectPoint(point.longitude, point.latitude, p))
}

/**
 * 沿全局航段方向（起点→终点）垂直偏移航点，
 * 并在端点处使用平滑衰减，使航线在港口处汇聚。
 * 这样可以避免因局部切线偏移导致的锯齿状走线。
 */
function applyGlobalOffset<T extends { x: number; y: number }>(points: T[], offset: number): T[] {
  if (points.length < 2) return points
  const first = points[0]!, last = points[points.length - 1]!
  const gdx = last.x - first.x, gdy = last.y - first.y
  const glen = Math.sqrt(gdx * gdx + gdy * gdy) || 1
  // 垂直于全局方向
  const nx = -gdy / glen, ny = gdx / glen

  return points.map((p, i) => {
    const total = points.length - 1
    // 基于正弦的平滑衰减：端点为0，中间为1
    const ratio = total > 0 ? i / total : 0
    const fade = Math.sin(ratio * Math.PI)
    return { ...p, x: p.x + nx * offset * fade, y: p.y + ny * offset * fade }
  })
}

/**
 * 使用 Catmull-Rom 到 Bezier 转换从航点构建 SVG 三次贝塞尔路径。
 * 生成原生 SVG C 曲线，实现无限平滑效果。
 */
function buildSmoothPathD(points: Array<{ x: number; y: number }>): string {
  if (points.length === 0) return ''
  if (points.length === 1) return `M ${points[0]!.x.toFixed(1)} ${points[0]!.y.toFixed(1)}`
  if (points.length === 2) {
    return `M ${points[0]!.x.toFixed(1)} ${points[0]!.y.toFixed(1)} L ${points[1]!.x.toFixed(1)} ${points[1]!.y.toFixed(1)}`
  }

  // Catmull-Rom → 三次 Bezier 控制点
  // 对于线段 P_i → P_{i+1}，控制点为：
  //   cp1 = P_i + (P_{i+1} - P_{i-1}) / 6
  //   cp2 = P_{i+1} - (P_{i+2} - P_i) / 6
  const f = (v: number) => v.toFixed(1)
  let d = `M ${f(points[0]!.x)} ${f(points[0]!.y)}`

  for (let i = 0; i < points.length - 1; i++) {
    const p0 = points[Math.max(0, i - 1)]!
    const p1 = points[i]!
    const p2 = points[i + 1]!
    const p3 = points[Math.min(points.length - 1, i + 2)]!

    const cp1x = p1.x + (p2.x - p0.x) / 6
    const cp1y = p1.y + (p2.y - p0.y) / 6
    const cp2x = p2.x - (p3.x - p1.x) / 6
    const cp2y = p2.y - (p3.y - p1.y) / 6

    d += ` C ${f(cp1x)} ${f(cp1y)}, ${f(cp2x)} ${f(cp2y)}, ${f(p2.x)} ${f(p2.y)}`
  }
  return d
}

/**
 * 在 SVG 三次贝塞尔路径上均匀采样 N 个点。
 * 用于航向箭头放置和陆地冲突计数。
 */
function sampleBezierPoints(points: Array<{ x: number; y: number }>, samplesPerSegment = 8): Array<{ x: number; y: number }> {
  if (points.length < 2) return [...points]
  const result: Array<{ x: number; y: number }> = []

  for (let i = 0; i < points.length - 1; i++) {
    const p0 = points[Math.max(0, i - 1)]!
    const p1 = points[i]!
    const p2 = points[i + 1]!
    const p3 = points[Math.min(points.length - 1, i + 2)]!

    const cp1x = p1.x + (p2.x - p0.x) / 6
    const cp1y = p1.y + (p2.y - p0.y) / 6
    const cp2x = p2.x - (p3.x - p1.x) / 6
    const cp2y = p2.y - (p3.y - p1.y) / 6

    for (let s = 0; s < samplesPerSegment; s++) {
      const t = s / samplesPerSegment
      const u = 1 - t
      const x = u * u * u * p1.x + 3 * u * u * t * cp1x + 3 * u * t * t * cp2x + t * t * t * p2.x
      const y = u * u * u * p1.y + 3 * u * u * t * cp1y + 3 * u * t * t * cp2y + t * t * t * p2.y
      result.push({ x, y })
    }
  }
  result.push(points[points.length - 1]!)
  return result
}

/**
 * 为航段生成偏移航点，选择陆地冲突较少的一侧（+/-）。
 * 返回偏移后的航点（不是密集采样后的点 —
 * SVG 路径请使用 buildSmoothPathD）。
 */
function generateOffsetWaypoints<T extends { x: number; y: number }>(
  points: T[], proj: MapProjection, preferredSign: 1 | -1 = 1, absOffset = 8,
): T[] {
  if (points.length < 2) return points
  const visiblePolys = getVisibleLandWithBBox(proj)
  const countLand = (pts: T[]) => {
    if (visiblePolys.length === 0) return 0
    const samples = sampleBezierPoints(pts)
    let count = 0
    for (const p of samples) {
      if (isInsideLandScreenFast(p.x, p.y, proj, visiblePolys)) count++
    }
    return count
  }

  // 首先尝试优选侧，然后尝试另一侧，偏移量递增
  const offsets = [absOffset, absOffset * 1.5, absOffset * 2]
  let bestPoints = points
  let bestLand = countLand(points)

  for (const off of offsets) {
    for (const sign of [preferredSign, -preferredSign as 1 | -1]) {
      const candidate = applyGlobalOffset(points, off * sign)
      const land = countLand(candidate)
      if (land < bestLand) {
        bestPoints = candidate
        bestLand = land
        if (land === 0) return bestPoints
      }
    }
  }
  return bestPoints
}

function buildPathD(points: Array<{ x: number; y: number }>) {
  return buildSmoothPathD(points)
}

function projectWithLookup(points: GeoPoint[], proj: MapProjection) {
  const projected = projectGeoPoints(points, proj)
  const byKey = new Map<string, { x: number; y: number }>()
  points.forEach((point, index) => {
    byKey.set(`${point.longitude.toFixed(6)}:${point.latitude.toFixed(6)}`, projected[index]!)
  })
  return byKey
}

function computeSegmentArrow(points: Array<{ x: number; y: number }>) {
  if (points.length < 2) return undefined
  let totalLen = 0
  const lengths: number[] = [0]
  for (let i = 1; i < points.length; i++) {
    const p1 = points[i - 1]!
    const p2 = points[i]!
    const d = Math.sqrt((p2.x - p1.x) ** 2 + (p2.y - p1.y) ** 2)
    totalLen += d
    lengths.push(totalLen)
  }
  const targetLen = totalLen / 2
  for (let i = 1; i < points.length; i++) {
    if (lengths[i]! >= targetLen || i === points.length - 1) {
      const p1 = points[i - 1]!
      const p2 = points[i]!
      const angle = Math.atan2(p2.y - p1.y, p2.x - p1.x) * 180 / Math.PI
      const segmentLen = lengths[i]! - lengths[i - 1]!
      const ratio = segmentLen > 0 ? (targetLen - lengths[i - 1]!) / segmentLen : 0.5
      return {
        x: p1.x + (p2.x - p1.x) * ratio,
        y: p1.y + (p2.y - p1.y) * ratio,
        angle
      }
    }
  }
  return undefined
}

function normalizeBackendRoute(routeMap?: BackendRouteMapInput) {
  if (!routeMap || !Array.isArray(routeMap.coordinates) || routeMap.coordinates.length === 0) return []
  return routeMap.coordinates
    .map((segment) => segment
      .filter((point) => Array.isArray(point) && point.length >= 2)
      .map((point) => ({ longitude: Number(point[0]), latitude: Number(point[1]) })))
    .filter((segment) => segment.length >= 2)
}

function buildGraph() {
  const graph = new Map<string, Array<{ to: string; weight: number }>>()
  for (const [from, to] of MARITIME_EDGES) {
    const fromPoint = MARITIME_NODES[from]!
    const toPoint = MARITIME_NODES[to]!
    const weight = geoDistance(fromPoint, toPoint)
    if (!graph.has(from)) graph.set(from, [])
    if (!graph.has(to)) graph.set(to, [])
    graph.get(from)!.push({ to, weight })
    graph.get(to)!.push({ to: from, weight })
  }
  return graph
}

const MARITIME_GRAPH = buildGraph()

function shortestNodePath(start: string, end: string) {
  if (start === end) return [start]
  const distances = new Map<string, number>()
  const previous = new Map<string, string | null>()
  const visited = new Set<string>()
  const queue = new Set<string>(Object.keys(MARITIME_NODES))

  for (const key of queue) {
    distances.set(key, key === start ? 0 : Number.POSITIVE_INFINITY)
    previous.set(key, null)
  }

  while (queue.size > 0) {
    let current: string | null = null
    let currentDistance = Number.POSITIVE_INFINITY
    for (const key of queue) {
      const distance = distances.get(key) ?? Number.POSITIVE_INFINITY
      if (distance < currentDistance) {
        current = key
        currentDistance = distance
      }
    }

    if (!current || current === end) break
    queue.delete(current)
    visited.add(current)

    for (const edge of MARITIME_GRAPH.get(current) || []) {
      if (visited.has(edge.to)) continue
      const candidate = currentDistance + edge.weight
      if (candidate < (distances.get(edge.to) ?? Number.POSITIVE_INFINITY)) {
        distances.set(edge.to, candidate)
        previous.set(edge.to, current)
      }
    }
  }

  const path: string[] = []
  let current: string | null = end
  while (current) {
    path.unshift(current)
    current = previous.get(current) ?? null
  }
  return path[0] === start ? path : [start, end]
}

function resolvePortNode(city: string, coordinate: GeoPoint) {
  const normalized = normalizeCityName(city)
  const explicit = PORT_CONNECTORS[normalized]
  if (explicit) return explicit

  let nearest = Object.keys(MARITIME_NODES)[0]!
  let nearestDistance = Number.POSITIVE_INFINITY
  for (const [nodeId, node] of Object.entries(MARITIME_NODES)) {
    const distance = geoDistance(coordinate, node)
    if (distance < nearestDistance) {
      nearest = nodeId
      nearestDistance = distance
    }
  }
  return nearest
}

function resolveRawStops(stops: RouteDisplayStop[]) {
  return stops.map((stop) => {
    const coordinate = resolvePortCoordinate(stop.city)
    const fallback = coordinate || fallbackCoordinate(stop.city)
    return {
      id: `${stop.dayNo}-${stop.city}`,
      city: stop.city,
      dayNo: stop.dayNo,
      latitude: fallback.latitude,
      longitude: fallback.longitude,
      hasCoordinate: Boolean(coordinate),
    } satisfies RawPoint
  })
}

function buildSegmentCoordinates(from: RawPoint, to: RawPoint) {
  const startNode = resolvePortNode(from.city, from)
  const endNode = resolvePortNode(to.city, to)
  const nodePath = shortestNodePath(startNode, endNode)
  const seaPoints = nodePath.map((nodeId) => ({ ...MARITIME_NODES[nodeId] }))
  const coordinates: GeoPoint[] = [{ latitude: from.latitude, longitude: from.longitude }, ...seaPoints, { latitude: to.latitude, longitude: to.longitude }]

  return coordinates.filter((point, index, list) => {
    if (index === 0) return true
    const previous = list[index - 1]!
    return previous.latitude !== point.latitude || previous.longitude !== point.longitude
  })
}

export function resolvePortCoordinate(city: string): PortCoordinate | null {
  const normalized = normalizeCityName(city)
  if (!normalized) return null
  return PORT_COORDINATES[normalized] || null
}

export function buildDisplayRouteStops(stops: RouteStopInput[]): RouteDisplayStop[] {
  const seen = new Set<string>()
  return stops
    .filter((stop) => !isSeaCruiseStop(stop.city, stop.summary))
    .filter((stop) => {
      const key = `${stop.dayNo}-${normalizeCityName(stop.city)}`
      if (seen.has(key)) return false
      seen.add(key)
      return true
    })
    .map((stop) => ({
      city: stop.city,
      dayNo: stop.dayNo,
      summary: stop.summary,
      etaTime: stop.etaTime,
      etdTime: stop.etdTime,
    }))
}

export function buildScheduleTableStops(stops: RouteStopInput[]): RouteDisplayStop[] {
  const seen = new Set<string>()
  return stops
    .filter((stop) => {
      const key = `${stop.dayNo}-${normalizeCityName(stop.city) || stop.city}`
      if (seen.has(key)) return false
      seen.add(key)
      return true
    })
    .map((stop) => {
      const isSea = isSeaCruiseStop(stop.city, stop.summary)
      const displayCity = normalizeCityName(stop.city) || stop.city
      return {
        city: displayCity,
        dayNo: stop.dayNo,
        summary: stop.summary,
        etaTime: isSea ? undefined : stop.etaTime,
        etdTime: isSea ? undefined : stop.etdTime,
      }
    })
}

export function buildRouteMapPoints(stops: RouteStopInput[]): RouteMapPoint[] {
  const rawStops = resolveRawStops(buildDisplayRouteStops(stops))
  if (rawStops.length === 0) return []
  const proj = computeProjection(rawStops)
  const projected = projectGeoPoints(rawStops, proj)
  return rawStops.map((stop, index) => ({
    city: stop.city,
    dayNo: stop.dayNo,
    latitude: stop.latitude,
    longitude: stop.longitude,
    hasCoordinate: stop.hasCoordinate,
    x: projected[index]!.x,
    y: projected[index]!.y,
  }))
}

export function buildRouteMapModel(stops: RouteStopInput[], backendRouteMap?: BackendRouteMapInput): RouteMapModel {
  const rawStops = resolveRawStops(buildDisplayRouteStops(stops))
  if (rawStops.length === 0) return { stops: [], segments: [], landPolygons: [], hasKnownCoordinates: false }

  const backendSegments = normalizeBackendRoute(backendRouteMap)
  if (backendSegments.length > 0) {
    const projectionPoints = [...rawStops, ...backendSegments.flat()]
    const proj = computeProjection(projectionPoints)
    const projectedMap = projectWithLookup(projectionPoints, proj)
    const projectedStops = rawStops.map((stop) => {
      const projected = projectedMap.get(`${stop.longitude.toFixed(6)}:${stop.latitude.toFixed(6)}`)!
      return {
        city: stop.city,
        dayNo: stop.dayNo,
        latitude: stop.latitude,
        longitude: stop.longitude,
        hasCoordinate: stop.hasCoordinate,
        x: projected.x,
        y: projected.y,
      } satisfies RouteMapPoint
    })

    const segments = backendSegments.map((segment, index) => {
      const projectedPoints = segment.map((point) => {
        const projected = projectedMap.get(`${point.longitude.toFixed(6)}:${point.latitude.toFixed(6)}`)!
        return {
          city: projectedStops[Math.min(index + 1, projectedStops.length - 1)]?.city || '',
          dayNo: projectedStops[Math.min(index + 1, projectedStops.length - 1)]?.dayNo || 0,
          latitude: point.latitude,
          longitude: point.longitude,
          hasCoordinate: true as boolean,
          x: projected.x,
          y: projected.y,
        } satisfies RouteMapPoint
      })
      projectedPoints[0] = { ...projectedPoints[0]!, city: projectedStops[Math.min(index, projectedStops.length - 1)]?.city || '', dayNo: projectedStops[Math.min(index, projectedStops.length - 1)]?.dayNo || 0 }
      projectedPoints[projectedPoints.length - 1] = { ...projectedPoints[projectedPoints.length - 1]!, city: projectedStops[Math.min(index + 1, projectedStops.length - 1)]?.city || '', dayNo: projectedStops[Math.min(index + 1, projectedStops.length - 1)]?.dayNo || 0 }

        const curvedPoints = generateOffsetWaypoints(projectedPoints, proj, index % 2 === 0 ? 1 : -1)

        return {
          id: `${projectedStops[Math.min(index, projectedStops.length - 1)]?.city || 'route'}-${index}`,
          from: projectedStops[Math.min(index, projectedStops.length - 1)]!,
          to: projectedStops[Math.min(index + 1, projectedStops.length - 1)]!,
          projectedPoints: curvedPoints,
          pathD: buildPathD(curvedPoints),
          arrow: computeSegmentArrow(sampleBezierPoints(curvedPoints)),
      } satisfies RouteMapSegment
    })

    const landPolygons = collectVisibleLandPolygons(proj)

    return {
      stops: projectedStops,
      segments,
      landPolygons,
      hasKnownCoordinates: true,
    }
  }

  const allCoordinates: GeoPoint[] = []
  const rawSegments: GeoPoint[][] = []
  for (let index = 0; index < rawStops.length - 1; index += 1) {
    const segmentCoordinates = buildSegmentCoordinates(rawStops[index]!, rawStops[index + 1]!)
    rawSegments.push(segmentCoordinates)
    allCoordinates.push(...segmentCoordinates)
  }
  if (allCoordinates.length === 0) {
    allCoordinates.push(...rawStops)
  }

  const proj = computeProjection(allCoordinates)
  const projectedCoords = projectGeoPoints(allCoordinates, proj)
  const coordinateKey = (point: GeoPoint) => `${point.longitude.toFixed(4)}:${point.latitude.toFixed(4)}`
  const projectedMap = new Map<string, { x: number; y: number }>()
  allCoordinates.forEach((point, index) => {
    projectedMap.set(coordinateKey(point), projectedCoords[index]!)
  })

  const projectedStops = rawStops.map((stop) => {
    const projected = projectedMap.get(coordinateKey(stop)) || projectGeoPoints([stop], proj)[0]!
    return {
      city: stop.city,
      dayNo: stop.dayNo,
      latitude: stop.latitude,
      longitude: stop.longitude,
      hasCoordinate: stop.hasCoordinate,
      x: projected.x,
      y: projected.y,
    } satisfies RouteMapPoint
  })

  const segments = rawSegments.map((segmentCoordinates, index) => {
    const projectedPoints = segmentCoordinates.map((point) => {
      const projected = projectedMap.get(coordinateKey(point))!
      return {
        city: projectedStops[Math.min(index + 1, projectedStops.length - 1)]!.city,
        dayNo: projectedStops[Math.min(index + 1, projectedStops.length - 1)]!.dayNo,
        latitude: point.latitude,
        longitude: point.longitude,
        hasCoordinate: true as boolean,
        x: projected.x,
        y: projected.y,
      } satisfies RouteMapPoint
    })

    projectedPoints[0] = { ...projectedPoints[0]!, city: projectedStops[index]!.city, dayNo: projectedStops[index]!.dayNo, hasCoordinate: projectedStops[index]!.hasCoordinate }
    projectedPoints[projectedPoints.length - 1] = { ...projectedPoints[projectedPoints.length - 1]!, city: projectedStops[index + 1]!.city, dayNo: projectedStops[index + 1]!.dayNo, hasCoordinate: projectedStops[index + 1]!.hasCoordinate }

    const curvedPoints = generateOffsetWaypoints(projectedPoints, proj, index % 2 === 0 ? 1 : -1)

    return {
      id: `${projectedStops[index]!.city}-${projectedStops[index + 1]!.city}-${index}`,
      from: projectedStops[index]!,
      to: projectedStops[index + 1]!,
      projectedPoints: curvedPoints,
      pathD: buildPathD(curvedPoints),
      arrow: computeSegmentArrow(sampleBezierPoints(curvedPoints)),
    } satisfies RouteMapSegment
  })

  const landPolygons = collectVisibleLandPolygons(proj)

  return {
    stops: projectedStops,
    segments,
    landPolygons,
    hasKnownCoordinates: projectedStops.some((stop) => stop.hasCoordinate),
  }
}