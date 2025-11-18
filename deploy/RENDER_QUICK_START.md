# Render 快速部署指南

## 前置要求

1. Render 帳號（[註冊](https://render.com/)）
2. GitHub 倉庫已連接
3. 外部 MQTT Broker（推薦使用 [CloudMQTT](https://www.cloudmqtt.com/) 免費方案）

## 快速步驟

### 1. 設置 MQTT Broker

**推薦：CloudMQTT**
1. 註冊並創建免費實例
2. 記錄以下信息：
   - Broker URL (例如: `m20.cloudmqtt.com`)
   - Port (通常是 `1883` 或 `8883` for TLS)
   - Username
   - Password

### 2. 部署到 Render

#### 選項 A: 使用 Blueprint (推薦)

1. 在 Render Dashboard 點擊 "New +" → "Blueprint"
2. 連接你的 GitHub 倉庫
3. Render 會自動檢測 `render.yaml`
4. 點擊 "Apply" 創建所有服務

#### 選項 B: 手動部署

按照 `RENDER_DEPLOYMENT.md` 中的說明手動創建每個服務。

### 3. 設置環境變數

#### 在所有需要 MQTT 的服務中設置：
- `MQTT_BROKER_URL`: 你的 MQTT broker URL
- `MQTT_USERNAME`: MQTT 用戶名（如果需要）
- `MQTT_PASSWORD`: MQTT 密碼（如果需要）

#### 在前端服務中設置（**重要**）：
部署完所有後端服務後，獲取它們的 URL，然後在前端服務中設置：

1. 找到 `feeder-ide-api` 服務的 URL（例如: `https://feeder-ide-api-xxx.onrender.com`）
2. 找到 `feeder-sim-engine` 服務的 URL（例如: `https://feeder-sim-engine-xxx.onrender.com`）
3. 在前端服務的環境變數中設置：
   - `VITE_API_BASE_URL`: `https://feeder-ide-api-xxx.onrender.com/api/v1`
   - `VITE_SIM_API_BASE_URL`: `https://feeder-sim-engine-xxx.onrender.com`
4. **重新部署前端服務**（環境變數變更需要重新構建）

### 4. 驗證部署

1. 訪問前端 URL
2. 打開瀏覽器開發者工具檢查是否有錯誤
3. 測試功能：
   - 創建拓樸
   - 運行模擬
   - 執行安全測試

## 常見問題

### Q: 前端無法連接到後端
**A:** 檢查：
1. 後端服務是否已成功部署
2. 前端環境變數是否正確設置
3. 是否已重新部署前端（環境變數變更需要重新構建）

### Q: MQTT 連接失敗
**A:** 檢查：
1. MQTT broker URL 是否正確
2. 端口是否正確（1883 或 8883）
3. 認證信息是否正確
4. 某些 MQTT 服務需要 TLS，確保使用正確的端口

### Q: 服務啟動失敗
**A:** 檢查：
1. Render 服務日誌
2. Dockerfile 路徑是否正確
3. 環境變數是否設置正確

## 成本優化

### 免費方案限制
- 每個服務 750 小時/月
- 15 分鐘無活動後休眠
- 休眠後首次請求有延遲

### 建議
- **開發/測試**: 使用免費方案
- **生產環境**: 升級到 Starter Plan ($7/月/服務)

## 下一步

部署完成後，可以：
1. 設置自定義域名
2. 配置自動部署（Git push 觸發）
3. 設置監控和告警
4. 配置 SSL 證書（Render 自動提供）

