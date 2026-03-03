-- Sprint 4.1 回滚：邮轮介绍模块字段补齐（Part A）
-- 警告：该回滚会删除新增列并丢失对应数据，执行前请先备份。

-- 索引回滚
DROP INDEX IF EXISTS idx_images_entity;
DROP INDEX IF EXISTS idx_cabin_types_code;
DROP INDEX IF EXISTS idx_cruises_status;
DROP INDEX IF EXISTS idx_cruises_code;

-- Image 回滚
ALTER TABLE images DROP COLUMN IF EXISTS is_primary;

-- FacilityCategory 回滚
ALTER TABLE facility_categories DROP COLUMN IF EXISTS status;

-- Facility 回滚
ALTER TABLE facilities DROP COLUMN IF EXISTS target_audience;
ALTER TABLE facilities DROP COLUMN IF EXISTS charge_price_tip;
ALTER TABLE facilities DROP COLUMN IF EXISTS extra_charge;
ALTER TABLE facilities DROP COLUMN IF EXISTS open_hours;

-- CabinType 回滚
ALTER TABLE cabin_types DROP COLUMN IF EXISTS floor_plan_url;
ALTER TABLE cabin_types DROP COLUMN IF EXISTS amenities;
ALTER TABLE cabin_types DROP COLUMN IF EXISTS tags;
ALTER TABLE cabin_types DROP COLUMN IF EXISTS bed_type;
ALTER TABLE cabin_types DROP COLUMN IF EXISTS max_capacity;
ALTER TABLE cabin_types DROP COLUMN IF EXISTS area_max;
ALTER TABLE cabin_types DROP COLUMN IF EXISTS area_min;
ALTER TABLE cabin_types DROP COLUMN IF EXISTS code;

-- Cruise 回滚
ALTER TABLE cruises DROP COLUMN IF EXISTS deck_count;
ALTER TABLE cruises DROP COLUMN IF EXISTS width;
ALTER TABLE cruises DROP COLUMN IF EXISTS length;
ALTER TABLE cruises DROP COLUMN IF EXISTS refurbish_year;
ALTER TABLE cruises DROP COLUMN IF EXISTS crew_count;
ALTER TABLE cruises DROP COLUMN IF EXISTS code;

-- CabinSKU 回滚
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS amenities;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS bed_type;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS has_balcony;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS has_window;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS orientation;
ALTER TABLE cabin_skus DROP COLUMN IF EXISTS position;

-- CabinInventory 回滚
ALTER TABLE cabin_inventories DROP COLUMN IF EXISTS alert_threshold;

-- CabinPrice 回滚
ALTER TABLE cabin_prices DROP COLUMN IF EXISTS price_type;
