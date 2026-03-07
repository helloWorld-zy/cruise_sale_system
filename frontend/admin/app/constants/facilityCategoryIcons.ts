export type FacilityCategoryIconOption = {
  value: string
  label: string
}

export const facilityCategoryIconOptions: FacilityCategoryIconOption[] = [
  { value: 'utensils', label: '餐饮' },
  { value: 'music', label: '演出' },
  { value: 'spa', label: '水疗' },
  { value: 'dumbbell', label: '健身' },
  { value: 'swimmer', label: '泳池' },
  { value: 'gamepad', label: '娱乐' },
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