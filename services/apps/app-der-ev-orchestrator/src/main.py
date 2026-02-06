"""
DER + EV Orchestrator
控制 DER 和 EV 充電樁以優化 feeder 負載
"""
import asyncio
import os
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from src.registry import AssetRegistry
from src.control_loop import ControlLoop
from src.mqtt_client import MQTTClient

app = FastAPI(
    title="DER + EV Orchestrator",
    description="Orchestrator for DER and EV charging control",
    version="0.1.0"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 初始化組件
registry = AssetRegistry()
mqtt_client = None
control_loop = None


def _sync_mqtt_connect(broker: str, port: int, feeder_id: str) -> MQTTClient:
    """同步 MQTT 連線，供 asyncio.to_thread 使用，避免阻塞 event loop"""
    import json
    import paho.mqtt.client as mqtt
    client = MQTTClient(broker, port, feeder_id)
    client.client = mqtt.Client(client_id=f"der-ev-orchestrator-{feeder_id}")

    def on_connect(c, u, f, rc):
        if rc == 0:
            print(f"MQTT connected to {broker}:{port}")
        else:
            print(f"MQTT connection failed: {rc}")

    def on_message(c, u, msg):
        topic = msg.topic
        if topic in client.subscriptions:
            try:
                payload = json.loads(msg.payload.decode())
                client.subscriptions[topic](topic, payload)
            except Exception as e:
                print(f"Error processing message: {e}")

    client.client.on_connect = on_connect
    client.client.on_message = on_message
    client.client.connect(broker, port, 60)
    client.client.loop_start()
    return client


@app.on_event("startup")
async def startup():
    """啟動時初始化 MQTT 和 control loop（非阻塞，避免 Render 部署時 No open ports）"""
    global mqtt_client, control_loop

    async def init_mqtt_and_loop():
        global mqtt_client, control_loop
        feeder_id = os.getenv("FEEDER_ID", "feeder-001")
        mqtt_broker = os.getenv("MQTT_BROKER", "mqtt")
        mqtt_port = int(os.getenv("MQTT_PORT", "1883"))

        max_retries = 10
        retry_delay = 2
        for i in range(max_retries):
            try:
                # 在執行緒中執行 MQTT 連線，避免阻塞 event loop 導致 Render 無法偵測 port
                mqtt_client = await asyncio.to_thread(_sync_mqtt_connect, mqtt_broker, mqtt_port, feeder_id)
                print(f"Successfully connected to MQTT broker at {mqtt_broker}:{mqtt_port}")
                break
            except Exception as e:
                if i < max_retries - 1:
                    print(f"Failed to connect to MQTT broker (attempt {i+1}/{max_retries}): {e}")
                    await asyncio.sleep(retry_delay)
                else:
                    print(f"Failed to connect to MQTT broker after {max_retries} attempts. Continuing without MQTT...")
                    mqtt_client = None

        if mqtt_client:
            control_loop = ControlLoop(registry, mqtt_client, feeder_id)
            asyncio.create_task(control_loop.run())
        else:
            print("Warning: Control loop not started due to MQTT connection failure")

    # 背景執行，不阻塞 startup 完成，讓 uvicorn 可立即 bind port
    asyncio.create_task(init_mqtt_and_loop())


@app.on_event("shutdown")
async def shutdown():
    """關閉時清理資源"""
    if mqtt_client:
        await mqtt_client.disconnect()


@app.get("/health")
async def health_check():
    """健康檢查"""
    return {"status": "ok"}


# 導入 API routes
from src.api import router, set_registry
set_registry(registry)
app.include_router(router, prefix="/api/v1", tags=["orchestrator"])


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8083)

