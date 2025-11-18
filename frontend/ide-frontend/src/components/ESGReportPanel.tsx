import { useTranslation } from 'react-i18next'
import { ESGCalculationResult } from '../api/esgApi'
import './ESGReportPanel.css'

interface ESGReportPanelProps {
  esgResult: ESGCalculationResult | null
}

function ESGReportPanel({ esgResult }: ESGReportPanelProps) {
  const { t } = useTranslation()

  if (!esgResult) {
    return (
      <div className="esg-report-panel">
        <h2 className="panel-title">{t('esg.report.title')}</h2>
        <p className="empty-message">{t('esg.report.empty')}</p>
      </div>
    )
  }

  return (
    <div className="esg-report-panel">
      <h2 className="panel-title">{t('esg.report.title')}</h2>

      <div className="esg-summary-section">
        <h3>{t('esg.report.summary')}</h3>
        <div className="summary-item">
          <span className="label">{t('esg.report.total_emissions')}:</span>
          <span className="value negative">
            {esgResult.total_emissions_ton_co2.toFixed(3)} ton CO₂
          </span>
        </div>
        <div className="summary-item">
          <span className="label">{t('esg.report.carbon_credits')}:</span>
          <span className={`value ${esgResult.carbon_credits_ton > 0 ? 'positive' : ''}`}>
            {esgResult.carbon_credits_ton.toFixed(3)} ton
          </span>
        </div>
        <div className="summary-item">
          <span className="label">{t('esg.report.carbon_value')}:</span>
          <span className={`value ${esgResult.carbon_credit_value_usd > 0 ? 'positive' : ''}`}>
            ${esgResult.carbon_credit_value_usd.toFixed(2)} USD
          </span>
        </div>
        <div className="summary-item">
          <span className="label">{t('esg.report.esg_score')}:</span>
          <span className={`value score-${Math.floor(esgResult.esg_score / 20)}`}>
            {esgResult.esg_score.toFixed(1)}/100
          </span>
        </div>
      </div>

      <div className="esg-node-emissions-section">
        <h3>{t('esg.report.node_emissions')}</h3>
        <div className="node-emissions-list">
          {esgResult.node_emissions.map((emission) => (
            <div key={emission.node_id} className="node-emission-item">
              <div className="node-emission-header">
                <span className="node-id">{emission.node_id}</span>
                <span className="node-type">{emission.node_type}</span>
              </div>
              <div className="node-emission-details">
                <span>
                  {t('esg.report.power')}: {emission.power_kw.toFixed(2)} kW
                </span>
                <span>
                  {t('esg.report.energy')}: {emission.energy_kwh.toFixed(2)} kWh
                </span>
                <span className={emission.emission_kg_co2 >= 0 ? 'negative' : 'positive'}>
                  {t('esg.report.emission')}: {emission.emission_kg_co2.toFixed(3)} kg CO₂
                </span>
              </div>
            </div>
          ))}
        </div>
      </div>

      {esgResult.recommendations.length > 0 && (
        <div className="esg-recommendations-section">
          <h3>{t('esg.report.recommendations')}</h3>
          <div className="recommendations-list">
            {esgResult.recommendations.map((rec, index) => (
              <div key={index} className={`recommendation-item priority-${rec.priority}`}>
                <div className="recommendation-header">
                  <span className="recommendation-title">{rec.title}</span>
                  <span className={`priority-badge priority-${rec.priority}`}>
                    {t(`esg.report.priority.${rec.priority}`)}
                  </span>
                </div>
                <p className="recommendation-description">{rec.description}</p>
                {rec.estimated_reduction_ton && (
                  <div className="recommendation-impact">
                    {t('esg.report.estimated_reduction')}:{' '}
                    {rec.estimated_reduction_ton.toFixed(3)} ton CO₂
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default ESGReportPanel

