package bus

import (
	"fmt"
	"log"
	"time"

	"github.com/feeder-platform/feeder-os-controller/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Bus 定義 message bus 介面
type Bus interface {
	Publish(topic string, payload []byte) error
	Subscribe(topic string, handler MessageHandler) error
	Unsubscribe(topic string) error
	Close() error
}

// MessageHandler 處理接收到的訊息
type MessageHandler func(topic string, payload []byte)

// MQTTBus MQTT 實作
type MQTTBus struct {
	client mqtt.Client
	cfg    config.MQTTConfig
}

// NewMQTTBus 建立新的 MQTT bus
func NewMQTTBus(cfg config.MQTTConfig) (*MQTTBus, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.Broker, cfg.Port))
	opts.SetClientID(cfg.ClientID)
	
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		log.Printf("Unexpected message on topic: %s", msg.Topic())
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Println("Connected to MQTT broker")
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
	})

	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	return &MQTTBus{
		client: client,
		cfg:    cfg,
	}, nil
}

// Publish 發布訊息
func (b *MQTTBus) Publish(topic string, payload []byte) error {
	token := b.client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to publish: %w", token.Error())
	}
	return nil
}

// Subscribe 訂閱主題
func (b *MQTTBus) Subscribe(topic string, handler MessageHandler) error {
	token := b.client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Topic(), msg.Payload())
	})
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe: %w", token.Error())
	}
	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// Unsubscribe 取消訂閱
func (b *MQTTBus) Unsubscribe(topic string) error {
	token := b.client.Unsubscribe(topic)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe: %w", token.Error())
	}
	log.Printf("Unsubscribed from topic: %s", topic)
	return nil
}

// Close 關閉連接
func (b *MQTTBus) Close() error {
	b.client.Disconnect(250)
	return nil
}

