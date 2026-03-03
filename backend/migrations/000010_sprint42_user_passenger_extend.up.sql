-- User extension fields
ALTER TABLE users ADD COLUMN IF NOT EXISTS alipay_uid VARCHAR(80) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS email VARCHAR(100);

-- Passenger extension fields
ALTER TABLE passengers ADD COLUMN IF NOT EXISTS english_name VARCHAR(100);
ALTER TABLE passengers ADD COLUMN IF NOT EXISTS phone VARCHAR(20);
ALTER TABLE passengers ADD COLUMN IF NOT EXISTS email VARCHAR(100);
ALTER TABLE passengers ADD COLUMN IF NOT EXISTS emergency_contact VARCHAR(50);
ALTER TABLE passengers ADD COLUMN IF NOT EXISTS emergency_phone VARCHAR(20);
ALTER TABLE passengers ADD COLUMN IF NOT EXISTS special_needs TEXT;
ALTER TABLE passengers ADD COLUMN IF NOT EXISTS is_favorite BOOLEAN DEFAULT FALSE;

CREATE INDEX IF NOT EXISTS idx_passengers_favorite ON passengers(user_id, is_favorite);
