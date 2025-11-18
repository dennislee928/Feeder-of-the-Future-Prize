# Render 部署指南

本指南說明如何將 Feeder of the Future Platform 部署到 Render。

## 重要注意事項

### MQTT Broker
Render 的 Background Worker 不適合長期運行的 MQTT broker。建議使用外部 MQTT 服務：
- **CloudMQTT** (免費方案可用)
- **HiveMQ Cloud** (免費方案可用)
- **AWS IoT Core**
- **EMQX Cloud**

### 服務間通訊
Render 會為每個服務分配不同的 URL。需要更新環境變數以指向正確的服務 URL。

## 部署步驟

### 1. 準備 MQTT Broker

#### 選項 A: 使用 CloudMQTT (推薦)
1. 註冊 [CloudMQTT](https://www.cloudmqtt.com/)
2. 創建一個免費實例
3. 獲取連接信息：
   - Broker URL
   - Port
   - Username
   - Password

#### 選項 B: 使用 Render Background Worker (不推薦，僅用於測試)
如果必須使用 Render，可以部署 MQTT broker 作為 Background Worker，但這不適合生產環境。

### 2. 在 Render 創建服務

#### 方法 A: 使用 render.yaml (推薦)
1. 將此倉庫連接到 Render
2. 在 Render Dashboard 中選擇 "New Blueprint"
3. 選擇此倉庫並導入 `render.yaml`
4. Render 會自動創建所有服務

#### 方法 B: 手動創建每個服務
為每個服務創建獨立的 Web Service：

1. **Frontend**
   - Type: Web Service
   - Environment: Docker
   - Dockerfile Path: `./frontend/ide-frontend/Dockerfile`
   - Docker Context: `./frontend/ide-frontend`
   - Environment Variables:
     - `VITE_API_BASE_URL`: `https://feeder-ide-api.onrender.com/api/v1`
     - `VITE_SIM_API_BASE_URL`: `https://feeder-sim-engine.onrender.com`

2. **Feeder IDE API**
   - Type: Web Service
   - Environment: Docker
   - Dockerfile Path: `./services/feeder-ide-api/Dockerfile`
   - Docker Context: `./services/feeder-ide-api`
   - Environment Variables:
     - `PORT`: `8080`

3. **Feeder Sim Engine**
   - Type: Web Service
   - Environment: Docker
   - Dockerfile Path: `./services/feeder-sim-engine/Dockerfile`
   - Docker Context: `./services/feeder-sim-engine`

4. **Feeder OS Controller**
   - Type: Web Service
   - Environment: Docker
   - Dockerfile Path: `./services/feeder-os-controller/Dockerfile`
   - Docker Context: `./services/feeder-os-controller`
   - Environment Variables:
     - `PORT`: `8082`
     - `MQTT_BROKER`: `<你的 MQTT broker URL>`
     - `MQTT_PORT`: `1883` (或你的 MQTT broker 端口)
     - `MQTT_CLIENT_ID`: `feeder-os-controller`

5. **DER + EV Orchestrator**
   - Type: Web Service
   - Environment: Docker
   - Dockerfile Path: `./services/apps/app-der-ev-orchestrator/Dockerfile`
   - Docker Context: `./services/apps/app-der-ev-orchestrator`
   - Environment Variables:
     - `FEEDER_ID`: `feeder-001`
     - `MQTT_BROKER`: `<你的 MQTT broker URL>`
     - `MQTT_PORT`: `1883`

6. **Rural Resilience Engine**
   - Type: Web Service
   - Environment: Docker
   - Dockerfile Path: `./services/apps/app-rural-resilience/Dockerfile`
   - Docker Context: `./services/apps/app-rural-resilience`

7. **Security Gateway**
   - Type: Web Service
   - Environment: Docker
   - Dockerfile Path: `./services/security-fabric/security-gateway/Dockerfile`
   - Docker Context: `./services/security-fabric/security-gateway`
   - Environment Variables:
     - `PORT`: `8443`

8. **Telemetry Collector**
   - Type: Web Service
   - Environment: Docker
   - Dockerfile Path: `./services/security-fabric/telemetry-collector/Dockerfile`
   - Docker Context: `./services/security-fabric/telemetry-collector`

### 3. 設置環境變數

在 Render Dashboard 中，為每個服務設置以下環境變數：

#### 共享環境變數（如果使用外部 MQTT）
- `MQTT_BROKER_URL`: 你的 MQTT broker URL (例如: `m20.cloudmqtt.com`)
- `MQTT_USERNAME`: MQTT 用戶名
- `MQTT_PASSWORD`: MQTT 密碼

#### 前端環境變數
部署完所有後端服務後，更新前端的環境變數：
- `VITE_API_BASE_URL`: `https://<你的-feeder-ide-api-url>/api/v1`
- `VITE_SIM_API_BASE_URL`: `https://<你的-feeder-sim-engine-url>`

### 4. 更新服務 URL

部署完成後，需要更新前端的環境變數以指向實際的服務 URL：

1. 在 Render Dashboard 中找到每個服務的 URL
2. 更新前端的環境變數：
   - `VITE_API_BASE_URL`: 指向 `feeder-ide-api` 的 URL
   - `VITE_SIM_API_BASE_URL`: 指向 `feeder-sim-engine` 的 URL
3. 重新部署前端服務

### 5. 驗證部署

1. 訪問前端 URL
2. 檢查瀏覽器控制台是否有錯誤
3. 測試各個功能：
   - 創建拓樸
   - 運行模擬
   - 執行安全測試

## 成本估算

Render 的免費方案限制：
- 每個服務 750 小時/月
- 服務在 15 分鐘無活動後會休眠
- 休眠後首次請求會有延遲（冷啟動）

建議的服務配置：
- **Starter Plan** ($7/月/服務): 適合生產環境
- **Free Plan**: 適合開發和測試

## 故障排除

### 服務無法啟動
- 檢查 Dockerfile 路徑是否正確
- 檢查環境變數是否設置正確
- 查看 Render 的服務日誌

### 前端無法連接到後端
- 確認後端服務已成功部署
- 檢查 CORS 設置
- 確認環境變數中的 URL 正確

### MQTT 連接失敗
- 確認 MQTT broker URL 正確
- 檢查 MQTT broker 的認證信息
- 確認端口是否正確（某些 MQTT 服務使用 8883 for TLS）

## 替代方案

如果 Render 不適合你的需求，可以考慮：
- **Railway**: 支持 docker-compose
- **Fly.io**: 支持多區域部署
- **AWS/GCP/Azure**: 完整的雲端平台
- **DigitalOcean App Platform**: 類似 Render 但支持更多功能

