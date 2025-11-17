"""
Control Loop
定時執行控制邏輯
"""
import asyncio
from typing import Dict
from src.registry import AssetRegistry
from src.mqtt_client import MQTTClient


class ControlLoop:
    """控制迴圈"""
    def __init__(self, registry: AssetRegistry, mqtt_client: MQTTClient, feeder_id: str):
        self.registry = registry
        self.mqtt_client = mqtt_client
        self.feeder_id = feeder_id
        self.interval_seconds = 300  # 5 分鐘
        self.loading_threshold = 0.85  # 85% 載流率
        self.measurements: Dict[str, Dict] = {}
        self.running = False
        
        # 訂閱所有 measurements（如果 MQTT client 可用）
        if mqtt_client and mqtt_client.client:
            measurements_topic = f"feeder/{feeder_id}/measurements/+"
            mqtt_client.subscribe(measurements_topic, self._on_measurement)
    
    def _on_measurement(self, topic: str, payload: Dict):
        """處理接收到的測量資料"""
        asset_id = payload.get("asset_id")
        if asset_id:
            self.measurements[asset_id] = payload
    
    async def run(self):
        """執行控制迴圈"""
        self.running = True
        print("Control loop started")
        
        while self.running:
            try:
                await self._control_step()
            except Exception as e:
                print(f"Error in control step: {e}")
            
            await asyncio.sleep(self.interval_seconds)
    
    async def _control_step(self):
        """執行一個控制步驟"""
        # 計算 feeder 總負載
        total_loading = self._calculate_feeder_loading()
        
        print(f"Feeder loading: {total_loading:.2%}")
        
        # 如果負載超過閾值，執行控制
        if total_loading > self.loading_threshold:
            print(f"Loading exceeds threshold ({self.loading_threshold:.2%}), applying control...")
            self._reduce_ev_charging()
        else:
            # 負載正常，可以恢復 EV 充電
            self._restore_ev_charging()
    
    def _calculate_feeder_loading(self) -> float:
        """計算 feeder 載流率（簡化版本）"""
        # 從 measurements 取得負載資料
        total_power = 0.0
        for asset_id, measurement in self.measurements.items():
            data = measurement.get("data", {})
            power = data.get("power_kw", 0.0)
            total_power += power
        
        # 假設 feeder 額定容量為 1000 kW（之後可以從 topology 取得）
        rated_capacity = 1000.0
        loading = total_power / rated_capacity if rated_capacity > 0 else 0.0
        
        return min(loading, 1.0)  # 限制在 0-1 之間
    
    def _reduce_ev_charging(self):
        """降低 EV 充電功率"""
        ev_chargers = self.registry.list_assets("ev_charger")
        
        # 計算需要降低的功率
        total_ev_power = sum(c.current_power_kw for c in ev_chargers)
        reduction_factor = 0.5  # 降低 50%
        
        for charger in ev_chargers:
            new_power = charger.current_power_kw * reduction_factor
            self.registry.update_asset_power(charger.id, new_power)
            
            # 發送控制命令（如果 MQTT client 可用）
            if self.mqtt_client and self.mqtt_client.client:
                command = {
                    "asset_id": charger.id,
                    "command": "set_power",
                    "power_kw": new_power,
                    "timestamp": asyncio.get_event_loop().time()
                }
                
                topic = self.mqtt_client.commands_topic(charger.id)
                self.mqtt_client.publish(topic, command)
                print(f"Reduced EV charger {charger.id} power to {new_power:.2f} kW")
            else:
                print(f"MQTT not available, would reduce EV charger {charger.id} power to {new_power:.2f} kW")
    
    def _restore_ev_charging(self):
        """恢復 EV 充電功率"""
        ev_chargers = self.registry.list_assets("ev_charger")
        
        for charger in ev_chargers:
            # 恢復到最大功率的 80%（避免立即過載）
            target_power = charger.max_power_kw * 0.8
            
            if charger.current_power_kw < target_power:
                self.registry.update_asset_power(charger.id, target_power)
                
                if self.mqtt_client and self.mqtt_client.client:
                    command = {
                        "asset_id": charger.id,
                        "command": "set_power",
                        "power_kw": target_power,
                        "timestamp": asyncio.get_event_loop().time()
                    }
                    
                    topic = self.mqtt_client.commands_topic(charger.id)
                    self.mqtt_client.publish(topic, command)
                    print(f"Restored EV charger {charger.id} power to {target_power:.2f} kW")
                else:
                    print(f"MQTT not available, would restore EV charger {charger.id} power to {target_power:.2f} kW")

