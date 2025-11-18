-- 刪除付費記錄表
DROP INDEX IF EXISTS idx_payments_provider_id;
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_subscription_id;
DROP INDEX IF EXISTS idx_payments_user_id;
DROP TABLE IF EXISTS payments;

