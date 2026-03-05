-- Rollback Sprint 4.3 company hardening

DROP INDEX IF EXISTS idx_cruise_companies_english_name;
DROP INDEX IF EXISTS idx_cruise_companies_name;
