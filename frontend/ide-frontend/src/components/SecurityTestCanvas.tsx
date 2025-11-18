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
  MarkerType,
} from 'reactflow'
import 'reactflow/dist/style.css'
import { ideApi } from '../api/ideApi'
import { penetrationApi, PenetrationTestResult } from '../api/penetrationApi'
import AttackScenarioPanel from './AttackScenarioPanel'
import SecurityReportPanel from './SecurityReportPanel'
import Palette from './Palette'
import CustomNode from './CustomNode'
import './SecurityTestCanvas.css'

const nodeTypes = {
  default: CustomNode,
}

interface SecurityTestCanvasProps {
  onNodeSelect: (nodeId: string | null) => void
  currentTopologyId: string | null
  onTopologyIdChange: (id: string | null) => void
}

function SecurityTestCanvas({
  onNodeSelect,
  currentTopologyId: _currentTopologyId,
  onTopologyIdChange,
}: SecurityTestCanvasProps) {
  const { t } = useTranslation()
  const [nodes, setNodes, onNodesChange] = useNodesState([])
  const [edges, setEdges, onEdgesChange] = useEdgesState([])
  const [isLoading, setIsLoading] = useState(false)
  const [isRunningTest, setIsRunningTest] = useState(false)
  const [selectedScenarios, setSelectedScenarios] = useState<string[]>([])
  const [testResult, setTestResult] = useState<PenetrationTestResult | null>(null)

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
        x: event.clientX - 300, // 調整左側欄位寬度
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
    [setNodes]
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
        alert(t('security.no_topology'))
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
      alert(t('security.topology_loaded'))
    } catch (error) {
      console.error('Load failed:', error)
      alert(t('security.load_failed'))
    } finally {
      setIsLoading(false)
    }
  }, [setNodes, setEdges, onTopologyIdChange, t])

  // 執行滲透測試
  const handleRunPenetrationTest = useCallback(async () => {
    if (nodes.length === 0) {
      alert(t('security.add_nodes_first'))
      return
    }

    if (selectedScenarios.length === 0) {
      alert(t('security.select_scenarios'))
      return
    }

    setIsRunningTest(true)
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
        profile_type: 'suburban',
      }

      const result = await penetrationApi.runPenetrationTest({
        topology,
        attack_scenarios: selectedScenarios,
      })

      setTestResult(result)
    } catch (error) {
      console.error('Penetration test failed:', error)
      alert(t('security.test_failed'))
    } finally {
      setIsRunningTest(false)
    }
  }, [nodes, edges, selectedScenarios, t])

  // 根據測試結果更新節點和線路視覺化
  useEffect(() => {
    if (!testResult) return

    // 收集所有受影響的節點和線路
    const affectedNodesSet = new Set<string>()
    const affectedLinesSet = new Set<string>()
    const attackPaths: Array<{ from: string; to: string }> = []

    testResult.attacks.forEach((attack) => {
      attack.affected_nodes.forEach((nodeId) => affectedNodesSet.add(nodeId))
      attack.affected_lines.forEach((lineId) => affectedLinesSet.add(lineId))
      attack.attack_path.forEach((path) => {
        attackPaths.push({ from: path.from, to: path.to })
      })
    })

    // 更新節點顏色（根據嚴重程度）
    setNodes((nds) =>
      nds.map((node) => {
        if (!affectedNodesSet.has(node.id)) {
          return node
        }

        // 找到影響此節點的最嚴重攻擊
        const affectingAttacks = testResult.attacks.filter((attack) =>
          attack.affected_nodes.includes(node.id)
        )
        const maxSeverity = affectingAttacks.reduce((max, attack) => {
          const severityOrder = { low: 1, medium: 2, high: 3, critical: 4 }
          return Math.max(max, severityOrder[attack.severity])
        }, 0)

        let backgroundColor = '#2a2a2a'
        let borderColor = '#444'

        if (maxSeverity >= 4) {
          // critical
          backgroundColor = '#dc2626'
          borderColor = '#991b1b'
        } else if (maxSeverity >= 3) {
          // high
          backgroundColor = '#ef4444'
          borderColor = '#dc2626'
        } else if (maxSeverity >= 2) {
          // medium
          backgroundColor = '#f59e0b'
          borderColor = '#d97706'
        } else {
          // low
          backgroundColor = '#fbbf24'
          borderColor = '#f59e0b'
        }

        return {
          ...node,
          style: {
            ...node.style,
            background: backgroundColor,
            border: `3px solid ${borderColor}`,
          },
          data: {
            ...node.data,
            isAffected: true,
            severity: maxSeverity >= 4 ? 'critical' : maxSeverity >= 3 ? 'high' : maxSeverity >= 2 ? 'medium' : 'low',
          },
        }
      })
    )

    // 更新線路顏色和標示攻擊路徑
    setEdges((eds) =>
      eds.map((edge) => {
        const isAffected = affectedLinesSet.has(edge.id)
        const isAttackPath = attackPaths.some(
          (path) => path.from === edge.source && path.to === edge.target
        )

        if (!isAffected && !isAttackPath) {
          return edge
        }

        return {
          ...edge,
          style: {
            ...edge.style,
            stroke: isAttackPath ? '#ef4444' : '#f59e0b',
            strokeWidth: isAttackPath ? 3 : 2,
          },
          markerEnd: isAttackPath
            ? {
                type: MarkerType.ArrowClosed,
                color: '#ef4444',
              }
            : undefined,
          animated: isAttackPath,
        }
      })
    )
  }, [testResult, setNodes, setEdges])

  return (
    <div className="security-test-container">
      <div className="security-sidebar-left">
        <div className="security-sidebar-section">
          <Palette />
        </div>
        <div className="security-sidebar-section">
          <AttackScenarioPanel
            selectedScenarios={selectedScenarios}
            onScenariosChange={setSelectedScenarios}
          />
        </div>
      </div>
      <div className="security-main" onDrop={onDrop} onDragOver={onDragOver}>
        <div className="security-canvas-toolbar">
          <div className="toolbar-group">
            <button
              className="toolbar-button"
              onClick={handleLoadTopology}
              disabled={isLoading}
            >
              {isLoading ? t('security.loading') : t('security.load_topology')}
            </button>
          </div>
          <div className="toolbar-group">
            <button
              className="penetration-test-button"
              onClick={handleRunPenetrationTest}
              disabled={isRunningTest || selectedScenarios.length === 0 || nodes.length === 0}
              title={
                nodes.length === 0
                  ? t('security.add_nodes_first')
                  : selectedScenarios.length === 0
                  ? t('security.select_scenarios')
                  : t('security.run_test')
              }
            >
              {isRunningTest ? t('security.running') : t('security.run_penetration_test')}
            </button>
            {testResult && (
              <div className="test-summary">
                <span>
                  {t('security.total_attacks')}: {testResult.summary.total_attacks}
                </span>
                <span className="success">
                  {t('security.successful')}: {testResult.summary.successful}
                </span>
                <span className="critical">
                  {t('security.critical')}: {testResult.summary.critical_vulnerabilities}
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
      <div className="security-sidebar-right">
        <SecurityReportPanel testResult={testResult} />
      </div>
    </div>
  )
}

export default SecurityTestCanvas

