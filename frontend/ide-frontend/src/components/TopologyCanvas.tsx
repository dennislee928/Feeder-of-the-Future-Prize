import { useCallback, useState, useEffect } from 'react'
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
  NodeTypes,
} from 'reactflow'
import 'reactflow/dist/style.css'
import { simApi, PowerflowResult } from '../api/simApi'
import './TopologyCanvas.css'

interface TopologyCanvasProps {
  onNodeSelect: (nodeId: string | null) => void
  simulationResult: PowerflowResult | null
  onSimulationComplete: (result: PowerflowResult | null) => void
}

function TopologyCanvas({ onNodeSelect, simulationResult, onSimulationComplete }: TopologyCanvasProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState([])
  const [edges, setEdges, onEdgesChange] = useEdgesState([])
  const [isRunningSimulation, setIsRunningSimulation] = useState(false)

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
          label: item.name,
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
    [setNodes]
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

  // 執行模擬
  const handleRunSimulation = useCallback(async () => {
    if (nodes.length === 0) {
      alert('請先建立拓樸節點')
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
        })),
        profile_type: 'suburban', // 預設值，之後可以從 UI 選擇
      }

      const result = await simApi.runPowerflow({ topology })
      onSimulationComplete(result)
    } catch (error) {
      console.error('Simulation failed:', error)
      alert('模擬執行失敗，請檢查網路連線與服務狀態')
    } finally {
      setIsRunningSimulation(false)
    }
  }, [nodes, edges, onSimulationComplete])

  return (
    <div className="topology-canvas" onDrop={onDrop} onDragOver={onDragOver}>
      <div className="canvas-toolbar">
        <button
          className="simulation-button"
          onClick={handleRunSimulation}
          disabled={isRunningSimulation || nodes.length === 0}
        >
          {isRunningSimulation ? '執行中...' : '執行模擬 (Powerflow)'}
        </button>
        {simulationResult && (
          <div className="simulation-summary">
            <span>平均電壓: {simulationResult.summary.average_voltage_pu.toFixed(4)} pu</span>
            <span>最大載流率: {simulationResult.summary.max_line_loading_percent.toFixed(2)}%</span>
          </div>
        )}
      </div>
      <ReactFlow
        nodes={nodes}
        edges={edges}
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

