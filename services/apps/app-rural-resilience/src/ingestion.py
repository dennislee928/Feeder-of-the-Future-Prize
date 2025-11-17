"""
Fault Log Ingestion
接收並儲存故障紀錄
"""
from typing import List, Dict, Optional
from datetime import datetime
from dataclasses import dataclass, field


@dataclass
class FaultLog:
    """故障紀錄"""
    id: str
    asset_id: str
    asset_type: str
    fault_type: str
    timestamp: datetime
    repair_time_minutes: float
    description: Optional[str] = None


class FaultLogIngestion:
    """故障紀錄接收器"""
    def __init__(self):
        self._faults: List[FaultLog] = []
        self._assets: Dict[str, Dict] = {}  # asset_id -> metadata
    
    def ingest_fault(self, fault_data: Dict) -> FaultLog:
        """接收故障紀錄"""
        fault = FaultLog(
            id=fault_data.get("id", f"fault-{len(self._faults)}"),
            asset_id=fault_data["asset_id"],
            asset_type=fault_data.get("asset_type", "unknown"),
            fault_type=fault_data.get("fault_type", "unknown"),
            timestamp=datetime.fromisoformat(fault_data["timestamp"]) if isinstance(fault_data.get("timestamp"), str) else fault_data.get("timestamp", datetime.now()),
            repair_time_minutes=fault_data.get("repair_time_minutes", 0.0),
            description=fault_data.get("description")
        )
        self._faults.append(fault)
        return fault
    
    def get_faults_by_asset(self, asset_id: str) -> List[FaultLog]:
        """取得特定資產的故障紀錄"""
        return [f for f in self._faults if f.asset_id == asset_id]
    
    def get_all_faults(self) -> List[FaultLog]:
        """取得所有故障紀錄"""
        return self._faults.copy()
    
    def get_fault_count_by_asset(self, asset_id: str) -> int:
        """取得資產的故障次數"""
        return len(self.get_faults_by_asset(asset_id))
    
    def register_asset(self, asset_id: str, metadata: Dict):
        """註冊資產 metadata"""
        self._assets[asset_id] = metadata
    
    def get_asset_metadata(self, asset_id: str) -> Optional[Dict]:
        """取得資產 metadata"""
        return self._assets.get(asset_id)

