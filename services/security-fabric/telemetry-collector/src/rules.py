"""
Anomaly Detection Rules
rule-based 異常偵測
"""
from typing import List, Dict
from datetime import datetime, timedelta
from collections import defaultdict


class AnomalyDetector:
    """異常偵測器（rule-based v0）"""
    def __init__(self):
        self.logs: List[Dict] = []
        self.command_counts: Dict[str, int] = defaultdict(int)
        self.config_changes: List[Dict] = []
        self.time_window = timedelta(hours=1)
        
        # 閾值
        self.max_commands_per_hour = 100
        self.max_config_changes_per_hour = 10
        self.late_night_start = 22  # 22:00
        self.late_night_end = 6     # 06:00
    
    def ingest_log(self, log_data: Dict):
        """接收 log"""
        log_entry = {
            "timestamp": datetime.now(),
            "source": log_data.get("source", "unknown"),
            "type": log_data.get("type", "info"),
            "message": log_data.get("message", ""),
            "metadata": log_data.get("metadata", {})
        }
        
        self.logs.append(log_entry)
        
        # 檢查異常
        anomalies = self.check_anomalies(log_entry)
        
        return {
            "log_id": len(self.logs),
            "anomalies": anomalies
        }
    
    def check_anomalies(self, log_entry: Dict) -> List[Dict]:
        """檢查異常"""
        anomalies = []
        
        # Rule 1: 非預期大量控制命令
        if log_entry["type"] == "command":
            source = log_entry["source"]
            self.command_counts[source] += 1
            
            # 檢查最近一小時的命令數
            recent_commands = sum(
                1 for log in self.logs[-100:]  # 檢查最近 100 條
                if log["type"] == "command" 
                and log["source"] == source
                and (log_entry["timestamp"] - log["timestamp"]) < self.time_window
            )
            
            if recent_commands > self.max_commands_per_hour:
                anomalies.append({
                    "type": "excessive_commands",
                    "severity": "high",
                    "message": f"Excessive commands from {source}: {recent_commands} in last hour",
                    "threshold": self.max_commands_per_hour
                })
        
        # Rule 2: 深夜大量 config 變更
        if log_entry["type"] == "config_change":
            hour = log_entry["timestamp"].hour
            is_late_night = hour >= self.late_night_start or hour < self.late_night_end
            
            if is_late_night:
                self.config_changes.append(log_entry)
                
                # 檢查最近一小時的變更數
                recent_changes = sum(
                    1 for change in self.config_changes
                    if (log_entry["timestamp"] - change["timestamp"]) < self.time_window
                )
                
                if recent_changes > self.max_config_changes_per_hour:
                    anomalies.append({
                        "type": "late_night_config_changes",
                        "severity": "medium",
                        "message": f"Multiple config changes during late night: {recent_changes}",
                        "threshold": self.max_config_changes_per_hour
                    })
        
        # 如果有異常，記錄 alert
        if anomalies:
            for anomaly in anomalies:
                self._alert(anomaly)
        
        return anomalies
    
    def _alert(self, anomaly: Dict):
        """發出 alert（簡化：先 log）"""
        print(f"[ALERT] {anomaly['severity'].upper()}: {anomaly['message']}")
        # TODO: 發送到 SIEM 或通知系統
    
    def get_recent_logs(self, limit: int = 100) -> List[Dict]:
        """取得最近的 logs"""
        return self.logs[-limit:]
    
    def get_anomaly_summary(self) -> Dict:
        """取得異常摘要"""
        return {
            "total_logs": len(self.logs),
            "total_commands": sum(self.command_counts.values()),
            "total_config_changes": len(self.config_changes),
            "recent_anomalies": len([a for log in self.logs[-100:] if log.get("anomalies")])
        }

