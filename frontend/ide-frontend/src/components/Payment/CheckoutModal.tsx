import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { paymentApi } from '../../api/paymentApi'
import './CheckoutModal.css'

interface CheckoutModalProps {
  isOpen: boolean
  onClose: () => void
  tier: 'premium'
  provider: 'stripe' | 'paypal'
}

const CheckoutModal: React.FC<CheckoutModalProps> = ({ isOpen, onClose, tier, provider }) => {
  const { t } = useTranslation()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleCheckout = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await paymentApi.createCheckout({ tier, provider })

      if (response.checkout_url) {
        // Stripe: 重定向到 checkout URL
        window.location.href = response.checkout_url
      } else if (response.subscription_id) {
        // PayPal: 可能需要不同的處理方式
        // 這裡簡化處理，實際應該根據 PayPal 的流程處理
        alert(t('payment.checkout.paypal_redirect'))
        onClose()
      } else {
        setError(t('payment.checkout.failed'))
      }
    } catch (err: any) {
      console.error('Checkout failed:', err)
      setError(err.response?.data?.error || t('payment.checkout.failed'))
      setIsLoading(false)
    }
  }

  if (!isOpen) return null

  return (
    <div className="checkout-modal-overlay" onClick={onClose}>
      <div className="checkout-modal" onClick={(e) => e.stopPropagation()}>
        <div className="checkout-modal-header">
          <h2>{t('payment.checkout.title')}</h2>
          <button className="checkout-modal-close" onClick={onClose}>
            ×
          </button>
        </div>
        <div className="checkout-modal-body">
          <div className="checkout-summary">
            <div className="checkout-summary-item">
              <span className="checkout-summary-label">{t('payment.checkout.plan')}:</span>
              <span className="checkout-summary-value">
                {t(`payment.plans.${tier}.name`)}
              </span>
            </div>
            <div className="checkout-summary-item">
              <span className="checkout-summary-label">{t('payment.checkout.provider')}:</span>
              <span className="checkout-summary-value">
                {provider === 'stripe' ? 'Stripe' : 'PayPal'}
              </span>
            </div>
          </div>
          {error && <div className="checkout-error">{error}</div>}
          <div className="checkout-actions">
            <button
              className="checkout-button checkout-button-cancel"
              onClick={onClose}
              disabled={isLoading}
            >
              {t('payment.checkout.cancel')}
            </button>
            <button
              className="checkout-button checkout-button-confirm"
              onClick={handleCheckout}
              disabled={isLoading}
            >
              {isLoading ? t('payment.checkout.processing') : t('payment.checkout.confirm')}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default CheckoutModal

