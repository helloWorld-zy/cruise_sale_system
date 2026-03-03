ALTER TABLE IF EXISTS user_bindings
    DROP CONSTRAINT IF EXISTS uq_user_bindings_provider_identifier;
