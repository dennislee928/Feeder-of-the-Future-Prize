import { useEffect, useState } from 'react'
import './PropertiesPanel.css'

interface PropertiesPanelProps {
  nodeId: string | null
}

interface NodeProperties {
  id: string
  name: string
  type: string
  [key: string]: any
}

function PropertiesPanel({ nodeId }: PropertiesPanelProps) {
  const [properties, setProperties] = useState<NodeProperties | null>(null)

  useEffect(() => {
    // TODO: 從 store 或 API 取得節點屬性
    if (nodeId) {
      // 暫時使用模擬資料
      setProperties({
        id: nodeId,
        name: `Node ${nodeId}`,
        type: 'bus',
      })
    } else {
      setProperties(null)
    }
  }, [nodeId])

  if (!nodeId || !properties) {
    return (
      <div className="properties-panel">
        <h2 className="properties-title">屬性面板</h2>
        <div className="properties-empty">
          <p>請選擇一個節點以查看屬性</p>
        </div>
      </div>
    )
  }

  return (
    <div className="properties-panel">
      <h2 className="properties-title">屬性面板</h2>
      <div className="properties-content">
        <div className="property-group">
          <label className="property-label">ID</label>
          <input
            type="text"
            className="property-input"
            value={properties.id}
            readOnly
          />
        </div>
        <div className="property-group">
          <label className="property-label">名稱</label>
          <input
            type="text"
            className="property-input"
            value={properties.name}
            onChange={(e) =>
              setProperties({ ...properties, name: e.target.value })
            }
          />
        </div>
        <div className="property-group">
          <label className="property-label">類型</label>
          <input
            type="text"
            className="property-input"
            value={properties.type}
            readOnly
          />
        </div>
      </div>
    </div>
  )
}

export default PropertiesPanel

