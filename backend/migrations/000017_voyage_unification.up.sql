-- 000017_voyage_unification.up.sql
-- Merge route management into voyage management and add structured itinerary rows.

ALTER TABLE voyages
  ADD COLUMN IF NOT EXISTS brief_info VARCHAR(300) NOT NULL DEFAULT '';

UPDATE voyages
SET brief_info = ''
WHERE brief_info IS NULL;

DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'chk_voyages_return_after_depart'
  ) THEN
    ALTER TABLE voyages
      ADD CONSTRAINT chk_voyages_return_after_depart CHECK (return_date >= depart_date);
  END IF;
END $$;

ALTER TABLE voyages
  DROP CONSTRAINT IF EXISTS voyages_route_id_fkey;

DROP INDEX IF EXISTS idx_voyages_route_id;

ALTER TABLE voyages
  DROP COLUMN IF EXISTS route_id;

DROP TABLE IF EXISTS routes;

CREATE TABLE IF NOT EXISTS voyage_itineraries (
  id BIGSERIAL PRIMARY KEY,
  voyage_id BIGINT NOT NULL REFERENCES voyages(id) ON DELETE CASCADE,
  day_no INT NOT NULL,
  stop_index INT NOT NULL,
  city VARCHAR(120) NOT NULL,
  summary TEXT NOT NULL DEFAULT '',
  eta_time TIME,
  etd_time TIME,
  has_breakfast BOOLEAN NOT NULL DEFAULT FALSE,
  has_lunch BOOLEAN NOT NULL DEFAULT FALSE,
  has_dinner BOOLEAN NOT NULL DEFAULT FALSE,
  has_accommodation BOOLEAN NOT NULL DEFAULT FALSE,
  accommodation_text VARCHAR(300) NOT NULL DEFAULT '',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_voyage_itinerary_day_stop UNIQUE (voyage_id, day_no, stop_index),
  CONSTRAINT chk_voyage_itinerary_day_no CHECK (day_no >= 1),
  CONSTRAINT chk_voyage_itinerary_stop_index CHECK (stop_index >= 1)
);

CREATE INDEX IF NOT EXISTS idx_voyage_itineraries_voyage_day_stop
  ON voyage_itineraries(voyage_id, day_no, stop_index);
