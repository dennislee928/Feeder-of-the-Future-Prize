import React from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../../contexts/AuthContext'
import './UserProfile.css'

const UserProfile: React.FC = () => {
  const { t } = useTranslation()
  const { user, subscription, quota } = useAuth()

  if (!user) {
    return null
  }

  const getTierDisplayName = (tier: string) => {
    switch (tier) {
      case 'demo':
        return t('auth.tier.demo')
      case 'free':
        return t('auth.tier.free')
      case 'premium':
        return t('auth.tier.premium')
      default:
        return tier
    }
  }

  return (
    <div className="user-profile">
      <div className="user-profile-header">
        {user.avatar_url && (
          <img src={user.avatar_url} alt={user.name || user.email} className="user-profile-avatar" />
        )}
        <div className="user-profile-info">
          <div className="user-profile-name">{user.name || user.email}</div>
          <div className="user-profile-email">{user.email}</div>
          <div className="user-profile-tier">
            {t('auth.tier.label')}: {getTierDisplayName(user.subscription_tier)}
          </div>
        </div>
      </div>
      {quota && (
        <div className="user-profile-quota">
          <div className="quota-item">
            <span className="quota-label">{t('auth.quota.topologies')}:</span>
            <span className="quota-value">
              {quota.used_topologies} / {quota.max_topologies === 999999 ? '∞' : quota.max_topologies}
            </span>
          </div>
          <div className="quota-item">
            <span className="quota-label">{t('auth.quota.simulations')}:</span>
            <span className="quota-value">
              {quota.used_simulations_today} / {quota.max_simulations_per_day === 999999 ? '∞' : quota.max_simulations_per_day}
            </span>
          </div>
        </div>
      )}
      {subscription && subscription.status === 'active' && (
        <div className="user-profile-subscription">
          {subscription.current_period_end && (
            <div className="subscription-expiry">
              {t('auth.subscription.expires')}: {new Date(subscription.current_period_end).toLocaleDateString()}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

export default UserProfile

