-- 可以在这里初始化额外的 database 或 roles
-- 根据 docker-compose 设置，DB `cruise_booking` 会自动创建

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建简单的测试数据或仅验证连接
SELECT 1;
