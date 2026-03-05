-- Sprint 4.3: company admin hardening
-- Ensure company display fields are clean and searchable for admin CRUD/selectors.

ALTER TABLE cruise_companies
    ALTER COLUMN name SET NOT NULL;

UPDATE cruise_companies
SET english_name = NULL
WHERE english_name = '';

UPDATE cruise_companies
SET logo_url = NULL
WHERE logo_url = '';

CREATE INDEX IF NOT EXISTS idx_cruise_companies_name ON cruise_companies(name);
CREATE INDEX IF NOT EXISTS idx_cruise_companies_english_name ON cruise_companies(english_name);
