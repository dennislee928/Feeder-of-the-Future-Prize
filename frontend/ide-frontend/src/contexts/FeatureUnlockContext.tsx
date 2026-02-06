import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'

const UNLOCK_CODE = 'mitake@123'
const STORAGE_KEY = 'fe_unlock_all'

interface FeatureUnlockContextType {
  isUnlocked: boolean
  setUnlockCode: (code: string) => void
}

const FeatureUnlockContext = createContext<FeatureUnlockContextType | undefined>(undefined)

export const useFeatureUnlock = () => {
  const context = useContext(FeatureUnlockContext)
  if (!context) {
    throw new Error('useFeatureUnlock must be used within FeatureUnlockProvider')
  }
  return context
}

interface FeatureUnlockProviderProps {
  children: ReactNode
}

export const FeatureUnlockProvider: React.FC<FeatureUnlockProviderProps> = ({ children }) => {
  const [isUnlocked, setIsUnlocked] = useState(() => {
    return localStorage.getItem(STORAGE_KEY) === '1'
  })

  useEffect(() => {
    if (isUnlocked) {
      localStorage.setItem(STORAGE_KEY, '1')
    } else {
      localStorage.removeItem(STORAGE_KEY)
    }
  }, [isUnlocked])

  const setUnlockCode = (code: string) => {
    if (code === UNLOCK_CODE) {
      setIsUnlocked(true)
    }
  }

  const value: FeatureUnlockContextType = {
    isUnlocked,
    setUnlockCode,
  }

  return (
    <FeatureUnlockContext.Provider value={value}>
      {children}
    </FeatureUnlockContext.Provider>
  )
}
