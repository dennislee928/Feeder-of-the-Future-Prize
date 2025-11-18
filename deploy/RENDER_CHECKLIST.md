# Render 部署檢查清單

## ✅ 已完成的修復

### 1. Dockerfile 修復
- ✅ `feeder-sim-engine/Dockerfile` - 使用 PORT 環境變數
- ✅ `app-der-ev-orchestrator/Dockerfile` - 使用 PORT 環境變數
- ✅ `app-rural-resilience/Dockerfile` - 使用 PORT 環境變數
- ✅ `telemetry-collector/Dockerfile` - 使用 PORT 環境變數
- ✅ 所有服務都創建了 `start.sh` 啟動腳本

### 2. 配置文件
- ✅ `render.yaml` - Render Blueprint 配置
- ✅ `README_RENDER.md` - 部署說明
- ✅ `deploy/RENDER_DEPLOYMENT.md` - 詳細指南
- ✅ `deploy/RENDER_MANUAL_SETUP.md` - 手動設置指南

## 📋 部署前檢查

### 必須準備
- [ ] Render 帳號
- [ ] GitHub 倉庫已連接
- [ ] 外部 MQTT Broker（推薦 CloudMQTT）

### 部署步驟

1. **設置 MQTT Broker**
   - [ ] 註冊 CloudMQTT 或類似服務
   - [ ] 記錄 Broker URL、Port、Username、Password

2. **部署後端服務**（按順序）
   - [ ] `feeder-ide-api`
   - [ ] `feeder-sim-engine`
   - [ ] `feeder-os-controller`（設置 MQTT 環境變數）
   - [ ] `app-der-ev-orchestrator`（設置 MQTT 環境變數）
   - [ ] `app-rural-resilience`
   - [ ] `security-gateway`
   - [ ] `telemetry-collector`

3. **記錄服務 URL**
   - [ ] `feeder-ide-api` URL: ________________
   - [ ] `feeder-sim-engine` URL: ________________

4. **部署前端服務**
   - [ ] 設置 `VITE_API_BASE_URL`
   - [ ] 設置 `VITE_SIM_API_BASE_URL`
   - [ ] 重新部署前端

5. **驗證**
   - [ ] 訪問前端 URL
   - [ ] 檢查瀏覽器控制台
   - [ ] 測試創建拓樸
   - [ ] 測試運行模擬
   - [ ] 測試安全測試模式

## 🔧 如果 render.yaml 無法使用

如果 Render Blueprint 無法直接導入，請：
1. 使用 `deploy/RENDER_MANUAL_SETUP.md` 手動創建每個服務
2. 按照檢查清單逐步部署
3. 確保每個服務的環境變數正確設置

## ⚠️ 常見問題

### 問題：服務無法啟動
**解決**：檢查 Render 日誌，確認 PORT 環境變數已設置

### 問題：前端無法連接後端
**解決**：
1. 確認後端服務已部署
2. 檢查前端環境變數中的 URL 是否正確
3. 確認已重新部署前端（環境變數變更需要重新構建）

### 問題：MQTT 連接失敗
**解決**：
1. 確認 MQTT broker URL 正確
2. 檢查端口（1883 或 8883）
3. 確認認證信息正確

