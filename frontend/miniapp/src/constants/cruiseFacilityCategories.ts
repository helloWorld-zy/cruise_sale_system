export type CruiseFacilityCategoryId =
  | 'included-dining'
  | 'specialty-dining'
  | 'bar-lounge'
  | 'leisure-entertainment'
  | 'kids-family'
  | 'suite-privilege'
  | 'sports-fitness'
  | 'other'

export type CruiseFacilityCategory = {
  id: CruiseFacilityCategoryId
  label: string
  icon: CruiseFacilityCategoryId
  exactNames: string[]
  keywords: string[]
}

export const cruiseFacilityCategories: CruiseFacilityCategory[] = [
  {
    id: 'included-dining',
    label: '免费餐厅',
    icon: 'included-dining',
    exactNames: ['免费餐厅', '主餐厅', '自助餐厅', '免费咖啡厅'],
    keywords: ['免费', '主餐', '自助', '餐厅', '咖啡厅'],
  },
  {
    id: 'specialty-dining',
    label: '特色餐厅',
    icon: 'specialty-dining',
    exactNames: ['特色餐厅', '收费餐厅', '牛排馆', '铁板烧', '日料'],
    keywords: ['特色', '收费', '牛排', '铁板', '火锅', '日料', '料理'],
  },
  {
    id: 'bar-lounge',
    label: '酒吧',
    icon: 'bar-lounge',
    exactNames: ['酒吧', '酒廊', '咖啡吧', '葡萄酒吧'],
    keywords: ['酒吧', '酒廊', '咖啡吧', '葡萄酒', '威士忌', '鸡尾酒'],
  },
  {
    id: 'leisure-entertainment',
    label: '休闲娱乐',
    icon: 'leisure-entertainment',
    exactNames: ['休闲娱乐', '剧院', '秀场', '赌场', 'KTV'],
    keywords: ['剧院', '秀', '娱乐', '赌场', 'ktv', '派对', '影院', '演艺'],
  },
  {
    id: 'kids-family',
    label: '亲子童趣',
    icon: 'kids-family',
    exactNames: ['亲子童趣', '儿童中心', '青少年中心', '托管中心'],
    keywords: ['儿童', '亲子', '青少年', '托管', '电玩', '游戏室', '家庭'],
  },
  {
    id: 'suite-privilege',
    label: '舒享舱房',
    icon: 'suite-privilege',
    exactNames: ['舒享舱房', '套房礼遇', '行政酒廊', '客房服务'],
    keywords: ['套房', '礼遇', '行政酒廊', '客房', '舱房', '贵宾', '私享'],
  },
  {
    id: 'sports-fitness',
    label: '运动健身',
    icon: 'sports-fitness',
    exactNames: ['运动健身', '泳池', 'SPA', '健身房', '球场'],
    keywords: ['泳池', 'spa', '健身', '球场', '跑道', '运动', '水疗', '冲浪'],
  },
  {
    id: 'other',
    label: '其它',
    icon: 'other',
    exactNames: ['其它', '其他'],
    keywords: [],
  },
]

function normalizeText(value: unknown) {
  return String(value ?? '').trim().toLowerCase()
}

function getMatchSource(item: Record<string, any>) {
  return [item.category_name, item.name, item.english_name]
    .map(normalizeText)
    .filter(Boolean)
}

export function resolveCruiseFacilityCategory(item: Record<string, any>) {
  const sources = getMatchSource(item)
  if (sources.length === 0) {
    return cruiseFacilityCategories[cruiseFacilityCategories.length - 1]
  }

  const exact = cruiseFacilityCategories.find((category) =>
    category.id !== 'other' && category.exactNames.some((name) => sources.includes(normalizeText(name))),
  )

  if (exact) {
    return exact
  }

  const partial = cruiseFacilityCategories.find((category) =>
    category.id !== 'other' && category.keywords.some((keyword) => sources.some((source) => source.includes(normalizeText(keyword)))),
  )

  return partial ?? cruiseFacilityCategories[cruiseFacilityCategories.length - 1]
}

export function groupCruiseFacilities(items: Record<string, any>[]) {
  const groups = new Map<CruiseFacilityCategoryId, Record<string, any>[]>()
  cruiseFacilityCategories.forEach((category) => groups.set(category.id, []))

  items.forEach((item) => {
    const category = resolveCruiseFacilityCategory(item)
    groups.get(category.id)?.push(item)
  })

  return groups
}