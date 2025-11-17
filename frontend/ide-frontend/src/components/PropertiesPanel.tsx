import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PowerflowResult } from '../api/simApi'
import './PropertiesPanel.css'

interface PropertiesPanelProps {
  nodeId: string | null
  simulationResult: PowerflowResult | null
}

interface NodeProperties {
  id: string
  name: string
  type: string
  [key: string]: any
}

function PropertiesPanel({ nodeId, simulationResult }: PropertiesPanelProps) {
  const { t } = useTranslation()
  const [properties, setProperties] = useState<NodeProperties | null>(null)
  const [nodeSimulationData, setNodeSimulationData] = useState<any>(null)

  useEffect(() => {
    // TODO: 從 store 或 API 取得節點屬性
    if (nodeId) {
      // 暫時使用模擬資料
      setProperties({
        id: nodeId,
        name: `Node ${nodeId}`,
        type: 'bus',
      })

      // 如果有模擬結果，取得該節點的資料
      if (simulationResult) {
        const nodeData = simulationResult.nodes.find((n) => n.node_id === nodeId)
        setNodeSimulationData(nodeData || null)
      } else {
        setNodeSimulationData(null)
      }
    } else {
      setProperties(null)
      setNodeSimulationData(null)
    }
  }, [nodeId, simulationResult])

  if (!nodeId || !properties) {
    return (
      <div className="properties-panel">
        <h2 className="properties-title">{t('properties.title')}</h2>
        <div className="properties-empty">
          <p>{t('properties.empty')}</p>
        </div>
      </div>
    )
  }

  return (
    <div className="properties-panel">
      <h2 className="properties-title">{t('properties.title')}</h2>
      <div className="properties-content">
        <div className="property-group">
          <label className="property-label">{t('properties.id')}</label>
          <input
            type="text"
            className="property-input"
            value={properties.id}
            readOnly
          />
        </div>
        <div className="property-group">
          <label className="property-label">{t('properties.name')}</label>
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
          <label className="property-label">{t('properties.type')}</label>
          <input
            type="text"
            className="property-input"
            value={properties.type}
            readOnly
          />
        </div>

        {nodeSimulationData && (
          <>
            <div className="property-divider">{t('properties.simulation_results')}</div>
            <div className="property-group">
              <label className="property-label">{t('properties.voltage')}</label>
              <input
                type="text"
                className="property-input"
                value={nodeSimulationData.voltage_pu.toFixed(4)}
                readOnly
              />
            </div>
            <div className="property-group">
              <label className="property-label">{t('properties.voltage_kv')}</label>
              <input
                type="text"
                className="property-input"
                value={nodeSimulationData.voltage_kv.toFixed(4)}
                readOnly
              />
            </div>
            <div className="property-group">
              <label className="property-label">{t('properties.voltage_deviation')}</label>
              <input
                type="text"
                className={`property-input ${
                  nodeSimulationData.status === 'critical'
                    ? 'status-critical'
                    : nodeSimulationData.status === 'warning'
                    ? 'status-warning'
                    : 'status-normal'
                }`}
                value={`${nodeSimulationData.voltage_deviation_percent > 0 ? '+' : ''}${nodeSimulationData.voltage_deviation_percent.toFixed(2)}%`}
                readOnly
              />
            </div>
            <div className="property-group">
              <label className="property-label">{t('properties.status')}</label>
              <input
                type="text"
                className={`property-input ${
                  nodeSimulationData.status === 'critical'
                    ? 'status-critical'
                    : nodeSimulationData.status === 'warning'
                    ? 'status-warning'
                    : 'status-normal'
                }`}
                value={
                  nodeSimulationData.status === 'critical'
                    ? t('properties.status_critical')
                    : nodeSimulationData.status === 'warning'
                    ? t('properties.status_warning')
                    : t('properties.status_normal')
                }
                readOnly
              />
            </div>
          </>
        )}
      </div>
    </div>
  )
}

export default PropertiesPanel

