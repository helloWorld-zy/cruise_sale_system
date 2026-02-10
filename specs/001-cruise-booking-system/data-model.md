# Data Model: Cruise Booking System

**Feature**: `001-cruise-booking-system`
**Database**: PostgreSQL 17

## Core Product Entities

### `cruises` (Cruise Ships)
*   `id`: UUID (PK)
*   `name_en`: VARCHAR
*   `name_cn`: VARCHAR
*   `code`: VARCHAR (Unique)
*   `company_id`: UUID (FK -> cruise_companies)
*   `tonnage`: INT
*   `capacity`: INT
*   `decks`: INT
*   `status`: ENUM ('active', 'maintenance', 'retired')
*   `gallery`: JSONB (Array of image URLs)
*   `description`: TEXT (HTML/JSON)

### `cabin_types`
*   `id`: UUID (PK)
*   `cruise_id`: UUID (FK -> cruises)
*   `name`: VARCHAR
*   `code`: VARCHAR
*   `base_area`: DECIMAL
*   `capacity`: INT
*   `facilities`: JSONB (List of amenities)
*   `images`: JSONB

### `facilities`
*   `id`: UUID (PK)
*   `cruise_id`: UUID
*   `category_id`: UUID
*   `name`: VARCHAR
*   `location`: VARCHAR
*   `is_paid`: BOOLEAN

## Inventory & Sales

### `voyages` (Sailings)
*   `id`: UUID (PK)
*   `cruise_id`: UUID
*   `route_id`: UUID
*   `departure_date`: TIMESTAMPTZ
*   `return_date`: TIMESTAMPTZ
*   `status`: ENUM ('scheduled', 'sailing', 'completed', 'cancelled')

### `cabins` (Physical Rooms)
*   `id`: UUID (PK)
*   `cruise_id`: UUID
*   `cabin_type_id`: UUID
*   `number`: VARCHAR
*   `deck`: INT
*   `location_type`: ENUM ('forward', 'mid', 'aft')

### `inventory` (Per Voyage/Cabin Type)
*   `id`: UUID (PK)
*   `voyage_id`: UUID
*   `cabin_type_id`: UUID
*   `total_qty`: INT
*   `available_qty`: INT
*   `reserved_qty`: INT (Locked but not paid)
*   `sold_qty`: INT

### `price_rules`
*   `id`: UUID (PK)
*   `voyage_id`: UUID
*   `cabin_type_id`: UUID
*   `date`: DATE (For daily pricing)
*   `price_adult`: DECIMAL
*   `price_child`: DECIMAL
*   `currency`: VARCHAR (Default 'CNY')

## Order Management

### `orders`
*   `id`: UUID (PK)
*   `order_no`: VARCHAR (Unique, Human readable)
*   `user_id`: UUID
*   `voyage_id`: UUID
*   `status`: ENUM ('pending', 'paid', 'confirmed', 'cancelled', 'refunded', 'completed')
*   `total_amount`: DECIMAL
*   `currency`: VARCHAR
*   `created_at`: TIMESTAMPTZ
*   `expires_at`: TIMESTAMPTZ (Auto-cancel time)

### `order_items`
*   `id`: UUID (PK)
*   `order_id`: UUID
*   `cabin_id`: UUID (Nullable if "Guaranteed" cabin not yet assigned)
*   `cabin_type_id`: UUID
*   `price_snapshot`: DECIMAL

### `passengers`
*   `id`: UUID (PK)
*   `order_id`: UUID
*   `name_cn`: VARCHAR
*   `name_en`: VARCHAR
*   `doc_type`: ENUM ('passport', 'id_card')
*   `doc_number`: VARCHAR
*   `phone`: VARCHAR

## User System

### `users`
*   `id`: UUID (PK)
*   `phone`: VARCHAR (Unique)
*   `wx_openid`: VARCHAR
*   `password_hash`: VARCHAR

### `staffs`
*   `id`: UUID (PK)
*   `username`: VARCHAR
*   `role_id`: UUID (Casbin Role)
