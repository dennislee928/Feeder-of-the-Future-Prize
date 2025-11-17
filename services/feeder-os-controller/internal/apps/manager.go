package apps

import (
	"fmt"
	"sync"
	"time"

	"github.com/feeder-platform/feeder-os-controller/internal/bus"
	"github.com/feeder-platform/feeder-os-controller/internal/config"
	"github.com/google/uuid"
)

// Manager 管理 app 的生命週期
type Manager struct {
	mu    sync.RWMutex
	apps  map[string]*App
	bus   bus.Bus
	cfg   *config.Config
}

// NewManager 建立新的 app manager
func NewManager(bus bus.Bus, cfg *config.Config) *Manager {
	return &Manager{
		apps: make(map[string]*App),
		bus:  bus,
		cfg:  cfg,
	}
}

// InstallApp 安裝 app
func (m *Manager) InstallApp(req *InstallAppRequest) (*App, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	app := &App{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Version:     req.Version,
		Status:      StatusInstalled,
		Image:       req.Image,
		Config:      req.Config,
		Topics:      req.Topics,
		InstalledAt: time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.apps[app.ID] = app
	return app, nil
}

// EnableApp 啟用 app
func (m *Manager) EnableApp(appID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	app, exists := m.apps[appID]
	if !exists {
		return fmt.Errorf("app not found: %s", appID)
	}

	if app.Status == StatusEnabled {
		return nil // 已經啟用
	}

	app.Status = StatusEnabled
	app.UpdatedAt = time.Now()

	// TODO: 實際啟動 app container 或 process
	// 這裡先只更新狀態

	return nil
}

// DisableApp 停用 app
func (m *Manager) DisableApp(appID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	app, exists := m.apps[appID]
	if !exists {
		return fmt.Errorf("app not found: %s", appID)
	}

	if app.Status == StatusDisabled {
		return nil // 已經停用
	}

	app.Status = StatusDisabled
	app.UpdatedAt = time.Now()

	// TODO: 實際停止 app container 或 process

	return nil
}

// GetApp 取得 app
func (m *Manager) GetApp(appID string) (*App, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	app, exists := m.apps[appID]
	if !exists {
		return nil, fmt.Errorf("app not found: %s", appID)
	}

	return app, nil
}

// ListApps 列出所有 apps
func (m *Manager) ListApps() ([]*App, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	apps := make([]*App, 0, len(m.apps))
	for _, app := range m.apps {
		apps = append(apps, app)
	}

	return apps, nil
}

// InstallAppRequest 安裝 app 的請求
type InstallAppRequest struct {
	Name    string            `json:"name" binding:"required"`
	Version string            `json:"version" binding:"required"`
	Image   string            `json:"image" binding:"required"`
	Config  map[string]string `json:"config"`
	Topics  []string          `json:"topics"`
}

