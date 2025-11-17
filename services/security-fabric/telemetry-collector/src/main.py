"""
Telemetry Collector
收集 logs 並進行異常偵測
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from src.rules import AnomalyDetector

app = FastAPI(
    title="Telemetry Collector",
    description="Collect telemetry and detect anomalies",
    version="0.1.0"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# 初始化異常偵測器
detector = AnomalyDetector()

# 導入 API routes
from src.api import router, set_detector
set_detector(detector)
app.include_router(router, prefix="/api/v1", tags=["telemetry"])


@app.get("/health")
async def health_check():
    """健康檢查"""
    return {"status": "ok"}


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8085)

