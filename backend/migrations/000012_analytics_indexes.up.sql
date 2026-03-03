-- Improve analytics query performance for dashboard-related aggregations.
CREATE INDEX IF NOT EXISTS idx_bookings_created_at ON bookings (created_at);
CREATE INDEX IF NOT EXISTS idx_bookings_status_created_at ON bookings (status, created_at);
CREATE INDEX IF NOT EXISTS idx_bookings_cabin_sku_id ON bookings (cabin_sku_id);

CREATE INDEX IF NOT EXISTS idx_payments_status_created_at ON payments (status, created_at);
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments (created_at);

CREATE INDEX IF NOT EXISTS idx_cabin_inventories_alert_threshold ON cabin_inventories (alert_threshold);
