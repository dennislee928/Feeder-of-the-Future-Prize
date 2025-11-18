package user

import (
	"fmt"
	"time"
)

// TopologyCounter 拓樸計數器介面（用於檢查配額）
type TopologyCounter interface {
	CountByUserID(userID *string) (int, error)
}

// Service 會員服務
type Service struct {
	repo            Repository
	topologyCounter TopologyCounter // 可選，用於檢查拓樸配額
}

// NewService 建立新的會員服務
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// SetTopologyCounter 設置拓樸計數器（用於檢查配額）
func (s *Service) SetTopologyCounter(counter TopologyCounter) {
	s.topologyCounter = counter
}

// GetUserTier 取得用戶等級
func (s *Service) GetUserTier(userID *string) (string, error) {
	if userID == nil {
		return "demo", nil
	}

	user, err := s.repo.GetUserByID(*userID)
	if err != nil {
		return "demo", err
	}

	return user.SubscriptionTier, nil
}

// GetUserQuota 取得或創建用戶配額
func (s *Service) GetUserQuota(userID string) (*UserQuota, error) {
	quota, err := s.repo.GetQuotaByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 如果沒有配額，創建默認配額
	if quota == nil {
		quota = s.getDefaultQuota(userID)
		if err := s.repo.CreateOrUpdateQuota(quota); err != nil {
			return nil, fmt.Errorf("failed to create quota: %w", err)
		}
	}

	// 檢查是否需要重置每日模擬計數
	now := time.Now()
	if quota.LastSimulationResetDate.Before(now.Truncate(24 * time.Hour)) {
		quota.UsedSimulationsToday = 0
		quota.LastSimulationResetDate = now.Truncate(24 * time.Hour)
		if err := s.repo.UpdateQuota(quota); err != nil {
			return nil, fmt.Errorf("failed to reset simulation count: %w", err)
		}
	}

	return quota, nil
}

// getDefaultQuota 根據用戶等級返回默認配額
func (s *Service) getDefaultQuota(userID string) *UserQuota {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		// 如果無法取得用戶，返回免費會員配額
		return &UserQuota{
			UserID:                  userID,
			MaxTopologies:           999999, // 無限
			UsedTopologies:          0,
			MaxSimulationsPerDay:    100,
			UsedSimulationsToday:    0,
			LastSimulationResetDate: time.Now().Truncate(24 * time.Hour),
			CanUse3DRendering:       true,
			CanUseAIPrediction:      true,
			CanUseAdvancedSecurity:   true,
			CanAccessAPI:             false,
		}
	}

	quota := &UserQuota{
		UserID:                  userID,
		UsedTopologies:          0,
		UsedSimulationsToday:    0,
		LastSimulationResetDate: time.Now().Truncate(24 * time.Hour),
	}

	switch user.SubscriptionTier {
	case "demo":
		quota.MaxTopologies = 3
		quota.MaxSimulationsPerDay = 10
		quota.CanUse3DRendering = false
		quota.CanUseAIPrediction = false
		quota.CanUseAdvancedSecurity = false
		quota.CanAccessAPI = false
	case "free":
		quota.MaxTopologies = 999999 // 無限
		quota.MaxSimulationsPerDay = 100
		quota.CanUse3DRendering = true
		quota.CanUseAIPrediction = true
		quota.CanUseAdvancedSecurity = true
		quota.CanAccessAPI = false
	case "premium":
		quota.MaxTopologies = 999999 // 無限
		quota.MaxSimulationsPerDay = 999999 // 無限
		quota.CanUse3DRendering = true
		quota.CanUseAIPrediction = true
		quota.CanUseAdvancedSecurity = true
		quota.CanAccessAPI = true
	default:
		// 默認使用免費會員配額
		quota.MaxTopologies = 999999
		quota.MaxSimulationsPerDay = 100
		quota.CanUse3DRendering = true
		quota.CanUseAIPrediction = true
		quota.CanUseAdvancedSecurity = true
		quota.CanAccessAPI = false
	}

	return quota
}

// CheckTopologyQuota 檢查拓樸配額
func (s *Service) CheckTopologyQuota(userID *string) (bool, int, int, error) {
	if userID == nil {
		// Demo 模式：最多 3 個拓樸
		if s.topologyCounter != nil {
			count, _ := s.topologyCounter.CountByUserID(nil)
			canCreate := count < 3
			return canCreate, count, 3, nil
		}
		return true, 0, 3, nil
	}

	quota, err := s.GetUserQuota(*userID)
	if err != nil {
		return false, 0, 0, err
	}

	// 更新已使用的拓樸數量
	var count int
	if s.topologyCounter != nil {
		count, err = s.topologyCounter.CountByUserID(userID)
		if err != nil {
			return false, 0, 0, err
		}
	}

	quota.UsedTopologies = count
	if err := s.repo.UpdateQuota(quota); err != nil {
		return false, 0, 0, err
	}

	canCreate := quota.UsedTopologies < quota.MaxTopologies
	return canCreate, quota.UsedTopologies, quota.MaxTopologies, nil
}

// CheckSimulationQuota 檢查模擬配額
func (s *Service) CheckSimulationQuota(userID *string) (bool, error) {
	if userID == nil {
		// Demo 模式：允許模擬
		return true, nil
	}

	quota, err := s.GetUserQuota(*userID)
	if err != nil {
		return false, err
	}

	canSimulate := quota.UsedSimulationsToday < quota.MaxSimulationsPerDay
	return canSimulate, nil
}

// IncrementSimulationCount 增加模擬計數
func (s *Service) IncrementSimulationCount(userID string) error {
	quota, err := s.GetUserQuota(userID)
	if err != nil {
		return err
	}

	quota.UsedSimulationsToday++
	if err := s.repo.UpdateQuota(quota); err != nil {
		return err
	}

	return nil
}

// CanUseFeature 檢查用戶是否可以使用特定功能
func (s *Service) CanUseFeature(userID *string, feature string) (bool, error) {
	if userID == nil {
		// Demo 模式功能限制
		switch feature {
		case "3d_rendering", "ai_prediction", "advanced_security", "api_access":
			return false, nil
		default:
			return true, nil
		}
	}

	quota, err := s.GetUserQuota(*userID)
	if err != nil {
		return false, err
	}

	switch feature {
	case "3d_rendering":
		return quota.CanUse3DRendering, nil
	case "ai_prediction":
		return quota.CanUseAIPrediction, nil
	case "advanced_security":
		return quota.CanUseAdvancedSecurity, nil
	case "api_access":
		return quota.CanAccessAPI, nil
	default:
		return true, nil
	}
}

// UpdateUserTier 更新用戶等級（用於付費後）
func (s *Service) UpdateUserTier(userID string, tier string) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.SubscriptionTier = tier
	user.SubscriptionStatus = "active"

	// 更新配額
	quota := s.getDefaultQuota(userID)
	if err := s.repo.CreateOrUpdateQuota(quota); err != nil {
		return fmt.Errorf("failed to update quota: %w", err)
	}

	return s.repo.UpdateUser(user)
}

// CountByUserID 統計用戶拓樸數量（需要從 topology repository 調用）
// 這個方法需要 topology repository，所以我們在 service 中不直接實現
// 而是在需要時通過參數傳入

