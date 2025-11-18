#!/bin/bash
# Render 部署設置腳本
# 此腳本用於在 Render 部署後設置環境變數

echo "Render 部署設置腳本"
echo "===================="
echo ""
echo "請按照以下步驟設置 Render 服務："
echo ""
echo "1. 在 Render Dashboard 中創建所有服務（使用 render.yaml 或手動創建）"
echo ""
echo "2. 設置共享環境變數（在每個需要 MQTT 的服務中）："
echo "   - MQTT_BROKER_URL: 你的 MQTT broker URL"
echo "   - MQTT_USERNAME: MQTT 用戶名（如果需要）"
echo "   - MQTT_PASSWORD: MQTT 密碼（如果需要）"
echo ""
echo "3. 部署所有後端服務後，獲取它們的 URL："
echo "   - feeder-ide-api URL"
echo "   - feeder-sim-engine URL"
echo ""
echo "4. 更新前端服務的環境變數："
echo "   - VITE_API_BASE_URL: https://<feeder-ide-api-url>/api/v1"
echo "   - VITE_SIM_API_BASE_URL: https://<feeder-sim-engine-url>"
echo ""
echo "5. 重新部署前端服務以應用新的環境變數"
echo ""
echo "注意：前端環境變數需要在構建時設置，所以更新後需要重新構建。"

