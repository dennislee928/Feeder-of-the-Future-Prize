"""
Penetration Test Stub
提供滲透測試模擬（Layer 7-1 全棧攻擊場景）
"""
import random
import uuid
from typing import List, Dict, Any, Optional
from datetime import datetime


class PenetrationStub:
    """滲透測試 stub - 模擬各種攻擊場景"""
    
    # 攻擊場景定義
    ATTACK_SCENARIOS = {
        # Layer 7 (應用層)
        "sql_injection": {
            "name": "SQL Injection on API",
            "layer": 7,
            "severity": "high",
            "description": "SQL注入攻擊，嘗試透過API參數注入惡意SQL指令"
        },
        "xss": {
            "name": "XSS in Web Interface",
            "layer": 7,
            "severity": "medium",
            "description": "跨站腳本攻擊，在Web介面中注入惡意腳本"
        },
        "unauthorized_access": {
            "name": "Unauthorized API Access",
            "layer": 7,
            "severity": "high",
            "description": "未授權API存取，嘗試存取需要權限的端點"
        },
        "command_injection_mqtt": {
            "name": "Command Injection via MQTT",
            "layer": 7,
            "severity": "critical",
            "description": "透過MQTT訊息注入系統命令"
        },
        # Layer 6 (表示層)
        "data_tampering": {
            "name": "Data Tampering",
            "layer": 6,
            "severity": "high",
            "description": "數據篡改攻擊，修改傳輸中的數據"
        },
        "encryption_bypass": {
            "name": "Encryption Bypass",
            "layer": 6,
            "severity": "critical",
            "description": "加密繞過，嘗試破解或繞過加密機制"
        },
        # Layer 5 (會話層)
        "session_hijacking": {
            "name": "Session Hijacking",
            "layer": 5,
            "severity": "high",
            "description": "會話劫持，竊取並使用有效會話"
        },
        "replay_attack": {
            "name": "Replay Attack",
            "layer": 5,
            "severity": "medium",
            "description": "重放攻擊，重複發送已截獲的訊息"
        },
        # Layer 4 (傳輸層)
        "tcp_syn_flood": {
            "name": "TCP SYN Flood",
            "layer": 4,
            "severity": "high",
            "description": "TCP SYN洪水攻擊，耗盡伺服器資源"
        },
        "udp_flood": {
            "name": "UDP Flood",
            "layer": 4,
            "severity": "high",
            "description": "UDP洪水攻擊，大量UDP封包造成服務中斷"
        },
        # Layer 3 (網路層)
        "icmp_flood": {
            "name": "ICMP Flood",
            "layer": 3,
            "severity": "medium",
            "description": "ICMP洪水攻擊，大量ping請求"
        },
        "route_poisoning": {
            "name": "Route Poisoning",
            "layer": 3,
            "severity": "critical",
            "description": "路由中毒，操縱路由表導致流量重定向"
        },
        # Layer 2 (數據鏈路層)
        "arp_spoofing": {
            "name": "ARP Spoofing",
            "layer": 2,
            "severity": "high",
            "description": "ARP欺騙，偽造MAC地址進行中間人攻擊"
        },
        "mac_flood": {
            "name": "MAC Flood",
            "layer": 2,
            "severity": "medium",
            "description": "MAC洪水攻擊，耗盡交換機MAC表"
        },
        # Layer 1 (物理層)
        "physical_access": {
            "name": "Physical Device Access",
            "layer": 1,
            "severity": "critical",
            "description": "物理設備存取，直接接觸硬體設備"
        },
        "line_disconnection": {
            "name": "Line Disconnection",
            "layer": 1,
            "severity": "critical",
            "description": "線路切斷，物理切斷電力線路"
        }
    }
    
    def __init__(self):
        pass
    
    def run_penetration_test(
        self,
        topology: Dict[str, Any],
        attack_scenarios: List[str],
        target_nodes: Optional[List[str]] = None
    ) -> Dict[str, Any]:
        """
        執行滲透測試
        
        Args:
            topology: 拓樸結構（包含nodes和lines）
            attack_scenarios: 要執行的攻擊場景列表
            target_nodes: 目標節點（可選，如果未指定則隨機選擇）
            
        Returns:
            滲透測試結果
        """
        nodes = topology.get("nodes", [])
        lines = topology.get("lines", [])
        
        if not nodes:
            return {
                "error": "No nodes in topology",
                "summary": {
                    "total_attacks": 0,
                    "successful": 0,
                    "failed": 0,
                    "critical_vulnerabilities": 0
                }
            }
        
        # 執行每個攻擊場景
        attack_results = []
        successful_attacks = 0
        failed_attacks = 0
        critical_vulns = 0
        
        for scenario_id in attack_scenarios:
            if scenario_id not in self.ATTACK_SCENARIOS:
                continue
            
            scenario_info = self.ATTACK_SCENARIOS[scenario_id]
            
            # 決定攻擊是否成功（根據嚴重程度調整成功率）
            success_rate = {
                "low": 0.3,
                "medium": 0.5,
                "high": 0.7,
                "critical": 0.9
            }.get(scenario_info["severity"], 0.5)
            
            is_successful = random.random() < success_rate
            
            # 選擇受影響的節點和線路
            affected_nodes = self._select_affected_nodes(nodes, target_nodes, scenario_info["layer"])
            affected_lines = self._select_affected_lines(lines, affected_nodes)
            
            # 生成攻擊路徑
            attack_path = self._generate_attack_path(nodes, lines, affected_nodes, scenario_info["layer"])
            
            # 生成影響描述和建議
            impact, recommendations = self._generate_impact_and_recommendations(
                scenario_id, scenario_info, is_successful, affected_nodes
            )
            
            if is_successful:
                successful_attacks += 1
                if scenario_info["severity"] == "critical":
                    critical_vulns += 1
            else:
                failed_attacks += 1
            
            attack_results.append({
                "attack_id": str(uuid.uuid4()),
                "scenario": scenario_id,
                "scenario_name": scenario_info["name"],
                "layer": scenario_info["layer"],
                "severity": scenario_info["severity"],
                "successful": is_successful,
                "affected_nodes": affected_nodes,
                "affected_lines": affected_lines,
                "attack_path": attack_path,
                "impact": impact,
                "recommendations": recommendations,
                "timestamp": datetime.now().isoformat()
            })
        
        return {
            "attacks": attack_results,
            "summary": {
                "total_attacks": len(attack_results),
                "successful": successful_attacks,
                "failed": failed_attacks,
                "critical_vulnerabilities": critical_vulns,
                "total_nodes": len(nodes),
                "total_lines": len(lines),
                "affected_nodes_count": len(set(
                    node for result in attack_results for node in result["affected_nodes"]
                )),
                "affected_lines_count": len(set(
                    line for result in attack_results for line in result["affected_lines"]
                ))
            }
        }
    
    def _select_affected_nodes(
        self,
        nodes: List[Dict],
        target_nodes: Optional[List[str]],
        layer: int
    ) -> List[str]:
        """選擇受影響的節點"""
        if target_nodes:
            # 只選擇存在的目標節點
            node_ids = [n["id"] for n in nodes]
            return [nid for nid in target_nodes if nid in node_ids]
        
        # 根據層級決定影響範圍
        # 高層攻擊（Layer 7-5）通常影響特定節點
        # 低層攻擊（Layer 4-1）可能影響多個節點
        if layer >= 5:
            num_affected = random.randint(1, min(3, len(nodes)))
        else:
            num_affected = random.randint(1, min(len(nodes), max(1, len(nodes) // 2)))
        
        selected = random.sample([n["id"] for n in nodes], num_affected)
        return selected
    
    def _select_affected_lines(
        self,
        lines: List[Dict],
        affected_nodes: List[str]
    ) -> List[str]:
        """選擇受影響的線路（連接到受影響節點的線路）"""
        affected_lines = []
        for line in lines:
            from_node = line.get("from_node_id", "")
            to_node = line.get("to_node_id", "")
            if from_node in affected_nodes or to_node in affected_nodes:
                affected_lines.append(line.get("id", ""))
        return affected_lines
    
    def _generate_attack_path(
        self,
        nodes: List[Dict],
        lines: List[Dict],
        affected_nodes: List[str],
        layer: int
    ) -> List[Dict[str, str]]:
        """生成攻擊路徑"""
        if not affected_nodes:
            return []
        
        # 簡化：從第一個節點開始，連接到其他受影響節點
        path = []
        if len(affected_nodes) > 1:
            for i in range(len(affected_nodes) - 1):
                path.append({
                    "from": affected_nodes[i],
                    "to": affected_nodes[i + 1],
                    "layer": str(layer)
                })
        else:
            # 單一節點攻擊，標示為入口點
            path.append({
                "from": "external",
                "to": affected_nodes[0],
                "layer": str(layer)
            })
        
        return path
    
    def _generate_impact_and_recommendations(
        self,
        scenario_id: str,
        scenario_info: Dict,
        is_successful: bool,
        affected_nodes: List[str]
    ) -> tuple[str, List[str]]:
        """生成影響描述和修復建議"""
        if not is_successful:
            impact = f"攻擊 {scenario_info['name']} 被成功阻擋。"
            recommendations = [
                "繼續監控類似攻擊嘗試",
                "確認安全防護機制正常運作"
            ]
            return impact, recommendations
        
        # 根據攻擊類型生成不同的影響和建議
        impact_templates = {
            "sql_injection": f"SQL注入成功，可能導致資料庫資料外洩。受影響節點: {', '.join(affected_nodes)}",
            "xss": f"XSS攻擊成功，可能竊取使用者會話。受影響節點: {', '.join(affected_nodes)}",
            "unauthorized_access": f"未授權存取成功，攻擊者獲得系統存取權限。受影響節點: {', '.join(affected_nodes)}",
            "command_injection_mqtt": f"命令注入成功，可能導致系統被完全控制。受影響節點: {', '.join(affected_nodes)}",
            "data_tampering": f"數據篡改成功，傳輸中的數據被修改。受影響節點: {', '.join(affected_nodes)}",
            "encryption_bypass": f"加密繞過成功，敏感資料可能被竊取。受影響節點: {', '.join(affected_nodes)}",
            "session_hijacking": f"會話劫持成功，攻擊者可以使用被竊取的會話。受影響節點: {', '.join(affected_nodes)}",
            "replay_attack": f"重放攻擊成功，舊的訊息被重複執行。受影響節點: {', '.join(affected_nodes)}",
            "tcp_syn_flood": f"TCP SYN洪水攻擊成功，服務可能暫時中斷。受影響節點: {', '.join(affected_nodes)}",
            "udp_flood": f"UDP洪水攻擊成功，網路頻寬被耗盡。受影響節點: {', '.join(affected_nodes)}",
            "icmp_flood": f"ICMP洪水攻擊成功，網路設備可能過載。受影響節點: {', '.join(affected_nodes)}",
            "route_poisoning": f"路由中毒成功，流量可能被重定向到惡意節點。受影響節點: {', '.join(affected_nodes)}",
            "arp_spoofing": f"ARP欺騙成功，可能進行中間人攻擊。受影響節點: {', '.join(affected_nodes)}",
            "mac_flood": f"MAC洪水攻擊成功，交換機可能無法正常運作。受影響節點: {', '.join(affected_nodes)}",
            "physical_access": f"物理存取成功，設備可能被直接操控。受影響節點: {', '.join(affected_nodes)}",
            "line_disconnection": f"線路切斷成功，電力供應中斷。受影響節點: {', '.join(affected_nodes)}"
        }
        
        recommendation_templates = {
            "sql_injection": [
                "實施參數化查詢",
                "使用ORM框架避免直接SQL",
                "加強輸入驗證和過濾"
            ],
            "xss": [
                "實施內容安全政策(CSP)",
                "對所有使用者輸入進行編碼",
                "使用安全的框架和函式庫"
            ],
            "unauthorized_access": [
                "實施強制存取控制(MAC)",
                "使用多因素認證",
                "定期審查存取權限"
            ],
            "command_injection_mqtt": [
                "驗證所有MQTT訊息內容",
                "實施訊息簽章機制",
                "限制MQTT主題存取權限"
            ],
            "data_tampering": [
                "使用數位簽章驗證數據完整性",
                "實施端到端加密",
                "監控數據傳輸異常"
            ],
            "encryption_bypass": [
                "升級加密演算法",
                "使用強加密金鑰",
                "實施金鑰輪換機制"
            ],
            "session_hijacking": [
                "使用HTTPS傳輸",
                "實施會話固定保護",
                "定期更新會話ID"
            ],
            "replay_attack": [
                "實施時間戳記驗證",
                "使用nonce防止重放",
                "實施訊息序號檢查"
            ],
            "tcp_syn_flood": [
                "配置SYN cookies",
                "實施速率限制",
                "使用DDoS防護服務"
            ],
            "udp_flood": [
                "實施UDP速率限制",
                "配置防火牆規則",
                "使用流量清洗服務"
            ],
            "icmp_flood": [
                "限制ICMP請求速率",
                "配置防火牆阻擋不必要的ICMP",
                "實施網路流量監控"
            ],
            "route_poisoning": [
                "實施路由認證機制",
                "使用安全路由協定",
                "監控路由表異常變更"
            ],
            "arp_spoofing": [
                "實施靜態ARP表",
                "使用ARP監控工具",
                "實施網路分段隔離"
            ],
            "mac_flood": [
                "配置交換機MAC限制",
                "實施端口安全",
                "監控MAC表使用率"
            ],
            "physical_access": [
                "加強實體安全措施",
                "實施設備鎖定機制",
                "監控設備存取記錄"
            ],
            "line_disconnection": [
                "實施線路監控系統",
                "配置自動故障切換",
                "加強實體線路保護"
            ]
        }
        
        impact = impact_templates.get(scenario_id, f"攻擊 {scenario_info['name']} 成功。受影響節點: {', '.join(affected_nodes)}")
        recommendations = recommendation_templates.get(scenario_id, [
            "加強安全防護措施",
            "實施深度防禦策略",
            "定期進行安全審計"
        ])
        
        return impact, recommendations
    
    def get_available_scenarios(self) -> List[Dict[str, Any]]:
        """取得所有可用的攻擊場景"""
        return [
            {
                "id": scenario_id,
                **info
            }
            for scenario_id, info in self.ATTACK_SCENARIOS.items()
        ]

