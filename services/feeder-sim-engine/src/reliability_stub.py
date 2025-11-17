"""
Reliability Stub
提供簡單的可靠度分析 stub（生成合理的 SAIDI/SAIFI 資料）
"""
import random
import math
from typing import Dict, Any, List


class ReliabilityStub:
    """可靠度分析 stub - 生成合理的 SAIDI/SAIFI 資料"""
    
    def __init__(self):
        # 預設參數
        self.default_fault_rate = 0.1  # 故障率（次/年/公里）
        self.default_repair_time = 240  # 平均修復時間（分鐘）
        
    def run_reliability_analysis(
        self, 
        topology: Dict[str, Any], 
        parameters: Dict[str, Any]
    ) -> Dict[str, Any]:
        """
        執行可靠度分析（stub）
        
        Args:
            topology: 拓樸資料
            parameters: 分析參數
            
        Returns:
            包含 SAIDI/SAIFI 的結果
        """
        nodes = topology.get("nodes", [])
        lines = topology.get("lines", [])
        profile_type = topology.get("profile_type", "suburban")
        
        # 從參數取得設定值，或使用預設值
        fault_rate = parameters.get("fault_rate", self.default_fault_rate)
        repair_time = parameters.get("repair_time", self.default_repair_time)
        
        # 估算 feeder 長度（簡化：假設每個 line 平均 1 km）
        estimated_length = len(lines) * 1.0
        
        # 計算預期故障次數
        expected_faults = estimated_length * fault_rate
        
        # 根據 profile 類型調整
        profile_multipliers = {
            "rural": 1.5,  # 鄉村地區故障率較高
            "suburban": 1.0,
            "urban": 0.7,  # 都市地區故障率較低
        }
        multiplier = profile_multipliers.get(profile_type, 1.0)
        expected_faults *= multiplier
        
        # 計算 SAIFI（System Average Interruption Frequency Index）
        # 假設平均每個故障影響 30% 的客戶
        avg_customers_affected = len(nodes) * 0.3
        saifi = expected_faults / len(nodes) if len(nodes) > 0 else 0
        
        # 計算 SAIDI（System Average Interruption Duration Index）
        # SAIDI = (總停電時間) / (總客戶數)
        total_outage_minutes = expected_faults * repair_time * avg_customers_affected
        saidi = total_outage_minutes / len(nodes) if len(nodes) > 0 else 0
        
        # 生成各節點的風險評分
        node_risks = []
        for node in nodes:
            node_id = node.get("id", "")
            node_type = node.get("type", "bus")
            
            # 根據節點類型與位置計算風險
            base_risk = random.uniform(0.1, 0.5)
            
            # 調整風險（變壓器、開關風險較高）
            if node_type == "transformer":
                base_risk *= 1.5
            elif node_type == "switch":
                base_risk *= 1.2
            
            node_risks.append({
                "node_id": node_id,
                "risk_score": round(base_risk, 3),
                "risk_level": self._get_risk_level(base_risk)
            })
        
        # 生成各線路的風險評分
        line_risks = []
        for line in lines:
            line_id = line.get("id", "")
            
            # 線路風險基於長度與位置
            base_risk = random.uniform(0.2, 0.6)
            
            line_risks.append({
                "line_id": line_id,
                "risk_score": round(base_risk, 3),
                "risk_level": self._get_risk_level(base_risk)
            })
        
        return {
            "saidi": round(saidi, 2),  # 分鐘/年
            "saifi": round(saifi, 2),  # 次/年
            "expected_faults_per_year": round(expected_faults, 2),
            "average_repair_time_minutes": repair_time,
            "node_risks": node_risks,
            "line_risks": line_risks,
            "summary": {
                "total_nodes": len(nodes),
                "total_lines": len(lines),
                "estimated_length_km": round(estimated_length, 2),
                "profile_type": profile_type
            }
        }
    
    def _get_risk_level(self, risk_score: float) -> str:
        """根據風險評分回傳風險等級"""
        if risk_score < 0.3:
            return "low"
        elif risk_score < 0.6:
            return "medium"
        else:
            return "high"

