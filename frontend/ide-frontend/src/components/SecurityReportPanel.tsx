import { useTranslation } from 'react-i18next'
import { PenetrationTestResult } from '../api/penetrationApi'
import './SecurityReportPanel.css'

interface SecurityReportPanelProps {
  testResult: PenetrationTestResult | null
}

function SecurityReportPanel({ testResult }: SecurityReportPanelProps) {
  const { t } = useTranslation()

  if (!testResult) {
    return (
      <div className="security-report-panel">
        <h2 className="panel-title">{t('security.report.title')}</h2>
        <div className="report-empty">
          <p>{t('security.report.empty')}</p>
        </div>
      </div>
    )
  }

  const { summary, attacks } = testResult

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

  const getSuccessIcon = (successful: boolean) => {
    return successful ? '✓' : '✗'
  }

  const getSuccessColor = (successful: boolean) => {
    return successful ? '#ef4444' : '#10b981'
  }

  // 收集所有受影響的節點和線路
  const affectedNodes = new Set<string>()
  const affectedLines = new Set<string>()
  attacks.forEach((attack) => {
    attack.affected_nodes.forEach((node) => affectedNodes.add(node))
    attack.affected_lines.forEach((line) => affectedLines.add(line))
  })

  // 收集所有建議
  const allRecommendations = new Set<string>()
  attacks.forEach((attack) => {
    attack.recommendations.forEach((rec) => allRecommendations.add(rec))
  })

  return (
    <div className="security-report-panel">
      <h2 className="panel-title">{t('security.report.title')}</h2>
      <div className="report-content">
        {/* 摘要 */}
        <div className="report-section">
          <h3 className="section-title">{t('security.report.summary')}</h3>
          <div className="summary-grid">
            <div className="summary-item">
              <span className="summary-label">{t('security.report.total_attacks')}</span>
              <span className="summary-value">{summary.total_attacks}</span>
            </div>
            <div className="summary-item">
              <span className="summary-label">{t('security.report.successful')}</span>
              <span className="summary-value success">{summary.successful}</span>
            </div>
            <div className="summary-item">
              <span className="summary-label">{t('security.report.failed')}</span>
              <span className="summary-value">{summary.failed}</span>
            </div>
            <div className="summary-item">
              <span className="summary-label">{t('security.report.critical_vulns')}</span>
              <span className="summary-value critical">{summary.critical_vulnerabilities}</span>
            </div>
            <div className="summary-item">
              <span className="summary-label">{t('security.report.affected_nodes')}</span>
              <span className="summary-value">{summary.affected_nodes_count}</span>
            </div>
            <div className="summary-item">
              <span className="summary-label">{t('security.report.affected_lines')}</span>
              <span className="summary-value">{summary.affected_lines_count}</span>
            </div>
          </div>
        </div>

        {/* 攻擊列表 */}
        <div className="report-section">
          <h3 className="section-title">{t('security.report.attacks')}</h3>
          <div className="attacks-list">
            {attacks.map((attack) => (
              <div key={attack.attack_id} className="attack-item">
                <div className="attack-header">
                  <div className="attack-info">
                    <span className="attack-name">{attack.scenario_name}</span>
                    <span
                      className="severity-badge"
                      style={{ backgroundColor: getSeverityColor(attack.severity) }}
                    >
                      {t(`security.severity.${attack.severity}`)}
                    </span>
                    <span
                      className="success-badge"
                      style={{ color: getSuccessColor(attack.successful) }}
                    >
                      {getSuccessIcon(attack.successful)}{' '}
                      {attack.successful
                        ? t('security.report.successful')
                        : t('security.report.blocked')}
                    </span>
                  </div>
                  <div className="attack-layer">Layer {attack.layer}</div>
                </div>
                <div className="attack-details">
                  <div className="attack-impact">
                    <strong>{t('security.report.impact')}:</strong> {attack.impact}
                  </div>
                  {attack.affected_nodes.length > 0 && (
                    <div className="attack-affected">
                      <strong>{t('security.report.affected_nodes')}:</strong>{' '}
                      {attack.affected_nodes.join(', ')}
                    </div>
                  )}
                  {attack.attack_path.length > 0 && (
                    <div className="attack-path">
                      <strong>{t('security.report.attack_path')}:</strong>
                      <div className="path-list">
                        {attack.attack_path.map((path, idx) => (
                          <span key={idx} className="path-item">
                            {path.from} → {path.to}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* 受影響資產 */}
        {(affectedNodes.size > 0 || affectedLines.size > 0) && (
          <div className="report-section">
            <h3 className="section-title">{t('security.report.affected_assets')}</h3>
            {affectedNodes.size > 0 && (
              <div className="assets-list">
                <strong>{t('security.report.nodes')}:</strong>
                <div className="assets-tags">
                  {Array.from(affectedNodes).map((node) => (
                    <span key={node} className="asset-tag">
                      {node}
                    </span>
                  ))}
                </div>
              </div>
            )}
            {affectedLines.size > 0 && (
              <div className="assets-list">
                <strong>{t('security.report.lines')}:</strong>
                <div className="assets-tags">
                  {Array.from(affectedLines).map((line) => (
                    <span key={line} className="asset-tag">
                      {line}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        {/* 修復建議 */}
        {allRecommendations.size > 0 && (
          <div className="report-section">
            <h3 className="section-title">{t('security.report.recommendations')}</h3>
            <ul className="recommendations-list">
              {Array.from(allRecommendations).map((rec, idx) => (
                <li key={idx} className="recommendation-item">
                  {rec}
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </div>
  )
}

export default SecurityReportPanel

