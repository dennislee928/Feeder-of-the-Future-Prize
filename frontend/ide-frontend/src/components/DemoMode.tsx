import React, { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from '../contexts/AuthContext'
import LoginModal from './Auth/LoginModal'
import './DemoMode.css'

interface DemoModeProps {
  onLoginClick?: () => void
}

const DemoMode: React.FC<DemoModeProps> = ({ onLoginClick }) => {
  const { t } = useTranslation()
  const { isAuthenticated } = useAuth()
  const [showLoginModal, setShowLoginModal] = useState(false)

  if (isAuthenticated) {
    return null // 已登入，不顯示 demo 提示
  }

  return (
    <>
      <div className="demo-mode-banner">
        <div className="demo-mode-content">
          <span className="demo-mode-icon">ℹ️</span>
          <span className="demo-mode-text">
            {t('demo.banner.text')}
          </span>
          <button
            className="demo-mode-login-button"
            onClick={() => {
              if (onLoginClick) {
                onLoginClick()
              } else {
                setShowLoginModal(true)
              }
            }}
          >
            {t('demo.banner.login')}
          </button>
        </div>
      </div>
      <LoginModal isOpen={showLoginModal} onClose={() => setShowLoginModal(false)} />
    </>
  )
}

export default DemoMode

