# Dummy App

示例 app，展示如何與 Feeder OS 溝通。

## Features

- 訂閱 measurements topic
- Log 接收到的資料

## Configuration

透過環境變數配置：

- `MQTT_BROKER` - MQTT broker address (default: localhost)
- `MQTT_PORT` - MQTT port (default: 1883)
- `FEEDER_ID` - Feeder ID (default: feeder-001)
- `CLIENT_ID` - MQTT client ID (default: dummy-app)

## Run locally

```bash
pip install -r requirements.txt
python src/main.py
```

## Docker

```bash
docker build -t dummy-app .
docker run -e MQTT_BROKER=localhost -e FEEDER_ID=feeder-001 dummy-app
```

