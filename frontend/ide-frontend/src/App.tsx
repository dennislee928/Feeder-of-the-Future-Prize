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
import { PowerflowResult } from './api/simApi'
import './App.css'

type AppMode = 'design' | 'security' | 'esg'

function App() {
  const { t } = useTranslation()
  const { isAuthenticated, logout } = useAuth()
  const [mode, setMode] = useState<AppMode>('design')
  const [selectedNode, setSelectedNode] = useState<string | null>(null)
  const [simulationResult, setSimulationResult] = useState<PowerflowResult | null>(null)
  const [currentTopologyId, setCurrentTopologyId] = useState<string | null>(null)
  const [showLoginModal, setShowLoginModal] = useState(false)

  return (
    <div className="app-container">
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
        ) : (
          <ESGSimulationCanvas
            onNodeSelect={setSelectedNode}
            currentTopologyId={currentTopologyId}
            onTopologyIdChange={setCurrentTopologyId}
          />
        )}
      </div>
      <LoginModal isOpen={showLoginModal} onClose={() => setShowLoginModal(false)} />
    </div>
  )
}

export default App

