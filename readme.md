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