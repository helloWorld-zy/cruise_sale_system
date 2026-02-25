-- 000003_route_voyage_cabin.down.sql
-- 回滚 Sprint 2 新增表：按依赖顺序删除航线、航次、舱房相关表。

DROP TABLE IF EXISTS inventory_logs;
DROP TABLE IF EXISTS cabin_inventories;
DROP TABLE IF EXISTS cabin_prices;
DROP TABLE IF EXISTS cabin_skus;
DROP TABLE IF EXISTS voyages;
DROP TABLE IF EXISTS routes;
