# Rural Resilience Engine

Predictive resilience analysis for rural feeders.

## Features

- Fault log ingestion
- Risk scoring (rule-based v0)
- Upgrade suggestions

## Development

### Prerequisites

- Python 3.11+

### Run locally

```bash
pip install -r requirements.txt
python -m uvicorn src.main:app --reload --port 8084
```

### API Endpoints

- `POST /api/v1/fault-logs` - Upload fault log
- `GET /api/v1/risk-scores` - Get risk scores
- `POST /api/v1/risk-scores` - Calculate risk scores with topology
- `GET /api/v1/upgrade-suggestions` - Get upgrade suggestions
- `POST /api/v1/upgrade-suggestions` - Generate suggestions with topology
- `GET /health` - Health check

