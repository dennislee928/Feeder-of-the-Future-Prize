-- 刪除 OAuth 關聯表
DROP INDEX IF EXISTS idx_user_oauth_provider;
DROP INDEX IF EXISTS idx_user_oauth_user_id;
DROP TABLE IF EXISTS user_oauth;

