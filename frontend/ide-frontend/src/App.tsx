import { useState } from 'react'
import TopologyCanvas from './components/TopologyCanvas'
import Palette from './components/Palette'
import PropertiesPanel from './components/PropertiesPanel'
import { PowerflowResult } from './api/simApi'
import './App.css'

function App() {
  const [selectedNode, setSelectedNode] = useState<string | null>(null)
  const [simulationResult, setSimulationResult] = useState<PowerflowResult | null>(null)

  return (
    <div className="app-container">
      <div className="app-header">
        <h1>Feeder IDE - Digital Twin & Design Platform</h1>
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

