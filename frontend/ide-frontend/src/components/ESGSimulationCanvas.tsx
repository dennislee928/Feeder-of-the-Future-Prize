import { useCallback, useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import ReactFlow, {
  Node,
  Edge,
  addEdge,
  Connection,
  useNodesState,
  useEdgesState,
  Background,
  Controls,
  MiniMap,
} from 'reactflow'
import 'reactflow/dist/style.css'
import { ideApi } from '../api/ideApi'
import { esgApi, ESGCalculationResult } from '../api/esgApi'
import { useFeaturePermission } from '../hooks/useFeaturePermission'
import ESGCalculationPanel from './ESGCalculationPanel'
import ESGReportPanel from './ESGReportPanel'
import Palette from './Palette'
import CustomNode from './CustomNode'
import './ESGSimulationCanvas.css'

const nodeTypes = {
  default: CustomNode,
}

interface ESGSimulationCanvasProps {
  onNodeSelect: (nodeId: string | null) => void
  currentTopologyId: string | null
  onTopologyIdChange: (id: string | null) => void
}

interface ESGParameters {
  time_hours: number
  ev_charging_hours: number
  solar_generation_hours: number
  battery_cycles: number
}

function ESGSimulationCanvas({
  onNodeSelect,
  currentTopologyId: _currentTopologyId,
  onTopologyIdChange,
}: ESGSimulationCanvasProps) {
  const { t } = useTranslation()
  const { canUseFeature } = useFeaturePermission()
  const [nodes, setNodes, onNodesChange] = useNodesState([])
  const [edges, setEdges, onEdgesChange] = useEdgesState([])
  const [isLoading, setIsLoading] = useState(false)
  const [isCalculating, setIsCalculating] = useState(false)
  const [esgResult, setEsgResult] = useState<ESGCalculationResult | null>(null)
  
  // ESG 功能需要註冊會員（免費或付費）
  const canUseESG = canUseFeature('advanced_security') // 使用相同的權限檢查
  const [parameters, setParameters] = useState<ESGParameters>({
    time_hours: 24,
    ev_charging_hours: 4,
    solar_generation_hours: 6,
    battery_cycles: 1,
  })

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge(params, eds)),
    [setEdges]
  )

  const onNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => {
      onNodeSelect(node.id)
    },
    [onNodeSelect]
  )

  const onPaneClick = useCallback(() => {
    onNodeSelect(null)
  }, [onNodeSelect])

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault()

      const data = event.dataTransfer.getData('application/json')
      if (!data) return

      const item = JSON.parse(data)
      const position = {
        x: event.clientX - 300,
        y: event.clientY - 100,
      }

      const newNode: Node = {
        id: `${item.type}-${Date.now()}`,
        type: 'default',
        position,
        data: {
          label: item.nameKey ? t(`palette.${item.nameKey}`) : item.type,
          type: item.type,
        },
        style: {
          background: '#2a2a2a',
          color: '#fff',
          border: '1px solid #444',
          borderRadius: '4px',
          padding: '10px',
        },
      }

      setNodes((nds) => nds.concat(newNode))
    },
    [setNodes, t]
  )

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
  }, [])

  // 載入拓樸
  const handleLoadTopology = useCallback(async () => {
    setIsLoading(true)
    try {
      const topologies = await ideApi.listTopologies()
      if (topologies.length === 0) {
        alert(t('esg.no_topology'))
        return
      }

      const topology = topologies[0]
      onTopologyIdChange(topology.id)

      const flowNodes: Node[] = topology.nodes.map((node) => ({
        id: node.id,
        type: 'default',
        position: node.position,
        data: {
          label: node.name,
          type: node.type,
          properties: node.properties || {},
        },
        style: {
          background: '#2a2a2a',
          color: '#fff',
          border: '1px solid #444',
          borderRadius: '4px',
          padding: '10px',
        },
      }))

      const flowEdges: Edge[] = topology.lines.map((line) => ({
        id: line.id,
        source: line.from_node_id,
        target: line.to_node_id,
        label: line.name,
      }))

      setNodes(flowNodes)
      setEdges(flowEdges)
      alert(t('esg.topology_loaded'))
    } catch (error) {
      console.error('Load failed:', error)
      alert(t('esg.load_failed'))
    } finally {
      setIsLoading(false)
    }
  }, [setNodes, setEdges, onTopologyIdChange, t])

  // 執行 ESG 計算
  const handleCalculateESG = useCallback(async () => {
    if (nodes.length === 0) {
      alert(t('esg.add_nodes_first'))
      return
    }

    // 檢查是否可以使用 ESG 功能（需要註冊會員）
    if (!canUseESG) {
      alert(t('esg.feature_locked'))
      return
    }

    setIsCalculating(true)
    try {
      const topology = {
        nodes: nodes.map((node) => ({
          id: node.id,
          type: node.data.type || 'bus',
          name: node.data.label || node.id,
          position: node.position,
          properties: node.data.properties || {},
        })),
        lines: edges.map((edge) => ({
          id: edge.id,
          from_node_id: edge.source,
          to_node_id: edge.target,
          name: typeof edge.label === 'string' ? edge.label : undefined,
        })),
        profile_type: 'suburban',
      }

      const result = await esgApi.calculateESG({
        topology,
        parameters,
      })

      setEsgResult(result)
    } catch (error) {
      console.error('ESG calculation failed:', error)
      alert(t('esg.calculation_failed'))
    } finally {
      setIsCalculating(false)
    }
  }, [nodes, edges, parameters, canUseESG, t])

  // 根據 ESG 結果更新節點視覺化
  useEffect(() => {
    if (!esgResult) return

    // 創建節點排放映射
    const nodeEmissionMap = new Map<string, number>()
    esgResult.node_emissions.forEach((emission) => {
      nodeEmissionMap.set(emission.node_id, emission.emission_kg_co2)
    })

    // 找到最大和最小排放值（用於顏色映射）
    const emissions = Array.from(nodeEmissionMap.values())
    const maxEmission = Math.max(...emissions, 1)
    const minEmission = Math.min(...emissions, -maxEmission)

    // 更新節點顏色
    setNodes((nds) =>
      nds.map((node) => {
        const emission = nodeEmissionMap.get(node.id) || 0

        // 計算顏色（紅色=高排放，綠色=負排放/減排）
        let backgroundColor = '#2a2a2a'
        let borderColor = '#444'

        if (emission > 0) {
          // 正排放：紅色系
          const intensity = Math.min(emission / maxEmission, 1)
          const red = Math.floor(200 + intensity * 55)
          const green = Math.floor(50 - intensity * 50)
          const blue = Math.floor(50 - intensity * 50)
          backgroundColor = `rgb(${red}, ${green}, ${blue})`
          borderColor = `rgb(${Math.min(red + 20, 255)}, ${Math.max(green - 20, 0)}, ${Math.max(blue - 20, 0)})`
        } else if (emission < 0) {
          // 負排放（減排）：綠色系
          const intensity = Math.min(Math.abs(emission) / Math.abs(minEmission), 1)
          const red = Math.floor(50 - intensity * 50)
          const green = Math.floor(200 + intensity * 55)
          const blue = Math.floor(50 - intensity * 50)
          backgroundColor = `rgb(${red}, ${green}, ${blue})`
          borderColor = `rgb(${Math.max(red - 20, 0)}, ${Math.min(green + 20, 255)}, ${Math.max(blue - 20, 0)})`
        }

        return {
          ...node,
          style: {
            ...node.style,
            background: backgroundColor,
            border: `2px solid ${borderColor}`,
          },
          data: {
            ...node.data,
            emission: emission,
          },
        }
      })
    )
  }, [esgResult, setNodes])

  return (
    <div className="esg-simulation-container">
      <div className="esg-sidebar-left">
        <div className="esg-sidebar-section">
          <Palette />
        </div>
        <div className="esg-sidebar-section">
          <ESGCalculationPanel
            parameters={parameters}
            onParametersChange={setParameters}
          />
        </div>
      </div>
      <div className="esg-main" onDrop={onDrop} onDragOver={onDragOver}>
        <div className="esg-canvas-toolbar">
          <div className="toolbar-group">
            <button
              className="toolbar-button"
              onClick={handleLoadTopology}
              disabled={isLoading}
            >
              {isLoading ? t('esg.loading') : t('esg.load_topology')}
            </button>
          </div>
          <div className="toolbar-group">
            <button
              className="esg-calculate-button"
              onClick={handleCalculateESG}
              disabled={isCalculating || nodes.length === 0}
              title={nodes.length === 0 ? t('esg.add_nodes_first') : t('esg.calculate')}
            >
              {isCalculating ? t('esg.calculating') : t('esg.calculate_esg')}
            </button>
            {esgResult && (
              <div className="esg-summary">
                <span>
                  {t('esg.total_emissions')}: {esgResult.total_emissions_ton_co2.toFixed(2)} ton CO₂
                </span>
                <span className={esgResult.carbon_credits_ton > 0 ? 'positive' : ''}>
                  {t('esg.carbon_credits')}: {esgResult.carbon_credits_ton.toFixed(2)} ton
                </span>
                <span>
                  {t('esg.esg_score')}: {esgResult.esg_score.toFixed(1)}/100
                </span>
              </div>
            )}
          </div>
        </div>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onNodeClick={onNodeClick}
          onPaneClick={onPaneClick}
          fitView
        >
          <Background />
          <Controls />
          <MiniMap />
        </ReactFlow>
      </div>
      <div className="esg-sidebar-right">
        <ESGReportPanel esgResult={esgResult} />
      </div>
    </div>
  )
}

export default ESGSimulationCanvas

