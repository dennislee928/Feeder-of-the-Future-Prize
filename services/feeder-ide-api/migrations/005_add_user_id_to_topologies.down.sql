-- 刪除拓樸表（如果需要回滾）
DROP INDEX IF EXISTS idx_topologies_created_at;
DROP INDEX IF EXISTS idx_topologies_user_id;
DROP TABLE IF EXISTS topologies;

