-- 000007_seed_admin_and_basics.down.sql
-- 回滚默认种子数据。

DELETE FROM users WHERE id = 1 AND phone = '13800138000';
DELETE FROM users WHERE id = 2 AND phone = '13800138001';
DELETE FROM cabin_prices WHERE cabin_sku_id = 1 AND occupancy = 2;
DELETE FROM cabin_inventories WHERE cabin_sku_id = 1;
DELETE FROM cabin_skus WHERE id = 1 AND code = 'SKU-DEFAULT';
DELETE FROM voyages WHERE id = 1 AND code = 'V-DEFAULT';
DELETE FROM routes WHERE id = 1 AND code = 'R-DEFAULT';
DELETE FROM cabin_types WHERE id = 1 AND name = '标准内舱';
DELETE FROM cruises WHERE id = 1 AND name = '默认测试邮轮';
DELETE FROM cruise_companies WHERE id = 1 AND name = '默认邮轮公司';
DELETE FROM staffs WHERE username = 'admin';
