# Feeder Simulation Engine

Simulation & Analysis API for feeder topologies.

## Features

- Powerflow simulation (stub)
- Reliability analysis (SAIDI/SAIFI)
- RESTful API with FastAPI

## Development

### Prerequisites

- Python 3.11+

### Run locally

```bash
pip install -r requirements.txt
python -m uvicorn src.app:app --reload --port 8081
```

Server will start on `http://localhost:8081`

### API Endpoints

- `POST /simulate/powerflow` - Run powerflow analysis
- `POST /simulate/reliability` - Run reliability analysis
- `GET /health` - Health check

### Example Request

```json
{
  "topology": {
    "nodes": [...],
    "lines": [...],
    "profile_type": "suburban"
  },
  "parameters": {
    "fault_rate": 0.1,
    "repair_time": 240
  }
}
```

