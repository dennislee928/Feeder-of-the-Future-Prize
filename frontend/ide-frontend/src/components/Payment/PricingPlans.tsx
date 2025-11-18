import React from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../../contexts/AuthContext'
import './PricingPlans.css'

interface PricingPlansProps {
  onSelectPlan: (tier: 'premium', provider: 'stripe' | 'paypal') => void
}

const PricingPlans: React.FC<PricingPlansProps> = ({ onSelectPlan }) => {
  const { t } = useTranslation()
  const { user } = useAuth()

  const plans = [
    {
      tier: 'free' as const,
      name: t('payment.plans.free.name'),
      price: t('payment.plans.free.price'),
      features: [
        t('payment.plans.free.features.unlimited_topologies'),
        t('payment.plans.free.features.all_simulations'),
        t('payment.plans.free.features.basic_3d'),
        t('payment.plans.free.features.basic_ai'),
        t('payment.plans.free.features.advanced_security'),
      ],
      current: user?.subscription_tier === 'free',
    },
    {
      tier: 'premium' as const,
      name: t('payment.plans.premium.name'),
      price: t('payment.plans.premium.price'),
      features: [
        t('payment.plans.premium.features.all_free'),
        t('payment.plans.premium.features.api_access'),
        t('payment.plans.premium.features.priority_support'),
        t('payment.plans.premium.features.advanced_3d'),
        t('payment.plans.premium.features.advanced_ai'),
        t('payment.plans.premium.features.bulk_export'),
        t('payment.plans.premium.features.collaboration'),
      ],
      current: user?.subscription_tier === 'premium',
      popular: true,
    },
  ]

  return (
    <div className="pricing-plans">
      <h2 className="pricing-plans-title">{t('payment.plans.title')}</h2>
      <p className="pricing-plans-description">{t('payment.plans.description')}</p>
      <div className="pricing-plans-grid">
        {plans.map((plan) => (
          <div
            key={plan.tier}
            className={`pricing-plan ${plan.current ? 'current' : ''} ${plan.popular ? 'popular' : ''}`}
          >
            {plan.popular && <div className="pricing-plan-badge">{t('payment.plans.popular')}</div>}
            <div className="pricing-plan-header">
              <h3 className="pricing-plan-name">{plan.name}</h3>
              <div className="pricing-plan-price">{plan.price}</div>
            </div>
            <ul className="pricing-plan-features">
              {plan.features.map((feature, index) => (
                <li key={index} className="pricing-plan-feature">
                  <svg className="feature-icon" viewBox="0 0 20 20" fill="currentColor">
                    <path
                      fillRule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clipRule="evenodd"
                    />
                  </svg>
                  {feature}
                </li>
              ))}
            </ul>
            {plan.current ? (
              <div className="pricing-plan-current">{t('payment.plans.current_plan')}</div>
            ) : plan.tier === 'premium' ? (
              <div className="pricing-plan-actions">
                <button
                  className="pricing-plan-button pricing-plan-button-stripe"
                  onClick={() => onSelectPlan('premium', 'stripe')}
                >
                  {t('payment.plans.pay_with_stripe')}
                </button>
                <button
                  className="pricing-plan-button pricing-plan-button-paypal"
                  onClick={() => onSelectPlan('premium', 'paypal')}
                >
                  {t('payment.plans.pay_with_paypal')}
                </button>
              </div>
            ) : null}
          </div>
        ))}
      </div>
    </div>
  )
}

export default PricingPlans

