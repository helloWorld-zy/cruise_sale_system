-- 为第三方账号绑定表增加唯一约束，防止同一标识符重复绑定不同账号。
-- 如尚未建立 user_bindings 表，请先确保前序迁移已创建。
ALTER TABLE IF EXISTS user_bindings
    ADD CONSTRAINT uq_user_bindings_provider_identifier UNIQUE (provider, identifier);
