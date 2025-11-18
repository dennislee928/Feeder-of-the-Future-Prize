-- 刪除訂閱表
DROP INDEX IF EXISTS idx_subscriptions_payment_subscription_id;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_user_id;
DROP TABLE IF EXISTS subscriptions;

