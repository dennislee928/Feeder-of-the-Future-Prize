# Tech Stack

> 盡量選你已經熟悉、好維護、好擴充的技術。

---

## Languages

- **Go**
  - Feeder OS controller
  - 高性能、易部署的 edge services
- **Python**
  - Simulation engine（未來接 pandapower / GridLAB-D / OpenDSS）
  - ML / data-driven modules（Rural risk, IDS 等）
- **TypeScript**
  - Web IDE frontend (React)

---

## Backend Frameworks

### Go

- HTTP/gRPC:
  - `chi` or `gorilla/mux` for REST
  - `grpc-go` for microservice comms（可後期導入）
- Messaging:
  - MQTT（e.g. Mosquitto as broker）or
  - NATS (輕量 pub/sub)

### Python

- Web:
  - **FastAPI**（型別支援佳，適合 simulation API / small services）
- Data & ML (future):
  - `pandas`, `numpy`, `scikit-learn` (for quick prototypes)

---

## Frontend

- **React + TypeScript + Vite**
- UI / Graph libs:
  - `React Flow` 或 `Cytoscape.js` 用來畫拓樸
- State management:
  - `Zustand` / `Redux Toolkit`（視喜好）

---

## Data Storage

- Config / Topology:
  - 簡單開始：`PostgreSQL` + `JSONB` 或純檔案（YAML/JSON）。
- Telemetry / Logs:
  - Phase 1–2：直接寫檔 / Postgres。
  - Future：`Prometheus` + `Loki` / `OpenSearch`。
- Time-series (future):
  - `TimescaleDB` 或 `InfluxDB`（如需更進階 analysis）。

---

## Messaging / Streaming

- 初版選一種：
  - **MQTT**（架構簡單，edge 友好）
    - Broker: `Eclipse Mosquitto`
  - or **NATS**（更通用、cloud-native）。

選擇準則：

- 想貼近 IoT / edge 生態 → MQTT。
- 想更 cloud/microservices → NATS。

---

## Deployment

- **Docker / Docker Compose**
  - 開發環境，一鍵載起所有 services。
- **k3s / Kubernetes (future)**
  - 模擬「Feeder OS」跑在 edge k3s node 上。
  - Cluster 模式可以放在 lab / cloud。

---

## Security & DevSecOps

- TLS/mTLS:
  - `step-ca` 或 `cfssl` 建簡單 internal CA。
- SAST / Dependency scanning（可逐步導入）：
  - Go: `gosec`
  - Python: `bandit`, `pip-audit`
  - JS: `npm audit`, `depcheck`
- Container:
  - `Trivy` for image scanning。
- CI/CD:
  - GitHub Actions：
    - Lint + Unit test。
    - Build images。
    - Basic security scan（Trivy / SAST）stub。

---

## Simulation & Power Systems (future)

- **pandapower**（Python）
- **GridLAB-D**（需另裝，視需求）
- **OpenDSS**（透過 Python bridge）

初期可以先 stub 掉。待架構穩定後再慢慢接真實電力庫。
