export interface CruiseCompany {
  id: number
  name: string
  english_name?: string
  description?: string
  logo_url?: string
  status: number
}

export interface Cruise {
  id: number
  company_id: number
  name: string
  english_name?: string
  status: number
}

export interface CabinType {
  id: number
  cruise_id: number
  name: string
  capacity: number
  area?: number
  status: number
}

export interface Route {
  id: number
  code: string
  name: string
  description?: string
  status: number
}

export interface Voyage {
  id: number
  route_id: number
  cruise_id: number
  code: string
  depart_date: string
  return_date: string
  status: number
}

export interface CabinSKU {
  id: number
  voyage_id: number
  cabin_type_id: number
  code: string
  deck?: string
  area?: number
  max_guests: number
  status: number
}

export interface CabinPrice {
  id: number
  cabin_sku_id: number
  date: string
  occupancy: number
  price_cents: number
}

export interface User {
  id: number
  phone?: string
  wx_open_id?: string
  nickname?: string
  avatar_url?: string
  status: number
}

export interface Booking {
  id: number
  user_id: number
  voyage_id: number
  cabin_sku_id: number
  status: string
  total_cents: number
}
