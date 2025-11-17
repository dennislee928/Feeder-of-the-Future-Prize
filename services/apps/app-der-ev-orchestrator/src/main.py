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


@app.on_event("startup")
async def startup():
    """啟動時初始化 MQTT 和 control loop"""
    global mqtt_client, control_loop
    
    feeder_id = os.getenv("FEEDER_ID", "feeder-001")
    mqtt_broker = os.getenv("MQTT_BROKER", "mqtt")  # 使用服務名稱
    mqtt_port = int(os.getenv("MQTT_PORT", "1883"))
    
    # 重試連接 MQTT
    max_retries = 10
    retry_delay = 2
    for i in range(max_retries):
        try:
            mqtt_client = MQTTClient(mqtt_broker, mqtt_port, feeder_id)
            await mqtt_client.connect()
            print(f"Successfully connected to MQTT broker at {mqtt_broker}:{mqtt_port}")
            break
        except Exception as e:
            if i < max_retries - 1:
                print(f"Failed to connect to MQTT broker (attempt {i+1}/{max_retries}): {e}")
                print(f"Retrying in {retry_delay} seconds...")
                await asyncio.sleep(retry_delay)
            else:
                print(f"Failed to connect to MQTT broker after {max_retries} attempts. Continuing without MQTT...")
                mqtt_client = None
    
    if mqtt_client:
        control_loop = ControlLoop(registry, mqtt_client, feeder_id)
        asyncio.create_task(control_loop.run())
    else:
        print("Warning: Control loop not started due to MQTT connection failure")


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

