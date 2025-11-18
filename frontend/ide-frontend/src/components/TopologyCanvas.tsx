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
import { simApi, PowerflowResult } from '../api/simApi'
import { ideApi } from '../api/ideApi'
import { useFeaturePermission } from '../hooks/useFeaturePermission'
import CustomNode from './CustomNode'
import './TopologyCanvas.css'

const nodeTypes = {
  default: CustomNode,
}

interface TopologyCanvasProps {
  onNodeSelect: (nodeId: string | null) => void
  simulationResult: PowerflowResult | null
  onSimulationComplete: (result: PowerflowResult | null) => void
  currentTopologyId: string | null
  onTopologyIdChange: (id: string | null) => void
}

function TopologyCanvas({ 
  onNodeSelect, 
  simulationResult, 
  onSimulationComplete,
  currentTopologyId,
  onTopologyIdChange
}: TopologyCanvasProps) {
  const { t } = useTranslation()
  const { canCreateTopology, canRunSimulation } = useFeaturePermission()
  const [nodes, setNodes, onNodesChange] = useNodesState([])
  const [edges, setEdges, onEdgesChange] = useEdgesState([])
  const [isRunningSimulation, setIsRunningSimulation] = useState(false)
  const [isSaving, setIsSaving] = useState(false)
  const [isLoading, setIsLoading] = useState(false)

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
        x: event.clientX - 250, // 調整左側欄位寬度
        y: event.clientY - 100, // 調整 header 高度
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

  // 根據模擬結果更新節點顏色
  useEffect(() => {
    if (!simulationResult) return

    setNodes((nds) =>
      nds.map((node) => {
        const nodeResult = simulationResult.nodes.find((n) => n.node_id === node.id)
        if (!nodeResult) return node

        // 根據電壓偏差設定顏色
        let backgroundColor = '#2a2a2a'
        let borderColor = '#444'

        if (nodeResult.status === 'warning') {
          backgroundColor = '#fbbf24' // yellow
          borderColor = '#f59e0b'
        } else if (nodeResult.status === 'critical') {
          backgroundColor = '#ef4444' // red
          borderColor = '#dc2626'
        } else {
          backgroundColor = '#10b981' // green
          borderColor = '#059669'
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
            voltage: nodeResult.voltage_pu,
            voltageDeviation: nodeResult.voltage_deviation_percent,
          },
        }
      })
    )
  }, [simulationResult, setNodes])

  // 儲存拓樸
  const handleSaveTopology = useCallback(async () => {
    if (nodes.length === 0) {
      alert(t('topology.add_nodes_first'))
      return
    }

    // 檢查配額（僅在創建新拓樸時）
    if (!currentTopologyId) {
      const quotaCheck = canCreateTopology()
      if (!quotaCheck.canCreate) {
        alert(
          t('topology.quota_exceeded', {
            used: quotaCheck.used,
            max: quotaCheck.max === Infinity ? '∞' : quotaCheck.max,
          })
        )
        return
      }
    }

    setIsSaving(true)
    try {
      const topologyData = {
        name: `Topology ${new Date().toLocaleString()}`,
        description: 'Created from IDE',
        profile_type: 'suburban' as const,
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
          name: typeof edge.label === 'string' ? edge.label : '',
          properties: {},
        })),
      }

      if (currentTopologyId) {
        // 更新現有拓樸
        await ideApi.updateTopology(currentTopologyId, topologyData)
        alert(t('topology.topology_updated'))
      } else {
        // 建立新拓樸
        const result = await ideApi.createTopology(topologyData)
        onTopologyIdChange(result.id)
        alert(t('topology.topology_saved'))
      }
    } catch (error) {
      console.error('Save failed:', error)
      alert(t('topology.save_failed'))
    } finally {
      setIsSaving(false)
    }
  }, [nodes, edges, currentTopologyId, onTopologyIdChange, t])

  // 載入拓樸
  const handleLoadTopology = useCallback(async () => {
    setIsLoading(true)
    try {
      const topologies = await ideApi.listTopologies()
      if (topologies.length === 0) {
        alert(t('topology.no_topology'))
        return
      }

      // 簡單選擇第一個（之後可以改成選擇對話框）
      const topology = topologies[0]
      onTopologyIdChange(topology.id)

      // 轉換為 React Flow 格式
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
      onTopologyIdChange(topology.id)
      alert(t('topology.topology_loaded'))
    } catch (error) {
      console.error('Load failed:', error)
      alert(t('topology.load_failed'))
    } finally {
      setIsLoading(false)
    }
  }, [setNodes, setEdges, onTopologyIdChange, t])

  // 執行模擬
  const handleRunSimulation = useCallback(async () => {
    if (nodes.length === 0) {
      alert(t('topology.add_nodes_first'))
      return
    }

    // 檢查模擬配額
    if (!canRunSimulation()) {
      alert(t('topology.simulation_quota_exceeded'))
      return
    }

    setIsRunningSimulation(true)

    try {
      const topology = {
        nodes: nodes.map((node) => ({
          id: node.id,
          type: node.data.type || 'bus',
          name: node.data.label || node.id,
          position: node.position,
        })),
        lines: edges.map((edge) => ({
          id: edge.id,
          from_node_id: edge.source,
          to_node_id: edge.target,
          name: typeof edge.label === 'string' ? edge.label : undefined,
        })),
        profile_type: 'suburban', // 預設值，之後可以從 UI 選擇
      }

      const result = await simApi.runPowerflow({ topology })
      onSimulationComplete(result)
    } catch (error) {
      console.error('Simulation failed:', error)
      alert(t('topology.simulation_failed'))
    } finally {
      setIsRunningSimulation(false)
    }
  }, [nodes, edges, onSimulationComplete, t])

  return (
    <div className="topology-canvas" onDrop={onDrop} onDragOver={onDragOver}>
      <div className="canvas-toolbar">
        <div className="toolbar-group">
          <button
            className="toolbar-button"
            onClick={handleSaveTopology}
            disabled={isSaving || nodes.length === 0}
          >
            {isSaving ? t('topology.saving') : t('topology.save')}
          </button>
          <button
            className="toolbar-button"
            onClick={handleLoadTopology}
            disabled={isLoading}
          >
            {isLoading ? t('topology.loading') : t('topology.load')}
          </button>
        </div>
        <div className="toolbar-group">
          <button
            className="simulation-button"
            onClick={handleRunSimulation}
            disabled={isRunningSimulation}
            title={nodes.length === 0 ? t('topology.simulation_tooltip_empty') : t('topology.simulation_tooltip')}
          >
            {isRunningSimulation ? t('topology.running') : t('topology.run_simulation')}
          </button>
          {simulationResult && (
            <div className="simulation-summary">
              <span>{t('topology.average_voltage')}: {simulationResult.summary.average_voltage_pu.toFixed(4)} pu</span>
              <span>{t('topology.max_line_loading')}: {simulationResult.summary.max_line_loading_percent.toFixed(2)}%</span>
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
  )
}

export default TopologyCanvas

