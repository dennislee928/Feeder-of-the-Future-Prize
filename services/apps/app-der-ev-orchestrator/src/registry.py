"""
Asset Registry
管理 EV chargers 和 PV batteries
"""
from typing import Dict, List, Optional
from datetime import datetime
import uuid


class Asset:
    """資產（EV charger 或 PV battery）"""
    def __init__(self, asset_id: str, asset_type: str, name: str, 
                 rated_power_kw: float, properties: Dict = None):
        self.id = asset_id
        self.type = asset_type  # "ev_charger" or "pv_battery"
        self.name = name
        self.rated_power_kw = rated_power_kw
        self.properties = properties or {}
        self.registered_at = datetime.now()
        self.is_controllable = True
        self.current_power_kw = 0.0
        self.max_power_kw = rated_power_kw
        self.min_power_kw = 0.0


class AssetRegistry:
    """資產註冊表"""
    def __init__(self):
        self._assets: Dict[str, Asset] = {}
    
    def register_ev_charger(self, name: str, rated_power_kw: float, 
                           properties: Dict = None) -> Asset:
        """註冊 EV 充電樁"""
        asset_id = f"ev-{uuid.uuid4().hex[:8]}"
        asset = Asset(asset_id, "ev_charger", name, rated_power_kw, properties)
        self._assets[asset_id] = asset
        return asset
    
    def register_pv_battery(self, name: str, rated_power_kw: float,
                           properties: Dict = None) -> Asset:
        """註冊 PV + Battery"""
        asset_id = f"pv-{uuid.uuid4().hex[:8]}"
        asset = Asset(asset_id, "pv_battery", name, rated_power_kw, properties)
        self._assets[asset_id] = asset
        return asset
    
    def get_asset(self, asset_id: str) -> Optional[Asset]:
        """取得資產"""
        return self._assets.get(asset_id)
    
    def list_assets(self, asset_type: Optional[str] = None) -> List[Asset]:
        """列出資產"""
        if asset_type:
            return [a for a in self._assets.values() if a.type == asset_type]
        return list(self._assets.values())
    
    def update_asset_power(self, asset_id: str, power_kw: float) -> bool:
        """更新資產功率"""
        asset = self._assets.get(asset_id)
        if not asset:
            return False
        
        # 確保功率在範圍內
        power_kw = max(asset.min_power_kw, min(asset.max_power_kw, power_kw))
        asset.current_power_kw = power_kw
        return True
    
    def get_total_power(self, asset_type: Optional[str] = None) -> float:
        """取得總功率"""
        assets = self.list_assets(asset_type)
        return sum(a.current_power_kw for a in assets)

