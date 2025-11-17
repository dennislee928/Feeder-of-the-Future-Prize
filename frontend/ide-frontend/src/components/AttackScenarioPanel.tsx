import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { penetrationApi, AttackScenario } from '../api/penetrationApi'
import './AttackScenarioPanel.css'

interface AttackScenarioPanelProps {
  selectedScenarios: string[]
  onScenariosChange: (scenarios: string[]) => void
}

function AttackScenarioPanel({ selectedScenarios, onScenariosChange }: AttackScenarioPanelProps) {
  const { t } = useTranslation()
  const [scenarios, setScenarios] = useState<AttackScenario[]>([])
  const [loading, setLoading] = useState(true)
  const [expandedLayers, setExpandedLayers] = useState<Set<number>>(new Set([7, 6, 5, 4, 3, 2, 1]))

  useEffect(() => {
    const fetchScenarios = async () => {
      try {
        const data = await penetrationApi.getAvailableScenarios()
        setScenarios(data)
      } catch (error) {
        console.error('Failed to fetch scenarios:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchScenarios()
  }, [])

  const toggleScenario = (scenarioId: string) => {
    if (selectedScenarios.includes(scenarioId)) {
      onScenariosChange(selectedScenarios.filter((id) => id !== scenarioId))
    } else {
      onScenariosChange([...selectedScenarios, scenarioId])
    }
  }

  const toggleLayer = (layer: number) => {
    const newExpanded = new Set(expandedLayers)
    if (newExpanded.has(layer)) {
      newExpanded.delete(layer)
    } else {
      newExpanded.add(layer)
    }
    setExpandedLayers(newExpanded)
  }

  const selectAll = () => {
    onScenariosChange(scenarios.map((s) => s.id))
  }

  const deselectAll = () => {
    onScenariosChange([])
  }

  // 按層級分組
  const scenariosByLayer: Record<number, AttackScenario[]> = {}
  scenarios.forEach((scenario) => {
    if (!scenariosByLayer[scenario.layer]) {
      scenariosByLayer[scenario.layer] = []
    }
    scenariosByLayer[scenario.layer].push(scenario)
  })

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return '#dc2626'
      case 'high':
        return '#ef4444'
      case 'medium':
        return '#f59e0b'
      case 'low':
        return '#fbbf24'
      default:
        return '#6b7280'
    }
  }

  const getLayerName = (layer: number) => {
    const layerNames: Record<number, string> = {
      7: t('security.layer.application'),
      6: t('security.layer.presentation'),
      5: t('security.layer.session'),
      4: t('security.layer.transport'),
      3: t('security.layer.network'),
      2: t('security.layer.datalink'),
      1: t('security.layer.physical'),
    }
    return layerNames[layer] || `Layer ${layer}`
  }

  if (loading) {
    return (
      <div className="attack-scenario-panel">
        <h2 className="panel-title">{t('security.scenarios.title')}</h2>
        <div className="loading">{t('security.scenarios.loading')}</div>
      </div>
    )
  }

  return (
    <div className="attack-scenario-panel">
      <div className="panel-header">
        <h2 className="panel-title">{t('security.scenarios.title')}</h2>
        <div className="panel-actions">
          <button className="action-button" onClick={selectAll}>
            {t('security.scenarios.select_all')}
          </button>
          <button className="action-button" onClick={deselectAll}>
            {t('security.scenarios.deselect_all')}
          </button>
        </div>
      </div>
      <div className="selected-count">
        {t('security.scenarios.selected')}: {selectedScenarios.length} / {scenarios.length}
      </div>
      <div className="scenarios-list">
        {[7, 6, 5, 4, 3, 2, 1].map((layer) => {
          const layerScenarios = scenariosByLayer[layer] || []
          if (layerScenarios.length === 0) return null

          const isExpanded = expandedLayers.has(layer)

          return (
            <div key={layer} className="layer-group">
              <div className="layer-header" onClick={() => toggleLayer(layer)}>
                <span className="layer-name">{getLayerName(layer)}</span>
                <span className="layer-count">({layerScenarios.length})</span>
                <span className="layer-toggle">{isExpanded ? '▼' : '▶'}</span>
              </div>
              {isExpanded && (
                <div className="layer-scenarios">
                  {layerScenarios.map((scenario) => {
                    const isSelected = selectedScenarios.includes(scenario.id)
                    return (
                      <div
                        key={scenario.id}
                        className={`scenario-item ${isSelected ? 'selected' : ''}`}
                        onClick={() => toggleScenario(scenario.id)}
                      >
                        <input
                          type="checkbox"
                          checked={isSelected}
                          onChange={() => toggleScenario(scenario.id)}
                          onClick={(e) => e.stopPropagation()}
                        />
                        <div className="scenario-content">
                          <div className="scenario-header">
                            <span className="scenario-name">{scenario.name}</span>
                            <span
                              className="severity-badge"
                              style={{ backgroundColor: getSeverityColor(scenario.severity) }}
                            >
                              {t(`security.severity.${scenario.severity}`)}
                            </span>
                          </div>
                          <div className="scenario-description">{scenario.description}</div>
                        </div>
                      </div>
                    )
                  })}
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}

export default AttackScenarioPanel

