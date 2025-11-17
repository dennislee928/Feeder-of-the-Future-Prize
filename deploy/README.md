# Deployment Guide

## Prerequisites

1. **Docker Desktop** 必須正在運行
   - Windows: 啟動 Docker Desktop 應用程式
   - 確認 Docker daemon 正在運行：`docker ps` 應該能正常執行

2. **檢查 Port 衝突**
   - 如果本地有 MQTT broker 運行在 port 1883，請先停止它
   - 或者修改 `docker-compose.yml` 中的 port 映射

3. 確保所有服務的 Dockerfile 都存在

## Quick Start

```bash
# 從專案根目錄
cd deploy

# 啟動所有服務
docker-compose up --build

# 或在背景執行
docker-compose up -d --build
```

## Troubleshooting

### Port 1883 已被佔用

如果看到錯誤：
```
Bind for 0.0.0.0:1883 failed: port is already allocated
```

**解決方案**：
1. 停止本地 MQTT broker（如果有的話）
2. 或修改 `docker-compose.yml` 中的 MQTT port 映射為其他 port（例如 `1884:1883`）

### Docker Desktop 未運行

如果看到錯誤：
```
error during connect: open //./pipe/dockerDesktopLinuxEngine: The system cannot find the file specified.
```

**解決方案**：
1. 啟動 Docker Desktop 應用程式
2. 等待 Docker Desktop 完全啟動（圖示不再閃爍）
3. 再次執行 `docker ps` 確認連接正常
4. 然後執行 `docker-compose up`

### 容器無法連接到 MQTT

如果看到錯誤：
```
socket.gaierror: [Errno -5] No address associated with hostname
```

**解決方案**：
1. 確保 MQTT 容器已成功啟動：`docker-compose ps`
2. 檢查 MQTT 容器的健康狀態
3. 確保所有服務都在同一個 network（feeder-network）

### 檢查服務狀態

```bash
# 查看所有容器狀態
docker-compose ps

# 查看特定服務的日誌
docker-compose logs feeder-ide-api
docker-compose logs mqtt

# 查看所有服務的日誌
docker-compose logs -f

# 查看特定服務的錯誤
docker-compose logs app-der-ev-orchestrator | grep -i error
```

### 停止服務

```bash
# 停止所有服務
docker-compose down

# 停止並移除 volumes
docker-compose down -v

# 停止並移除所有相關資源
docker-compose down --remove-orphans
```

## Services

啟動後，以下服務將可用：

- **IDE Frontend**: http://localhost:3000
- **Feeder IDE API**: http://localhost:8080
- **Simulation Engine**: http://localhost:8081
- **Feeder OS Controller**: http://localhost:8082
- **DER + EV Orchestrator**: http://localhost:8083
- **Rural Resilience Engine**: http://localhost:8084
- **Security Gateway**: http://localhost:8443
- **Telemetry Collector**: http://localhost:8085
- **MQTT Broker**: localhost:1884 (外部訪問) 或 mqtt:1883 (容器內部)

## Network Notes

- 容器之間使用服務名稱（如 `mqtt`）進行通信
- 外部訪問使用 `localhost` 和映射的 port
- 所有服務都在 `feeder-network` 網路中
