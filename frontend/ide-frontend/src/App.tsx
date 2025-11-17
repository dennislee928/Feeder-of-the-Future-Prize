import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import TopologyCanvas from './components/TopologyCanvas'
import Palette from './components/Palette'
import PropertiesPanel from './components/PropertiesPanel'
import LanguageSwitcher from './components/LanguageSwitcher'
import { PowerflowResult } from './api/simApi'
import './App.css'

function App() {
  const { t } = useTranslation()
  const [selectedNode, setSelectedNode] = useState<string | null>(null)
  const [simulationResult, setSimulationResult] = useState<PowerflowResult | null>(null)
  const [currentTopologyId, setCurrentTopologyId] = useState<string | null>(null)

  return (
    <div className="app-container">
      <div className="app-header">
        <h1>{t('app.title')}</h1>
        <LanguageSwitcher />
      </div>
      <div className="app-content">
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
      </div>
    </div>
  )
}

export default App

