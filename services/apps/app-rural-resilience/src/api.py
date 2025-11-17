"""
API Routes
"""
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import Optional, Dict, List
from datetime import datetime
from src.ingestion import FaultLogIngestion
from src.risk_model import RiskModel

router = APIRouter()

# 取得實例（應該從 main 注入）
ingestion: Optional[FaultLogIngestion] = None
risk_model: Optional[RiskModel] = None


def set_components(ing: FaultLogIngestion, rm: RiskModel):
    """設定組件"""
    global ingestion, risk_model
    ingestion = ing
    risk_model = rm


class FaultLogRequest(BaseModel):
    asset_id: str
    asset_type: Optional[str] = "line"
    fault_type: str
    timestamp: Optional[str] = None
    repair_time_minutes: float
    description: Optional[str] = None


class TopologyRequest(BaseModel):
    topology: Optional[Dict] = None


@router.post("/fault-logs")
async def upload_fault_log(request: FaultLogRequest):
    """上傳故障紀錄"""
    if not ingestion:
        raise HTTPException(status_code=500, detail="Ingestion not initialized")
    
    fault_data = {
        "asset_id": request.asset_id,
        "asset_type": request.asset_type,
        "fault_type": request.fault_type,
        "timestamp": request.timestamp or datetime.now().isoformat(),
        "repair_time_minutes": request.repair_time_minutes,
        "description": request.description
    }
    
    fault = ingestion.ingest_fault(fault_data)
    
    return {
        "id": fault.id,
        "asset_id": fault.asset_id,
        "timestamp": fault.timestamp.isoformat(),
        "status": "ingested"
    }


@router.get("/risk-scores")
async def get_risk_scores(topology: Optional[str] = None):
    """取得風險評分"""
    if not risk_model:
        raise HTTPException(status_code=500, detail="Risk model not initialized")
    
    # 解析 topology（簡化：這裡應該從參數或 body 取得）
    topology_data = None
    
    risk_scores = risk_model.calculate_risk_scores(topology_data)
    
    return {
        "risk_scores": risk_scores,
        "timestamp": datetime.now().isoformat()
    }


@router.post("/risk-scores")
async def calculate_risk_scores(request: TopologyRequest):
    """計算風險評分（使用拓樸資料）"""
    if not risk_model:
        raise HTTPException(status_code=500, detail="Risk model not initialized")
    
    risk_scores = risk_model.calculate_risk_scores(request.topology)
    
    return {
        "risk_scores": risk_scores,
        "timestamp": datetime.now().isoformat()
    }


@router.get("/upgrade-suggestions")
async def get_upgrade_suggestions(topology: Optional[str] = None):
    """取得升級建議"""
    if not risk_model:
        raise HTTPException(status_code=500, detail="Risk model not initialized")
    
    topology_data = None
    suggestions = risk_model.generate_upgrade_suggestions(topology_data)
    
    return {
        "suggestions": suggestions,
        "count": len(suggestions),
        "timestamp": datetime.now().isoformat()
    }


@router.post("/upgrade-suggestions")
async def generate_upgrade_suggestions(request: TopologyRequest):
    """生成升級建議（使用拓樸資料）"""
    if not risk_model:
        raise HTTPException(status_code=500, detail="Risk model not initialized")
    
    suggestions = risk_model.generate_upgrade_suggestions(request.topology)
    
    return {
        "suggestions": suggestions,
        "count": len(suggestions),
        "timestamp": datetime.now().isoformat()
    }

