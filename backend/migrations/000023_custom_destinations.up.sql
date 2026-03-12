-- 自定义目的地表（私属岛屿等搜索不到的特殊目的地）
CREATE TABLE IF NOT EXISTS custom_destinations (
    id            BIGSERIAL    PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,           -- 目的地名称（如"可可湾"）
    country       VARCHAR(100) NOT NULL DEFAULT '', -- 所属国家/地区（如"巴哈马"）
    latitude      DOUBLE PRECISION,                 -- 纬度
    longitude     DOUBLE PRECISION,                 -- 经度
    keywords      TEXT         NOT NULL DEFAULT '', -- 搜索关键词（逗号分隔，如"CocoCay,Perfect Day"）
    description   TEXT         NOT NULL DEFAULT '', -- 备注描述
    status        SMALLINT     NOT NULL DEFAULT 1,  -- 1=启用, 0=停用
    sort_order    INT          NOT NULL DEFAULT 0,  -- 排序权重
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_custom_destinations_deleted_at ON custom_destinations (deleted_at);
CREATE UNIQUE INDEX idx_custom_destinations_name_country ON custom_destinations (name, country) WHERE deleted_at IS NULL;
