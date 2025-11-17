import { useState } from 'react'
import './Palette.css'

interface PaletteItem {
  id: string
  name: string
  icon: string
  type: string
}

const paletteItems: PaletteItem[] = [
  { id: 'bus', name: 'Bus', icon: '⚡', type: 'bus' },
  { id: 'transformer', name: 'Transformer', icon: '🔌', type: 'transformer' },
  { id: 'switch', name: 'Switch', icon: '🔀', type: 'switch' },
  { id: 'line', name: 'Line', icon: '📏', type: 'line' },
  { id: 'ev_charger', name: 'EV Charger', icon: '🔋', type: 'ev_charger' },
  { id: 'der', name: 'DER', icon: '☀️', type: 'der' },
]

function Palette() {
  const [selectedItem, setSelectedItem] = useState<string | null>(null)

  return (
    <div className="palette">
      <h2 className="palette-title">資產面板</h2>
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
            <span className="palette-item-name">{item.name}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

export default Palette

