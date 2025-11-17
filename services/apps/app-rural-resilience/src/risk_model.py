"""
Risk Model
計算風險評分與升級建議
"""
from typing import List, Dict, Optional
from src.ingestion import FaultLogIngestion


class RiskModel:
    """風險模型（v0: rule-based）"""
    def __init__(self, ingestion: FaultLogIngestion):
        self.ingestion = ingestion
        self.fault_rate_threshold = 0.1  # 每年 0.1 次/公里
        self.distance_threshold = 20.0  # 20 公里
    
    def calculate_risk_scores(self, topology: Optional[Dict] = None) -> Dict[str, float]:
        """
        計算每個 asset 的風險評分
        
        Args:
            topology: 拓樸資料（可選，用於計算距離）
        
        Returns:
            asset_id -> risk_score (0-1)
        """
        risk_scores = {}
        
        # 取得所有資產
        assets = self._get_assets_from_topology(topology) if topology else []
        
        for asset in assets:
            asset_id = asset.get("id")
            if not asset_id:
                continue
            
            # 計算風險評分
            risk_score = self._calculate_asset_risk(asset_id, asset, topology)
            risk_scores[asset_id] = risk_score
        
        return risk_scores
    
    def _get_assets_from_topology(self, topology: Dict) -> List[Dict]:
        """從拓樸取得資產列表"""
        nodes = topology.get("nodes", [])
        lines = topology.get("lines", [])
        
        assets = []
        assets.extend(nodes)
        assets.extend(lines)
        return assets
    
    def _calculate_asset_risk(self, asset_id: str, asset: Dict, topology: Optional[Dict]) -> float:
        """計算單一資產的風險"""
        risk_factors = []
        
        # Factor 1: 故障頻率
        fault_count = self.ingestion.get_fault_count_by_asset(asset_id)
        fault_rate = fault_count / 10.0  # 假設 10 年歷史
        fault_risk = min(fault_rate / self.fault_rate_threshold, 1.0)
        risk_factors.append(fault_risk * 0.5)  # 權重 50%
        
        # Factor 2: 距離變電站（簡化：假設遠離的風險較高）
        distance_risk = 0.0
        if topology:
            # 簡化：假設節點數越多，距離越遠
            nodes = topology.get("nodes", [])
            node_index = next((i for i, n in enumerate(nodes) if n.get("id") == asset_id), -1)
            if node_index >= 0:
                # 假設距離與節點索引成正比
                normalized_distance = min(node_index / len(nodes), 1.0) if nodes else 0.0
                distance_risk = normalized_distance * 0.3  # 權重 30%
        
        risk_factors.append(distance_risk)
        
        # Factor 3: 單一路徑（簡化：假設只有一個連接的風險較高）
        path_risk = 0.0
        if topology:
            lines = topology.get("lines", [])
            connections = sum(1 for line in lines 
                            if line.get("from_node_id") == asset_id or 
                               line.get("to_node_id") == asset_id)
            if connections <= 1:
                path_risk = 0.2  # 權重 20%
        
        risk_factors.append(path_risk)
        
        # 總風險 = 加權平均
        total_risk = sum(risk_factors)
        return min(total_risk, 1.0)
    
    def generate_upgrade_suggestions(self, topology: Optional[Dict] = None) -> List[Dict]:
        """
        生成升級建議
        
        Returns:
            建議列表，每個建議包含位置和類型
        """
        suggestions = []
        risk_scores = self.calculate_risk_scores(topology)
        
        # 找出高風險區域
        high_risk_assets = [
            (asset_id, score) 
            for asset_id, score in risk_scores.items() 
            if score > 0.6
        ]
        
        # 對每個高風險資產生成建議
        for asset_id, risk_score in high_risk_assets:
            asset_metadata = self.ingestion.get_asset_metadata(asset_id)
            
            # 建議 1: 如果故障率高，建議加 sectionalizer
            fault_count = self.ingestion.get_fault_count_by_asset(asset_id)
            if fault_count > 3:
                suggestions.append({
                    "type": "sectionalizer",
                    "location": asset_id,
                    "reason": f"High fault rate ({fault_count} faults)",
                    "priority": "high",
                    "estimated_cost": 50000
                })
            
            # 建議 2: 如果距離遠且單一路徑，建議 parallel line
            if risk_score > 0.7:
                suggestions.append({
                    "type": "parallel_line",
                    "location": asset_id,
                    "reason": "Long distance with single path",
                    "priority": "medium",
                    "estimated_cost": 200000
                })
        
        return suggestions

