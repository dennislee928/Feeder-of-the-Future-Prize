"""
MQTT Client
處理與 message bus 的連接
"""
import json
import asyncio
from typing import Callable, Optional, Dict
import paho.mqtt.client as mqtt


class MQTTClient:
    """MQTT 客戶端"""
    def __init__(self, broker: str, port: int, feeder_id: str):
        self.broker = broker
        self.port = port
        self.feeder_id = feeder_id
        self.client: Optional[mqtt.Client] = None
        self.subscriptions: Dict[str, Callable] = {}
    
    async def connect(self):
        """連接到 MQTT broker"""
        self.client = mqtt.Client(client_id=f"der-ev-orchestrator-{self.feeder_id}")
        
        def on_connect(client, userdata, flags, rc):
            if rc == 0:
                print(f"MQTT connected to {self.broker}:{self.port}")
            else:
                print(f"MQTT connection failed: {rc}")
        
        def on_message(client, userdata, msg):
            topic = msg.topic
            if topic in self.subscriptions:
                try:
                    payload = json.loads(msg.payload.decode())
                    self.subscriptions[topic](topic, payload)
                except Exception as e:
                    print(f"Error processing message: {e}")
        
        self.client.on_connect = on_connect
        self.client.on_message = on_message
        
        # 使用同步連接（paho-mqtt 是同步的）
        try:
            result = self.client.connect(self.broker, self.port, 60)
            if result != mqtt.MQTT_ERR_SUCCESS:
                raise Exception(f"MQTT connection failed with code: {result}")
            self.client.loop_start()
        except Exception as e:
            raise Exception(f"Failed to connect to MQTT broker {self.broker}:{self.port}: {e}")
    
    async def disconnect(self):
        """斷開連接"""
        if self.client:
            self.client.loop_stop()
            self.client.disconnect()
    
    def subscribe(self, topic: str, callback: Callable):
        """訂閱 topic"""
        if self.client:
            self.client.subscribe(topic, 0)
            self.subscriptions[topic] = callback
    
    def publish(self, topic: str, payload: dict):
        """發布訊息"""
        if self.client:
            self.client.publish(topic, json.dumps(payload), 0)
    
    def measurements_topic(self, asset_id: str) -> str:
        """建立 measurements topic"""
        return f"feeder/{self.feeder_id}/measurements/{asset_id}"
    
    def commands_topic(self, asset_id: str) -> str:
        """建立 commands topic"""
        return f"feeder/{self.feeder_id}/commands/{asset_id}"

