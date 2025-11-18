import { useTranslation } from 'react-i18next'
import './ESGCalculationPanel.css'

interface ESGParameters {
  time_hours: number
  ev_charging_hours: number
  solar_generation_hours: number
  battery_cycles: number
}

interface ESGCalculationPanelProps {
  parameters: ESGParameters
  onParametersChange: (params: ESGParameters) => void
}

function ESGCalculationPanel({ parameters, onParametersChange }: ESGCalculationPanelProps) {
  const { t } = useTranslation()

  const handleChange = (key: keyof ESGParameters, value: number) => {
    onParametersChange({
      ...parameters,
      [key]: value,
    })
  }

  return (
    <div className="esg-calculation-panel">
      <h2 className="panel-title">{t('esg.calculation.title')}</h2>
      <div className="parameter-group">
        <label>
          {t('esg.calculation.time_hours')}
          <input
            type="number"
            min="1"
            max="8760"
            value={parameters.time_hours}
            onChange={(e) => handleChange('time_hours', parseFloat(e.target.value) || 24)}
          />
        </label>
      </div>
      <div className="parameter-group">
        <label>
          {t('esg.calculation.ev_charging_hours')}
          <input
            type="number"
            min="0"
            max="24"
            step="0.5"
            value={parameters.ev_charging_hours}
            onChange={(e) => handleChange('ev_charging_hours', parseFloat(e.target.value) || 4)}
          />
        </label>
      </div>
      <div className="parameter-group">
        <label>
          {t('esg.calculation.solar_generation_hours')}
          <input
            type="number"
            min="0"
            max="24"
            step="0.5"
            value={parameters.solar_generation_hours}
            onChange={(e) => handleChange('solar_generation_hours', parseFloat(e.target.value) || 6)}
          />
        </label>
      </div>
      <div className="parameter-group">
        <label>
          {t('esg.calculation.battery_cycles')}
          <input
            type="number"
            min="0"
            max="10"
            step="0.1"
            value={parameters.battery_cycles}
            onChange={(e) => handleChange('battery_cycles', parseFloat(e.target.value) || 1)}
          />
        </label>
      </div>
    </div>
  )
}

export default ESGCalculationPanel

