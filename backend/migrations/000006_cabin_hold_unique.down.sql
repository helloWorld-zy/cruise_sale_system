-- 000006_cabin_hold_unique.down.sql
-- 回滚：删除 cabin_holds 的唯一索引约束。

DROP INDEX IF EXISTS uq_cabin_holds_sku_user;
