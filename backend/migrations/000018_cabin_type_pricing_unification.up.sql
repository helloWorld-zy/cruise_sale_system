-- 000018_cabin_type_pricing_unification.up.sql
-- Force-rebuild cabin-type domain data and introduce category/media/pricing-version schema.

-- Strong rebuild: clear existing cabin-related business data.
TRUNCATE TABLE
  order_status_logs,
  refunds,
  payments,
  cabin_holds,
  booking_passengers,
  bookings,
  inventory_logs,
  cabin_inventories,
  cabin_prices,
  cabin_skus,
  cabin_types
RESTART IDENTITY CASCADE;

CREATE TABLE IF NOT EXISTS cabin_type_categories (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(32) NOT NULL UNIQUE,
  status SMALLINT NOT NULL DEFAULT 1,
  sort_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

INSERT INTO cabin_type_categories (name, code, status, sort_order)
VALUES
  ('内舱房', 'INNER', 1, 10),
  ('海景房', 'OCEANVIEW', 1, 20),
  ('阳台房', 'BALCONY', 1, 30),
  ('套房', 'SUITE', 1, 40)
ON CONFLICT (code) DO NOTHING;

ALTER TABLE cabin_types
  ADD COLUMN IF NOT EXISTS category_id BIGINT,
  ADD COLUMN IF NOT EXISTS occupancy INT,
  ADD COLUMN IF NOT EXISTS intro TEXT NOT NULL DEFAULT '';

UPDATE cabin_types
SET occupancy = COALESCE(occupancy, capacity, 2)
WHERE occupancy IS NULL;

UPDATE cabin_types
SET category_id = (SELECT id FROM cabin_type_categories WHERE code = 'INNER' LIMIT 1)
WHERE category_id IS NULL;

ALTER TABLE cabin_types
  ALTER COLUMN occupancy SET NOT NULL,
  ALTER COLUMN category_id SET NOT NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'fk_cabin_types_category_id'
  ) THEN
    ALTER TABLE cabin_types
      ADD CONSTRAINT fk_cabin_types_category_id FOREIGN KEY (category_id) REFERENCES cabin_type_categories(id);
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS cabin_type_cruise_bindings (
  id BIGSERIAL PRIMARY KEY,
  cabin_type_id BIGINT NOT NULL REFERENCES cabin_types(id) ON DELETE CASCADE,
  cruise_id BIGINT NOT NULL REFERENCES cruises(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_cabin_type_cruise_binding UNIQUE (cabin_type_id, cruise_id)
);

-- Seed binding from legacy cruise_id if it exists.
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns
    WHERE table_name = 'cabin_types' AND column_name = 'cruise_id'
  ) THEN
    INSERT INTO cabin_type_cruise_bindings (cabin_type_id, cruise_id)
    SELECT id, cruise_id
    FROM cabin_types
    WHERE cruise_id IS NOT NULL
    ON CONFLICT (cabin_type_id, cruise_id) DO NOTHING;
  END IF;
END $$;

CREATE TABLE IF NOT EXISTS cabin_type_media (
  id BIGSERIAL PRIMARY KEY,
  cabin_type_id BIGINT NOT NULL REFERENCES cabin_types(id) ON DELETE CASCADE,
  media_type VARCHAR(20) NOT NULL,
  url TEXT NOT NULL,
  title VARCHAR(120) NOT NULL DEFAULT '',
  sort_order INT NOT NULL DEFAULT 0,
  is_primary BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT chk_cabin_type_media_type CHECK (media_type IN ('image', 'floor_plan'))
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_cabin_type_media_primary
  ON cabin_type_media(cabin_type_id, media_type)
  WHERE is_primary = TRUE AND deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS voyage_cabin_type_price_versions (
  id BIGSERIAL PRIMARY KEY,
  voyage_id BIGINT NOT NULL REFERENCES voyages(id) ON DELETE CASCADE,
  cabin_type_id BIGINT NOT NULL REFERENCES cabin_types(id) ON DELETE CASCADE,
  inventory_total INT NOT NULL,
  settlement_price_cents BIGINT NOT NULL,
  sale_price_cents BIGINT NOT NULL,
  effective_at TIMESTAMPTZ NOT NULL,
  created_by BIGINT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vct_price_versions_effective
  ON voyage_cabin_type_price_versions(voyage_id, cabin_type_id, effective_at DESC);

CREATE TABLE IF NOT EXISTS voyage_cabin_type_current (
  voyage_id BIGINT NOT NULL REFERENCES voyages(id) ON DELETE CASCADE,
  cabin_type_id BIGINT NOT NULL REFERENCES cabin_types(id) ON DELETE CASCADE,
  inventory_total INT NOT NULL,
  settlement_price_cents BIGINT NOT NULL,
  sale_price_cents BIGINT NOT NULL,
  effective_at TIMESTAMPTZ NOT NULL,
  version_id BIGINT NOT NULL REFERENCES voyage_cabin_type_price_versions(id) ON DELETE RESTRICT,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY (voyage_id, cabin_type_id)
);
