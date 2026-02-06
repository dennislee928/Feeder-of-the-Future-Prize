import React, { useState, useCallback } from 'react'
import { useFeatureUnlock } from '../contexts/FeatureUnlockContext'
import './FeatureUnlockInput.css'

export const FeatureUnlockInput: React.FC = () => {
  const { isUnlocked, setUnlockCode } = useFeatureUnlock()
  const [value, setValue] = useState('')

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value
      setValue(newValue)
      setUnlockCode(newValue)
    },
    [setUnlockCode]
  )

  if (isUnlocked) {
    return (
      <span className="feature-unlock-badge" title="All features unlocked">
        ✓ Unlocked
      </span>
    )
  }

  return (
    <input
      type="text"
      className="feature-unlock-input"
      placeholder="Code"
      value={value}
      onChange={handleChange}
      aria-label="Feature unlock code"
    />
  )
}
