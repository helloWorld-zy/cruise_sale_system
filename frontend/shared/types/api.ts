export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PaginatedData<T> {
  list: T[]
  pagination: Pagination
}

export interface Pagination {
  page: number
  page_size: number
  total: number
  total_pages: number
}

export type PaginatedResponse<T> = ApiResponse<PaginatedData<T>>

export interface ListQuery {
  page?: number
  page_size?: number
  keyword?: string
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}
