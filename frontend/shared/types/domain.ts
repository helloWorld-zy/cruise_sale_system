// shared/types/domain.ts — 前后端共享的领域模型类型定义
// 与后端 Go domain 结构体保持一一对应

/** 邮轮公司 */
export interface CruiseCompany {
  id: number              // 主键 ID
  name: string            // 公司名称（中文）
  english_name?: string   // 公司英文名称
  description?: string    // 公司简介
  logo_url?: string       // Logo 图片地址
  status: number          // 状态：1=启用，0=停用
}

/** 邮轮 */
export interface Cruise {
  id: number                   // 主键 ID
  company_id: number           // 所属公司 ID
  name: string                 // 邮轮名称（中文）
  english_name: string         // 邮轮英文名称
  code: string                 // 邮轮代码
  crew_count: number           // 船员人数
  build_year: number           // 建造年份
  refurbish_year: number       // 翻新年份
  tonnage: number              // 吨位
  passenger_capacity: number   // 载客量
  room_count: number           // 舱房数量
  length: number               // 长度（米）
  width: number                // 宽度（米）
  deck_count: number           // 甲板层数
  description: string          // 描述
  status: number               // 状态：1=上架，0=下架
  sort_order: number           // 排序
  created_at: string           // 创建时间
  updated_at: string           // 更新时间
}

/** 舱房类型 */
export interface CabinType {
  id: number                 // 主键 ID
  cruise_id: number          // 所属邮轮 ID
  name: string               // 舱房类型名称
  english_name: string       // 英文名称
  code: string               // 舱型代码
  area_min: number           // 面积最小值
  area_max: number           // 面积最大值
  area: number               // 面积（兼容旧字段）
  capacity: number           // 默认容纳人数（兼容旧字段）
  max_capacity: number       // 最大容纳人数
  bed_type: string           // 床型
  tags: string               // 标签
  amenities: string          // 设施
  floor_plan_url: string     // 户型图 URL
  deck: string               // 甲板说明
  description: string        // 描述
  status: number             // 状态：1=启用，0=停用
  sort_order: number         // 排序
}

/** 设施分类 */
export interface FacilityCategory {
  id: number                // 主键 ID
  name: string              // 分类名
  icon: string              // 分类图标
  sort_order: number        // 排序
  status: number            // 状态：1=启用，0=停用
}

/** 航线 */
export interface Route {
  id: number              // 主键 ID
  code: string            // 航线编码（唯一）
  name: string            // 航线名称
  description?: string    // 航线描述
  status: number          // 状态：1=启用，0=停用
}

/** 航次 */
export interface Voyage {
  id: number              // 主键 ID
  route_id: number        // 所属航线 ID
  cruise_id: number       // 执行邮轮 ID
  code: string            // 航次编码（唯一）
  depart_date: string     // 出发日期（ISO 格式）
  return_date: string     // 返航日期（ISO 格式）
  status: number          // 状态：1=开放预订，0=关闭
}

/** 舱房 SKU（最小可售单元） */
export interface CabinSKU {
  id: number               // 主键 ID
  voyage_id: number        // 所属航次 ID
  cabin_type_id: number    // 所属舱房类型 ID
  code: string             // 舱房编号（唯一）
  deck: string             // 所在甲板层
  area: number             // 面积（平方米）
  max_guests: number       // 最大入住人数
  position: string         // 位置（前中后）
  orientation: string      // 朝向
  has_window: boolean      // 是否有窗
  has_balcony: boolean     // 是否有阳台
  bed_type: string         // 床型
  amenities: string        // 设施
  status: number           // 状态：1=上架，0=下架
}

/** 舱房价格（按日期和入住人数定价） */
export interface CabinPrice {
  id: number              // 主键 ID
  cabin_sku_id: number    // 关联的舱房 SKU ID
  date: string            // 价格生效日期（ISO 格式）
  occupancy: number       // 入住人数
  price_cents: number     // 价格（单位：分）
  price_type: string      // 价格类型
}

/** 舱房库存 */
export interface CabinInventory {
  id: number             // 主键 ID
  cabin_sku_id: number   // 舱位 ID
  total: number          // 总库存
  locked: number         // 锁定库存
  sold: number           // 已售库存
  alert_threshold: number // 预警阈值
}

/** 设施 */
export interface Facility {
  id: number                 // 主键 ID
  category_id: number        // 分类 ID
  cruise_id: number          // 邮轮 ID
  name: string               // 名称
  english_name: string       // 英文名称
  location: string           // 位置
  open_hours: string         // 开放时间
  extra_charge: boolean      // 是否收费
  charge_price_tip: string   // 收费说明
  target_audience: string    // 适合人群
  description: string        // 描述
  status: number             // 状态
  sort_order: number         // 排序
}

/** 图片 */
export interface AppImage {
  id: number                 // 主键 ID
  entity_type: string        // 关联实体类型
  entity_id: number          // 关联实体 ID
  url: string                // 图片地址
  sort_order: number         // 排序
  is_primary: boolean        // 是否主图
}

/** C 端用户 */
export interface User {
  id: number              // 主键 ID
  phone?: string          // 手机号
  wx_open_id?: string     // 微信 OpenID
  nickname?: string       // 用户昵称
  avatar_url?: string     // 头像图片地址
  status: number          // 状态：1=启用，0=停用
}

/** 预订订单 */
export interface Booking {
  id: number              // 主键 ID
  user_id: number         // 下单用户 ID
  voyage_id: number       // 所属航次 ID
  cabin_sku_id: number    // 预订的舱房 SKU ID
  status: string          // 订单状态（created/paid/cancelled）
  total_cents: number     // 订单总金额（单位：分）
}
