# Feeder OS Controller

Feeder OS + App Runtime controller service.

## Features

- App lifecycle management (install, enable, disable)
- MQTT message bus integration
- App contract definition

## Development

### Prerequisites

- Go 1.21+
- MQTT broker (e.g., Mosquitto)

### Run locally

```bash
# Start MQTT broker (using Docker)
docker run -it -p 1883:1883 eclipse-mosquitto

# Run the service
go run cmd/feeder-os/main.go
```

Server will start on `http://localhost:8082`

### Environment Variables

- `PORT` - Server port (default: 8082)
- `MQTT_BROKER` - MQTT broker address (default: localhost)
- `MQTT_CLIENT_ID` - MQTT client ID (default: feeder-os-controller)
- `MQTT_USERNAME` - MQTT username (optional)
- `MQTT_PASSWORD` - MQTT password (optional)
- `APPS_STORAGE_PATH` - Apps storage path (default: ./apps)

### API Endpoints

- `POST /api/v1/apps/install` - Install app
- `POST /api/v1/apps/enable` - Enable app
- `POST /api/v1/apps/disable` - Disable app
- `GET /api/v1/apps` - List all apps
- `GET /api/v1/apps/:id` - Get app by ID
- `GET /health` - Health check

### Topic Naming Convention

- `feeder/<feeder_id>/measurements/<asset_id>` - Measurement data
- `feeder/<feeder_id>/commands/<asset_id>` - Control commands
- `feeder/<feeder_id>/events/<severity>` - Event notifications

