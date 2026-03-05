ALTER TABLE facility_categories
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_facility_categories_deleted_at
  ON facility_categories(deleted_at);
