# Feeder-of-the-Future Platform

> A software-defined distribution feeder platform with a digital twin IDE, edge **Feeder OS + App Runtime**, and domain apps for DER/EV orchestration, rural resilience, and cyber-physical security.

This project is a **code-first** exploration of next-generation distribution feeder design, inspired by the Feeder of the Future Prize but not limited to the competition.

Core idea:  
Treat a feeder like a **software platform**：

- You **design & simulate** feeders in a **Digital Twin & Design IDE**.
- You **deploy control logic** as apps on an edge **Feeder OS + App Runtime**.
- You plug in domain apps:
  - **Suburban / Urban DER + EV Orchestrator**
  - **Rural Predictive Resilience Engine**
- Everything sits on top of a **Cyber-Physical Security Fabric** that gives you observability + defense across OT/IT.

---

## High-Level Architecture

```mermaid
flowchart TD
  subgraph L0[Physical Layer]
    GridAssets[Lines / Switches / Transformers / DER / EV Chargers]
  end

  subgraph L1[Cyber-Physical Security Fabric]
    SecGateway[Security Gateway]
    SecCollector[Telemetry & IDS/IPS]
  end

  subgraph L2[Feeder OS + App Runtime]
    FeederOS[Feeder OS\nEdge Controller + App Runtime]
    AppStore[App Store & Lifecycle Manager]
  end

  subgraph L3[Domain Apps]
    AppDER[DER + EV Orchestrator]
    AppRural[Rural Predictive Resilience]
    AppCore[Core Apps\nVolt/VAR / FLISR / Metrics]
  end

  subgraph L4[Digital Twin & Design IDE]
    IDEUI[Web IDE (Topology Editor)]
    SimAPI[Simulation & Analysis API]
  end

  GridAssets <--> SecGateway
  SecGateway <--> FeederOS
  FeederOS <--> AppDER
  FeederOS <--> AppRural
  FeederOS <--> AppCore
  AppDER <--> SimAPI
  AppRural <--> SimAPI
  IDEUI <--> SimAPI
  IDEUI <--> AppStore

Main Components
1. Digital Twin & Design IDE

A web IDE + API backend to:

Draw / import feeder topology (nodes, lines, switches, DER, EV chargers).

Parameterize Rural / Suburban / Urban profiles.

Run power-flow / reliability simulations (stub initially, then integrate pandapower / GridLAB-D / OpenDSS).

Export configs for Feeder OS + Apps.

2. Feeder OS + App Runtime

An edge controller that runs at the feeder head-end or substation:

Minimal Linux + container runtime (k3s / containerd).

gRPC / MQTT control bus for apps.

App lifecycle:

install / upgrade / rollback

configuration via GitOps / API

Pluggable domain apps (e.g. DER orchestrator, rural resilience).

3. Suburban / Urban DER + EV Orchestrator

A domain app that:

Connects to EV chargers, PV inverters, home batteries (via protocol adaptors).

Runs feeder-level OPF / simple heuristics / later MPC/RL to:

Avoid transformer / line overload

Control voltage

Improve DER hosting capacity

Focuses on suburban / urban scenarios with dense BTM assets and EV charging.

4. Rural Feeder Predictive Resilience Engine

A domain app targeting rural feeders:

Ingests fault history, asset meta, weather / GIS data.

Runs reliability simulations (SAIDI/SAIFI, outage risk).

Proposes:

Where to add automated sectionalizers / reclosers

Alternative routing / parallel paths

Can round-trip with the IDE to suggest topology changes.

5. Cyber-Physical Security Fabric

A horizontal layer that:

Adds security gateways / proxies at OT-IT boundaries.

Collects telemetry (NetFlow, protocol logs, app logs).

Applies:

mTLS / zero-trust auth

Anomaly detection for feeder operations

Integrates conceptually with a central Unified Security & Infrastructure Platform (can reuse your existing project).

Status

This is an experimental, research-grade project:

✅ Planning / architecture docs

🚧 Initial scaffolding (services, Docker, minimal APIs)

⏳ Advanced simulation and control algorithms

Focus is on:

Clean modular architecture

Hackable codebase for experiments

Strong DevSecOps practices from day one

Status

This is an experimental, research-grade project:

✅ Planning / architecture docs

🚧 Initial scaffolding (services, Docker, minimal APIs)

⏳ Advanced simulation and control algorithms

Focus is on:

Clean modular architecture

Hackable codebase for experiments

Strong DevSecOps practices from day one

Quick Start (planned)

⚠️ Until code exists, this is aspirational. Adjust as implementation evolves.

# 1. Clone
git clone https://github.com/<your-username>/feeder-of-the-future-platform.git
cd feeder-of-the-future-platform

# 2. Start dev stack (backend + frontend + minimal edge runtime)
docker-compose up --build

# 3. Open IDE
# Visit http://localhost:3000 for the Digital Twin & Design IDE

Repository Layout (preview)

See project_structure.md
 for a more detailed description.

.
├── README.md
├── spec.md
├── implementation_steps.md
├── tech_stack.md
├── project_structure.md
├── services/
│   ├── feeder-ide-api/
│   ├── feeder-sim-engine/
│   ├── feeder-os-controller/
│   ├── apps/
│   │   ├── app-der-ev-orchestrator/
│   │   └── app-rural-resilience/
│   └── security-fabric/
│       ├── security-gateway/
│       └── telemetry-collector/
├── frontend/
│   └── ide-frontend/
└── deploy/
    ├── docker-compose.yml
    └── k8s/

License

TBD (MIT is recommended for maximum reuse).

Contributing

Right now this is a personal R&D playground.
If it ever opens to contributors:

Fork → create feature branch → PR

Keep changes small and composable


---

## `spec.md`

```markdown
# Feeder-of-the-Future Platform – System Specification

## 1. Goals & Non-Goals

### 1.1 Goals

- 提供一個 **軟體定義的配電 feeder 平台**，讓你可以：
  - 在 IDE 裡設計 / 模擬 Rural / Suburban / Urban feeders。
  - 把控制邏輯包成 app，部署到 Edge「Feeder OS」上。
  - 針對不同情境（EV/DER, rural reliability）安裝不同 domain apps。
- 內建 **Cyber-Physical Security** 概念，從設計 → 部署全程考慮安全性。
- 保持架構 **模組化**，方便快速實驗新的：
  - 控制策略（OPF, MPC, RL）
  - 資安偵測模型
  - 拓樸 / 裝置類型

### 1.2 Non-Goals

- 不打算一開始就達到 **utility-grade、可直接上線的 SCADA/DMS**。
- 不追求完整支援所有電力協定（IEC 61850, DNP3, etc.）— 先用 mock / simplified API。
- 不以競賽文件為主軸；**以 hackable code / architecture exploration 為優先**。

---

## 2. Functional Requirements

### 2.1 Digital Twin & Design IDE

**ID: IDE-01 – Topology Modeling**

- 使用者可以：
  - 在 web UI 上新增 / 編輯：
    - Bus / node
    - Line / cable
    - Transformer
    - Switch / breaker
    - DER (PV, battery)
    - EV chargers
  - 設定基本電氣參數（額定電壓、阻抗、額定容量等）。
- 後端以 JSON / YAML 格式儲存拓樸。

**ID: IDE-02 – Track Profiles**

- 系統提供三種 profile：
  - Rural
  - Suburban
  - Urban
- 每個 profile 包含：
  - 典型負載成分（住宅 / 輕工業 / 商業比例）
  - 典型 feeder 長度 / 節點數
  - 目標可靠度指標（SAIDI/SAIFI upper bound）
  - 典型 DER / EV 滲透率 range

**ID: IDE-03 – Simulation API**

- 提供 REST/gRPC API：
  - `/simulate/powerflow`
  - `/simulate/reliability`
  - `/simulate/scenario-run`
- 結果包含：
  - 節點電壓 / 線路載流率
  - 估計 SAIDI/SAIFI
  - 簡單 cost 指標（設備數量、估算 CapEx）

**ID: IDE-04 – Export to Feeder OS**

- 一個拓樸可以被「編譯」成：
  - Feeder OS config（topic 名稱、資產 ID、metrics 來源）
  - Domain app 的初始 config（例如 DER/EV app 的 charger list）。

---

### 2.2 Feeder OS + App Runtime

**ID: FOS-01 – App Lifecycle**

- 提供 API / CLI：
  - 安裝 app：`fos app install app-der-ev-orchestrator`
  - 更新 app：`fos app upgrade app-der-ev-orchestrator`
  - 停用 / 啟用 / 回滾 app
- App 應該打包為 container image（含版本標記）。

**ID: FOS-02 – Message Bus**

- 在 Feeder OS 上運行 message bus（MQTT 或 NATS 或 Kafka）。
- 所有 apps 使用統一主題規則，例：
  - `feeder/<id>/measurements/<asset-id>`
  - `feeder/<id>/commands/<asset-id>`
  - `feeder/<id>/events/<severity>`

**ID: FOS-03 – Config Management**

- Feeder OS 能從 Git / S3 / API 拉取 config bundles：
  - 拓樸描述
  - app-specific config
- 支援 versioned config（方便回滾）。

---

### 2.3 Suburban / Urban DER + EV Orchestrator

**ID: DER-01 – Asset Registry**

- 透過 REST / MQTT 註冊以下 assets：
  - EV chargers
  - Residential PV + battery
  - Commercial loads
- 儲存：
  - 額定功率
  - 可接受調度範圍（例如 EV 充電可 delay 到幾點）

**ID: DER-02 – Feeder-Level Control**

- 每個控制週期（例如每 5 分鐘）：
  - 從 measurements topic 取得 feeder 狀態。
  - 跑簡單 OPF / heuristic：
    - 優先限制過載 / 電壓違規。
    - 在可行範圍內平滑峰值負載。
  - 發布控制指令給 chargers / DER（功率 setpoint / charging window）。

**ID: DER-03 – Scenario Hooks**

- 可以在 offline 模式下，掛到 simulation engine：
  - 用相同控制邏輯跑在 digital twin 上。
  - 用於測試 / training（未來可以放 RL）。

---

### 2.4 Rural Feeder Predictive Resilience Engine

**ID: RUR-01 – Data Ingestion**

- 接收：
  - 歷史 fault log（時間、設備、故障類型、修復時間）。
  - 資產 meta（設備年齡、型號、安裝位置）。
  - weather / GIS（可先用 mock JSON）。
- 儲存在簡單的 relational DB 或時序 DB。

**ID: RUR-02 – Risk Scoring**

- 為每個 line / transformer / switch 算一個 risk score：
  - 依據故障頻率、距離、天氣暴露程度。
- 將高風險 asset 可視化。

**ID: RUR-03 – Upgrade Suggestions (v1 heuristic)**

- 初版可以用 rule-based：
  - 若某 section fault rate > threshold：
    - 建議：在兩側加 sectionalizer。
  - 若某 feeder section 長度 > threshold 且只有單一路徑：
    - 建議：評估 parallel line。
- 將建議以 JSON 返回 IDE，IDE 在拓樸圖上標出「建議新增點」。

---

### 2.5 Cyber-Physical Security Fabric

**ID: SEC-01 – Security Gateway**

- 作為 OT/IT 邊界的 proxy：
  - 統一 terminate TLS。
  - 做 mTLS 認證 Feeder OS / apps。
- Logging：
  - 所有控制命令與設定修改都被記錄並送往 telemetry collector。

**ID: SEC-02 – Telemetry Collector**

- 收集：
  - Feeder OS / app logs
  - Network metadata（例如各 topic / endpoint 行為）
- 初版可以：
  - 以 rule-based 偵測異常（例如超出正常頻率的開關操作）。
  - 將重要事件送到 external SIEM（可以指向你既有 security platform）。

---

## 3. Non-Functional Requirements

- **Modularity**
  - 每個 service / app 獨立 repo 子資料夾，有清楚 API 介面。
- **Security**
  - 預設 all internal comms with TLS / mTLS。
  - 容器鏡像需有 basic SBOM / image signing（可先 stub）。
- **Observability**
  - 基礎 metrics + structured logging（Feeder OS & apps）。
- **Deployability**
  - 可在 local Docker Compose 起 basic dev stack。
  - 未來支援 K8s（k3s）部署。

---

## 4. Future Extensions

- 真正接軌：
  - IEC 61850 / DNP3 透過 protocol gateway。
- 更進階的 control：
  - MPC / RL based feeder control。
- UI 進階：
  - 更完善的 feeder layout 視覺化與 scenario compare。