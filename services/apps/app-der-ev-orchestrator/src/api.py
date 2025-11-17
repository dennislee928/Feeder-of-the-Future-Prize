"""
API Routes
提供 REST API 給外部調用
"""
from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import Optional, Dict
from src.registry import AssetRegistry

router = APIRouter()

# 取得 registry 實例（應該從 main 注入，這裡簡化）
registry: Optional[AssetRegistry] = None


def set_registry(reg: AssetRegistry):
    """設定 registry"""
    global registry
    registry = reg


class RegisterEVChargerRequest(BaseModel):
    name: str
    rated_power_kw: float
    properties: Optional[Dict] = None


class RegisterPVBatteryRequest(BaseModel):
    name: str
    rated_power_kw: float
    properties: Optional[Dict] = None


@router.post("/assets/ev-chargers")
async def register_ev_charger(request: RegisterEVChargerRequest):
    """註冊 EV 充電樁"""
    if not registry:
        raise HTTPException(status_code=500, detail="Registry not initialized")
    
    asset = registry.register_ev_charger(
        request.name,
        request.rated_power_kw,
        request.properties
    )
    
    return {
        "asset_id": asset.id,
        "name": asset.name,
        "type": asset.type,
        "rated_power_kw": asset.rated_power_kw,
        "registered_at": asset.registered_at.isoformat()
    }


@router.post("/assets/pv-batteries")
async def register_pv_battery(request: RegisterPVBatteryRequest):
    """註冊 PV + Battery"""
    if not registry:
        raise HTTPException(status_code=500, detail="Registry not initialized")
    
    asset = registry.register_pv_battery(
        request.name,
        request.rated_power_kw,
        request.properties
    )
    
    return {
        "asset_id": asset.id,
        "name": asset.name,
        "type": asset.type,
        "rated_power_kw": asset.rated_power_kw,
        "registered_at": asset.registered_at.isoformat()
    }


@router.get("/assets")
async def list_assets(asset_type: Optional[str] = None):
    """列出所有資產"""
    if not registry:
        raise HTTPException(status_code=500, detail="Registry not initialized")
    
    assets = registry.list_assets(asset_type)
    
    return {
        "assets": [
            {
                "id": a.id,
                "name": a.name,
                "type": a.type,
                "rated_power_kw": a.rated_power_kw,
                "current_power_kw": a.current_power_kw,
                "registered_at": a.registered_at.isoformat()
            }
            for a in assets
        ]
    }

