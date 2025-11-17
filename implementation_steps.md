# Implementation Steps

> 目標：先有可跑的骨架，再慢慢補電力/AI/資安細節。

---

## Phase 0 – Repo Bootstrap

1. 建新 repo：`feeder-of-the-future-platform`
2. 建立基本目錄結構（見 `project_structure.md`）。
3. 加入：
   - `.editorconfig`
   - `.gitignore`
   - 基本 CI（lint / build）stub（GitHub Actions）。

---

## Phase 1 – Minimal Digital Twin & IDE

### Backend (feeder-ide-api)

- [ ] 初始化專案（Python + FastAPI 或 Go + chi/grpc）。
- [ ] 建立基本 API：
  - [ ] `POST /topologies` – 新增拓樸。
  - [ ] `GET /topologies/{id}` – 取得拓樸。
  - [ ] `PUT /topologies/{id}` – 更新拓樸。
- [ ] Data model 設計：
  - [ ] Node / Line / Transformer / Switch / DER / EV Charger。
  - [ ] Track profile metadata。

### Frontend (ide-frontend)

- [ ] 初始化 React + TypeScript + Vite。
- [ ] 建 UI 架構：
  - [ ] 左側：資產 Palette。
  - [ ] 中央：拓樸 canvas（先用簡單 SVG / React Flow）。
  - [ ] 右側：屬性編輯面板。
- [ ] 串接 backend：
  - [ ] 儲存 / 載入拓樸。

---

## Phase 2 – Simulation Stub

### Sim Engine (feeder-sim-engine)

- [ ] 建立 service（Python 推薦，方便之後接電力套件）。
- [ ] 定義 API：
  - [ ] `POST /simulate/powerflow` – 接收拓樸 JSON，回傳隨機但合理的 voltage / loading。
  - [ ] `POST /simulate/reliability` – 接收拓樸 JSON + 假設參數，回傳 dummy SAIDI/SAIFI。
- [ ] 未接真實電力庫前，可用：
  - [ ] 簡單 algebra + heuristics 生成 deterministic fake data。

### IDE 整合

- [ ] 在 IDE UI 加上「Run Simulation」按鈕。
- [ ] 顯示簡易結果（例如節點顏色代表電壓偏差）。

---

## Phase 3 – Feeder OS + App Runtime Skeleton

### Feeder OS Controller (feeder-os-controller)

- [ ] 使用 Go 建立 service。
- [ ] 提供 REST/gRPC API：
  - [ ] `POST /apps/install`
  - [ ] `POST /apps/enable`
  - [ ] `POST /apps/disable`
  - [ ] `GET /apps`
- [ ] 整合 message bus（MQTT 或 NATS）：
  - [ ] 啟動 broker（可內嵌或 docker 另一 container）。
  - [ ] 提供簡單 publish/subscribe helper。

### App Contract

- [ ] 定義 app 與 Feeder OS 溝通介面：
  - [ ] Config 方式（環境變數 / config volume / gRPC）。
  - [ ] topics 命名規則。
- [ ] 建立「dummy app」示例：
  - [ ] 訂閱 measurements topic，log 資料即可。

---

## Phase 4 – DER + EV Orchestrator v0

### App Skeleton

- [ ] 建立 `app-der-ev-orchestrator` service（Python / Go 擇一）。
- [ ] 支援：
  - [ ] Asset registry：簡單 REST：
    - `POST /assets/ev-chargers`
    - `POST /assets/pv-batteries`
  - [ ] 控制 loop（定時 job）：
    - [ ] 從 message bus 讀取 measurements。
    - [ ] 執行簡單 heuristic：
      - 若 feeder loading > threshold → 降低部分 EV 充電功率。
    - [ ] 發 publish 控制命令。

### IDE / Sim 整合

- [ ] Offline 模式：
  - [ ] 用 simulation engine 假資料餵給 orchestrator。
  - [ ] 確認控制邏輯合理。

---

## Phase 5 – Rural Resilience Engine v0

### App Skeleton

- [ ] 建立 `app-rural-resilience` service。
- [ ] API：
  - [ ] `POST /fault-logs` – 上傳故障紀錄（可先用 JSON）。
  - [ ] `GET /risk-scores` – 回傳每個 asset 的 risk score。
  - [ ] `GET /upgrade-suggestions` – 回傳建議（新增裝置位置）。
- [ ] Algorithm v0：
  - [ ] 基於簡單 rule：
    - 故障次數多 → 高風險。
    - 遠離變電站且單路徑 → 高風險。

### IDE 整合

- [ ] IDE 取得 `upgrade-suggestions`。
- [ ] 在拓樸圖上用小 icon 標示「建議新增開關」。

---

## Phase 6 – Security Fabric v0

### Security Gateway

- [ ] 建立 `security-gateway` service：
  - [ ] 作為反向 proxy / API gateway。
  - [ ] 統一 terminate TLS。
- [ ] 實作：
  - [ ] 簡單 mTLS 驗證 app / feeder OS（可先用自簽 CA）。
  - [ ] Log 所有敏感 API 呼叫。

### Telemetry Collector

- [ ] 建立 `telemetry-collector` service：
  - [ ] 收 app / feeder OS logs（HTTP / gRPC / syslog 皆可）。
  - [ ] 存入 TS DB 或簡單文件 DB。
  - [ ] 針對:
    - 非預期大量控制命令
    - 深夜大量 config 變更
    - 發出 basic alert（先 log / print）。

---

## Phase 7 – DevSecOps & Polish

- [ ] Dockerfile for each service。
- [ ] `docker-compose.yml` 一鍵起 local stack：
  - IDE + Sim + Feeder OS + dummy apps + security gateway。
- [ ] GitHub Actions：
  - [ ] Lint（golangci-lint, flake8, eslint 等）。
  - [ ] Build images。
  - [ ] Basic unit tests。
- [ ] README 更新：
  - [ ] 加實際啟動步驟與 demo 流程。

---

## Suggested Implementation Order (Minimal Path)

1. Phase 0–1：有 IDE + backend，可以畫圖 / 存拓樸。
2. Phase 2：加 simulation stub，有基本「按下去 → 出結果」感覺。
3. Phase 3：Feeder OS skeleton + message bus + dummy app。
4. Phase 4：DER/EV app v0，形成第一個完整閉環（測量 → 控制）。
5. Phase 5–6：Rural + Security 作為進階模組慢慢補。
