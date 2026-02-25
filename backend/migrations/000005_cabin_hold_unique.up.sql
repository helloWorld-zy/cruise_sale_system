DELETE FROM cabin_holds
WHERE id NOT IN (
  SELECT MAX(id)
  FROM cabin_holds
  GROUP BY cabin_sku_id, user_id
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_cabin_holds_sku_user ON cabin_holds(cabin_sku_id, user_id);
