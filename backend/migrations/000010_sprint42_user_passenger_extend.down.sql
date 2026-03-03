DROP INDEX IF EXISTS idx_passengers_favorite;

ALTER TABLE passengers DROP COLUMN IF EXISTS is_favorite;
ALTER TABLE passengers DROP COLUMN IF EXISTS special_needs;
ALTER TABLE passengers DROP COLUMN IF EXISTS emergency_phone;
ALTER TABLE passengers DROP COLUMN IF EXISTS emergency_contact;
ALTER TABLE passengers DROP COLUMN IF EXISTS email;
ALTER TABLE passengers DROP COLUMN IF EXISTS phone;
ALTER TABLE passengers DROP COLUMN IF EXISTS english_name;

ALTER TABLE users DROP COLUMN IF EXISTS email;
ALTER TABLE users DROP COLUMN IF EXISTS alipay_uid;
