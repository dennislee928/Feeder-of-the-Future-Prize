import React, { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { paymentApi, Payment } from '../../api/paymentApi'
import './PaymentHistory.css'

const PaymentHistory: React.FC = () => {
  const { t } = useTranslation()
  const [payments, setPayments] = useState<Payment[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadPayments()
  }, [])

  const loadPayments = async () => {
    setIsLoading(true)
    setError(null)
    try {
      const response = await paymentApi.getPaymentHistory()
      setPayments(response.payments)
    } catch (err: any) {
      console.error('Failed to load payment history:', err)
      setError(err.response?.data?.error || t('payment.history.load_failed'))
    } finally {
      setIsLoading(false)
    }
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return '#34a853'
      case 'pending':
        return '#fbbc04'
      case 'failed':
        return '#ea4335'
      case 'refunded':
        return '#9aa0a6'
      default:
        return '#666'
    }
  }

  if (isLoading) {
    return <div className="payment-history-loading">{t('payment.history.loading')}</div>
  }

  if (error) {
    return <div className="payment-history-error">{error}</div>
  }

  if (payments.length === 0) {
    return <div className="payment-history-empty">{t('payment.history.empty')}</div>
  }

  return (
    <div className="payment-history">
      <h2 className="payment-history-title">{t('payment.history.title')}</h2>
      <div className="payment-history-table">
        <table>
          <thead>
            <tr>
              <th>{t('payment.history.date')}</th>
              <th>{t('payment.history.amount')}</th>
              <th>{t('payment.history.provider')}</th>
              <th>{t('payment.history.status')}</th>
            </tr>
          </thead>
          <tbody>
            {payments.map((payment) => (
              <tr key={payment.id}>
                <td>{new Date(payment.created_at).toLocaleDateString()}</td>
                <td>
                  {payment.amount.toFixed(2)} {payment.currency.toUpperCase()}
                </td>
                <td>{payment.payment_provider.toUpperCase()}</td>
                <td>
                  <span
                    className="payment-status"
                    style={{ color: getStatusColor(payment.status) }}
                  >
                    {t(`payment.status.${payment.status}`)}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default PaymentHistory

