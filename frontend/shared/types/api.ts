// shared/types/api.ts — 前后端共享的 API 通用类型定义

// ApiResponse 是后端统一返回的 JSON 信封结构，泛型 T 表示 data 字段的实际类型。
export interface ApiResponse<T = unknown> {
  code: number      // 业务状态码（0 表示成功）
  message: string   // 响应消息
  data: T           // 响应数据
}

// PaginatedData 表示分页查询的返回数据，包含列表和分页信息。
export interface PaginatedData<T> {
  list: T[]                  // 当前页的数据列表
  pagination: Pagination     // 分页元信息
}

// Pagination 包含分页的元信息。
export interface Pagination {
  page: number        // 当前页码
  page_size: number   // 每页数量
  total: number       // 总记录数
  total_pages: number // 总页数
}

// PaginatedResponse 是 ApiResponse 与 PaginatedData 的组合类型，用于分页接口的返回值。
export type PaginatedResponse<T> = ApiResponse<PaginatedData<T>>

// ListQuery 定义列表查询的通用参数结构。
export interface ListQuery {
  page?: number                    // 页码（默认 1）
  page_size?: number               // 每页数量（默认 10）
  keyword?: string                 // 关键词搜索
  sort_by?: string                 // 排序字段
  sort_order?: 'asc' | 'desc'     // 排序方向
}
