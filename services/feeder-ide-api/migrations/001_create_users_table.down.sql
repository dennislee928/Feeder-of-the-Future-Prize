-- 刪除用戶表
DROP INDEX IF EXISTS idx_users_api_key;
DROP INDEX IF EXISTS idx_users_subscription_tier;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

