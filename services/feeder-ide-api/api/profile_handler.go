package api

import (
	"net/http"

	"github.com/feeder-platform/feeder-ide-api/internal/profiles"
	"github.com/gin-gonic/gin"
)

// ProfileHandler 處理 profile 相關的 HTTP 請求
type ProfileHandler struct {
	repo profiles.Repository
}

// NewProfileHandler 建立新的 ProfileHandler
func NewProfileHandler(repo profiles.Repository) *ProfileHandler {
	return &ProfileHandler{repo: repo}
}

// ListProfiles 列出所有 profiles
// @Summary 列出所有 profiles
// @Description 取得所有可用的 feeder profiles（Rural/Suburban/Urban）
// @Tags profiles
// @Produce json
// @Success 200 {array} profiles.Profile
// @Router /api/v1/profiles [get]
func (h *ProfileHandler) ListProfiles(c *gin.Context) {
	profileList, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profileList)
}

// GetProfile 取得特定 profile
// @Summary 取得特定 profile
// @Description 根據類型取得 profile（rural/suburban/urban）
// @Tags profiles
// @Produce json
// @Param type path string true "Profile 類型" Enums(rural, suburban, urban)
// @Success 200 {object} profiles.Profile
// @Failure 404 {object} map[string]string
// @Router /api/v1/profiles/{type} [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	profileType := c.Param("type")

	profile, err := h.repo.GetByType(profileType)
	if err != nil {
		if err == profiles.ErrProfileNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

