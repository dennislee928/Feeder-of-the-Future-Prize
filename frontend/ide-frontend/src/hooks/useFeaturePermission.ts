import { useAuth } from '../contexts/AuthContext'

export type Feature = '3d_rendering' | 'ai_prediction' | 'advanced_security' | 'api_access'

export const useFeaturePermission = () => {
  const { user, quota } = useAuth()

  const canUseFeature = (feature: Feature): boolean => {
    // Demo 模式（未登入）功能限制
    if (!user) {
      switch (feature) {
        case '3d_rendering':
        case 'ai_prediction':
        case 'advanced_security':
        case 'api_access':
          return false
        default:
          return true
      }
    }

    // 已登入用戶，檢查配額
    if (!quota) {
      // 如果沒有配額資訊，根據用戶等級判斷
      if (user?.subscription_tier === 'premium') {
        return true
      }
      if (user?.subscription_tier === 'free') {
        return feature !== 'api_access'
      }
      return false
    }

    // 根據配額檢查
    switch (feature) {
      case '3d_rendering':
        return quota.can_use_3d_rendering
      case 'ai_prediction':
        return quota.can_use_ai_prediction
      case 'advanced_security':
        return quota.can_use_advanced_security
      case 'api_access':
        return quota.can_access_api
      default:
        return true
    }
  }

  const canCreateTopology = (): { canCreate: boolean; used: number; max: number } => {
    if (!user) {
      // Demo 模式：最多 3 個拓樸
      return { canCreate: true, used: 0, max: 3 }
    }

    if (!quota) {
      return { canCreate: true, used: 0, max: 999999 }
    }

    const canCreate = quota.used_topologies < quota.max_topologies
    return {
      canCreate,
      used: quota.used_topologies,
      max: quota.max_topologies === 999999 ? Infinity : quota.max_topologies,
    }
  }

  const canRunSimulation = (): boolean => {
    if (!user) {
      // Demo 模式：允許模擬
      return true
    }

    if (!quota) {
      return true
    }

    return quota.used_simulations_today < quota.max_simulations_per_day
  }

  return {
    canUseFeature,
    canCreateTopology,
    canRunSimulation,
    userTier: user?.subscription_tier || 'demo',
    quota,
  }
}

