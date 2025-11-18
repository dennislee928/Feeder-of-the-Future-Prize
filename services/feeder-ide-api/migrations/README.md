# Database Migrations

本目錄包含資料庫 migration 檔案。

## Migration 檔案命名規則

- `{序號}_{描述}.up.sql` - 升級 migration
- `{序號}_{描述}.down.sql` - 降級 migration

## 執行 Migration

### 使用 golang-migrate

```bash
# 安裝 golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 執行所有 migration
migrate -path ./migrations -database "postgres://feeder_user:feeder_password@localhost:5432/feeder_db?sslmode=disable" up

# 回滾一個版本
migrate -path ./migrations -database "postgres://feeder_user:feeder_password@localhost:5432/feeder_db?sslmode=disable" down 1
```

### 使用 Docker

```bash
# 在容器中執行 migration
docker exec -i feeder-postgres psql -U feeder_user -d feeder_db < migrations/001_create_users_table.up.sql
```

## Migration 順序

1. `001_create_users_table` - 創建用戶表
2. `002_create_oauth_table` - 創建 OAuth 關聯表
3. `003_create_subscriptions_table` - 創建訂閱表
4. `004_create_payments_table` - 創建付費記錄表
5. `005_add_user_id_to_topologies` - 為拓樸表添加用戶關聯
6. `006_create_user_quotas_table` - 創建用戶配額表

