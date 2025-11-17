# Project Structure

> 初版 monorepo，方便跨 service 重用程式碼與 CI/CD。

```text
.
├── README.md
├── spec.md
├── implementation_steps.md
├── tech_stack.md
├── project_structure.md
│
├── services/
│   ├── feeder-ide-api/             # Backend for the Digital Twin & Design IDE
│   │   ├── cmd/
│   │   │   └── ide-api/           # main() entrypoint
│   │   ├── internal/
│   │   │   ├── topology/          # topology models & CRUD logic
│   │   │   ├── profiles/          # rural/suburban/urban profiles
│   │   │   └── simclient/         # client to feeder-sim-engine
│   │   ├── api/                   # HTTP/gRPC handlers
│   │   └── Dockerfile
│   │
│   ├── feeder-sim-engine/         # Simulation & analysis service
│   │   ├── src/
│   │   │   ├── app.py             # FastAPI entry
│   │   │   ├── powerflow_stub.py  # powerflow stubs / wrappers
│   │   │   └── reliability_stub.py# SAIDI/SAIFI stubs
│   │   └── Dockerfile
│   │
│   ├── feeder-os-controller/      # Feeder OS + App Runtime
│   │   ├── cmd/
│   │   │   └── feeder-os/
│   │   ├── internal/
│   │   │   ├── apps/              # app lifecycle mgmt
│   │   │   ├── bus/               # MQTT/NATS abstraction
│   │   │   └── config/            # config loading/versioning
│   │   ├── api/
│   │   └── Dockerfile
│   │
│   ├── apps/                      # Domain apps living on Feeder OS
│   │   ├── app-der-ev-orchestrator/
│   │   │   ├── src/
│   │   │   │   ├── main.py        # main loop
│   │   │   │   ├── registry.py    # asset registry API
│   │   │   │   └── control_loop.py# simple OPF/heuristic
│   │   │   └── Dockerfile
│   │   │
│   │   └── app-rural-resilience/
│   │       ├── src/
│   │       │   ├── main.py
│   │       │   ├── ingestion.py   # fault/history ingestion
│   │       │   └── risk_model.py  # risk scoring + suggestions
│   │       └── Dockerfile
│   │
│   └── security-fabric/           # Cyber-Physical Security Fabric
│       ├── security-gateway/
│       │   ├── cmd/
│       │   │   └── sec-gateway/
│       │   ├── internal/
│       │   │   ├── proxy/         # reverse proxy logic
│       │   │   └── mTLS/          # certificate handling
│       │   └── Dockerfile
│       │
│       └── telemetry-collector/
│           ├── src/
│           │   ├── main.py
│           │   └── rules.py       # simple anomaly rules
│           └── Dockerfile
│
├── frontend/
│   └── ide-frontend/
│       ├── src/
│       │   ├── App.tsx
│       │   ├── components/
│       │   │   ├── TopologyCanvas.tsx
│       │   │   ├── Palette.tsx
│       │   │   └── PropertiesPanel.tsx
│       │   └── api/
│       │       └── ideApi.ts      # typed client for feeder-ide-api
│       ├── public/
│       └── package.json
│
├── deploy/
│   ├── docker-compose.yml         # dev stack: IDE + sim + feeder OS + apps + security
│   └── k8s/
│       ├── feeder-ide-api.yaml
│       ├── feeder-sim-engine.yaml
│       ├── feeder-os-controller.yaml
│       ├── app-der-ev-orchestrator.yaml
│       ├── app-rural-resilience.yaml
│       ├── security-gateway.yaml
│       └── telemetry-collector.yaml
│
└── .github/
    └── workflows/
        ├── ci-backend.yml         # lint/test/build for services/*
        ├── ci-frontend.yml        # lint/build for ide-frontend
        └── security-scan.yml      # optional Trivy / SAST stubs
