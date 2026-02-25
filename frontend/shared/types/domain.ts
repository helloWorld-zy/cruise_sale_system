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
  id: number              // 主键 ID
  company_id: number      // 所属公司 ID
  name: string            // 邮轮名称（中文）
  english_name?: string   // 邮轮英文名称
  status: number          // 状态：1=上架，0=下架
}

/** 舱房类型 */
export interface CabinType {
  id: number              // 主键 ID
  cruise_id: number       // 所属邮轮 ID
  name: string            // 舱房类型名称
  capacity: number        // 默认容纳人数
  area?: number           // 面积（平方米）
  status: number          // 状态：1=启用，0=停用
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
  id: number              // 主键 ID
  voyage_id: number       // 所属航次 ID
  cabin_type_id: number   // 所属舱房类型 ID
  code: string            // 舱房编号（唯一）
  deck?: string           // 所在甲板层
  area?: number           // 面积（平方米）
  max_guests: number      // 最大入住人数
  status: number          // 状态：1=上架，0=下架
}

/** 舱房价格（按日期和入住人数定价） */
export interface CabinPrice {
  id: number              // 主键 ID
  cabin_sku_id: number    // 关联的舱房 SKU ID
  date: string            // 价格生效日期（ISO 格式）
  occupancy: number       // 入住人数
  price_cents: number     // 价格（单位：分）
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
