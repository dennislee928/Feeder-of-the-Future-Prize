"""
Dummy App Example
訂閱 measurements topic，log 資料即可
"""
import json
import os
import sys
import socket
import time
import paho.mqtt.client as mqtt

# 從環境變數取得配置
MQTT_BROKER = os.getenv("MQTT_BROKER", "localhost")
MQTT_PORT = int(os.getenv("MQTT_PORT", "1883"))
FEEDER_ID = os.getenv("FEEDER_ID", "feeder-001")
CLIENT_ID = os.getenv("CLIENT_ID", "dummy-app")

# 訂閱的 topic
MEASUREMENTS_TOPIC = f"feeder/{FEEDER_ID}/measurements/+"


def on_connect(client, userdata, flags, rc):
    """連線回調"""
    if rc == 0:
        print(f"Connected to MQTT broker at {MQTT_BROKER}:{MQTT_PORT}")
        # 訂閱 measurements topic
        client.subscribe(MEASUREMENTS_TOPIC)
        print(f"Subscribed to: {MEASUREMENTS_TOPIC}")
    else:
        print(f"Failed to connect, return code {rc}")


def on_message(client, userdata, msg):
    """訊息接收回調"""
    try:
        payload = json.loads(msg.payload.decode())
        print(f"[{time.strftime('%Y-%m-%d %H:%M:%S')}] Received on {msg.topic}:")
        print(f"  Asset ID: {payload.get('asset_id', 'unknown')}")
        print(f"  Timestamp: {payload.get('timestamp', 'unknown')}")
        print(f"  Data: {json.dumps(payload.get('data', {}), indent=2)}")
        print("-" * 50)
    except json.JSONDecodeError:
        print(f"Failed to parse JSON: {msg.payload.decode()}")
    except Exception as e:
        print(f"Error processing message: {e}")


def main():
    """主程式"""
    print("Starting Dummy App...")
    print(f"Configuration:")
    print(f"  MQTT Broker: {MQTT_BROKER}:{MQTT_PORT}")
    print(f"  Feeder ID: {FEEDER_ID}")
    print(f"  Client ID: {CLIENT_ID}")
    print(f"  Topic: {MEASUREMENTS_TOPIC}")

    # 建立 MQTT client
    client = mqtt.Client(client_id=CLIENT_ID)
    client.on_connect = on_connect
    client.on_message = on_message

    # 連線到 broker（重試邏輯）
    max_retries = 10
    retry_delay = 2
    connected = False
    
    for attempt in range(max_retries):
        try:
            print(f"Attempting to connect to MQTT broker (attempt {attempt + 1}/{max_retries})...")
            client.connect(MQTT_BROKER, MQTT_PORT, 60)
            print(f"Successfully connected to MQTT broker!")
            connected = True
            break
        except socket.gaierror as e:
            if attempt < max_retries - 1:
                print(f"DNS resolution failed: {e}")
                print(f"Retrying in {retry_delay} seconds...")
                time.sleep(retry_delay)
            else:
                print(f"Failed to connect after {max_retries} attempts. Exiting...")
                sys.exit(1)
        except Exception as e:
            if attempt < max_retries - 1:
                print(f"Connection failed: {e}")
                print(f"Retrying in {retry_delay} seconds...")
                time.sleep(retry_delay)
            else:
                print(f"Failed to connect after {max_retries} attempts: {e}")
                sys.exit(1)
    
    if connected:
        try:
            client.loop_forever()
        except KeyboardInterrupt:
            print("\nShutting down...")
            client.disconnect()
        except Exception as e:
            print(f"Error during operation: {e}")


if __name__ == "__main__":
    main()

