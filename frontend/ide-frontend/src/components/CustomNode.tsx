import { memo } from 'react'
import { Handle, Position, NodeProps } from 'reactflow'
import './CustomNode.css'

interface CustomNodeData {
  label: string
  type?: string
  properties?: Record<string, any>
  voltage?: number
  voltageDeviation?: number
  isAffected?: boolean
  severity?: 'low' | 'medium' | 'high' | 'critical'
}

function CustomNode({ data }: NodeProps<CustomNodeData>) {
  return (
    <div
      className="custom-node"
      style={{
        background: data.style?.background || '#2a2a2a',
        border: data.style?.border || '1px solid #444',
        color: data.style?.color || '#fff',
        padding: '10px',
        borderRadius: '4px',
        minWidth: '120px',
        textAlign: 'center',
      }}
    >
      <Handle type="target" position={Position.Top} />
      <div className="node-label">{data.label || 'Node'}</div>
      {data.type && (
        <div className="node-type" style={{ fontSize: '0.75rem', color: '#aaa', marginTop: '4px' }}>
          {data.type}
        </div>
      )}
      {data.voltage !== undefined && (
        <div className="node-voltage" style={{ fontSize: '0.7rem', color: '#aaa', marginTop: '2px' }}>
          {data.voltage.toFixed(3)} pu
        </div>
      )}
      <Handle type="source" position={Position.Bottom} />
    </div>
  )
}

export default memo(CustomNode)

