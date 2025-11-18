"""
ESG / Carbon Credit Stub
提供 ESG 和碳權計算模擬
"""
import random
import math
from typing import List, Dict, Any, Optional
from datetime import datetime


class ESGStub:
    """ESG 和碳權計算 stub"""
    
    # 碳排放係數（kg CO2/kWh）
    EMISSION_FACTORS = {
        "grid": 0.5,  # 電網平均碳排放係數（台灣約 0.5 kg CO2/kWh）
        "solar": 0.05,  # 太陽能（含製造過程）
        "wind": 0.02,  # 風力發電
        "battery": 0.1,  # 電池（充放電損失）
        "ev_charging": 0.5,  # EV 充電（假設來自電網）
    }
    
    # 碳權價格（USD/ton CO2）
    CARBON_CREDIT_PRICE = 50.0
    
    def __init__(self):
        pass
    
    def calculate_emissions(
        self,
        topology: Dict[str, Any],
        parameters: Optional[Dict[str, Any]] = None
    ) -> Dict[str, Any]:
        """
        計算碳排放
        
        Args:
            topology: 拓樸結構
            parameters: 計算參數（時間範圍、負載等）
            
        Returns:
            碳排放計算結果
        """
        if parameters is None:
            parameters = {}
        
        nodes = topology.get("nodes", [])
        lines = topology.get("lines", [])
        
        # 計算時間範圍（小時）
        time_hours = parameters.get("time_hours", 24)
        
        # 計算各節點的碳排放
        node_emissions = []
        total_emissions = 0.0
        
        for node in nodes:
            node_id = node.get("id", "")
            node_type = node.get("type", "bus")
            properties = node.get("properties", {})
            
            # 根據節點類型計算碳排放
            emission = 0.0
            power_kw = 0.0
            
            if node_type == "ev_charger":
                # EV 充電樁
                power_kw = properties.get("power_rating", 7.0)  # 預設 7kW
                charging_hours = parameters.get("ev_charging_hours", 4.0)
                energy_kwh = power_kw * charging_hours
                emission = energy_kwh * self.EMISSION_FACTORS["ev_charging"]
                
            elif node_type == "solar":
                # 太陽能發電（負排放）
                power_kw = properties.get("capacity", 5.0)  # 預設 5kW
                generation_hours = parameters.get("solar_generation_hours", 6.0)
                energy_kwh = power_kw * generation_hours
                emission = -energy_kwh * self.EMISSION_FACTORS["solar"]  # 負值表示減排
                
            elif node_type == "battery":
                # 電池（充放電）
                power_kw = properties.get("capacity", 10.0)  # 預設 10kWh
                cycles = parameters.get("battery_cycles", 1.0)
                energy_kwh = power_kw * cycles
                emission = energy_kwh * self.EMISSION_FACTORS["battery"]
                
            elif node_type == "transformer":
                # 變壓器（損耗）
                power_kw = properties.get("rating", 500.0)  # 預設 500kVA
                load_factor = properties.get("load_factor", 0.7)
                loss_factor = 0.02  # 2% 損耗
                energy_kwh = power_kw * load_factor * loss_factor * time_hours
                emission = energy_kwh * self.EMISSION_FACTORS["grid"]
                
            else:
                # 一般負載節點
                power_kw = properties.get("load", 10.0)  # 預設 10kW
                energy_kwh = power_kw * time_hours
                emission = energy_kwh * self.EMISSION_FACTORS["grid"]
            
            node_emissions.append({
                "node_id": node_id,
                "node_type": node_type,
                "power_kw": power_kw,
                "energy_kwh": energy_kwh if node_type != "solar" else -energy_kwh,
                "emission_kg_co2": emission,
            })
            
            total_emissions += emission
        
        # 計算碳權價值
        total_emissions_ton = total_emissions / 1000.0
        carbon_credits = max(0, -total_emissions_ton)  # 只有負排放（減排）才能獲得碳權
        carbon_credit_value = carbon_credits * self.CARBON_CREDIT_PRICE
        
        # 計算 ESG 分數（0-100）
        esg_score = self._calculate_esg_score(total_emissions_ton, nodes, topology)
        
        return {
            "timestamp": datetime.now().isoformat(),
            "time_hours": time_hours,
            "total_emissions_kg_co2": total_emissions,
            "total_emissions_ton_co2": total_emissions_ton,
            "carbon_credits_ton": carbon_credits,
            "carbon_credit_value_usd": carbon_credit_value,
            "esg_score": esg_score,
            "node_emissions": node_emissions,
            "recommendations": self._generate_recommendations(total_emissions_ton, nodes),
        }
    
    def _calculate_esg_score(
        self,
        total_emissions_ton: float,
        nodes: List[Dict],
        topology: Dict[str, Any]
    ) -> float:
        """
        計算 ESG 分數（0-100）
        
        考慮因素：
        - 總碳排放量
        - 再生能源比例
        - DER 滲透率
        """
        # 基礎分數
        base_score = 50.0
        
        # 計算再生能源比例
        renewable_count = sum(1 for n in nodes if n.get("type") in ["solar", "wind", "battery"])
        total_assets = len([n for n in nodes if n.get("type") != "bus"])
        renewable_ratio = renewable_count / max(total_assets, 1)
        
        # 碳排放影響（越低越好）
        emission_penalty = min(total_emissions_ton * 10, 30)  # 最多扣 30 分
        
        # 再生能源加分
        renewable_bonus = renewable_ratio * 20  # 最多加 20 分
        
        # DER 滲透率加分
        der_count = sum(1 for n in nodes if n.get("type") in ["solar", "battery", "ev_charger"])
        der_ratio = der_count / max(len(nodes), 1)
        der_bonus = der_ratio * 10  # 最多加 10 分
        
        # 計算最終分數
        score = base_score - emission_penalty + renewable_bonus + der_bonus
        return max(0, min(100, score))
    
    def _generate_recommendations(
        self,
        total_emissions_ton: float,
        nodes: List[Dict]
    ) -> List[Dict[str, Any]]:
        """生成減排建議"""
        recommendations = []
        
        # 檢查是否有太陽能
        has_solar = any(n.get("type") == "solar" for n in nodes)
        if not has_solar and total_emissions_ton > 0:
            recommendations.append({
                "type": "add_solar",
                "priority": "high",
                "title": "建議新增太陽能發電",
                "description": "新增太陽能發電可大幅降低碳排放，並獲得碳權",
                "estimated_reduction_ton": total_emissions_ton * 0.3,
            })
        
        # 檢查是否有電池儲能
        has_battery = any(n.get("type") == "battery" for n in nodes)
        if not has_battery and total_emissions_ton > 0:
            recommendations.append({
                "type": "add_battery",
                "priority": "medium",
                "title": "建議新增電池儲能系統",
                "description": "電池儲能可優化再生能源使用，減少電網依賴",
                "estimated_reduction_ton": total_emissions_ton * 0.1,
            })
        
        # 檢查 EV 充電樁
        ev_count = sum(1 for n in nodes if n.get("type") == "ev_charger")
        if ev_count > 0 and total_emissions_ton > 0:
            recommendations.append({
                "type": "optimize_ev_charging",
                "priority": "medium",
                "title": "優化 EV 充電策略",
                "description": "在再生能源發電高峰時段充電，可降低碳排放",
                "estimated_reduction_ton": total_emissions_ton * 0.15,
            })
        
        return recommendations

