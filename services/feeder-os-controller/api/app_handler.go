package api

import (
	"net/http"

	"github.com/feeder-platform/feeder-os-controller/internal/apps"
	"github.com/gin-gonic/gin"
)

// AppHandler 處理 app 相關的 HTTP 請求
type AppHandler struct {
	manager *apps.Manager
}

// NewAppHandler 建立新的 AppHandler
func NewAppHandler(manager *apps.Manager) *AppHandler {
	return &AppHandler{manager: manager}
}

// InstallAppRequest 安裝 app 的請求
type InstallAppRequest struct {
	Name    string            `json:"name" binding:"required"`
	Version string            `json:"version" binding:"required"`
	Image   string            `json:"image" binding:"required"`
	Config  map[string]string `json:"config"`
	Topics  []string          `json:"topics"`
}

// InstallApp 安裝 app
// @Summary 安裝 app
// @Description 安裝一個新的 app 到 Feeder OS
// @Tags apps
// @Accept json
// @Produce json
// @Param app body InstallAppRequest true "App 資料"
// @Success 201 {object} apps.App
// @Failure 400 {object} map[string]string
// @Router /api/v1/apps/install [post]
func (h *AppHandler) InstallApp(c *gin.Context) {
	var req InstallAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := h.manager.InstallApp(&apps.InstallAppRequest{
		Name:    req.Name,
		Version: req.Version,
		Image:   req.Image,
		Config:  req.Config,
		Topics:  req.Topics,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, app)
}

// EnableAppRequest 啟用 app 的請求
type EnableAppRequest struct {
	AppID string `json:"app_id" binding:"required"`
}

// EnableApp 啟用 app
// @Summary 啟用 app
// @Description 啟用一個已安裝的 app
// @Tags apps
// @Accept json
// @Produce json
// @Param request body EnableAppRequest true "App ID"
// @Success 200 {object} apps.App
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/apps/enable [post]
func (h *AppHandler) EnableApp(c *gin.Context) {
	var req EnableAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.manager.EnableApp(req.AppID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	app, err := h.manager.GetApp(req.AppID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

// DisableAppRequest 停用 app 的請求
type DisableAppRequest struct {
	AppID string `json:"app_id" binding:"required"`
}

// DisableApp 停用 app
// @Summary 停用 app
// @Description 停用一個已啟用的 app
// @Tags apps
// @Accept json
// @Produce json
// @Param request body DisableAppRequest true "App ID"
// @Success 200 {object} apps.App
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/apps/disable [post]
func (h *AppHandler) DisableApp(c *gin.Context) {
	var req DisableAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.manager.DisableApp(req.AppID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	app, err := h.manager.GetApp(req.AppID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

// GetApp 取得 app
// @Summary 取得 app
// @Description 根據 ID 取得 app 資訊
// @Tags apps
// @Produce json
// @Param id path string true "App ID"
// @Success 200 {object} apps.App
// @Failure 404 {object} map[string]string
// @Router /api/v1/apps/{id} [get]
func (h *AppHandler) GetApp(c *gin.Context) {
	appID := c.Param("id")

	app, err := h.manager.GetApp(appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, app)
}

// ListApps 列出所有 apps
// @Summary 列出所有 apps
// @Description 取得所有已安裝的 apps
// @Tags apps
// @Produce json
// @Success 200 {array} apps.App
// @Router /api/v1/apps [get]
func (h *AppHandler) ListApps(c *gin.Context) {
	apps, err := h.manager.ListApps()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, apps)
}

