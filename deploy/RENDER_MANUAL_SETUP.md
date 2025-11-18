# Render 手動部署指南

由於 Render Blueprint 格式可能因版本而異，如果 `render.yaml` 無法直接使用，請按照以下步驟手動創建每個服務。

## 服務列表

需要創建以下 8 個 Web Service：

1. `feeder-ide-frontend` - 前端
2. `feeder-ide-api` - IDE API
3. `feeder-sim-engine` - 模擬引擎
4. `feeder-os-controller` - Feeder OS 控制器
5. `app-der-ev-orchestrator` - DER/EV 編排器
6. `app-rural-resilience` - 農村韌性引擎
7. `security-gateway` - 安全網關
8. `telemetry-collector` - 遙測收集器

## 手動創建服務步驟

### 1. 前端服務 (feeder-ide-frontend)

1. 在 Render Dashboard 點擊 "New +" → "Web Service"
2. 連接你的 GitHub 倉庫
3. 設置：
   - **Name**: `feeder-ide-frontend`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./frontend/ide-frontend/Dockerfile`
   - **Docker Context**: `./frontend/ide-frontend`
   - **Plan**: `Starter` ($7/月) 或 `Free` (測試用)
4. 環境變數（**先設置臨時值，部署後端後再更新**）：
   - `VITE_API_BASE_URL`: `https://feeder-ide-api.onrender.com/api/v1`
   - `VITE_SIM_API_BASE_URL`: `https://feeder-sim-engine.onrender.com`
5. 點擊 "Create Web Service"

### 2. Feeder IDE API (feeder-ide-api)

1. "New +" → "Web Service"
2. 連接 GitHub 倉庫
3. 設置：
   - **Name**: `feeder-ide-api`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./services/feeder-ide-api/Dockerfile`
   - **Docker Context**: `./services/feeder-ide-api`
   - **Plan**: `Starter` 或 `Free`
4. 環境變數：
   - `PORT`: `8080`
5. 點擊 "Create Web Service"
6. **記錄服務 URL**（例如: `https://feeder-ide-api-xxx.onrender.com`）

### 3. Feeder Sim Engine (feeder-sim-engine)

1. "New +" → "Web Service"
2. 連接 GitHub 倉庫
3. 設置：
   - **Name**: `feeder-sim-engine`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./services/feeder-sim-engine/Dockerfile`
   - **Docker Context**: `./services/feeder-sim-engine`
   - **Plan**: `Starter` 或 `Free`
4. 點擊 "Create Web Service"
5. **記錄服務 URL**（例如: `https://feeder-sim-engine-xxx.onrender.com`）

### 4. Feeder OS Controller (feeder-os-controller)

1. "New +" → "Web Service"
2. 連接 GitHub 倉庫
3. 設置：
   - **Name**: `feeder-os-controller`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./services/feeder-os-controller/Dockerfile`
   - **Docker Context**: `./services/feeder-os-controller`
   - **Plan**: `Starter` 或 `Free`
4. 環境變數：
   - `PORT`: `8082`
   - `MQTT_BROKER`: `<你的 MQTT broker URL>`（例如: `m20.cloudmqtt.com`）
   - `MQTT_PORT`: `1883` 或 `8883`（取決於你的 MQTT 服務）
   - `MQTT_CLIENT_ID`: `feeder-os-controller`
   - `MQTT_USERNAME`: `<MQTT 用戶名>`（如果需要）
   - `MQTT_PASSWORD`: `<MQTT 密碼>`（如果需要）
5. 點擊 "Create Web Service"

### 5. DER + EV Orchestrator (app-der-ev-orchestrator)

1. "New +" → "Web Service"
2. 連接 GitHub 倉庫
3. 設置：
   - **Name**: `app-der-ev-orchestrator`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./services/apps/app-der-ev-orchestrator/Dockerfile`
   - **Docker Context**: `./services/apps/app-der-ev-orchestrator`
   - **Plan**: `Starter` 或 `Free`
4. 環境變數：
   - `FEEDER_ID`: `feeder-001`
   - `MQTT_BROKER`: `<你的 MQTT broker URL>`
   - `MQTT_PORT`: `1883` 或 `8883`
   - `MQTT_USERNAME`: `<MQTT 用戶名>`（如果需要）
   - `MQTT_PASSWORD`: `<MQTT 密碼>`（如果需要）
5. 點擊 "Create Web Service"

### 6. Rural Resilience Engine (app-rural-resilience)

1. "New +" → "Web Service"
2. 連接 GitHub 倉庫
3. 設置：
   - **Name**: `app-rural-resilience`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./services/apps/app-rural-resilience/Dockerfile`
   - **Docker Context**: `./services/apps/app-rural-resilience`
   - **Plan**: `Starter` 或 `Free`
4. 點擊 "Create Web Service"

### 7. Security Gateway (security-gateway)

1. "New +" → "Web Service"
2. 連接 GitHub 倉庫
3. 設置：
   - **Name**: `security-gateway`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./services/security-fabric/security-gateway/Dockerfile`
   - **Docker Context**: `./services/security-fabric/security-gateway`
   - **Plan**: `Starter` 或 `Free`
4. 環境變數：
   - `PORT`: `8443`
5. 點擊 "Create Web Service"

### 8. Telemetry Collector (telemetry-collector)

1. "New +" → "Web Service"
2. 連接 GitHub 倉庫
3. 設置：
   - **Name**: `telemetry-collector`
   - **Environment**: `Docker`
   - **Dockerfile Path**: `./services/security-fabric/telemetry-collector/Dockerfile`
   - **Docker Context**: `./services/security-fabric/telemetry-collector`
   - **Plan**: `Starter` 或 `Free`
4. 點擊 "Create Web Service"

## 重要步驟：更新前端環境變數

**在所有後端服務部署完成後**：

1. 在 Render Dashboard 中找到每個後端服務的實際 URL
2. 進入 `feeder-ide-frontend` 服務設置
3. 更新環境變數：
   - `VITE_API_BASE_URL`: `https://<實際的-feeder-ide-api-url>/api/v1`
   - `VITE_SIM_API_BASE_URL`: `https://<實際的-feeder-sim-engine-url>`
4. **保存並重新部署前端服務**（環境變數變更需要重新構建）

## 驗證部署

1. 訪問前端 URL
2. 打開瀏覽器開發者工具
3. 檢查控制台是否有錯誤
4. 測試功能：
   - 創建拓樸
   - 運行模擬
   - 執行安全測試

## 故障排除

如果遇到問題，請查看：
- Render 服務的構建日誌
- Render 服務的運行日誌
- 瀏覽器控制台的錯誤信息

