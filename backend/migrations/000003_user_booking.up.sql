CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  phone VARCHAR(20) UNIQUE,
  wx_open_id VARCHAR(80) UNIQUE,
  nickname VARCHAR(50),
  avatar_url VARCHAR(500),
  status SMALLINT DEFAULT 1,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS passengers (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  name VARCHAR(50) NOT NULL,
  id_type VARCHAR(20) NOT NULL,
  id_number VARCHAR(50) NOT NULL,
  birthday TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bookings (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id),
  voyage_id BIGINT NOT NULL,
  cabin_sku_id BIGINT NOT NULL,
  status VARCHAR(20) NOT NULL,
  total_cents BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS booking_passengers (
  id BIGSERIAL PRIMARY KEY,
  booking_id BIGINT NOT NULL REFERENCES bookings(id),
  passenger_id BIGINT NOT NULL REFERENCES passengers(id)
);

CREATE TABLE IF NOT EXISTS cabin_holds (
  id BIGSERIAL PRIMARY KEY,
  cabin_sku_id BIGINT NOT NULL,
  user_id BIGINT NOT NULL,
  qty INT NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_cabin_holds_sku ON cabin_holds(cabin_sku_id);
