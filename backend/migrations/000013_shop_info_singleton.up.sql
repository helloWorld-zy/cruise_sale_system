-- Enforce singleton row for shop_info: only id=1 is allowed.
ALTER TABLE shop_info
ADD CONSTRAINT ck_shop_info_singleton_id CHECK (id = 1);
