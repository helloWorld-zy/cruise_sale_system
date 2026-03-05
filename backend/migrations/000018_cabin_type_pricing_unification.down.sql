-- Rollback 000018: drop cabin-type unification artifacts.

DROP TABLE IF EXISTS voyage_cabin_type_current;
DROP TABLE IF EXISTS voyage_cabin_type_price_versions;
DROP TABLE IF EXISTS cabin_type_media;
DROP TABLE IF EXISTS cabin_type_cruise_bindings;

ALTER TABLE cabin_types
  DROP CONSTRAINT IF EXISTS fk_cabin_types_category_id;

ALTER TABLE cabin_types
  DROP COLUMN IF EXISTS intro,
  DROP COLUMN IF EXISTS occupancy,
  DROP COLUMN IF EXISTS category_id;

DROP TABLE IF EXISTS cabin_type_categories;
