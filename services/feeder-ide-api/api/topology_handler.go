package api

import (
	"net/http"
	"time"

	"github.com/feeder-platform/feeder-ide-api/internal/auth"
	"github.com/feeder-platform/feeder-ide-api/internal/middleware"
	"github.com/feeder-platform/feeder-ide-api/internal/topology"
	"github.com/feeder-platform/feeder-ide-api/internal/user"
	"github.com/gin-gonic/gin"
)

// TopologyHandler 處理拓樸相關的 HTTP 請求
type TopologyHandler struct {
	repo        topology.Repository
	userService *user.Service
}

// NewTopologyHandler 建立新的 TopologyHandler
func NewTopologyHandler(repo topology.Repository, userService *user.Service) *TopologyHandler {
	return &TopologyHandler{
		repo:        repo,
		userService: userService,
	}
}

// CreateTopologyRequest 建立拓樸的請求
type CreateTopologyRequest struct {
	Name        string                 `json:"name" binding:"required"`
	Description string                 `json:"description,omitempty"`
	ProfileType string                 `json:"profile_type" binding:"required,oneof=rural suburban urban"`
	Nodes       []topology.Node        `json:"nodes"`
	Lines       []topology.Line        `json:"lines"`
}

// CreateTopology 建立新拓樸
// @Summary 建立新拓樸
// @Description 建立一個新的配電 feeder 拓樸
// @Tags topologies
// @Accept json
// @Produce json
// @Param topology body CreateTopologyRequest true "拓樸資料"
// @Success 201 {object} topology.Topology
// @Failure 400 {object} map[string]string
// @Router /api/v1/topologies [post]
func (h *TopologyHandler) CreateTopology(c *gin.Context) {
	var req CreateTopologyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 檢查配額（如果 userService 可用）
	userID := auth.GetUserID(c)
	if h.userService != nil {
		canCreate, used, max, err := h.userService.CheckTopologyQuota(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check quota: " + err.Error()})
			return
		}
		if !canCreate {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Topology quota exceeded",
				"used": used,
				"max":  max,
			})
			return
		}
	}

	topo := &topology.Topology{
		UserID:      userID, // 設置用戶ID（如果已登入）
		Name:        req.Name,
		Description: req.Description,
		ProfileType: req.ProfileType,
		Nodes:       req.Nodes,
		Lines:       req.Lines,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.repo.Create(topo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, topo)
}

// GetTopology 取得拓樸
// @Summary 取得拓樸
// @Description 根據 ID 取得拓樸
// @Tags topologies
// @Produce json
// @Param id path string true "拓樸 ID"
// @Success 200 {object} topology.Topology
// @Failure 404 {object} map[string]string
// @Router /api/v1/topologies/{id} [get]
func (h *TopologyHandler) GetTopology(c *gin.Context) {
	id := c.Param("id")
	userID := auth.GetUserID(c)

	// 使用 GetByIDAndUserID 確保用戶只能訪問自己的拓樸（或 demo 拓樸）
	topo, err := h.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		if err == topology.ErrTopologyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topo)
}

// UpdateTopologyRequest 更新拓樸的請求
type UpdateTopologyRequest struct {
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	ProfileType string                 `json:"profile_type,omitempty" binding:"omitempty,oneof=rural suburban urban"`
	Nodes       []topology.Node        `json:"nodes,omitempty"`
	Lines       []topology.Line        `json:"lines,omitempty"`
}

// UpdateTopology 更新拓樸
// @Summary 更新拓樸
// @Description 更新現有拓樸
// @Tags topologies
// @Accept json
// @Produce json
// @Param id path string true "拓樸 ID"
// @Param topology body UpdateTopologyRequest true "拓樸資料"
// @Success 200 {object} topology.Topology
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/topologies/{id} [put]
func (h *TopologyHandler) UpdateTopology(c *gin.Context) {
	id := c.Param("id")
	userID := auth.GetUserID(c)

	var req UpdateTopologyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 取得現有拓樸（檢查權限）
	existing, err := h.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		if err == topology.ErrTopologyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新欄位
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.ProfileType != "" {
		existing.ProfileType = req.ProfileType
	}
	if req.Nodes != nil {
		existing.Nodes = req.Nodes
	}
	if req.Lines != nil {
		existing.Lines = req.Lines
	}
	existing.UpdatedAt = time.Now()

	if err := h.repo.Update(id, existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

// DeleteTopology 刪除拓樸
// @Summary 刪除拓樸
// @Description 根據 ID 刪除拓樸
// @Tags topologies
// @Param id path string true "拓樸 ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/topologies/{id} [delete]
func (h *TopologyHandler) DeleteTopology(c *gin.Context) {
	id := c.Param("id")
	userID := auth.GetUserID(c)

	// 檢查權限：確保用戶只能刪除自己的拓樸
	existing, err := h.repo.GetByIDAndUserID(id, userID)
	if err != nil {
		if err == topology.ErrTopologyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 確保是該用戶的拓樸（demo 模式允許刪除無 userID 的拓樸）
	if userID != nil && existing.UserID != nil && *existing.UserID != *userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this topology"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListTopologies 列出所有拓樸
// @Summary 列出所有拓樸
// @Description 取得所有拓樸列表
// @Tags topologies
// @Produce json
// @Success 200 {array} topology.Topology
// @Router /api/v1/topologies [get]
func (h *TopologyHandler) ListTopologies(c *gin.Context) {
	userID := auth.GetUserID(c)

	// 根據用戶ID列出拓樸
	topologies, err := h.repo.ListByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topologies)
}

