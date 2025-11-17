"""
Powerflow Stub
提供簡單的潮流分析 stub（生成合理的假資料）
"""
import random
import math
from typing import List, Dict, Any


class PowerflowStub:
    """潮流分析 stub - 生成合理的電壓與載流率資料"""
    
    def __init__(self):
        # 基準電壓（kV）
        self.base_voltage = 12.47  # 典型配電電壓
        
    def run_powerflow(self, nodes: List[Dict], lines: List[Dict]) -> Dict[str, Any]:
        """
        執行潮流分析（stub）
        
        Args:
            nodes: 節點列表
            lines: 線路列表
            
        Returns:
            包含節點電壓與線路載流率的結果
        """
        # 生成節點電壓（pu，標么值）
        node_results = []
        for node in nodes:
            node_id = node.get("id", "")
            node_type = node.get("type", "bus")
            
            # 根據節點類型生成不同的電壓範圍
            if node_type == "transformer":
                # 變壓器節點電壓較穩定
                voltage_pu = random.uniform(0.98, 1.02)
            elif node_type == "ev_charger":
                # EV 充電樁可能造成電壓降
                voltage_pu = random.uniform(0.95, 1.00)
            else:
                # 一般節點
                voltage_pu = random.uniform(0.96, 1.01)
            
            # 計算實際電壓（kV）
            voltage_kv = voltage_pu * self.base_voltage
            
            # 計算電壓偏差（%）
            voltage_deviation = (voltage_pu - 1.0) * 100
            
            node_results.append({
                "node_id": node_id,
                "voltage_pu": round(voltage_pu, 4),
                "voltage_kv": round(voltage_kv, 4),
                "voltage_deviation_percent": round(voltage_deviation, 2),
                "status": "normal" if abs(voltage_deviation) < 5 else "warning"
            })
        
        # 生成線路載流率
        line_results = []
        for line in lines:
            line_id = line.get("id", "")
            
            # 生成載流率（0-100%）
            loading_percent = random.uniform(20, 85)
            
            # 根據載流率判斷狀態
            if loading_percent > 80:
                status = "warning"
            elif loading_percent > 90:
                status = "critical"
            else:
                status = "normal"
            
            line_results.append({
                "line_id": line_id,
                "loading_percent": round(loading_percent, 2),
                "status": status
            })
        
        # 計算整體統計
        avg_voltage_pu = sum(n["voltage_pu"] for n in node_results) / len(node_results) if node_results else 1.0
        max_loading = max((l["loading_percent"] for l in line_results), default=0)
        
        return {
            "nodes": node_results,
            "lines": line_results,
            "summary": {
                "average_voltage_pu": round(avg_voltage_pu, 4),
                "max_line_loading_percent": round(max_loading, 2),
                "total_nodes": len(node_results),
                "total_lines": len(line_results),
                "converged": True
            }
        }

