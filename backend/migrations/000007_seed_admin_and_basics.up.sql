-- 000007_seed_admin_and_basics.up.sql
-- 初始化最小可用后台数据：管理员账号、基础公司/邮轮/航线/航次/舱房与 C 端测试用户。

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 默认管理员：admin / admin123
INSERT INTO staffs (username, password_hash, real_name, status, created_at, updated_at)
SELECT 'admin', crypt('admin123', gen_salt('bf')), '系统管理员', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM staffs WHERE username = 'admin');

-- 最小公司数据，便于创建邮轮时引用 company_id。
INSERT INTO cruise_companies (id, name, created_at, updated_at)
SELECT 1, '默认邮轮公司', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM cruise_companies WHERE id = 1);

-- 最小邮轮数据，便于创建航次时引用 cruise_id。
INSERT INTO cruises (id, company_id, name, created_at, updated_at)
SELECT 1, 1, '默认测试邮轮', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM cruises WHERE id = 1);

-- 最小舱型数据，便于创建舱房 SKU 时引用 cabin_type_id。
INSERT INTO cabin_types (id, cruise_id, name, created_at, updated_at)
SELECT 1, 1, '标准内舱', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM cabin_types WHERE id = 1);

-- 最小航线与航次数据，便于创建舱房与订单测试。
INSERT INTO routes (id, code, name, created_at, updated_at)
SELECT 1, 'R-DEFAULT', '默认航线', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM routes WHERE id = 1);

INSERT INTO voyages (id, route_id, cruise_id, code, depart_date, return_date, created_at, updated_at)
SELECT 1, 1, 1, 'V-DEFAULT', NOW() + INTERVAL '7 day', NOW() + INTERVAL '10 day', NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM voyages WHERE id = 1);

INSERT INTO cabin_skus (id, voyage_id, cabin_type_id, code, max_guests, created_at, updated_at)
SELECT 1, 1, 1, 'SKU-DEFAULT', 2, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM cabin_skus WHERE id = 1);

INSERT INTO cabin_inventories (cabin_sku_id, total, locked, sold, updated_at)
SELECT 1, 20, 0, 0, NOW()
WHERE NOT EXISTS (SELECT 1 FROM cabin_inventories WHERE cabin_sku_id = 1);

INSERT INTO cabin_prices (cabin_sku_id, date, occupancy, price_cents, created_at, updated_at)
SELECT 1, NOW()::date, 2, 19900, NOW(), NOW()
WHERE NOT EXISTS (
	SELECT 1 FROM cabin_prices WHERE cabin_sku_id = 1 AND occupancy = 2 AND date::date = NOW()::date
);

-- C 端测试用户，便于管理后台测试“新建订单”链路。
INSERT INTO users (id, phone, nickname, status, created_at, updated_at)
SELECT 1, '13800138000', '测试用户', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 1);

INSERT INTO users (id, phone, nickname, status, created_at, updated_at)
SELECT 2, '13800138001', '测试用户-2', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM users WHERE id = 2);

-- 修复序列，避免显式 ID 插入后下次自增冲突。
SELECT setval('cruise_companies_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM cruise_companies), 1));
SELECT setval('cruises_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM cruises), 1));
SELECT setval('cabin_types_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM cabin_types), 1));
SELECT setval('routes_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM routes), 1));
SELECT setval('voyages_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM voyages), 1));
SELECT setval('cabin_skus_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM cabin_skus), 1));
SELECT setval('users_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM users), 1));
