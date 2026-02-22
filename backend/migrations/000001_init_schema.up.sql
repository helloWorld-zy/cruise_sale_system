CREATE TABLE cruise_companies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    english_name VARCHAR(100),
    description TEXT,
    logo_url VARCHAR(500),
    status SMALLINT DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_cruise_companies_deleted_at ON cruise_companies(deleted_at);

CREATE TABLE cruises (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES cruise_companies(id),
    name VARCHAR(100) NOT NULL,
    english_name VARCHAR(100),
    build_year INTEGER,
    tonnage DOUBLE PRECISION,
    passenger_capacity INTEGER,
    room_count INTEGER,
    description TEXT,
    status SMALLINT DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_cruises_company_id ON cruises(company_id);
CREATE INDEX idx_cruises_deleted_at ON cruises(deleted_at);

CREATE TABLE cabin_types (
    id BIGSERIAL PRIMARY KEY,
    cruise_id BIGINT NOT NULL REFERENCES cruises(id),
    name VARCHAR(100) NOT NULL,
    english_name VARCHAR(100),
    capacity INTEGER DEFAULT 2,
    area DOUBLE PRECISION,
    deck VARCHAR(50),
    description TEXT,
    status SMALLINT DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_cabin_types_cruise_id ON cabin_types(cruise_id);
CREATE INDEX idx_cabin_types_deleted_at ON cabin_types(deleted_at);

CREATE TABLE facility_categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    icon VARCHAR(255),
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE facilities (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES facility_categories(id),
    cruise_id BIGINT NOT NULL REFERENCES cruises(id),
    name VARCHAR(100) NOT NULL,
    english_name VARCHAR(100),
    location VARCHAR(100),
    description TEXT,
    status SMALLINT DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_facilities_category_id ON facilities(category_id);
CREATE INDEX idx_facilities_cruise_id ON facilities(cruise_id);
CREATE INDEX idx_facilities_deleted_at ON facilities(deleted_at);

CREATE TABLE images (
    id BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id BIGINT NOT NULL,
    url VARCHAR(500) NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_images_entity_type_id ON images(entity_type, entity_id);

CREATE TABLE staffs (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    real_name VARCHAR(50),
    phone VARCHAR(20),
    email VARCHAR(100),
    avatar_url VARCHAR(500),
    status SMALLINT DEFAULT 1,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_staffs_deleted_at ON staffs(deleted_at);

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
