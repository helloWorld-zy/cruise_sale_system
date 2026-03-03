-- Route 扩展字段
ALTER TABLE routes ADD COLUMN IF NOT EXISTS departure_port VARCHAR(100);
ALTER TABLE routes ADD COLUMN IF NOT EXISTS arrival_port VARCHAR(100);
ALTER TABLE routes ADD COLUMN IF NOT EXISTS stops TEXT;
ALTER TABLE routes ADD COLUMN IF NOT EXISTS sort_order INT DEFAULT 0;

-- CabinSKU 扩展字段
ALTER TABLE cabin_skus ADD COLUMN IF NOT EXISTS position VARCHAR(20);
ALTER TABLE cabin_skus ADD COLUMN IF NOT EXISTS orientation VARCHAR(20);
ALTER TABLE cabin_skus ADD COLUMN IF NOT EXISTS has_window BOOLEAN DEFAULT FALSE;
ALTER TABLE cabin_skus ADD COLUMN IF NOT EXISTS has_balcony BOOLEAN DEFAULT FALSE;
ALTER TABLE cabin_skus ADD COLUMN IF NOT EXISTS bed_type VARCHAR(100);
ALTER TABLE cabin_skus ADD COLUMN IF NOT EXISTS amenities TEXT;
ALTER TABLE cabin_skus ADD COLUMN IF NOT EXISTS grade VARCHAR(50);

-- CabinPrice 扩展字段
ALTER TABLE cabin_prices ADD COLUMN IF NOT EXISTS child_price_cents BIGINT DEFAULT 0;
ALTER TABLE cabin_prices ADD COLUMN IF NOT EXISTS single_supplement_cents BIGINT DEFAULT 0;
ALTER TABLE cabin_prices ADD COLUMN IF NOT EXISTS price_type VARCHAR(20) DEFAULT 'base';

-- CabinInventory 扩展字段
ALTER TABLE cabin_inventories ADD COLUMN IF NOT EXISTS alert_threshold INT DEFAULT 0;

-- 索引
CREATE INDEX IF NOT EXISTS idx_cabin_prices_type ON cabin_prices(price_type);
CREATE INDEX IF NOT EXISTS idx_routes_departure ON routes(departure_port);
