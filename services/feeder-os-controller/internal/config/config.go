package config

import (
	"os"
)

// Config 應用程式配置
type Config struct {
	MQTT MQTTConfig
	Apps AppsConfig
}

// MQTTConfig MQTT broker 配置
type MQTTConfig struct {
	Broker   string
	Port     int
	ClientID string
	Username string
	Password string
}

// AppsConfig App 管理配置
type AppsConfig struct {
	StoragePath string
}

// Load 載入配置（從環境變數）
func Load() *Config {
	return &Config{
		MQTT: MQTTConfig{
			Broker:   getEnv("MQTT_BROKER", "localhost"),
			Port:     1883,
			ClientID: getEnv("MQTT_CLIENT_ID", "feeder-os-controller"),
			Username: getEnv("MQTT_USERNAME", ""),
			Password: getEnv("MQTT_PASSWORD", ""),
		},
		Apps: AppsConfig{
			StoragePath: getEnv("APPS_STORAGE_PATH", "./apps"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

