import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { useAuth } from './contexts/AuthContext'
import TopologyCanvas from './components/TopologyCanvas'
import SecurityTestCanvas from './components/SecurityTestCanvas'
import ESGSimulationCanvas from './components/ESGSimulationCanvas'
import Palette from './components/Palette'
import PropertiesPanel from './components/PropertiesPanel'
import LanguageSwitcher from './components/LanguageSwitcher'
import LoginModal from './components/Auth/LoginModal'
import UserProfile from './components/Auth/UserProfile'
import PricingPlans from './components/Payment/PricingPlans'
import CheckoutModal from './components/Payment/CheckoutModal'
import PaymentHistory from './components/Payment/PaymentHistory'
import DemoMode from './components/DemoMode'
import { PowerflowResult } from './api/simApi'
import './App.css'

type AppMode = 'design' | 'security' | 'esg' | 'pricing'

function App() {
  const { t } = useTranslation()
  const { isAuthenticated, logout } = useAuth()
  const [mode, setMode] = useState<AppMode>('design')
  const [selectedNode, setSelectedNode] = useState<string | null>(null)
  const [simulationResult, setSimulationResult] = useState<PowerflowResult | null>(null)
  const [currentTopologyId, setCurrentTopologyId] = useState<string | null>(null)
  const [showLoginModal, setShowLoginModal] = useState(false)
  const [showCheckoutModal, setShowCheckoutModal] = useState(false)
  const [checkoutTier, setCheckoutTier] = useState<'premium'>('premium')
  const [checkoutProvider, setCheckoutProvider] = useState<'stripe' | 'paypal'>('stripe')

  return (
    <div className="app-container">
      <DemoMode onLoginClick={() => setShowLoginModal(true)} />
      <div className="app-header">
        <h1>{t('app.title')}</h1>
        <div className="app-header-right">
          <div className="mode-tabs">
            <button
              className={`mode-tab ${mode === 'design' ? 'active' : ''}`}
              onClick={() => setMode('design')}
            >
              {t('app.mode.design')}
            </button>
            <button
              className={`mode-tab ${mode === 'security' ? 'active' : ''}`}
              onClick={() => setMode('security')}
            >
              {t('app.mode.security')}
            </button>
            <button
              className={`mode-tab ${mode === 'esg' ? 'active' : ''}`}
              onClick={() => setMode('esg')}
            >
              {t('app.mode.esg')}
            </button>
            {isAuthenticated && (
              <button
                className={`mode-tab ${mode === 'pricing' ? 'active' : ''}`}
                onClick={() => setMode('pricing')}
              >
                {t('app.mode.pricing')}
              </button>
            )}
          </div>
          {isAuthenticated ? (
            <div className="user-menu">
              <UserProfile />
              <button className="logout-button" onClick={logout}>
                {t('auth.logout')}
              </button>
            </div>
          ) : (
            <button className="login-button-header" onClick={() => setShowLoginModal(true)}>
              {t('auth.login.button')}
            </button>
          )}
          <LanguageSwitcher />
        </div>
      </div>
      <div className="app-content">
        {mode === 'design' ? (
          <>
            <div className="app-sidebar-left">
              <Palette />
            </div>
            <div className="app-main">
              <TopologyCanvas 
                onNodeSelect={setSelectedNode}
                simulationResult={simulationResult}
                onSimulationComplete={setSimulationResult}
                currentTopologyId={currentTopologyId}
                onTopologyIdChange={setCurrentTopologyId}
              />
            </div>
            <div className="app-sidebar-right">
              <PropertiesPanel nodeId={selectedNode} simulationResult={simulationResult} />
            </div>
          </>
        ) : mode === 'security' ? (
          <SecurityTestCanvas
            onNodeSelect={setSelectedNode}
            currentTopologyId={currentTopologyId}
            onTopologyIdChange={setCurrentTopologyId}
          />
        ) : mode === 'pricing' ? (
          <div className="pricing-page">
            <PricingPlans
              onSelectPlan={(tier, provider) => {
                setCheckoutTier(tier)
                setCheckoutProvider(provider)
                setShowCheckoutModal(true)
              }}
            />
            {isAuthenticated && (
              <div className="payment-history-section">
                <PaymentHistory />
              </div>
            )}
          </div>
        ) : (
          <ESGSimulationCanvas
            onNodeSelect={setSelectedNode}
            currentTopologyId={currentTopologyId}
            onTopologyIdChange={setCurrentTopologyId}
          />
        )}
      </div>
      <LoginModal isOpen={showLoginModal} onClose={() => setShowLoginModal(false)} />
      <CheckoutModal
        isOpen={showCheckoutModal}
        onClose={() => setShowCheckoutModal(false)}
        tier={checkoutTier}
        provider={checkoutProvider}
      />
    </div>
  )
}

export default App

