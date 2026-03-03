export type OrderStatus =
  | 'created'
  | 'pending_payment'
  | 'paid'
  | 'confirmed'
  | 'pending_travel'
  | 'traveling'
  | 'completed'
  | 'cancelled'
  | 'refunding'
  | 'refunded'

export interface OrderSummary {
  id: number
  booking_no?: string
  status: OrderStatus | string
  total_cents: number
  created_at?: string
  voyage_id?: number
  route_name?: string
  departure_port?: string
  travel_date?: string
}
