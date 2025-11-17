import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import TopologyCanvas from './components/TopologyCanvas'
import SecurityTestCanvas from './components/SecurityTestCanvas'
import Palette from './components/Palette'
import PropertiesPanel from './components/PropertiesPanel'
import LanguageSwitcher from './components/LanguageSwitcher'
import { PowerflowResult } from './api/simApi'
import './App.css'

type AppMode = 'design' | 'security'

function App() {
  const { t } = useTranslation()
  const [mode, setMode] = useState<AppMode>('design')
  const [selectedNode, setSelectedNode] = useState<string | null>(null)
  const [simulationResult, setSimulationResult] = useState<PowerflowResult | null>(null)
  const [currentTopologyId, setCurrentTopologyId] = useState<string | null>(null)

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
          </div>
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
        ) : (
          <SecurityTestCanvas
            onNodeSelect={setSelectedNode}
            currentTopologyId={currentTopologyId}
            onTopologyIdChange={setCurrentTopologyId}
          />
        )}
      </div>
    </div>
  )
}

export default App

