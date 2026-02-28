-- 000006_cabin_hold_unique.up.sql
-- 为 cabin_holds 表添加唯一索引，确保同一用户同一 SKU 只能存在一条占座记录。
-- 先清理可能的重复数据，再创建唯一约束。

DELETE FROM cabin_holds
WHERE id NOT IN (
  SELECT MAX(id)
  FROM cabin_holds
  GROUP BY cabin_sku_id, user_id
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_cabin_holds_sku_user ON cabin_holds(cabin_sku_id, user_id);
