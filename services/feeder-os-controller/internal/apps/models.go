package apps

import (
	"fmt"
	"time"
)

// App 代表一個安裝的 app
type App struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Status      AppStatus         `json:"status"`
	Image       string            `json:"image"`       // Docker image
	Config      map[string]string `json:"config"`       // 環境變數配置
	Topics      []string          `json:"topics"`      // 訂閱的 topics
	InstalledAt time.Time         `json:"installed_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// AppStatus app 狀態
type AppStatus string

const (
	StatusInstalled AppStatus = "installed"
	StatusEnabled   AppStatus = "enabled"
	StatusDisabled  AppStatus = "disabled"
	StatusError     AppStatus = "error"
)

// AppContract 定義 app 與 Feeder OS 的溝通介面
type AppContract struct {
	// Topics 命名規則
	// feeder/<feeder_id>/measurements/<asset_id> - 測量資料
	// feeder/<feeder_id>/commands/<asset_id> - 控制命令
	// feeder/<feeder_id>/events/<severity> - 事件通知
	
	// Config 方式
	// 1. 環境變數：透過 Config map 傳遞
	// 2. Config volume：未來可支援
	// 3. gRPC：未來可支援
}

// TopicNaming 定義 topic 命名規則
func TopicNaming(feederID, category, resourceID string) string {
	return fmt.Sprintf("feeder/%s/%s/%s", feederID, category, resourceID)
}

// MeasurementsTopic 建立 measurements topic
func MeasurementsTopic(feederID, assetID string) string {
	return TopicNaming(feederID, "measurements", assetID)
}

// CommandsTopic 建立 commands topic
func CommandsTopic(feederID, assetID string) string {
	return TopicNaming(feederID, "commands", assetID)
}

// EventsTopic 建立 events topic
func EventsTopic(feederID, severity string) string {
	return TopicNaming(feederID, "events", severity)
}

