DROP INDEX IF EXISTS idx_facility_categories_deleted_at;

ALTER TABLE facility_categories
  DROP COLUMN IF EXISTS deleted_at;
