export type FacilityCategoryIconOption = {
  value: string
  label: string
}

export const facilityCategoryIconOptions: FacilityCategoryIconOption[] = [
  { value: 'included-dining', label: '免费餐厅' },
  { value: 'specialty-dining', label: '特色餐厅' },
  { value: 'bar-lounge', label: '酒吧' },
  { value: 'leisure-entertainment', label: '休闲娱乐' },
  { value: 'kids-family', label: '亲子童趣' },
  { value: 'suite-privilege', label: '舒享舱房' },
  { value: 'sports-fitness', label: '运动健身' },
  { value: 'other', label: '其它' },
]

export function ensureFacilityCategoryIconOptions(currentIcon: string) {
  const normalizedIcon = currentIcon.trim()
  if (!normalizedIcon) {
    return facilityCategoryIconOptions
  }
  if (facilityCategoryIconOptions.some((item) => item.value === normalizedIcon)) {
    return facilityCategoryIconOptions
  }
  return [{ value: normalizedIcon, label: normalizedIcon }, ...facilityCategoryIconOptions]
}

export function getFacilityCategoryIconLabel(icon: string) {
  const normalizedIcon = icon.trim()
  if (!normalizedIcon) {
    return '请选择图标'
  }
  return facilityCategoryIconOptions.find((item) => item.value === normalizedIcon)?.label ?? normalizedIcon
}