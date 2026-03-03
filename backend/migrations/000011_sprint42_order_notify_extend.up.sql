-- 订单状态转换日志
CREATE TABLE IF NOT EXISTS order_status_logs (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES bookings(id),
    from_status VARCHAR(30) NOT NULL,
    to_status VARCHAR(30) NOT NULL,
    operator_id BIGINT,
    remark TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_order_status_logs_order ON order_status_logs(order_id);

-- 通知模板
CREATE TABLE IF NOT EXISTS notification_templates (
    id BIGSERIAL PRIMARY KEY,
    event_type VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    template TEXT NOT NULL,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_notification_templates_event ON notification_templates(event_type);

-- 店铺信息
CREATE TABLE IF NOT EXISTS shop_info (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100),
    logo VARCHAR(500),
    contact_phone VARCHAR(20),
    contact_email VARCHAR(100),
    company_desc TEXT,
    service_desc TEXT,
    icp_number VARCHAR(50),
    business_license VARCHAR(100),
    address VARCHAR(200),
    wechat VARCHAR(50)
);

-- 退款规则
CREATE TABLE IF NOT EXISTS refund_rules (
    id BIGSERIAL PRIMARY KEY,
    min_days INTEGER NOT NULL,
    max_days INTEGER NOT NULL,
    refund_rate INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 财务对账记录
CREATE TABLE IF NOT EXISTS reconciliations (
    id BIGSERIAL PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    total_payments BIGINT DEFAULT 0,
    total_payment_amount BIGINT DEFAULT 0,
    total_refund_amount BIGINT DEFAULT 0,
    discrepancy_count BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_reconciliations_date ON reconciliations(date);

-- 操作日志
CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    staff_id BIGINT,
    operation VARCHAR(50),
    resource VARCHAR(50),
    resource_id BIGINT,
    details TEXT,
    ip_address VARCHAR(50),
    user_agent VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_operation_logs_staff ON operation_logs(staff_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_created ON operation_logs(created_at);

-- 员工角色
ALTER TABLE staff ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'operator';
