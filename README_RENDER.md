# Render 部署說明

本項目已配置好 Render 部署。按照以下步驟即可部署到 Render。

## 快速開始

### 1. 準備 MQTT Broker

**推薦使用 CloudMQTT（免費方案）**：
1. 訪問 [CloudMQTT](https://www.cloudmqtt.com/)
2. 註冊並創建免費實例
3. 記錄連接信息：
   - Broker URL
   - Port (通常是 1883 或 8883)
   - Username
   - Password

### 2. 部署到 Render

#### 方法 A: 使用 Blueprint（推薦）

1. 登入 [Render Dashboard](https://dashboard.render.com/)
2. 點擊 "New +" → "Blueprint"
3. 連接你的 GitHub 倉庫
4. Render 會自動檢測 `render.yaml`
5. 點擊 "Apply" 創建所有服務

#### 方法 B: 手動部署

參考 `deploy/RENDER_DEPLOYMENT.md` 中的詳細說明。

### 3. 設置環境變數

#### 步驟 1: 設置 MQTT 環境變數

在以下服務中設置 MQTT 相關環境變數：
- `feeder-os-controller`
- `app-der-ev-orchestrator`

環境變數：
- `MQTT_BROKER`: 你的 MQTT broker URL（例如: `m20.cloudmqtt.com`）
- `MQTT_PORT`: MQTT 端口（通常是 `1883` 或 `8883`）
- `MQTT_USERNAME`: MQTT 用戶名（如果需要）
- `MQTT_PASSWORD`: MQTT 密碼（如果需要）

#### 步驟 2: 設置前端環境變數（重要）

**必須在所有後端服務部署完成後進行**：

1. 在 Render Dashboard 中找到每個後端服務的 URL：
   - `feeder-ide-api` 的 URL（例如: `https://feeder-ide-api-xxx.onrender.com`）
   - `feeder-sim-engine` 的 URL（例如: `https://feeder-sim-engine-xxx.onrender.com`）

2. 在 `feeder-ide-frontend` 服務中設置環境變數：
   - `VITE_API_BASE_URL`: `https://feeder-ide-api-xxx.onrender.com/api/v1`
   - `VITE_SIM_API_BASE_URL`: `https://feeder-sim-engine-xxx.onrender.com`

3. **重新部署前端服務**（環境變數變更需要重新構建）

### 4. 驗證部署

1. 訪問前端 URL
2. 打開瀏覽器開發者工具檢查控制台
3. 測試功能：
   - 創建拓樸
   - 運行模擬
   - 執行安全測試

## 服務列表

部署後會創建以下服務：

| 服務名稱 | 類型 | 端口 | 說明 |
|---------|------|------|------|
| feeder-ide-frontend | Web | 80 | 前端應用 |
| feeder-ide-api | Web | 8080 | IDE API 後端 |
| feeder-sim-engine | Web | 8081 | 模擬引擎 |
| feeder-os-controller | Web | 8082 | Feeder OS 控制器 |
| app-der-ev-orchestrator | Web | 8083 | DER/EV 編排器 |
| app-rural-resilience | Web | 8084 | 農村韌性引擎 |
| security-gateway | Web | 8443 | 安全網關 |
| telemetry-collector | Web | 8085 | 遙測收集器 |

## 重要提示

### 前端環境變數
- `VITE_*` 環境變數需要在**構建時**設置
- 更新環境變數後必須**重新部署**前端服務
- 建議先部署所有後端服務，獲取 URL 後再設置前端環境變數

### MQTT Broker
- Render 的 Background Worker 不適合長期運行的 MQTT broker
- **強烈建議使用外部 MQTT 服務**（如 CloudMQTT）
- 如果必須使用 Render，可以部署 MQTT broker 作為 Background Worker，但不推薦用於生產環境

### 服務間通訊
- Render 會為每個服務分配不同的 URL
- 服務間通訊使用 HTTPS
- 確保 CORS 設置正確

## 故障排除

### 前端無法連接到後端
1. 檢查後端服務是否已成功部署
2. 檢查前端環境變數是否正確
3. 確認已重新部署前端（環境變數變更需要重新構建）
4. 檢查瀏覽器控制台的錯誤信息

### MQTT 連接失敗
1. 確認 MQTT broker URL 正確
2. 檢查端口是否正確（1883 或 8883）
3. 確認認證信息正確
4. 某些 MQTT 服務需要 TLS，確保使用正確的端口

### 服務啟動失敗
1. 查看 Render 的服務日誌
2. 檢查 Dockerfile 路徑是否正確
3. 確認環境變數設置正確
4. 檢查構建日誌中的錯誤信息

## 成本估算

### 免費方案
- 每個服務 750 小時/月
- 15 分鐘無活動後休眠
- 休眠後首次請求有延遲（冷啟動）

### 付費方案
- **Starter Plan**: $7/月/服務
- 適合生產環境使用
- 無休眠，性能更好

## 相關文檔

- `deploy/RENDER_DEPLOYMENT.md` - 詳細部署指南
- `deploy/RENDER_QUICK_START.md` - 快速開始指南
- `render.yaml` - Render Blueprint 配置

