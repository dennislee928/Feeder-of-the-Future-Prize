-- 創建用戶配額表
CREATE TABLE IF NOT EXISTS user_quotas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    max_topologies INTEGER DEFAULT 3,
    used_topologies INTEGER DEFAULT 0,
    max_simulations_per_day INTEGER DEFAULT 10,
    used_simulations_today INTEGER DEFAULT 0,
    last_simulation_reset_date DATE DEFAULT CURRENT_DATE,
    can_use_3d_rendering BOOLEAN DEFAULT FALSE,
    can_use_ai_prediction BOOLEAN DEFAULT FALSE,
    can_use_advanced_security BOOLEAN DEFAULT FALSE,
    can_access_api BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_quotas_user_id ON user_quotas(user_id);

