"""
API Routes
"""
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import Optional, Dict, List
from src.rules import AnomalyDetector

router = APIRouter()

# 取得實例（應該從 main 注入）
detector: Optional[AnomalyDetector] = None


def set_detector(det: AnomalyDetector):
    """設定 detector"""
    global detector
    detector = det


class LogRequest(BaseModel):
    source: str
    type: str
    message: str
    metadata: Optional[Dict] = None


@router.post("/logs")
async def ingest_log(request: LogRequest):
    """接收 log"""
    if not detector:
        raise HTTPException(status_code=500, detail="Detector not initialized")
    
    log_data = {
        "source": request.source,
        "type": request.type,
        "message": request.message,
        "metadata": request.metadata or {}
    }
    
    result = detector.ingest_log(log_data)
    
    return result


@router.get("/logs")
async def get_logs(limit: int = 100):
    """取得 logs"""
    if not detector:
        raise HTTPException(status_code=500, detail="Detector not initialized")
    
    logs = detector.get_recent_logs(limit)
    
    return {
        "logs": logs,
        "count": len(logs)
    }


@router.get("/anomalies/summary")
async def get_anomaly_summary():
    """取得異常摘要"""
    if not detector:
        raise HTTPException(status_code=500, detail="Detector not initialized")
    
    summary = detector.get_anomaly_summary()
    
    return summary

