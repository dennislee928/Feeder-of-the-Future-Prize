-- 為拓樸表添加 user_id 欄位
-- 注意：如果 topologies 表不存在，需要先創建

-- 創建拓樸表（如果不存在）
CREATE TABLE IF NOT EXISTS topologies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    profile_type VARCHAR(50) DEFAULT 'suburban' CHECK (profile_type IN ('rural', 'suburban', 'urban')),
    nodes JSONB NOT NULL DEFAULT '[]'::jsonb,
    lines JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 添加索引
CREATE INDEX IF NOT EXISTS idx_topologies_user_id ON topologies(user_id);
CREATE INDEX IF NOT EXISTS idx_topologies_created_at ON topologies(created_at);

