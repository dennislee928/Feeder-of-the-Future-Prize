# Feeder IDE API

Backend API for the Digital Twin & Design IDE.

## Features

- Topology CRUD operations
- Profile management (Rural/Suburban/Urban)
- RESTful API with Gin framework

## Development

### Prerequisites

- Go 1.21+

### Run locally

```bash
go run cmd/ide-api/main.go
```

Server will start on `http://localhost:8080`

### API Endpoints

- `POST /api/v1/topologies` - Create topology
- `GET /api/v1/topologies/:id` - Get topology
- `PUT /api/v1/topologies/:id` - Update topology
- `DELETE /api/v1/topologies/:id` - Delete topology
- `GET /api/v1/topologies` - List all topologies
- `GET /api/v1/profiles` - List all profiles
- `GET /api/v1/profiles/:type` - Get profile by type
- `GET /health` - Health check

