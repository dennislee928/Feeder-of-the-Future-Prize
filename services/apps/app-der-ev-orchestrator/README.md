# DER + EV Orchestrator

Domain app for orchestrating DER and EV charging to optimize feeder load.

## Features

- Asset registry (EV chargers, PV batteries)
- Control loop with heuristic-based load management
- MQTT integration for measurements and commands

## Development

### Prerequisites

- Python 3.11+
- MQTT broker

### Run locally

```bash
pip install -r requirements.txt
python -m uvicorn src.main:app --reload --port 8083
```

### Environment Variables

- `FEEDER_ID` - Feeder ID (default: feeder-001)
- `MQTT_BROKER` - MQTT broker address (default: localhost)
- `MQTT_PORT` - MQTT port (default: 1883)

### API Endpoints

- `POST /api/v1/assets/ev-chargers` - Register EV charger
- `POST /api/v1/assets/pv-batteries` - Register PV battery
- `GET /api/v1/assets` - List all assets
- `GET /health` - Health check

