# Telemetry Collector

Collect telemetry and detect anomalies.

## Features

- Log ingestion
- Rule-based anomaly detection
- Alert generation

## Development

### Prerequisites

- Python 3.11+

### Run locally

```bash
pip install -r requirements.txt
python -m uvicorn src.main:app --reload --port 8085
```

### API Endpoints

- `POST /api/v1/logs` - Ingest log
- `GET /api/v1/logs` - Get logs
- `GET /api/v1/anomalies/summary` - Get anomaly summary
- `GET /health` - Health check

