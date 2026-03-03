-- Route 扩展字段回滚
ALTER TABLE routes DROP COLUMN IF EXISTS departure_port;
ALTER TABLE routes DROP COLUMN IF EXISTS arrival_port;
ALTER TABLE routes DROP COLUMN IF EXISTS stops;
ALTER TABLE routes DROP COLUMN IF EXISTS sort_order;

-- CabinSKU 扩展字段回滚
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS position;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS orientation;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS has_window;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS has_balcony;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS bed_type;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS amenities;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS grade;

-- CabinPrice 扩展字段回滚
ALTER TABLE cabin_prices DROP COLUMN IF EXISTS child_price_cents;
ALTER TABLE cabin_prices DROP COLUMN IF EXISTS single_supplement_cents;
ALTER TABLE cabin_prices DROP COLUMN IF EXISTS price_type;

-- CabinInventory 扩展字段回滚
ALTER TABLE cabin_inventories DROP COLUMN IF EXISTS alert_threshold;

-- 索引回滚
DROP INDEX IF EXISTS idx_cabin_prices_type;
DROP INDEX IF EXISTS idx_routes_departure;
