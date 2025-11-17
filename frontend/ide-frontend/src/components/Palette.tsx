import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import './Palette.css'

interface PaletteItem {
  id: string
  nameKey: string
  icon: string
  type: string
}

const paletteItems: PaletteItem[] = [
  { id: 'bus', nameKey: 'bus', icon: '⚡', type: 'bus' },
  { id: 'transformer', nameKey: 'transformer', icon: '🔌', type: 'transformer' },
  { id: 'switch', nameKey: 'switch', icon: '🔀', type: 'switch' },
  { id: 'line', nameKey: 'line', icon: '📏', type: 'line' },
  { id: 'ev_charger', nameKey: 'ev_charger', icon: '🔋', type: 'ev_charger' },
  { id: 'der', nameKey: 'der', icon: '☀️', type: 'der' },
]

function Palette() {
  const { t } = useTranslation()
  const [selectedItem, setSelectedItem] = useState<string | null>(null)

  return (
    <div className="palette">
      <h2 className="palette-title">{t('palette.title')}</h2>
      <div className="palette-items">
        {paletteItems.map((item) => (
          <div
            key={item.id}
            className={`palette-item ${selectedItem === item.id ? 'selected' : ''}`}
            onClick={() => setSelectedItem(item.id)}
            draggable
            onDragStart={(e) => {
              e.dataTransfer.setData('application/json', JSON.stringify(item))
            }}
          >
            <span className="palette-item-icon">{item.icon}</span>
            <span className="palette-item-name">{t(`palette.${item.nameKey}`)}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

export default Palette

